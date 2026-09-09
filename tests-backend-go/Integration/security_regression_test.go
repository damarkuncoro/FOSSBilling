package Integration

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
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
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
)

// TestBug01_NegativePriceValidation (Resolved)
func TestBug01_NegativePriceValidation(t *testing.T) {
	ctx := context.Background()
	invRepo := memory.NewMockInvoiceRepository()
	clientRepo := memory.NewMockClientRepository()
	prodRepo := memory.NewMockProductRepository()
	promoRepo := memory.NewMockPromoRepository()
	orderRepo := memory.NewMockOrderRepository()

	invService := billing.NewInvoiceService(invRepo, clientRepo, nil, nil)
	cartService := cart.NewCartService(cart.NewPromoCalculator(promoRepo), promoRepo, orderRepo, prodRepo, clientRepo, nil, nil, invService, events.NewEventBus())

	_ = prodRepo.Create(ctx, &domain.Product{ID: 101, Name: "Test"})

	badCart := &cart.Cart{
		ClientID: 1,
		Items: []cart.CartItem{
			{ProductID: 101, Price: decimal.FromFloat(-100.00), Quantity: 1},
		},
	}

	_, err := cartService.Checkout(ctx, badCart)
	if err == nil {
		t.Fatal("SECURITY FAILURE: Allowed checkout with negative price (BUG-01)")
	}
}

// TestBug05_IDORInvoicePayment (Resolved)
func TestBug05_IDORInvoicePayment(t *testing.T) {
	ctx := context.Background()
	invRepo := memory.NewMockInvoiceRepository()
	clientRepo := memory.NewMockClientRepository()
	invService := billing.NewInvoiceService(invRepo, clientRepo, nil, nil)

	userA_ID, userB_ID := int64(1), int64(2)
	invA := &domain.Invoice{ID: 100, ClientID: userA_ID, Total: decimal.FromFloat(50.00), Status: domain.InvoiceStatusUnpaid}
	_ = invRepo.Create(ctx, invA, nil)

	_ = clientRepo.AddBalanceTransaction(ctx, &domain.ClientBalance{ClientID: userB_ID, Amount: decimal.FromFloat(100.00), Type: domain.BalanceTypeCredit})

	_, err := invService.PayWithBalance(ctx, userB_ID, invA.ID)
	if err == nil {
		t.Fatal("SECURITY FAILURE: User B paid User A's invoice (BUG-05 IDOR)")
	}
}

// TestBug06_PriceInjectionZeroID (Resolved)
func TestBug06_PriceInjectionZeroID(t *testing.T) {
	ctx := context.Background()
	cartService := cart.NewCartService(nil, nil, nil, nil, nil, nil, nil, nil, nil)

	badCart := &cart.Cart{
		ClientID: 1,
		Items: []cart.CartItem{
			{ProductID: 0, Price: decimal.FromFloat(1.00), Quantity: 1},
		},
	}

	_, err := cartService.Checkout(ctx, badCart)
	if err == nil || !strings.Contains(err.Error(), "invalid product ID") {
		t.Fatal("SECURITY FAILURE: Allowed checkout with ProductID=0 (BUG-06)")
	}
}

// TestBug07_PromoMaxUsesRaceCondition (Resolved)
func TestBug07_PromoMaxUsesRaceCondition(t *testing.T) {
	ctx := context.Background()
	promoRepo := memory.NewMockPromoRepository()
	prodRepo := memory.NewMockProductRepository()
	orderRepo := memory.NewMockOrderRepository()
	invRepo := memory.NewMockInvoiceRepository()
	clientRepo := memory.NewMockClientRepository()

	invService := billing.NewInvoiceService(invRepo, clientRepo, nil, nil)
	cartService := cart.NewCartService(cart.NewPromoCalculator(promoRepo), promoRepo, orderRepo, prodRepo, clientRepo, nil, nil, invService, events.NewEventBus())

	_ = promoRepo.Create(ctx, &domain.Promo{ID: 1, Code: "LIMIT1", MaxUses: 1, Active: true, Value: 100})
	_ = prodRepo.Create(ctx, &domain.Product{ID: 101, PriceMonthly: 1000})

	var wg sync.WaitGroup
	numReqs := 5
	results := make(chan error, numReqs)

	for i := 0; i < numReqs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := cartService.Checkout(ctx, &cart.Cart{ClientID: 1, PromoCode: "LIMIT1", Items: []cart.CartItem{{ProductID: 101, Quantity: 1}}})
			results <- err
		}()
	}
	wg.Wait()
	close(results)

	successCount := 0
	for err := range results {
		if err == nil {
			successCount++
		}
	}

	if successCount > 1 {
		t.Fatalf("INTEGRITY FAILURE: Promo limit bypassed! %d users succeeded (BUG-07)", successCount)
	}
}

