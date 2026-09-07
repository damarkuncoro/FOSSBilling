package main

import (
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/config"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/listener"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/notification"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/payment"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/payment/gateways"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/provisioning"
	activityUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/activity"
	antispamUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/antispam"
	apikeyUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/apikey"
	authUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/auth"
	billingUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/billing"
	cartUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/cart"
	catalogUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/catalog"
	companyUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/company"
	cookieconsentUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/cookieconsent"
	currencyUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/currency"
	domainUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/domain"
	downloadableUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/downloadable"
	extensionUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/extension"
	formbuilderUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/formbuilder"
	licenseUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/license"
	massmailUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/massmail"
	newsUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/news"
	notificationUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/notification"
	orderUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/order"
	pageUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/page"
	paymentUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/payment"
	redirectUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/redirect"
	seoUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/seo"
	staffUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	statsUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/stats"
	supportUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/support"
	systemUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/system"
	themeUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/theme"
	widgetUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/widget"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/cache"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/mailer"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/security"
)

// Services holds all application domain use cases and services
type Services struct {
	Auth          *authUsecase.AuthUsecase
	Password      *authUsecase.PasswordUsecase
	Order         *orderUsecase.OrderService
	Invoice       *billingUsecase.InvoiceService
	Cart          *cartUsecase.CartService
	Webhook       *paymentUsecase.WebhookService
	Payment       *paymentUsecase.PaymentService
	Support       *supportUsecase.SupportService
	Staff         *staffUsecase.StaffService
	Stats         *statsUsecase.StatsService
	Company       companyUsecase.CompanyService
	Currency      *currencyUsecase.CurrencyService
	News          *newsUsecase.NewsService
	Downloadable  *downloadableUsecase.DownloadableService
	Domain        *domainUsecase.DomainService
	License       *licenseUsecase.LicenseService
	APIKey        *apikeyUsecase.APIKeyService
	MassMail      *massmailUsecase.MassMailService
	Product       *catalogUsecase.ProductService
	System        *systemUsecase.SystemService
	Page          *pageUsecase.PageService
	Activity      *activityUsecase.ActivityService
	Notification  *notificationUsecase.NotificationService
	Antispam      *antispamUsecase.AntispamService
	Formbuilder   *formbuilderUsecase.FormbuilderService
	Extension     *extensionUsecase.ExtensionService
	Redirect      *redirectUsecase.RedirectService
	CookieConsent *cookieconsentUsecase.CookieConsentService
	Theme         *themeUsecase.ThemeService
	SEO           *seoUsecase.SEOService
	Widget        *widgetUsecase.WidgetService
	Gateways      *payment.GatewayRegistry
	Cache         cache.Cache
	EventBus      *events.EventBus
}

