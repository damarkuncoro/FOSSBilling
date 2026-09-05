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
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/mailer"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pgPool, err := postgres.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err == nil {
		defer pgPool.Close()
	}

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

	seedMemoryData(ctx, memClient, memStaff, memPromo, memNews, memCurr, memInv)

	taxCalculator := billingUsecase.NewTaxCalculator(nil)
	promoCalc := cartUsecase.NewPromoCalculator(memPromo)
	orderService := orderUsecase.NewOrderService(memOrder)
	invoiceService := billingUsecase.NewInvoiceService(memInv, memClient, taxCalculator)
	cartService := cartUsecase.NewCartService(promoCalc, memPromo, memOrder, invoiceService)
	webhookService := paymentUsecase.NewWebhookService(memTxn, memInv, orderService, memOrder)
	supportService := supportUsecase.NewSupportService(memSupport, memClient)
	staffService := staffUsecase.NewStaffService(memStaff, cfg.JWTSecret)
	statsService := statsUsecase.NewStatsService(memClient, memOrder, memInv, memSupport)
	authUc := authUsecase.NewAuthUsecase(memClient, cfg.JWTSecret)
	passwordUc := authUsecase.NewPasswordUsecase(memClient)

	companyService := companyUsecase.NewCompanyService(memCompany)
	currencyService := currencyUsecase.NewCurrencyService(memCurr)
	newsService := newsUsecase.NewNewsService(memNews)
	downloadService := downloadableUsecase.NewDownloadableService(memDl, memOrder, cfg.JWTSecret)
	apiKeyService := apikeyUsecase.NewAPIKeyService(memKey)
	massMailService := massmailUsecase.NewMassMailService(memMail, memClient, mailer.NewMockMailer(), "admin@fossbilling.org", "FOSSBilling")

	handlers := &AppHandlers{
		GuestAuth:      guest.NewAuthHandler(authUc),
		GuestCart:      guest.NewCartHandler(cartService),
		GuestWebhook:   guest.NewWebhookHandler(webhookService),
		GuestCurrency:  guest.NewCurrencyHandler(currencyService),
		GuestNews:      guest.NewNewsHandler(newsService),
		GuestCompany:   guest.NewCompanyHandler(companyService),
		ClientProfile:  client.NewProfileHandler(authUc, passwordUc),
		ClientOrder:    client.NewOrderHandler(memOrder),
		ClientInvoice:  client.NewInvoiceHandler(memInv, memClient, invoiceService),
		ClientDeposit:  client.NewDepositHandler(invoiceService),
		ClientSupport:  client.NewSupportHandler(supportService),
		ClientDownload: client.NewDownloadHandler(downloadService),
		ClientAPIKey:   client.NewAPIKeyHandler(apiKeyService),
		AdminAuth:      admin.NewStaffAuthHandler(staffService),
		AdminStaff:     admin.NewStaffManagementHandler(staffService, memClient, memOrder, orderService, supportService),
		AdminClient:    admin.NewClientManagementHandler(staffService, memClient),
		AdminInvoice:   admin.NewInvoiceManagementHandler(memInv, memClient, invoiceService),
		AdminStats:     admin.NewStatsHandler(statsService, staffService),
		AdminCurrency:  admin.NewCurrencyHandler(currencyService, staffService),
		AdminNews:      admin.NewNewsHandler(newsService, staffService),
		AdminMassMail:  admin.NewMassMailHandler(massMailService, staffService),
		AdminCompany:   admin.NewCompanyHandler(companyService, staffService),
		AdminCatalog:   admin.NewCatalogHandler(),
		AdminBilling:   admin.NewBillingModuleHandler(),
		AdminSystem:    admin.NewSystemModuleHandler(),
	}

	rateLimiter := middleware.NewRateLimiter(60, time.Second)
	router := setupRoutes(cfg, handlers, rateLimiter)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		log.Printf("🚀 FOSSBilling API Server running on port %s in %s mode...", cfg.Port, cfg.AppEnv)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	_ = server.Shutdown(shutdownCtx)
	log.Println("✅ Server exited cleanly.")
}
