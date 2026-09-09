package http_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/http/admin"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/http/client"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/http/guest"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/listener"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/notification"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/payment"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/provisioning"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/antispam"
	authUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/auth"
	billingUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/billing"
	cartUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/cart"
	formbuilderUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/formbuilder"
	orderUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/order"
	paymentUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/payment"
	staffUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	supportUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/support"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/mailer"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/plugins"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

func setupTestServer() (*httptest.Server, *memory.MockPromoRepository, *memory.MockStaffRepository) {
	jwtSecret := "test-ultra-secret-key-123456789012"
	eventBus := events.NewEventBus()

	clientRepo := memory.NewMockClientRepository()
	orderRepo := memory.NewMockOrderRepository()
	invoiceRepo := memory.NewMockInvoiceRepository()
	txnRepo := memory.NewMockTransactionRepository()
	promoRepo := memory.NewMockPromoRepository()
	supportRepo := memory.NewMockSupportRepository()
	staffRepo := memory.NewMockStaffRepository()
	productRepo := memory.NewMockProductRepository()
	taxRepo := memory.NewMockTaxRepository()

	mockMailer := mailer.NewMockMailer()
	emailService := notification.NewEmailService(mockMailer, "admin@fossbilling.org", "FOSSBilling")

	taxCalc := billingUsecase.NewTaxCalculator(taxRepo)
	promoCalc := cartUsecase.NewPromoCalculator(promoRepo)
	formService := formbuilderUsecase.NewFormbuilderService(memory.NewMockFormbuilderRepository())

	// Registries
	regRegistry := provisioning.NewRegistrarRegistry()
	regRegistry.Register("rdap", provisioning.NewMockRegistrarDriver())
	provRegistry := provisioning.NewProvisionerRegistry()

	orderService := orderUsecase.NewOrderService(orderRepo, productRepo, provRegistry, regRegistry, eventBus)
	invoiceService := billingUsecase.NewInvoiceService(invoiceRepo, clientRepo, taxCalc, plugins.NewHookManager(), eventBus)
	cartService := cartUsecase.NewCartService(promoCalc, promoRepo, orderRepo, productRepo, clientRepo, formService, taxCalc, invoiceService, eventBus)
	gatewayRegistry := paymentService.NewGatewayRegistry()
	webhookService := paymentUsecase.NewWebhookService(txnRepo, invoiceRepo, gatewayRegistry, eventBus)

	gatewayRegistry := payment.NewGatewayRegistry()
	paymentService := paymentUsecase.NewPaymentService(gatewayRegistry, invoiceRepo, clientRepo)

	supportService := supportUsecase.NewSupportService(supportRepo, clientRepo, eventBus)
	staffService := staffUsecase.NewStaffService(staffRepo, jwtSecret)

	antispamService := antispam.NewAntispamService(memory.NewMockAntispamRepository(), nil, nil, nil)
	authUc := authUsecase.NewAuthUsecase(clientRepo, antispamService, jwtSecret)
	passwordUc := authUsecase.NewPasswordUsecase(clientRepo)

	orderListener := listener.NewOrderListener(emailService, orderRepo, productRepo, clientRepo, orderService, regRegistry, provRegistry)
	eventBus.Subscribe(events.EventInvoicePaid, orderListener.HandleInvoicePaid)
	eventBus.Subscribe(events.EventOrderActivated, orderListener.HandleOrderActivated)

	guestAuthHandler := guest.NewAuthHandler(authUc)
	guestCartHandler := guest.NewCartHandler(cartService)
	guestWebhookHandler := guest.NewWebhookHandler(webhookService)

	clientProfileHandler := client.NewProfileHandler(authUc, passwordUc)
	clientOrderHandler := client.NewOrderHandler(orderService)
	clientInvoiceHandler := client.NewInvoiceHandler(invoiceRepo, clientRepo, invoiceService, paymentService)
	clientDepositHandler := client.NewDepositHandler(invoiceService, paymentService)
	clientSupportHandler := client.NewSupportHandler(supportService)

	adminStaffAuthHandler := admin.NewStaffAuthHandler(staffService)
	adminStaffMgmtHandler := admin.NewStaffManagementHandler(staffService, clientRepo, orderRepo, orderService, supportService)
	adminInvHandler := admin.NewInvoiceManagementHandler(invoiceRepo, clientRepo, invoiceService)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]string{"status": "ok"}, nil)
	})

	mux.HandleFunc("POST /api/v1/guest/auth/register", guestAuthHandler.Register)
	mux.HandleFunc("POST /api/v1/guest/auth/login", guestAuthHandler.Login)
	mux.HandleFunc("POST /api/v1/guest/cart/calculate", guestCartHandler.Calculate)
	mux.HandleFunc("POST /api/v1/guest/cart/checkout", guestCartHandler.Checkout)
	mux.HandleFunc("POST /api/v1/guest/gateways/{gateway}/webhook", guestWebhookHandler.HandleGatewayWebhook)

	mux.HandleFunc("POST /api/v1/admin/auth/login", adminStaffAuthHandler.Login)

	clientAuthMiddleware := middleware.RequireAuth(jwtSecret, "client", "admin", "superadmin")
	mux.Handle("GET /api/v1/client/profile", clientAuthMiddleware(http.HandlerFunc(clientProfileHandler.GetProfile)))
	mux.Handle("POST /api/v1/client/profile/change-password", clientAuthMiddleware(http.HandlerFunc(clientProfileHandler.ChangePassword)))
	mux.Handle("GET /api/v1/client/orders", clientAuthMiddleware(http.HandlerFunc(clientOrderHandler.ListOrders)))
	mux.Handle("GET /api/v1/client/orders/{id}", clientAuthMiddleware(http.HandlerFunc(clientOrderHandler.GetOrder)))
	mux.Handle("GET /api/v1/client/invoices", clientAuthMiddleware(http.HandlerFunc(clientInvoiceHandler.ListInvoices)))
	mux.Handle("GET /api/v1/client/invoices/{id}", clientAuthMiddleware(http.HandlerFunc(clientInvoiceHandler.GetInvoice)))
	mux.Handle("GET /api/v1/client/invoices/{id}/pdf", clientAuthMiddleware(http.HandlerFunc(clientInvoiceHandler.DownloadPDF)))
	mux.Handle("POST /api/v1/client/invoices/{id}/pay-balance", clientAuthMiddleware(http.HandlerFunc(clientInvoiceHandler.PayWithBalance)))
	mux.Handle("POST /api/v1/client/funds/deposit", clientAuthMiddleware(http.HandlerFunc(clientDepositHandler.DepositFunds)))

	mux.Handle("POST /api/v1/client/support/tickets", clientAuthMiddleware(http.HandlerFunc(clientSupportHandler.OpenTicket)))
	mux.Handle("GET /api/v1/client/support/tickets", clientAuthMiddleware(http.HandlerFunc(clientSupportHandler.ListTickets)))
	mux.Handle("GET /api/v1/client/support/tickets/{id}", clientAuthMiddleware(http.HandlerFunc(clientSupportHandler.GetTicket)))
	mux.Handle("POST /api/v1/client/support/tickets/{id}/reply", clientAuthMiddleware(http.HandlerFunc(clientSupportHandler.ReplyTicket)))
	mux.Handle("POST /api/v1/client/support/tickets/{id}/close", clientAuthMiddleware(http.HandlerFunc(clientSupportHandler.CloseTicket)))

	adminAuthMiddleware := middleware.RequireAuth(jwtSecret, "admin", "superadmin", "support", "billing")
	mux.Handle("GET /api/v1/admin/clients", adminAuthMiddleware(http.HandlerFunc(adminStaffMgmtHandler.ListClients)))
	mux.Handle("GET /api/v1/admin/orders", adminAuthMiddleware(http.HandlerFunc(adminStaffMgmtHandler.ListOrders)))
	mux.Handle("POST /api/v1/admin/orders/{id}/suspend", adminAuthMiddleware(http.HandlerFunc(adminStaffMgmtHandler.SuspendOrder)))
	mux.Handle("POST /api/v1/admin/orders/{id}/unsuspend", adminAuthMiddleware(http.HandlerFunc(adminStaffMgmtHandler.UnsuspendOrder)))
	mux.Handle("POST /api/v1/admin/orders/{id}/activate", adminAuthMiddleware(http.HandlerFunc(adminStaffMgmtHandler.ActivateOrder)))
	mux.Handle("GET /api/v1/admin/invoices", adminAuthMiddleware(http.HandlerFunc(adminInvHandler.ListInvoices)))
	mux.Handle("POST /api/v1/admin/invoices", adminAuthMiddleware(http.HandlerFunc(adminInvHandler.CreateInvoice)))
	mux.Handle("GET /api/v1/admin/support/tickets", adminAuthMiddleware(http.HandlerFunc(adminStaffMgmtHandler.ListTickets)))
	mux.Handle("POST /api/v1/admin/support/tickets/{id}/reply", adminAuthMiddleware(http.HandlerFunc(adminStaffMgmtHandler.ReplyTicket)))
	mux.Handle("GET /api/v1/admin/audit-logs", adminAuthMiddleware(http.HandlerFunc(adminStaffAuthHandler.GetAuditLogs)))

	server := httptest.NewServer(middleware.Logger(middleware.CORS(mux)))
	return server, promoRepo, staffRepo
}