// InitServices instantiates and configures all domain services with their respective repository dependencies
func InitServices(cfg *config.Config, repos *Repositories, eventBus *events.EventBus, appCache cache.Cache) *Services {
	mockMailer := mailer.NewMockMailer()
	emailService := notification.NewEmailService(mockMailer, "admin@fossbilling.org", "FOSSBilling")

	taxCalculator := billingUsecase.NewTaxCalculator(nil)
	promoCalc := cartUsecase.NewPromoCalculator(repos.Promo)

	authUc := authUsecase.NewAuthUsecase(repos.Client, cfg.JWTSecret)
	passwordUc := authUsecase.NewPasswordUsecase(repos.Client)

	orderService := orderUsecase.NewOrderService(repos.Order, eventBus)
	invoiceService := billingUsecase.NewInvoiceService(repos.Invoice, repos.Client, taxCalculator, eventBus)
	cartService := cartUsecase.NewCartService(promoCalc, repos.Promo, repos.Order, repos.Client, taxCalculator, invoiceService)

	gatewayRegistry := payment.NewGatewayRegistry()
	gatewayRegistry.Register(gateways.NewStripeGateway("sk_test", "pk_test", "whsec_test"))
	gatewayRegistry.Register(gateways.NewPayPalGateway("client_id", "secret", false))
	gatewayRegistry.Register(gateways.NewMidtransGateway("server_key", "client_key", false))
	gatewayRegistry.Register(gateways.NewBankTransferGateway("Bank Mandiri", "1234567890", "FOSSBilling Indonesia"))
	gatewayRegistry.Register(gateways.NewCustomGateway())

	webhookService := paymentUsecase.NewWebhookService(repos.Transaction, repos.Invoice, eventBus)
	paymentService := paymentUsecase.NewPaymentService(gatewayRegistry, repos.Invoice, repos.Client)

	supportService := supportUsecase.NewSupportService(repos.Support, repos.Client, eventBus)
	staffService := staffUsecase.NewStaffService(repos.Staff, cfg.JWTSecret)
	statsService := statsUsecase.NewStatsService(repos.Client, repos.Order, repos.Invoice, repos.Support)

	companyService := companyUsecase.NewCompanyService(repos.Company)
	currencyService := currencyUsecase.NewCurrencyService(repos.Currency)
	newsService := newsUsecase.NewNewsService(repos.News)
	downloadService := downloadableUsecase.NewDownloadableService(repos.Downloadable, repos.Order, cfg.JWTSecret)

	provisionerRegistry := provisioning.NewProvisionerRegistry()
	provisionerRegistry.Register("cpanel", provisioning.NewCpanelProvisioner(provisioning.CpanelConfig{
		Host:     "cpanel.fossbilling.org",
		Username: "root",
		APIToken: "MOCK_TOKEN_123",
		Insecure: true,
	}))
	provisionerRegistry.Register("directadmin", provisioning.NewDirectAdminProvisioner("da.fossbilling.org", 2222, "admin", "pass"))
	provisionerRegistry.Register("plesk", provisioning.NewPleskProvisioner(provisioning.PleskConfig{
		Host:     "plesk.fossbilling.org",
		Port:     8443,
		APIKey:   "MOCK_KEY_456",
		Insecure: true,
	}))
	provisionerRegistry.Register("hestia", provisioning.NewHestiaProvisioner(provisioning.HestiaConfig{
		Host:      "hestia.fossbilling.org",
		Port:      8083,
		AccessKey: "admin",
		SecretKey: "MOCK_HESTIA_KEY_789",
		Insecure:  true,
	}))
	provisionerRegistry.Register("cwp", provisioning.NewCWPProvisioner(provisioning.CWPConfig{
		Host:     "cwp.fossbilling.org",
		Port:     2304,
		APIKey:   "MOCK_CWP_KEY_101",
		Insecure: true,
	}))
	provisionerRegistry.Register("custom", provisioning.NewCustomServerProvisioner(provisioning.CustomServerConfig{
		EndpointURL: "https://webhooks.fossbilling.org/server",
		AuthToken:   "MOCK_CUSTOM_AUTH_202",
	}))

	registrarRegistry := provisioning.NewRegistrarRegistry()
	registrarRegistry.Register("rdap", provisioning.NewRDAPRegistrarDriver())
	registrarRegistry.Register("email", provisioning.NewEmailRegistrarDriver(emailService, "admin@fossbilling.org"))
	registrarRegistry.Register("custom", provisioning.NewCustomRegistrarDriver())

	// LogicBoxes compatible registrars
	lbConfig := provisioning.ResellerClubConfig{IsTest: true}
	registrarRegistry.Register("resellerclub", provisioning.NewResellerClubRegistrarDriver(lbConfig))
	registrarRegistry.Register("resellerid", provisioning.NewResellerClubRegistrarDriver(lbConfig))
	registrarRegistry.Register("resellbiz", provisioning.NewResellerClubRegistrarDriver(lbConfig))
	registrarRegistry.Register("netearthone", provisioning.NewResellerClubRegistrarDriver(lbConfig))

	registrarRegistry.Register("namecheap", provisioning.NewNamecheapRegistrarDriver(provisioning.NamecheapConfig{IsSandbox: true}))
	registrarRegistry.Register("internetbs", provisioning.NewInternetbsRegistrarDriver(provisioning.InternetbsConfig{IsTest: true}))

	domainService := domainUsecase.NewDomainService(repos.Order, registrarRegistry)
	licenseService := licenseUsecase.NewLicenseService(repos.Order)
	apiKeyService := apikeyUsecase.NewAPIKeyService(repos.APIKey)
	productService := catalogUsecase.NewProductService(repos.Product, appCache)
	systemService := systemUsecase.NewSystemService(repos.System)
	pageService := pageUsecase.NewPageService(repos.Page)
	activityService := activityUsecase.NewActivityService(repos.Activity)
	notificationService := notificationUsecase.NewNotificationService(repos.Notification)
	massMailService := massmailUsecase.NewMassMailService(repos.MassMail, repos.Client, mockMailer, "admin@fossbilling.org", "FOSSBilling")

	turnstileVerifier := security.NewTurnstileVerifier("")
	emailChecker := security.NewDisposableEmailChecker()
	sfsChecker := security.NewStopForumSpamChecker()
	antispamService := antispamUsecase.NewAntispamService(repos.Antispam, turnstileVerifier, emailChecker, sfsChecker)

	// Register Event Listeners
	orderListener := listener.NewOrderListener(emailService, repos.Order, repos.Product, repos.Client, orderService, registrarRegistry, provisionerRegistry)
	eventBus.Subscribe(events.EventOrderActivated, orderListener.HandleOrderActivated)
	eventBus.Subscribe(events.EventInvoicePaid, orderListener.HandleInvoicePaid)

	invoiceListener := listener.NewInvoiceListener(emailService, repos.Invoice, repos.Client)
	eventBus.Subscribe(events.EventInvoicePaid, invoiceListener.HandleInvoicePaid)

	activityListener := listener.NewActivityListener(activityService)
	eventBus.Subscribe(events.EventOrderActivated, activityListener.HandleOrderActivated)
	eventBus.Subscribe(events.EventInvoicePaid, activityListener.HandleInvoicePaid)
	eventBus.Subscribe(events.EventTicketOpened, activityListener.HandleTicketOpened)

	notificationListener := listener.NewNotificationListener(notificationService)
	eventBus.Subscribe(events.EventInvoicePaid, notificationListener.HandleInvoicePaid)
	eventBus.Subscribe(events.EventOrderActivated, notificationListener.HandleOrderActivated)

	return &Services{
		Auth:          authUc,
		Password:      passwordUc,
		Order:         orderService,
		Invoice:       invoiceService,
		Cart:          cartService,
		Webhook:       webhookService,
		Payment:       paymentService,
		Support:       supportService,
		Staff:         staffService,
		Stats:         statsService,
		Company:       companyService,
		Currency:      currencyService,
		News:          newsService,
		Downloadable:  downloadService,
		Domain:        domainService,
		License:       licenseService,
		APIKey:        apiKeyService,
		Product:       productService,
		System:        systemService,
		Page:          pageService,
		Activity:      activityService,
		Notification:  notificationService,
		Antispam:      antispamService,
		Formbuilder:   formbuilderUsecase.NewFormbuilderService(repos.Formbuilder),
		Extension:     extensionUsecase.NewExtensionService(repos.Extension),
		Redirect:      redirectUsecase.NewRedirectService(repos.Redirect),
		CookieConsent: cookieconsentUsecase.NewCookieConsentService(repos.Extension),
		Theme:         themeUsecase.NewThemeService(repos.Theme),
		SEO:           seoUsecase.NewSEOService(repos.Page, repos.News, repos.Product),
		Widget:        widgetUsecase.NewWidgetService(),
		Gateways:      gatewayRegistry,
		Cache:         appCache,
		MassMail:      massMailService,
		EventBus:      eventBus,
	}
}