// TestBug08_StoredXSSProfile (Resolved)
func TestBug08_StoredXSSProfile(t *testing.T) {
	ctx := context.Background()
	clientRepo := memory.NewMockClientRepository()
	authUc := auth.NewAuthUsecase(clientRepo, nil, "secret")

	_ = clientRepo.Create(ctx, &domain.Client{ID: 1, Email: "test@test.com"})

	xss := "<script>alert(1)</script>"
	_, _ = authUc.UpdateProfile(ctx, 1, auth.UpdateProfileDTO{FirstName: xss})

	updated, _ := clientRepo.GetByID(ctx, 1)
	if strings.Contains(updated.FirstName, "<script>") {
		t.Fatal("SECURITY FAILURE: Stored XSS in Profile fields (BUG-08)")
	}
}

// TestBug09_StoredXSSSupport (Resolved)
func TestBug09_StoredXSSSupport(t *testing.T) {
	ctx := context.Background()
	supportRepo := memory.NewMockSupportRepository()
	clientRepo := memory.NewMockClientRepository()
	supportUc := support.NewSupportService(supportRepo, clientRepo)

	_ = clientRepo.Create(ctx, &domain.Client{ID: 1})
	xss := "<img src=x onerror=alert(1)>"
	ticket, _ := supportUc.OpenTicket(ctx, support.CreateTicketDTO{ClientID: 1, Subject: "help", Message: xss})

	msgs, _ := supportRepo.GetMessages(ctx, ticket.ID)
	if strings.Contains(msgs[0].Content, "onerror") {
		t.Fatal("SECURITY FAILURE: Stored XSS in Support Ticket (BUG-09)")
	}
}

// TestBug12_WeakPasswordPolicy (Resolved)
func TestBug12_WeakPasswordPolicy(t *testing.T) {
	ctx := context.Background()
	clientRepo := memory.NewMockClientRepository()
	authUc := auth.NewAuthUsecase(clientRepo, nil, "secret")

	_, _, err := authUc.Register(ctx, auth.RegisterDTO{Email: "weak@test.com", Password: "123", FirstName: "A", LastName: "B"}, "1.1.1.1")
	if err == nil {
		t.Fatal("SECURITY FAILURE: Weak password accepted (BUG-12)")
	}
}

// TestBug15_ClientIDEnumeration (Resolved)
func TestBug15_ClientIDEnumeration(t *testing.T) {
	clientRepo := memory.NewMockClientRepository()
	taxRepo := memory.NewMockTaxRepository()
	_ = clientRepo.Create(context.Background(), &domain.Client{ID: 55, Country: "ID"})

	taxCalc := billing.NewTaxCalculator(taxRepo)
	cartService := cart.NewCartService(nil, nil, nil, nil, clientRepo, nil, taxCalc, nil, nil)
	handler := guest.NewCartHandler(cartService)

	// Simulate guest request with injected client_id
	body := `{"client_id": 55, "items": [{"product_id": 1, "price": 100, "quantity": 1}]}`
	req := httptest.NewRequest("POST", "/api/v1/guest/cart/calculate", strings.NewReader(body))
	rr := httptest.NewRecorder()

	handler.Calculate(rr, req)

	if strings.Contains(rr.Body.String(), `"client_id":55`) {
		t.Fatal("SECURITY FAILURE: ClientID leaked/reflected in Guest API (BUG-15)")
	}
}

// TestBug16_17_PrivilegeEscalation (Resolved)
func TestBug16_17_PrivilegeEscalation(t *testing.T) {
	ctx := context.Background()
	staffRepo := memory.NewMockStaffRepository()
	orderRepo := memory.NewMockOrderRepository()

	staffUc := staff.NewStaffService(staffRepo, "secret")
	orderUc := order.NewOrderService(orderRepo, nil, provisioning.NewProvisionerRegistry(), provisioning.NewRegistrarRegistry(), nil)
	handler := admin.NewStaffManagementHandler(staffUc, nil, orderRepo, orderUc, nil)

	// Group with only 'read' permission
	_ = staffRepo.CreateGroup(ctx, &domain.AdminGroup{ID: 2, Permissions: map[string][]string{"orders": {"read"}}})
	_ = staffRepo.Create(ctx, &domain.Staff{ID: 10, GroupID: 2, Role: domain.StaffRoleSupport})

	req := httptest.NewRequest("POST", "/api/v1/admin/orders/1/unsuspend", nil)
	req.SetPathValue("id", "1")
	req = req.WithContext(middleware.WithClientID(req.Context(), 10))

	rr := httptest.NewRecorder()
	handler.UnsuspendOrder(rr, req)

	if rr.Code == http.StatusOK {
		t.Fatal("SECURITY FAILURE: Staff with read-only access performed write action (BUG-16/17)")
	}
}

