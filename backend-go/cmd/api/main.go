package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fossbilling/backend-go/internal/config"
	"github.com/fossbilling/backend-go/internal/handler/http/admin"
	"github.com/fossbilling/backend-go/internal/handler/http/client"
	"github.com/fossbilling/backend-go/internal/handler/http/guest"
	"github.com/fossbilling/backend-go/internal/handler/middleware"
	"github.com/fossbilling/backend-go/internal/repository/postgres"
	apikeyUsecase "github.com/fossbilling/backend-go/internal/usecase/apikey"
	authUsecase "github.com/fossbilling/backend-go/internal/usecase/auth"
	billingUsecase "github.com/fossbilling/backend-go/internal/usecase/billing"
	cartUsecase "github.com/fossbilling/backend-go/internal/usecase/cart"
	currencyUsecase "github.com/fossbilling/backend-go/internal/usecase/currency"
	downloadableUsecase "github.com/fossbilling/backend-go/internal/usecase/downloadable"
	massmailUsecase "github.com/fossbilling/backend-go/internal/usecase/massmail"
	newsUsecase "github.com/fossbilling/backend-go/internal/usecase/news"
	orderUsecase "github.com/fossbilling/backend-go/internal/usecase/order"
	paymentUsecase "github.com/fossbilling/backend-go/internal/usecase/payment"
	staffUsecase "github.com/fossbilling/backend-go/internal/usecase/staff"
	statsUsecase "github.com/fossbilling/backend-go/internal/usecase/stats"
	supportUsecase "github.com/fossbilling/backend-go/internal/usecase/support"
	"github.com/fossbilling/backend-go/pkg/mailer"
	"github.com/fossbilling/backend-go/pkg/response"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Initialize PostgreSQL Connection Pool
	pgPool, err := postgres.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Printf("⚠️ Database connection failed (%v). Running with pool handle.", err)
	} else {
		defer pgPool.Close()
		log.Println("✅ Connected to PostgreSQL database pool.")
	}

	// 2. Initialize Repositories
	clientRepo := postgres.NewClientRepository(pgPool)
	_ = postgres.NewProductRepository(pgPool)
	orderRepo := postgres.NewOrderRepository(pgPool)
	invoiceRepo := postgres.NewInvoiceRepository(pgPool)
	txnRepo := postgres.NewTransactionRepository(pgPool)
	promoRepo := postgres.NewPromoRepository(pgPool)
	supportRepo := postgres.NewSupportRepository(pgPool)
	staffRepo := postgres.NewStaffRepository(pgPool)
	currencyRepo := postgres.NewCurrencyRepository(pgPool)
	newsRepo := postgres.NewNewsRepository(pgPool)
	downloadRepo := postgres.NewDownloadableRepository(pgPool)
	apiKeyRepo := postgres.NewAPIKeyRepository(pgPool)
	massMailRepo := postgres.NewMassMailRepository(pgPool)

	// 3. Initialize Services & Engines
	taxCalculator := billingUsecase.NewTaxCalculator(nil)
	promoCalc := cartUsecase.NewPromoCalculator(promoRepo)
	orderService := orderUsecase.NewOrderService(orderRepo)
	invoiceService := billingUsecase.NewInvoiceService(invoiceRepo, clientRepo, taxCalculator)

	cartService := cartUsecase.NewCartService(promoCalc, promoRepo, orderRepo, invoiceService)
	webhookService := paymentUsecase.NewWebhookService(txnRepo, invoiceRepo, orderService, orderRepo)
	supportService := supportUsecase.NewSupportService(supportRepo, clientRepo)
	staffService := staffUsecase.NewStaffService(staffRepo, cfg.JWTSecret)
	statsService := statsUsecase.NewStatsService(clientRepo, orderRepo, invoiceRepo, supportRepo)
	authUc := authUsecase.NewAuthUsecase(clientRepo, cfg.JWTSecret)

	currencyService := currencyUsecase.NewCurrencyService(currencyRepo)
	newsService := newsUsecase.NewNewsService(newsRepo)
	downloadService := downloadableUsecase.NewDownloadableService(downloadRepo, orderRepo, cfg.JWTSecret)
	apiKeyService := apikeyUsecase.NewAPIKeyService(apiKeyRepo)
	appMailer := mailer.NewMockMailer()
	massMailService := massmailUsecase.NewMassMailService(massMailRepo, clientRepo, appMailer, "admin@fossbilling.org", "FOSSBilling")

	// 4. Initialize Handlers
	guestAuthHandler := guest.NewAuthHandler(authUc)
	guestCartHandler := guest.NewCartHandler(cartService)
	guestWebhookHandler := guest.NewWebhookHandler(webhookService)
	guestCurrencyHandler := guest.NewCurrencyHandler(currencyService)
	guestNewsHandler := guest.NewNewsHandler(newsService)

	clientProfileHandler := client.NewProfileHandler(authUc)
	clientOrderHandler := client.NewOrderHandler(orderRepo)
	clientInvoiceHandler := client.NewInvoiceHandler(invoiceRepo, clientRepo, invoiceService)
	clientSupportHandler := client.NewSupportHandler(supportService)
	clientDownloadHandler := client.NewDownloadHandler(downloadService)
	clientAPIKeyHandler := client.NewAPIKeyHandler(apiKeyService)

	adminStaffHandler := admin.NewStaffHandler(staffService, clientRepo, orderRepo, orderService, supportService)
	adminStatsHandler := admin.NewStatsHandler(statsService, staffService)
	adminCurrencyHandler := admin.NewCurrencyHandler(currencyService, staffService)
	adminNewsHandler := admin.NewNewsHandler(newsService, staffService)
	adminMassMailHandler := admin.NewMassMailHandler(massMailService, staffService)

	// 5. Rate Limiter for public endpoints (60 req / min)
	rateLimiter := middleware.NewRateLimiter(60, time.Second)

	// 6. Setup HTTP Router (Go 1.22 enhanced ServeMux)
	mux := http.NewServeMux()

	// System & Health Endpoints
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]interface{}{
			"status":      "ok",
			"environment": cfg.AppEnv,
			"timestamp":   time.Now().UTC().Format(time.RFC3339),
		}, nil)
	})

	mux.HandleFunc("GET /api/v1", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]string{
			"message": "FOSSBilling Next-Gen API v1 (Golang)",
			"docs":    "/docs",
		}, nil)
	})

	// Public / Guest Routes (Rate Limited)
	mux.Handle("POST /api/v1/guest/auth/register", rateLimiter.RateLimit(http.HandlerFunc(guestAuthHandler.Register)))
	mux.Handle("POST /api/v1/guest/auth/login", rateLimiter.RateLimit(http.HandlerFunc(guestAuthHandler.Login)))
	mux.Handle("POST /api/v1/guest/cart/calculate", rateLimiter.RateLimit(http.HandlerFunc(guestCartHandler.Calculate)))
	mux.Handle("POST /api/v1/guest/cart/checkout", rateLimiter.RateLimit(http.HandlerFunc(guestCartHandler.Checkout)))
	mux.Handle("POST /api/v1/guest/gateways/{gateway}/webhook", http.HandlerFunc(guestWebhookHandler.HandleGatewayWebhook))
	mux.Handle("GET /api/v1/guest/currencies", rateLimiter.RateLimit(http.HandlerFunc(guestCurrencyHandler.List)))
	mux.Handle("GET /api/v1/guest/news", rateLimiter.RateLimit(http.HandlerFunc(guestNewsHandler.List)))
	mux.Handle("GET /api/v1/guest/news/{slug}", rateLimiter.RateLimit(http.HandlerFunc(guestNewsHandler.Get)))

	// Admin Auth (Public Login)
	mux.Handle("POST /api/v1/admin/auth/login", rateLimiter.RateLimit(http.HandlerFunc(adminStaffHandler.Login)))

	// Protected Client Routes
	clientAuthMiddleware := middleware.RequireAuth(cfg.JWTSecret, "client", "admin", "superadmin")
	mux.Handle("GET /api/v1/client/profile", clientAuthMiddleware(http.HandlerFunc(clientProfileHandler.GetProfile)))
	mux.Handle("PUT /api/v1/client/profile", clientAuthMiddleware(http.HandlerFunc(clientProfileHandler.UpdateProfile)))

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

	mux.Handle("GET /api/v1/client/downloads/{id}/link", clientAuthMiddleware(http.HandlerFunc(clientDownloadHandler.GenerateLink)))
	mux.Handle("GET /api/v1/client/downloads/{id}/file", http.HandlerFunc(clientDownloadHandler.StreamFile))
	mux.Handle("GET /api/v1/client/api-keys", clientAuthMiddleware(http.HandlerFunc(clientAPIKeyHandler.List)))
	mux.Handle("POST /api/v1/client/api-keys", clientAuthMiddleware(http.HandlerFunc(clientAPIKeyHandler.Generate)))
	mux.Handle("DELETE /api/v1/client/api-keys/{id}", clientAuthMiddleware(http.HandlerFunc(clientAPIKeyHandler.Revoke)))

	// Protected Admin Routes
	adminAuthMiddleware := middleware.RequireAuth(cfg.JWTSecret, "admin", "superadmin", "support", "billing")
	mux.Handle("GET /api/v1/admin/stats/dashboard", adminAuthMiddleware(http.HandlerFunc(adminStatsHandler.GetDashboard)))
	mux.Handle("GET /api/v1/admin/clients", adminAuthMiddleware(http.HandlerFunc(adminStaffHandler.ListClients)))
	mux.Handle("GET /api/v1/admin/orders", adminAuthMiddleware(http.HandlerFunc(adminStaffHandler.ListOrders)))
	mux.Handle("POST /api/v1/admin/orders/{id}/suspend", adminAuthMiddleware(http.HandlerFunc(adminStaffHandler.SuspendOrder)))
	mux.Handle("POST /api/v1/admin/orders/{id}/unsuspend", adminAuthMiddleware(http.HandlerFunc(adminStaffHandler.UnsuspendOrder)))
	mux.Handle("POST /api/v1/admin/orders/{id}/activate", adminAuthMiddleware(http.HandlerFunc(adminStaffHandler.ActivateOrder)))
	mux.Handle("GET /api/v1/admin/support/tickets", adminAuthMiddleware(http.HandlerFunc(adminStaffHandler.ListTickets)))
	mux.Handle("POST /api/v1/admin/support/tickets/{id}/reply", adminAuthMiddleware(http.HandlerFunc(adminStaffHandler.ReplyTicket)))
	mux.Handle("GET /api/v1/admin/audit-logs", adminAuthMiddleware(http.HandlerFunc(adminStaffHandler.GetAuditLogs)))

	mux.Handle("GET /api/v1/admin/currencies", adminAuthMiddleware(http.HandlerFunc(adminCurrencyHandler.List)))
	mux.Handle("POST /api/v1/admin/currencies", adminAuthMiddleware(http.HandlerFunc(adminCurrencyHandler.Create)))
	mux.Handle("PUT /api/v1/admin/currencies/{code}", adminAuthMiddleware(http.HandlerFunc(adminCurrencyHandler.Update)))
	mux.Handle("DELETE /api/v1/admin/currencies/{code}", adminAuthMiddleware(http.HandlerFunc(adminCurrencyHandler.Delete)))
	mux.Handle("POST /api/v1/admin/currencies/{code}/default", adminAuthMiddleware(http.HandlerFunc(adminCurrencyHandler.SetDefault)))

	mux.Handle("GET /api/v1/admin/news", adminAuthMiddleware(http.HandlerFunc(adminNewsHandler.List)))
	mux.Handle("POST /api/v1/admin/news", adminAuthMiddleware(http.HandlerFunc(adminNewsHandler.Create)))
	mux.Handle("PUT /api/v1/admin/news/{id}", adminAuthMiddleware(http.HandlerFunc(adminNewsHandler.Update)))
	mux.Handle("DELETE /api/v1/admin/news/{id}", adminAuthMiddleware(http.HandlerFunc(adminNewsHandler.Delete)))

	mux.Handle("GET /api/v1/admin/mass-mail", adminAuthMiddleware(http.HandlerFunc(adminMassMailHandler.List)))
	mux.Handle("POST /api/v1/admin/mass-mail", adminAuthMiddleware(http.HandlerFunc(adminMassMailHandler.Create)))
	mux.Handle("POST /api/v1/admin/mass-mail/{id}/send", adminAuthMiddleware(http.HandlerFunc(adminMassMailHandler.Send)))

	// Wrap with Global Middlewares: Logger -> CORS -> Mux
	handler := middleware.Logger(middleware.CORS(mux))

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful Shutdown Channel
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("🚀 FOSSBilling API Server running on port %s in %s mode...", cfg.Port, cfg.AppEnv)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	case sig := <-shutdown:
		log.Printf("🛑 Signal %v received. Starting graceful shutdown...", sig)
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)
			_ = server.Close()
		}
		log.Println("✅ Server exited cleanly.")
	}
}
