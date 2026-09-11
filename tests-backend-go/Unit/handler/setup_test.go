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
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/payment/gateways"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/provisioning"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/activity"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/affiliate"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/antispam"
	authUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/auth"
	billingUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/billing"
	cartUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/cart"
	formbuilderUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/formbuilder"
	orderUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/order"
	paymentUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/payment"
	staffUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	statsUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/stats"
	supportUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/support"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/cache"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/mailer"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/plugins"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

func setupTestServer() (*httptest.Server, *memory.MockPromoRepository, *memory.MockStaffRepository, *memory.MockProductRepository) {
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
	companyRepo := memory.NewMockCompanyRepository()
	tplRepo := memory.NewMockEmailTemplateRepository()

	mockMailer := mailer.NewMockMailer()
	emailService := notification.NewEmailService(mockMailer, tplRepo, "admin@fossbilling.org", "FOSSBilling")

	taxCalc := billingUsecase.NewTaxCalculator(taxRepo)
	promoCalc := cartUsecase.NewPromoCalculator(promoRepo)
	formService := formbuilderUsecase.NewFormbuilderService(memory.NewMockFormbuilderRepository())

	regRegistry := provisioning.NewRegistrarRegistry()
	regRegistry.Register("rdap", provisioning.NewMockRegistrarDriver())
	provRegistry := provisioning.NewProvisionerRegistry()

	orderService := orderUsecase.NewOrderService(orderRepo, productRepo, provRegistry, regRegistry, eventBus)
	invoiceService := billingUsecase.NewInvoiceService(invoiceRepo, clientRepo, companyRepo, taxCalc, plugins.NewHookManager(), eventBus)
	cartService := cartUsecase.NewCartService(promoCalc, promoRepo, orderRepo, productRepo, clientRepo, formService, taxCalc, invoiceService, nil, eventBus)

	gatewayRegistry := payment.NewGatewayRegistry()
	gatewayRegistry.Register(gateways.NewMidtransGateway("", "", false))
	gatewayRegistry.Register(gateways.NewStripeGateway("", "", "http://localhost:8080"))

	webhookService := paymentUsecase.NewWebhookService(txnRepo, invoiceRepo, gatewayRegistry, eventBus)
	paymentService := paymentUsecase.NewPaymentService(gatewayRegistry, invoiceRepo, clientRepo)

	supportService := supportUsecase.NewSupportService(supportRepo, clientRepo, eventBus)
	affiliateService := affiliate.NewAffiliateService(memory.NewMockAffiliateRepository(), clientRepo, orderRepo)
	staffService := staffUsecase.NewStaffService(staffRepo, jwtSecret, "FOSSBilling")
	activityService := activity.NewActivityService(memory.NewMockActivityRepository())

	antispamService := antispam.NewAntispamService(memory.NewMockAntispamRepository(), nil, nil, nil)
	authUc := authUsecase.NewAuthUsecase(clientRepo, antispamService, jwtSecret, "FOSSBilling", activityService)
	passwordUc := authUsecase.NewPasswordUsecase(clientRepo)

	orderListener := listener.NewOrderListener(emailService, orderRepo, productRepo, clientRepo, orderService, affiliateService, regRegistry, provRegistry)
	eventBus.Subscribe(events.EventInvoicePaid, orderListener.HandleInvoicePaid)
	eventBus.Subscribe(events.EventOrderActivated, orderListener.HandleOrderActivated)

	guestAuthHandler := guest.NewAuthHandler(authUc)
	guestCartHandler := guest.NewCartHandler(cartService)
	guestWebhookHandler := guest.NewWebhookHandler(webhookService)

	clientProfileHandler := client.NewProfileHandler(authUc, passwordUc)
	clientOrderHandler := client.NewOrderHandler(orderService)
	clientInvoiceHandler := client.NewInvoiceHandler(invoiceRepo, clientRepo, memory.NewMockCompanyRepository(), invoiceService, paymentService)
	clientDepositHandler := client.NewDepositHandler(invoiceService, paymentService, clientRepo)
	clientSupportHandler := client.NewSupportHandler(supportService)

	adminStaffAuthHandler := admin.NewStaffAuthHandler(staffService)
	adminStaffMgmtHandler := admin.NewStaffManagementHandler(staffService, clientRepo, orderRepo, orderService, supportService)
	adminInvHandler := admin.NewInvoiceManagementHandler(staffService, invoiceRepo, clientRepo, invoiceService)
	statsService := statsUsecase.NewStatsService(clientRepo, orderRepo, invoiceRepo, supportRepo, cache.NewMemoryCache())
	adminBillingHandler := admin.NewBillingModuleHandler(staffService, statsService, taxCalc, promoRepo)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) { response.JSON(w, 200, map[string]string{"status": "ok"}, nil) })
	mux.HandleFunc("POST /api/v1/guest/auth/register", guestAuthHandler.Register)
	mux.HandleFunc("POST /api/v1/guest/auth/login", guestAuthHandler.Login)
	mux.HandleFunc("POST /api/v1/guest/cart/calculate", guestCartHandler.Calculate)
	mux.HandleFunc("POST /api/v1/guest/cart/checkout", guestCartHandler.Checkout)
	mux.HandleFunc("POST /api/v1/guest/gateways/{gateway}/webhook", guestWebhookHandler.HandleGatewayWebhook)
	mux.HandleFunc("POST /api/v1/admin/auth/login", adminStaffAuthHandler.Login)

	ca := middleware.RequireAuth(jwtSecret, "client", "admin", "superadmin")
	mux.Handle("GET /api/v1/client/profile", ca(http.HandlerFunc(clientProfileHandler.GetProfile)))
	mux.Handle("POST /api/v1/client/profile/change-password", ca(http.HandlerFunc(clientProfileHandler.ChangePassword)))
	mux.Handle("GET /api/v1/client/orders", ca(http.HandlerFunc(clientOrderHandler.ListOrders)))
	mux.Handle("GET /api/v1/client/orders/{id}", ca(http.HandlerFunc(clientOrderHandler.GetOrder)))
	mux.Handle("GET /api/v1/client/invoices", ca(http.HandlerFunc(clientInvoiceHandler.ListInvoices)))
	mux.Handle("GET /api/v1/client/invoices/{id}", ca(http.HandlerFunc(clientInvoiceHandler.GetInvoice)))
	mux.Handle("GET /api/v1/client/invoices/{id}/pdf", ca(http.HandlerFunc(clientInvoiceHandler.DownloadPDF)))
	mux.Handle("POST /api/v1/client/invoices/{id}/pay-balance", ca(http.HandlerFunc(clientInvoiceHandler.PayWithBalance)))
	mux.Handle("POST /api/v1/client/funds/deposit", ca(http.HandlerFunc(clientDepositHandler.DepositFunds)))
	mux.Handle("POST /api/v1/client/support/tickets", ca(http.HandlerFunc(clientSupportHandler.OpenTicket)))
	mux.Handle("GET /api/v1/client/support/tickets", ca(http.HandlerFunc(clientSupportHandler.ListTickets)))
	mux.Handle("GET /api/v1/client/support/tickets/{id}", ca(http.HandlerFunc(clientSupportHandler.GetTicket)))
	mux.Handle("POST /api/v1/client/support/tickets/{id}/reply", ca(http.HandlerFunc(clientSupportHandler.ReplyTicket)))
	mux.Handle("POST /api/v1/client/support/tickets/{id}/close", ca(http.HandlerFunc(clientSupportHandler.CloseTicket)))

	aa := middleware.RequireAuth(jwtSecret, "admin", "superadmin", "support", "billing")
	mux.Handle("GET /api/v1/admin/clients", aa(http.HandlerFunc(adminStaffMgmtHandler.ListClients)))
	mux.Handle("GET /api/v1/admin/orders", aa(http.HandlerFunc(adminStaffMgmtHandler.ListOrders)))
	mux.Handle("POST /api/v1/admin/orders/{id}/suspend", aa(http.HandlerFunc(adminStaffMgmtHandler.SuspendOrder)))
	mux.Handle("POST /api/v1/admin/orders/{id}/unsuspend", aa(http.HandlerFunc(adminStaffMgmtHandler.UnsuspendOrder)))
	mux.Handle("POST /api/v1/admin/orders/{id}/activate", aa(http.HandlerFunc(adminStaffMgmtHandler.ActivateOrder)))
	mux.Handle("GET /api/v1/admin/invoices", aa(http.HandlerFunc(adminInvHandler.ListInvoices)))
	mux.Handle("POST /api/v1/admin/invoices", aa(http.HandlerFunc(adminInvHandler.CreateInvoice)))
	mux.Handle("GET /api/v1/admin/tax-rules", aa(http.HandlerFunc(adminBillingHandler.ListTaxRules)))
	mux.Handle("POST /api/v1/admin/tax-rules", aa(http.HandlerFunc(adminBillingHandler.CreateTaxRule)))
	mux.Handle("PUT /api/v1/admin/tax-rules/{id}", aa(http.HandlerFunc(adminBillingHandler.UpdateTaxRule)))
	mux.Handle("DELETE /api/v1/admin/tax-rules/{id}", aa(http.HandlerFunc(adminBillingHandler.DeleteTaxRule)))
	mux.Handle("GET /api/v1/admin/taxes", aa(http.HandlerFunc(adminBillingHandler.ListTaxRules)))
	mux.Handle("POST /api/v1/admin/taxes", aa(http.HandlerFunc(adminBillingHandler.CreateTaxRule)))
	mux.Handle("PUT /api/v1/admin/taxes/{id}", aa(http.HandlerFunc(adminBillingHandler.UpdateTaxRule)))
	mux.Handle("DELETE /api/v1/admin/taxes/{id}", aa(http.HandlerFunc(adminBillingHandler.DeleteTaxRule)))
	mux.Handle("GET /api/v1/admin/support/tickets", aa(http.HandlerFunc(adminStaffMgmtHandler.ListTickets)))
	mux.Handle("POST /api/v1/admin/support/tickets/{id}/reply", aa(http.HandlerFunc(adminStaffMgmtHandler.ReplyTicket)))
	mux.Handle("GET /api/v1/admin/audit-logs", aa(http.HandlerFunc(adminStaffAuthHandler.GetAuditLogs)))

	s := httptest.NewServer(middleware.Logger(middleware.CORS("*")(mux)))
	return s, promoRepo, staffRepo, productRepo
}
