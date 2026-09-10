package integration_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/http/admin"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/http/guest"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/listener"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/provisioning"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/auth"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/billing"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/cart"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/order"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/support"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
	"net/http/httptest"
)

func TestBug01_NegativePriceValidation(t *testing.T) {
	ctx := context.Background(); ir, cr, pr, pmr, or, cor := memory.NewMockInvoiceRepository(), memory.NewMockClientRepository(), memory.NewMockProductRepository(), memory.NewMockPromoRepository(), memory.NewMockOrderRepository(), memory.NewMockCompanyRepository()
	is := billing.NewInvoiceService(ir, cr, cor, nil, nil); cs := cart.NewCartService(nil, pmr, or, pr, cr, nil, nil, is, nil, nil)
	_ = pr.Create(ctx, &domain.Product{ID: 101, Title: "T"}); bad := &cart.Cart{ClientID: 1, Items: []cart.CartItem{{ProductID: 101, Price: -1000000, Quantity: 1}}}
	if _, err := cs.Checkout(ctx, bad, ""); err == nil { t.Fatal("BUG-01") }
}

func TestBug05_IDORInvoicePayment(t *testing.T) {
	ctx := context.Background(); ir, cr, cor := memory.NewMockInvoiceRepository(), memory.NewMockClientRepository(), memory.NewMockCompanyRepository(); is := billing.NewInvoiceService(ir, cr, cor, nil, nil)
	_ = ir.Create(ctx, &domain.Invoice{ID: 100, ClientID: 1, Total: 500000, Status: domain.InvoiceStatusUnpaid}, nil)
	_ = cr.AddBalanceTransaction(ctx, &domain.ClientBalance{ClientID: 2, Amount: 1000000, Type: domain.BalanceTypeCredit})
	if _, err := is.PayWithBalance(ctx, 2, 100); err == nil { t.Fatal("BUG-05") }
}

func TestBug08_StoredXSSProfile(t *testing.T) {
	ctx := context.Background(); cr := memory.NewMockClientRepository(); uc := auth.NewAuthUsecase(cr, nil, "s", "F", nil)
	_ = cr.Create(ctx, &domain.Client{ID: 1, Email: "t@t.com"}); xss := "<script>alert(1)</script>"
	_, _ = uc.UpdateProfile(ctx, 1, auth.UpdateProfileDTO{FirstName: xss})
	up, _ := cr.GetByID(ctx, 1); if strings.Contains(up.FirstName, "<script>") { t.Fatal("BUG-08") }
}

func TestBug09_StoredXSSSupport(t *testing.T) {
	ctx := context.Background(); sr, cr := memory.NewMockSupportRepository(), memory.NewMockClientRepository(); su := support.NewSupportService(sr, cr, nil)
	_ = cr.Create(ctx, &domain.Client{ID: 1}); xss := "<img src=x onerror=alert(1)>"
	tk, _ := su.OpenTicket(ctx, support.CreateTicketDTO{ClientID: 1, Subject: "h", Message: xss}); ms, _ := sr.GetMessages(ctx, tk.ID)
	if strings.Contains(ms[0].Content, "onerror") { t.Fatal("BUG-09") }
}

func TestBug15_ClientIDEnumeration(t *testing.T) {
	cr, pr := memory.NewMockClientRepository(), memory.NewMockProductRepository(); tr := memory.NewMockTaxRepository(); _ = cr.Create(context.Background(), &domain.Client{ID: 55, Country: "ID"})
	cs := cart.NewCartService(nil, nil, nil, pr, cr, nil, billing.NewTaxCalculator(tr), nil, nil, nil); h := guest.NewCartHandler(cs)
	_ = pr.Create(context.Background(), &domain.Product{ID: 1, Title: "P"}); req := httptest.NewRequest("POST", "/", strings.NewReader(`{"client_id": 55, "items": [{"product_id": 1, "price": 100, "quantity": 1}]}`)); rr := httptest.NewRecorder()
	h.Calculate(rr, req); if strings.Contains(rr.Body.String(), `"client_id":55`) { t.Fatal("BUG-15") }
}

func TestBug16_17_PrivilegeEscalation(t *testing.T) {
	ctx := context.Background(); sr, or := memory.NewMockStaffRepository(), memory.NewMockOrderRepository()
	su, ou := staff.NewStaffService(sr, "s", "F"), order.NewOrderService(or, nil, nil, nil, nil)
	h := admin.NewStaffManagementHandler(su, nil, or, ou, nil)
	_ = sr.CreateGroup(ctx, &domain.AdminGroup{ID: 2, Name: "R", Permissions: map[string][]string{"orders": {"read"}}})
	_ = sr.Create(ctx, &domain.Staff{ID: 10, GroupID: 2, Role: domain.StaffRoleSupport, Status: "active"})
	req := httptest.NewRequest("POST", "/", nil); req.SetPathValue("id", "1")
	req = req.WithContext(middleware.WithClientID(req.Context(), 10)); rr := httptest.NewRecorder(); h.UnsuspendOrder(rr, req)
	if rr.Code == 200 { t.Fatal("BUG-16/17") }
}

func TestBug20_ProvisioningFailureHandling(t *testing.T) {
	ctx := context.Background(); or, pr, cr := memory.NewMockOrderRepository(), memory.NewMockProductRepository(), memory.NewMockClientRepository()
	pReg := provisioning.NewProvisionerRegistry(); pReg.Register("fail", &MockFailingProv{})
	ou := order.NewOrderService(or, pr, pReg, nil, nil); ls := listener.NewOrderListener(nil, or, pr, cr, ou, nil, pReg)
	_ = pr.Create(ctx, &domain.Product{ID: 101, Title: "T", Type: domain.ProductTypeHosting})
	_ = or.Create(ctx, &domain.Order{ID: 1, ProductID: 101, Status: domain.OrderStatusActive, Config: []byte(`{"server_type":"fail"}`)})
	_ = ls.HandleOrderActivated(ctx, events.Event{Payload: domain.OrderActivatedPayload{OrderID: 1}})
	fo, _ := or.GetByID(ctx, 1); if fo.Status == domain.OrderStatusActive { t.Fatal("BUG-20") }
}

type MockFailingProv struct{}
func (m *MockFailingProv) Type() domain.ProductType { return domain.ProductTypeHosting }
func (m *MockFailingProv) Create(ctx context.Context, o *domain.Order) (*domain.ProvisionResult, error) { return nil, fmt.Errorf("FAIL") }
func (m *MockFailingProv) Suspend(ctx context.Context, o *domain.Order, r string) error { return nil }
func (m *MockFailingProv) Unsuspend(ctx context.Context, o *domain.Order) error { return nil }
func (m *MockFailingProv) Renew(ctx context.Context, o *domain.Order) error { return nil }
func (m *MockFailingProv) Terminate(ctx context.Context, o *domain.Order) error { return nil }
func (m *MockFailingProv) Sync(ctx context.Context, o *domain.Order) (*domain.ServiceStatus, error) { return &domain.ServiceStatus{IsActive: true}, nil }
func (m *MockFailingProv) ChangePassword(ctx context.Context, o *domain.Order, p string) error { return nil }
func (m *MockFailingProv) TestConnection(ctx context.Context) error { return nil }