// TestBug18_PromoOncePerClientRaceCondition (New Bug Hunt)
func TestBug18_PromoOncePerClientRaceCondition(t *testing.T) {
	ctx := context.Background()
	promoRepo := memory.NewMockPromoRepository()
	prodRepo := memory.NewMockProductRepository()
	orderRepo := memory.NewMockOrderRepository()
	invRepo := memory.NewMockInvoiceRepository()
	clientRepo := memory.NewMockClientRepository()

	invService := billing.NewInvoiceService(invRepo, clientRepo, nil, nil)
	cartService := cart.NewCartService(cart.NewPromoCalculator(promoRepo), promoRepo, orderRepo, prodRepo, clientRepo, nil, nil, invService, events.NewEventBus())

	// Promo yang hanya boleh digunakan 1x oleh tiap user
	_ = promoRepo.Create(ctx, &domain.Promo{ID: 1, Code: "ONCE", OncePerClient: true, Active: true, Value: 100})
	_ = prodRepo.Create(ctx, &domain.Product{ID: 101, PriceMonthly: 1000})
	_ = clientRepo.Create(ctx, &domain.Client{ID: 1})

	var wg sync.WaitGroup
	numReqs := 5
	results := make(chan error, numReqs)

	// Simulasi 5 request checkout secara bersamaan oleh 1 User yang sama
	for i := 0; i < numReqs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := cartService.Checkout(ctx, &cart.Cart{ClientID: 1, PromoCode: "ONCE", Items: []cart.CartItem{{ProductID: 101, Quantity: 1}}})
			results <- err
		}()
	}
	wg.Wait()
	close(results)

	successCount := 0
	for err := range results {
		if err == nil {
			successCount++
		}
	}

	// Seharusnya hanya 1 yang sukses, sisanya error "Promo already used"
	if successCount > 1 {
		t.Fatalf("INTEGRITY FAILURE: User used 'once_per_client' promo %d times! (BUG-18)", successCount)
	}
}

// TestBug20_ProvisioningFailureHandling memastikan kegagalan provisioning mengubah status order kembali
func TestBug20_ProvisioningFailureHandling(t *testing.T) {
	ctx := context.Background()
	orderRepo := memory.NewMockOrderRepository()
	productRepo := memory.NewMockProductRepository()
	clientRepo := memory.NewMockClientRepository()

	// Mock Provisioner that ALWAYS fails
	provRegistry := provisioning.NewProvisionerRegistry()
	provRegistry.Register("fail_driver", &MockProvisionerFailure{})

	orderUc := order.NewOrderService(orderRepo, productRepo, provRegistry, nil, nil)
	listener := listener.NewOrderListener(nil, orderRepo, productRepo, clientRepo, orderUc, nil, provRegistry)

	// Setup: Order for a product that uses the failing driver
	_ = productRepo.Create(ctx, &domain.Product{ID: 101, Type: domain.ProductTypeHosting})
	ord := &domain.Order{ID: 1, ProductID: 101, Status: domain.OrderStatusActive, Config: []byte(`{"server_type":"fail_driver"}`)}
	_ = orderRepo.Create(ctx, ord)

	// Simulate event trigger
	_ = listener.HandleOrderActivated(ctx, events.Event{Payload: domain.OrderActivatedPayload{OrderID: 1}})

	// Check final status
	finalOrd, _ := orderRepo.GetByID(ctx, 1)
	if finalOrd.Status == domain.OrderStatusActive {
		t.Fatal("LOGIC FAILURE: Order remained ACTIVE even though provisioning failed (BUG-20)")
	}
}

type MockProvisioner struct{ PType domain.ProductType }
func (m *MockProvisioner) Type() domain.ProductType { return m.PType }
func (m *MockProvisioner) Create(ctx context.Context, o *domain.Order) (*domain.ProvisionResult, error) {
	return &domain.ProvisionResult{Success: true}, nil
}
func (m *MockProvisioner) Suspend(ctx context.Context, o *domain.Order, r string) error { return nil }
func (m *MockProvisioner) Unsuspend(ctx context.Context, o *domain.Order) error           { return nil }
func (m *MockProvisioner) Renew(ctx context.Context, o *domain.Order) error               { return nil }
func (m *MockProvisioner) Terminate(ctx context.Context, o *domain.Order) error           { return nil }
func (m *MockProvisioner) Sync(ctx context.Context, o *domain.Order) (*domain.ServiceStatus, error) {
	return &domain.ServiceStatus{IsActive: true}, nil
}
func (m *MockProvisioner) ChangePassword(ctx context.Context, o *domain.Order, p string) error { return nil }
func (m *MockProvisioner) TestConnection(ctx context.Context) error                          { return nil }

type MockProvisionerFailure struct{ MockProvisioner }
func (m *MockProvisionerFailure) Create(ctx context.Context, o *domain.Order) (*domain.ProvisionResult, error) {
	return nil, fmt.Errorf("CONNECTION REFUSED")
}
func (m *MockProvisionerFailure) Type() domain.ProductType { return domain.ProductTypeHosting }
