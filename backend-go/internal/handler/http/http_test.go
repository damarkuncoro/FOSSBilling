package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/handler/http/admin"
	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/handler/http/client"
	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/handler/http/guest"
	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/repository/memory"
	authUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/internal/usecase/auth"
	billingUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/internal/usecase/billing"
	cartUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/internal/usecase/cart"
	orderUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/internal/usecase/order"
	paymentUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/internal/usecase/payment"
	staffUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/internal/usecase/staff"
	supportUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/internal/usecase/support"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/auth"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

func setupTestServer() (*httptest.Server, *memory.MockPromoRepository, *memory.MockStaffRepository) {
	jwtSecret := "test-ultra-secret-key-123456789012"

	clientRepo := memory.NewMockClientRepository()
	orderRepo := memory.NewMockOrderRepository()
	invoiceRepo := memory.NewMockInvoiceRepository()
	txnRepo := memory.NewMockTransactionRepository()
	promoRepo := memory.NewMockPromoRepository()
	supportRepo := memory.NewMockSupportRepository()
	staffRepo := memory.NewMockStaffRepository()

	taxCalc := billingUsecase.NewTaxCalculator([]billingUsecase.TaxRule{
		{Name: "PPN", Country: "ID", Rate: 11.0},
	})
	promoCalc := cartUsecase.NewPromoCalculator(promoRepo)
	orderService := orderUsecase.NewOrderService(orderRepo)
	invoiceService := billingUsecase.NewInvoiceService(invoiceRepo, clientRepo, taxCalc)
	cartService := cartUsecase.NewCartService(promoCalc, promoRepo, orderRepo, invoiceService)
	webhookService := paymentUsecase.NewWebhookService(txnRepo, invoiceRepo, orderService, orderRepo)
	supportService := supportUsecase.NewSupportService(supportRepo, clientRepo)
	staffService := staffUsecase.NewStaffService(staffRepo, jwtSecret)
	authUc := authUsecase.NewAuthUsecase(clientRepo, jwtSecret)

	guestAuthHandler := guest.NewAuthHandler(authUc)
	guestCartHandler := guest.NewCartHandler(cartService)
	guestWebhookHandler := guest.NewWebhookHandler(webhookService)

	clientProfileHandler := client.NewProfileHandler(authUc)
	clientOrderHandler := client.NewOrderHandler(orderRepo)
	clientInvoiceHandler := client.NewInvoiceHandler(invoiceRepo, clientRepo, invoiceService)
	clientSupportHandler := client.NewSupportHandler(supportService)


	adminStaffHandler := admin.NewStaffHandler(staffService, clientRepo, orderRepo, orderService, supportService)

	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]string{"status": "ok"}, nil)
	})

	// Guest
	mux.HandleFunc("POST /api/v1/guest/auth/register", guestAuthHandler.Register)
	mux.HandleFunc("POST /api/v1/guest/auth/login", guestAuthHandler.Login)
	mux.HandleFunc("POST /api/v1/guest/cart/calculate", guestCartHandler.Calculate)
	mux.HandleFunc("POST /api/v1/guest/cart/checkout", guestCartHandler.Checkout)
	mux.HandleFunc("POST /api/v1/guest/gateways/{gateway}/webhook", guestWebhookHandler.HandleGatewayWebhook)

	// Admin Auth
	mux.HandleFunc("POST /api/v1/admin/auth/login", adminStaffHandler.Login)

	// Client protected
	clientAuthMiddleware := middleware.RequireAuth(jwtSecret, "client", "admin", "superadmin")
	mux.Handle("GET /api/v1/client/profile", clientAuthMiddleware(http.HandlerFunc(clientProfileHandler.GetProfile)))
	mux.Handle("GET /api/v1/client/orders", clientAuthMiddleware(http.HandlerFunc(clientOrderHandler.ListOrders)))
	mux.Handle("GET /api/v1/client/orders/{id}", clientAuthMiddleware(http.HandlerFunc(clientOrderHandler.GetOrder)))
	mux.Handle("GET /api/v1/client/invoices", clientAuthMiddleware(http.HandlerFunc(clientInvoiceHandler.ListInvoices)))
	mux.Handle("GET /api/v1/client/invoices/{id}", clientAuthMiddleware(http.HandlerFunc(clientInvoiceHandler.GetInvoice)))
	mux.Handle("GET /api/v1/client/invoices/{id}/pdf", clientAuthMiddleware(http.HandlerFunc(clientInvoiceHandler.DownloadPDF)))
	mux.Handle("POST /api/v1/client/invoices/{id}/pay-balance", clientAuthMiddleware(http.HandlerFunc(clientInvoiceHandler.PayWithBalance)))

	mux.Handle("POST /api/v1/client/support/tickets", clientAuthMiddleware(http.HandlerFunc(clientSupportHandler.OpenTicket)))
	mux.Handle("GET /api/v1/client/support/tickets", clientAuthMiddleware(http.HandlerFunc(clientSupportHandler.ListTickets)))
	mux.Handle("GET /api/v1/client/support/tickets/{id}", clientAuthMiddleware(http.HandlerFunc(clientSupportHandler.GetTicket)))
	mux.Handle("POST /api/v1/client/support/tickets/{id}/reply", clientAuthMiddleware(http.HandlerFunc(clientSupportHandler.ReplyTicket)))
	mux.Handle("POST /api/v1/client/support/tickets/{id}/close", clientAuthMiddleware(http.HandlerFunc(clientSupportHandler.CloseTicket)))

	// Admin protected
	adminAuthMiddleware := middleware.RequireAuth(jwtSecret, "admin", "superadmin", "support", "billing")
	mux.Handle("GET /api/v1/admin/clients", adminAuthMiddleware(http.HandlerFunc(adminStaffHandler.ListClients)))
	mux.Handle("GET /api/v1/admin/orders", adminAuthMiddleware(http.HandlerFunc(adminStaffHandler.ListOrders)))
	mux.Handle("POST /api/v1/admin/orders/{id}/suspend", adminAuthMiddleware(http.HandlerFunc(adminStaffHandler.SuspendOrder)))
	mux.Handle("POST /api/v1/admin/orders/{id}/unsuspend", adminAuthMiddleware(http.HandlerFunc(adminStaffHandler.UnsuspendOrder)))
	mux.Handle("POST /api/v1/admin/orders/{id}/activate", adminAuthMiddleware(http.HandlerFunc(adminStaffHandler.ActivateOrder)))
	mux.Handle("GET /api/v1/admin/support/tickets", adminAuthMiddleware(http.HandlerFunc(adminStaffHandler.ListTickets)))
	mux.Handle("POST /api/v1/admin/support/tickets/{id}/reply", adminAuthMiddleware(http.HandlerFunc(adminStaffHandler.ReplyTicket)))
	mux.Handle("GET /api/v1/admin/audit-logs", adminAuthMiddleware(http.HandlerFunc(adminStaffHandler.GetAuditLogs)))

	server := httptest.NewServer(middleware.Logger(middleware.CORS(mux)))
	return server, promoRepo, staffRepo
}

