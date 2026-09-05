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

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/config"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/http/admin"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/http/client"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/http/guest"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/postgres"
	apikeyUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/apikey"
	authUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/auth"
	billingUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/billing"
	cartUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/cart"
	companyUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/company"
	currencyUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/currency"
	downloadableUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/downloadable"
	massmailUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/massmail"
	newsUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/news"
	orderUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/order"
	paymentUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/payment"
	staffUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	statsUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/stats"
	supportUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/support"
	authPkg "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/auth"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/mailer"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
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

	// 2. Initialize Repositories (with Graceful In-Memory Fallback if Postgres is offline)
	var (
		clientRepo   domain.ClientRepository
		orderRepo    domain.OrderRepository
		invoiceRepo  domain.InvoiceRepository
		txnRepo      domain.TransactionRepository
		promoRepo    domain.PromoRepository
		supportRepo  domain.SupportRepository
		staffRepo    domain.StaffRepository
		currencyRepo domain.CurrencyRepository
		newsRepo     domain.NewsRepository
		downloadRepo domain.DownloadableRepository
		apiKeyRepo   domain.APIKeyRepository
		massMailRepo domain.MassMailRepository
		companyRepo  domain.CompanyRepository
	)

	if pgPool != nil && err == nil {
		clientRepo = postgres.NewClientRepository(pgPool)
		orderRepo = postgres.NewOrderRepository(pgPool)
		invoiceRepo = postgres.NewInvoiceRepository(pgPool)
		txnRepo = postgres.NewTransactionRepository(pgPool)
		promoRepo = postgres.NewPromoRepository(pgPool)
		supportRepo = postgres.NewSupportRepository(pgPool)
		staffRepo = postgres.NewStaffRepository(pgPool)
		currencyRepo = postgres.NewCurrencyRepository(pgPool)
		newsRepo = postgres.NewNewsRepository(pgPool)
		downloadRepo = postgres.NewDownloadableRepository(pgPool)
		apiKeyRepo = postgres.NewAPIKeyRepository(pgPool)
		massMailRepo = postgres.NewMassMailRepository(pgPool)
		companyRepo = postgres.NewCompanyRepository(pgPool)
	} else {
		log.Println("💡 [Fallback Mode] Running with fast In-Memory storage (seeded with initial sample data).")
		memClient := memory.NewMockClientRepository()
		memOrder := memory.NewMockOrderRepository()
		memInv := memory.NewMockInvoiceRepository()
		memTxn := memory.NewMockTransactionRepository()
		memPromo := memory.NewMockPromoRepository()
		memSupport := memory.NewMockSupportRepository()
		memStaff := memory.NewMockStaffRepository()
		memCurr := memory.NewMockCurrencyRepository()
		memNews := memory.NewMockNewsRepository()
		memDl := memory.NewMockDownloadableRepository()
		memKey := memory.NewMockAPIKeyRepository()
		memMail := memory.NewMockMassMailRepository()
		memCompany := memory.NewMockCompanyRepository()

		// Seed initial sample news
		now := time.Now().UTC()
		_ = memNews.Create(ctx, &domain.NewsPost{
			AdminID:     1,
			Title:       "Selamat Datang di FOSSBilling Golang Edition",
			Slug:        "selamat-datang-di-fossbilling-golang-edition",
			Content:     "Backend FOSSBilling kini hadir dengan arsitektur modern Go bertenaga tinggi, Clean Architecture, dan REST API terstandarisasi.",
			Status:      domain.NewsStatusPublished,
			PublishedAt: &now,
		})
		_ = memNews.Create(ctx, &domain.NewsPost{
			AdminID:     1,
			Title:       "Pembaruan Infrastruktur Cloud & Server 2026",
			Slug:        "pembaruan-infrastruktur-cloud-server-2026",
			Content:     "Kapasitas jaringan bandwidth server telah di-upgrade menjadi 10Gbps unmetered.",
			Status:      domain.NewsStatusPublished,
			PublishedAt: &now,
		})

		// Seed Default Superadmin Staff (admin@fossbilling.org / admin123)
		adminPassHash, _ := authPkg.HashPassword("admin123")
		_ = memStaff.CreateGroup(ctx, &domain.AdminGroup{
			Name: "superadmin",
			Permissions: map[string][]string{
				"system":     {"read", "write"},
				"clients":    {"read", "write"},
				"orders":     {"read", "write"},
				"support":    {"read", "write"},
				"billing":    {"read", "write"},
				"currencies": {"read", "write"},
				"news":       {"read", "write"},
			},
		})
		_ = memStaff.Create(ctx, &domain.Staff{
			GroupID:      1,
			Email:        "admin@fossbilling.org",
			PasswordHash: adminPassHash,
			Name:         "Super Administrator",
			Role:         "superadmin",
			Status:       "active",
		})

		// Seed Default Client (client@fossbilling.org / client123)
		clientPassHash, _ := authPkg.HashPassword("client123")
		_ = memClient.Create(ctx, &domain.Client{
			Email:        "client@fossbilling.org",
			PasswordHash: clientPassHash,
			FirstName:    "Demo",
			LastName:     "Customer",
			Country:      "ID",
			Currency:     "IDR",
			Status:       "active",
		})

		// Seed Sample Product & Downloadable
		_ = memDl.Create(ctx, &domain.DownloadableFile{
			ProductID:   1,
			Filename:    "fossbilling-starter-pack.zip",
			FilePath:    "/data/files/starter-pack.zip",
			FileSize:    1024 * 1024 * 25,
			ContentType: "application/zip",
			Version:     "1.0.0",
		})

		clientRepo = memClient
		orderRepo = memOrder
		invoiceRepo = memInv
		txnRepo = memTxn
		promoRepo = memPromo
		supportRepo = memSupport
		staffRepo = memStaff
		currencyRepo = memCurr
		newsRepo = memNews
		downloadRepo = memDl
		apiKeyRepo = memKey
		massMailRepo = memMail
		companyRepo = memCompany
	}

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

	companyService := companyUsecase.NewCompanyService(companyRepo)
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
	guestCompanyHandler := guest.NewCompanyHandler(companyService)

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
	adminCompanyHandler := admin.NewCompanyHandler(companyService, staffService)

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

	mux.HandleFunc("GET /openapi.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		data, err := os.ReadFile("docs/openapi.json")
		if err != nil {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "OpenAPI spec not found", nil)
			return
		}
		_, _ = w.Write(data)
	})

	mux.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		html := `<!doctype html>
<html>
  <head>
    <title>FOSSBilling Next-Gen API Documentation</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <script id="api-reference" data-url="/openapi.json"></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`
		_, _ = w.Write([]byte(html))
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
	mux.Handle("GET /api/v1/guest/company", rateLimiter.RateLimit(http.HandlerFunc(guestCompanyHandler.GetCompany)))

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

	mux.Handle("GET /api/v1/admin/company", adminAuthMiddleware(http.HandlerFunc(adminCompanyHandler.GetCompany)))
	mux.Handle("PUT /api/v1/admin/company", adminAuthMiddleware(http.HandlerFunc(adminCompanyHandler.UpdateCompany)))

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
