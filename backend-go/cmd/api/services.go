package main

import (
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/config"
	apikeyUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/apikey"
	authUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/auth"
	billingUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/billing"
	cartUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/cart"
	companyUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/company"
	currencyUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/currency"
	downloadableUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/downloadable"
	domainUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/domain"
	licenseUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/license"
	massmailUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/massmail"
	newsUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/news"
	orderUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/order"
	paymentUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/payment"
	staffUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	statsUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/stats"
	supportUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/support"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/provisioning"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/mailer"
)

// Services holds all application domain use cases and services
type Services struct {
	Auth         *authUsecase.AuthUsecase
	Password     *authUsecase.PasswordUsecase
	Order        *orderUsecase.OrderService
	Invoice      *billingUsecase.InvoiceService
	Cart         *cartUsecase.CartService
	Webhook      *paymentUsecase.WebhookService
	Support      *supportUsecase.SupportService
	Staff        *staffUsecase.StaffService
	Stats        *statsUsecase.StatsService
	Company      companyUsecase.CompanyService
	Currency     *currencyUsecase.CurrencyService
	News         *newsUsecase.NewsService
	Downloadable *downloadableUsecase.DownloadableService
	Domain       *domainUsecase.DomainService
	License      *licenseUsecase.LicenseService
	APIKey       *apikeyUsecase.APIKeyService
	MassMail     *massmailUsecase.MassMailService
}

// InitServices instantiates and configures all domain services with their respective repository dependencies
func InitServices(cfg *config.Config, repos *Repositories) *Services {
	taxCalculator := billingUsecase.NewTaxCalculator(nil)
	promoCalc := cartUsecase.NewPromoCalculator(repos.Promo)
	orderService := orderUsecase.NewOrderService(repos.Order)
	invoiceService := billingUsecase.NewInvoiceService(repos.Invoice, repos.Client, taxCalculator)
	cartService := cartUsecase.NewCartService(promoCalc, repos.Promo, repos.Order, invoiceService)
	webhookService := paymentUsecase.NewWebhookService(repos.Transaction, repos.Invoice, orderService, repos.Order)
	supportService := supportUsecase.NewSupportService(repos.Support, repos.Client)
	staffService := staffUsecase.NewStaffService(repos.Staff, cfg.JWTSecret)
	statsService := statsUsecase.NewStatsService(repos.Client, repos.Order, repos.Invoice, repos.Support)
	authUc := authUsecase.NewAuthUsecase(repos.Client, cfg.JWTSecret)
	passwordUc := authUsecase.NewPasswordUsecase(repos.Client)

	companyService := companyUsecase.NewCompanyService(repos.Company)
	currencyService := currencyUsecase.NewCurrencyService(repos.Currency)
	newsService := newsUsecase.NewNewsService(repos.News)
	downloadService := downloadableUsecase.NewDownloadableService(repos.Downloadable, repos.Order, cfg.JWTSecret)
	domainService := domainUsecase.NewDomainService(repos.Order, provisioning.NewRDAPRegistrarDriver())
	licenseService := licenseUsecase.NewLicenseService(repos.Order)
	apiKeyService := apikeyUsecase.NewAPIKeyService(repos.APIKey)
	massMailService := massmailUsecase.NewMassMailService(repos.MassMail, repos.Client, mailer.NewMockMailer(), "admin@fossbilling.org", "FOSSBilling")

	return &Services{
		Auth:         authUc,
		Password:     passwordUc,
		Order:        orderService,
		Invoice:      invoiceService,
		Cart:         cartService,
		Webhook:      webhookService,
		Support:      supportService,
		Staff:        staffService,
		Stats:        statsService,
		Company:      companyService,
		Currency:     currencyService,
		News:         newsService,
		Downloadable: downloadService,
		Domain:       domainService,
		License:      licenseService,
		APIKey:       apiKeyService,
		MassMail:     massMailService,
	}
}