func TestHTTP_FullLifecycleJourney(t *testing.T) {
	ts, promoRepo, staffRepo := setupTestServer()
	defer ts.Close()

	ctx := context.Background()

	// Setup Promo
	_ = promoRepo.Create(ctx, &domain.Promo{
		Code:        "MERDEKA20",
		Description: "20% discount",
		Type:        domain.PromoTypePercentage,
		Value:       200000,
		Active:      true,
	})

	// Setup Admin Group and Staff
	group := &domain.AdminGroup{
		Name: "Administrators",
		Permissions: map[string][]string{
			"clients": {"*"},
			"orders":  {"*"},
			"support": {"*"},
			"system":  {"*"},
		},
	}
	_ = staffRepo.CreateGroup(ctx, group)
	passHash, _ := auth.HashPassword("SuperSecretAdmin123!")
	_ = staffRepo.Create(ctx, &domain.Staff{
		GroupID:      group.ID,
		Email:        "admin@fossbilling.org",
		PasswordHash: passHash,
		Name:         "Super Admin",
		Role:         domain.StaffRoleSuperAdmin,
		Status:       "active",
	})

	// 1. Health Check
	res, err := http.Get(ts.URL + "/health")
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("Health check failed: %v, status: %d", err, res.StatusCode)
	}

	// 2. Client Registration
	regBody, _ := json.Marshal(map[string]interface{}{
		"email":      "ahmad.dhani@example.com",
		"password":   "Password123!",
		"first_name": "Ahmad",
		"last_name":  "Dhani",
		"country":    "ID",
		"currency":   "USD",
	})
	regResp, err := http.Post(ts.URL+"/api/v1/guest/auth/register", "application/json", bytes.NewBuffer(regBody))
	if err != nil || regResp.StatusCode != http.StatusCreated {
		t.Fatalf("Register failed: %v, status: %d", err, regResp.StatusCode)
	}

	var regData struct {
		Data struct {
			Token  string `json:"token"`
			Client struct {
				ID int64 `json:"id"`
			} `json:"client"`
		} `json:"data"`
	}
	_ = json.NewDecoder(regResp.Body).Decode(&regData)
	clientToken := regData.Data.Token
	clientID := regData.Data.Client.ID
	if clientToken == "" || clientID == 0 {
		t.Fatalf("Invalid register data: token=%s, clientID=%d", clientToken, clientID)
	}

	// 3. Client Cart Calculation & Checkout
	cartBody, _ := json.Marshal(map[string]interface{}{
		"client_id":  clientID,
		"promo_code": "MERDEKA20",
		"items": []map[string]interface{}{
			{
				"product_id": 1,
				"title":      "Cloud Hosting Pro",
				"period":     "1M",
				"price":      1000000, // $100.00
				"quantity":   1,
			},
		},
	})
	checkoutResp, err := http.Post(ts.URL+"/api/v1/guest/cart/checkout", "application/json", bytes.NewBuffer(cartBody))
	if err != nil || checkoutResp.StatusCode != http.StatusCreated {
		t.Fatalf("Checkout failed: %v, status: %d", err, checkoutResp.StatusCode)
	}

	var checkoutData struct {
		Data struct {
			Orders  []domain.Order `json:"orders"`
			Invoice domain.Invoice `json:"invoice"`
		} `json:"data"`
	}
	_ = json.NewDecoder(checkoutResp.Body).Decode(&checkoutData)
	if len(checkoutData.Data.Orders) == 0 {
		t.Fatal("Expected orders created in checkout")
	}
	invoiceID := checkoutData.Data.Invoice.ID
	orderID := checkoutData.Data.Orders[0].ID

	// 4. Client Get Orders
	req, _ := http.NewRequest("GET", ts.URL+"/api/v1/client/orders", nil)
	req.Header.Set("Authorization", "Bearer "+clientToken)
	ordersResp, err := http.DefaultClient.Do(req)
	if err != nil || ordersResp.StatusCode != http.StatusOK {
		t.Fatalf("Get client orders failed: %v, status: %d", err, ordersResp.StatusCode)
	}

	// 5. Payment Webhook (Midtrans / Stripe)
	webhookBody, _ := json.Marshal(map[string]interface{}{
		"gateway_id": "midtrans",
		"txn_id":     fmt.Sprintf("TRX-%d", time.Now().UnixNano()),
		"invoice_id": invoiceID,
		"amount":     checkoutData.Data.Invoice.Total,
		"currency":   "USD",
	})
	webhookResp, err := http.Post(ts.URL+"/api/v1/guest/gateways/midtrans/webhook", "application/json", bytes.NewBuffer(webhookBody))
	if err != nil || webhookResp.StatusCode != http.StatusOK {
		t.Fatalf("Webhook failed: %v, status: %d", err, webhookResp.StatusCode)
	}

	// Verify Order is now active
	req, _ = http.NewRequest("GET", fmt.Sprintf("%s/api/v1/client/orders/%d", ts.URL, orderID), nil)
	req.Header.Set("Authorization", "Bearer "+clientToken)
	getOrdResp, _ := http.DefaultClient.Do(req)
	var ordData struct {
		Data domain.Order `json:"data"`
	}
	_ = json.NewDecoder(getOrdResp.Body).Decode(&ordData)
	if ordData.Data.Status != domain.OrderStatusActive {
		t.Errorf("Expected order status active after payment webhook, got: %s", ordData.Data.Status)
	}

	// 6. Support Ticket Lifecycle
	ticketBody, _ := json.Marshal(map[string]interface{}{
		"helpdesk_id": 1,
		"subject":     "Need help with SSH keys",
		"message":     "How do I add my public key?",
		"priority":    "high",
	})
	req, _ = http.NewRequest("POST", ts.URL+"/api/v1/client/support/tickets", bytes.NewBuffer(ticketBody))
	req.Header.Set("Authorization", "Bearer "+clientToken)
	ticketResp, err := http.DefaultClient.Do(req)
	if err != nil || ticketResp.StatusCode != http.StatusCreated {
		t.Fatalf("Open ticket failed: %v, status: %d", err, ticketResp.StatusCode)
	}

	var ticketData struct {
		Data domain.Ticket `json:"data"`
	}
	_ = json.NewDecoder(ticketResp.Body).Decode(&ticketData)
	ticketID := ticketData.Data.ID

	// 7. Admin Staff Login & Ticket Reply
	adminLoginBody, _ := json.Marshal(map[string]string{
		"email":    "admin@fossbilling.org",
		"password": "SuperSecretAdmin123!",
	})
	adminLoginResp, err := http.Post(ts.URL+"/api/v1/admin/auth/login", "application/json", bytes.NewBuffer(adminLoginBody))
	if err != nil || adminLoginResp.StatusCode != http.StatusOK {
		t.Fatalf("Admin login failed: %v, status: %d", err, adminLoginResp.StatusCode)
	}
	var adminAuthData struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	_ = json.NewDecoder(adminLoginResp.Body).Decode(&adminAuthData)
	adminToken := adminAuthData.Data.Token

	// Staff replies to ticket
	staffReplyBody, _ := json.Marshal(map[string]string{
		"message": "You can paste your SSH public key in the cPanel security section.",
	})
	req, _ = http.NewRequest("POST", fmt.Sprintf("%s/api/v1/admin/support/tickets/%d/reply", ts.URL, ticketID), bytes.NewBuffer(staffReplyBody))
	req.Header.Set("Authorization", "Bearer "+adminToken)
	staffReplyResp, err := http.DefaultClient.Do(req)
	if err != nil || staffReplyResp.StatusCode != http.StatusCreated {
		t.Fatalf("Staff reply failed: %v, status: %d", err, staffReplyResp.StatusCode)
	}

	// 8. Client closes ticket
	req, _ = http.NewRequest("POST", fmt.Sprintf("%s/api/v1/client/support/tickets/%d/close", ts.URL, ticketID), nil)
	req.Header.Set("Authorization", "Bearer "+clientToken)
	closeResp, err := http.DefaultClient.Do(req)
	if err != nil || closeResp.StatusCode != http.StatusOK {
		t.Fatalf("Close ticket failed: %v, status: %d", err, closeResp.StatusCode)
	}

	// 9. Admin checks Audit Logs
	req, _ = http.NewRequest("GET", ts.URL+"/api/v1/admin/audit-logs", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	auditResp, err := http.DefaultClient.Do(req)
	if err != nil || auditResp.StatusCode != http.StatusOK {
		t.Fatalf("Get audit logs failed: %v, status: %d", err, auditResp.StatusCode)
	}

	var auditData struct {
		Data []domain.AuditLog `json:"data"`
	}
	_ = json.NewDecoder(auditResp.Body).Decode(&auditData)
	if len(auditData.Data) == 0 {
		t.Error("Expected at least 1 audit log entry")
	}
}
