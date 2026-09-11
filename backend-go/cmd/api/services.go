package main

import (
	"github.com/jackc/pgx/v5/pgxpool"
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
	knowledgebase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/knowledgebase"
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
	affiliateUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/affiliate"
	themeUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/theme"
	widgetUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/widget"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/cache"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/centralalerts"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/fraud"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/mailer"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/notifications"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/plugins"
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
	Affiliate     *affiliateUsecase.AffiliateService
	Company       companyUsecase.CompanyService
	Currency      *currencyUsecase.CurrencyService
	News          *newsUsecase.NewsService
	Downloadable  *downloadableUsecase.DownloadableService
	Domain        *domainUsecase.DomainService
	License       *licenseUsecase.LicenseService
	APIKey        *apikeyUsecase.APIKeyService
	MassMail      *massmailUsecase.MassMailService
	KB            *knowledgebase.Service
	Product       *catalogUsecase.ProductService
	Server        *catalogUsecase.ServerService
	System        *systemUsecase.SystemService
	Page          *pageUsecase.PageService
	Activity      *activityUsecase.ActivityService
	Notification  *notificationUsecase.NotificationService
	AdminNotif    *notificationUsecase.AdminNotificationService
	Antispam      *antispamUsecase.AntispamService
	Fraud         fraud.FraudChecker
	Formbuilder   *formbuilderUsecase.FormbuilderService
	Extension     *extensionUsecase.ExtensionService
	Redirect      *redirectUsecase.RedirectService
	CookieConsent *cookieconsentUsecase.CookieConsentService
	Theme         *themeUsecase.ThemeService
	SEO           *seoUsecase.SEOService
	Widget        *widgetUsecase.WidgetService
	Tax           *billingUsecase.TaxCalculator
	Health        *systemUsecase.HealthUsecase
	Gateways      *payment.GatewayRegistry
	Hooks         *plugins.HookManager
	Cache         cache.Cache
	WSHub         *notifications.WSHub
	EventBus      *events.EventBus
}

// InitServices instantiates and configures all domain services with their respective repository dependencies
func InitServices(cfg *config.Config, repos *Repositories, pool *pgxpool.Pool, eventBus *events.EventBus, appCache cache.Cache, hookManager *plugins.HookManager) *Services {
	var appMailer mailer.Mailer
	if cfg.MailDriver == "smtp" {
		appMailer = mailer.NewSMTPMailer(cfg.MailHost, cfg.MailPort, cfg.MailUser, cfg.MailPass, cfg.MailFromAddr)
	} else {
		appMailer = mailer.NewMockMailer()
	}

	emailService := notification.NewEmailService(appMailer, repos.EmailTemplate, cfg.MailFromAddr, cfg.MailFromName)
	wsHub := notifications.NewWSHub()

	taxCalculator := billingUsecase.NewTaxCalculator(repos.Tax)
	promoCalc := cartUsecase.NewPromoCalculator(repos.Promo)
	fraudChecker := &fraud.MockFraudChecker{}

	// 1. Antispam & Security
	turnstileVerifier := security.NewTurnstileVerifier("")
	emailChecker := security.NewDisposableEmailChecker()
	sfsChecker := security.NewStopForumSpamChecker()
	antispamService := antispamUsecase.NewAntispamService(repos.Antispam, turnstileVerifier, emailChecker, sfsChecker)

	activityService := activityUsecase.NewActivityService(repos.Activity)

	// 2. Auth
	authUc := authUsecase.NewAuthUsecase(repos.Client, antispamService, cfg.JWTSecret, cfg.CompanyName, activityService)
	passwordUc := authUsecase.NewPasswordUsecase(repos.Client)

	// 3. Provisioning Registries
	provisionerFactory := provisioning.NewProvisionerFactory()
	provisionerRegistry := provisioning.NewProvisionerRegistry()
	// Real world provisioners are matched dynamically. Mocks for demo:
	provisionerRegistry.Register("cpanel", provisioning.NewCpanelProvisioner(provisioning.CpanelConfig{Insecure: true}))
	provisionerRegistry.Register("directadmin", provisioning.NewDirectAdminProvisioner("localhost", 2222, "admin", ""))

	registrarRegistry := provisioning.NewRegistrarRegistry()
	registrarRegistry.Register("rdap", provisioning.NewRDAPRegistrarDriver())
	registrarRegistry.Register("email", provisioning.NewEmailRegistrarDriver(emailService, cfg.MailFromAddr))
	registrarRegistry.Register("custom", provisioning.NewCustomRegistrarDriver())

	if cfg.NamecheapAPIUser != "" {
		registrarRegistry.Register("namecheap", provisioning.NewNamecheapRegistrarDriver(provisioning.NamecheapConfig{
			ApiUser:   cfg.NamecheapAPIUser,
			ApiKey:    cfg.NamecheapAPIKey,
			IsSandbox: cfg.AppEnv != "production",
		}))
	}

	dnsRegistry := provisioning.NewDNSProviderRegistry()
	if cfg.CloudflareToken != "" {
		dnsRegistry.Register("cloudflare", provisioning.NewCloudflareDNSProvider(cfg.CloudflareToken))
	}

	// 4. Core Business Logic
	orderService := orderUsecase.NewOrderService(repos.Order, repos.Product, provisionerRegistry, registrarRegistry, eventBus)
	invoiceService := billingUsecase.NewInvoiceService(repos.Invoice, repos.Client, repos.Company, taxCalculator, hookManager, eventBus)
	formbuilderService := formbuilderUsecase.NewFormbuilderService(repos.Formbuilder)
	cartService := cartUsecase.NewCartService(promoCalc, repos.Promo, repos.Order, repos.Product, repos.Client, formbuilderService, taxCalculator, invoiceService, fraudChecker, eventBus)

	// 5. Gateways
	gatewayRegistry := payment.NewGatewayRegistry()
	gatewayRegistry.Register(gateways.NewStripeGateway(cfg.StripeSecretKey, cfg.StripePublicKey, ""))
	gatewayRegistry.Register(gateways.NewMidtransGateway(cfg.MidtransServerKey, cfg.MidtransClientKey, cfg.AppEnv == "production"))
	gatewayRegistry.Register(gateways.NewBankTransferGateway("Bank Transfer", "", cfg.CompanyName))

	// 6. Rest of services
	webhookService := paymentUsecase.NewWebhookService(repos.Transaction, repos.Invoice, gatewayRegistry, eventBus)
	paymentService := paymentUsecase.NewPaymentService(gatewayRegistry, repos.Invoice, repos.Client)
	affiliateService := affiliateUsecase.NewAffiliateService(repos.Affiliate, repos.Client, repos.Order)
	supportService := supportUsecase.NewSupportService(repos.Support, repos.Client, eventBus)
	staffService := staffUsecase.NewStaffService(repos.Staff, cfg.JWTSecret, cfg.CompanyName)
	statsService := statsUsecase.NewStatsService(repos.Client, repos.Order, repos.Invoice, repos.Support, appCache)
	companyService := companyUsecase.NewCompanyService(repos.Company, repos.System)
	currencyService := currencyUsecase.NewCurrencyService(repos.Currency)
	newsService := newsUsecase.NewNewsService(repos.News)
	downloadService := downloadableUsecase.NewDownloadableService(repos.Downloadable, repos.Order, cfg.JWTSecret)
	domainService := domainUsecase.NewDomainService(repos.Order, registrarRegistry, dnsRegistry)
	licenseService := licenseUsecase.NewLicenseService(repos.Order)
	apiKeyService := apikeyUsecase.NewAPIKeyService(repos.APIKey)
	kbService := knowledgebase.NewKBService(repos.KB)
	productService := catalogUsecase.NewProductService(repos.Product, appCache)
	serverService := catalogUsecase.NewServerService(repos.Catalog, provisionerFactory)
	systemService := systemUsecase.NewSystemService(repos.System)
	pageService := pageUsecase.NewPageService(repos.Page)
	notificationService := notificationUsecase.NewNotificationService(repos.Notification)
	adminNotifService := notificationUsecase.NewAdminNotificationService(repos.AdminNotification)
	healthService := systemUsecase.NewHealthUsecase(pool, appCache)
	massMailService := massmailUsecase.NewMassMailService(repos.MassMail, repos.Client, appMailer, cfg.MailFromAddr, cfg.MailFromName)

	// 7. Event Listeners
	orderListener := listener.NewOrderListener(emailService, repos.Order, repos.Product, repos.Client, orderService, affiliateService, registrarRegistry, provisionerRegistry)
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

	telegramService := centralalerts.NewTelegramService(cfg.TelegramBotToken, cfg.TelegramChatID)
	adminAlertListener := listener.NewAdminAlertListener(adminNotifService, telegramService, wsHub)
	eventBus.Subscribe(events.EventInvoicePaid, adminAlertListener.HandleInvoicePaid)
	eventBus.Subscribe(events.EventTicketOpened, adminAlertListener.HandleTicketOpened)
	eventBus.Subscribe(events.EventOrderProvisioningFailed, adminAlertListener.HandleOrderProvisioningFailed)

	systemListener := listener.NewSystemListener(emailService, cfg.MailFromAddr)
	eventBus.Subscribe(events.EventLowStock, systemListener.HandleLowStock)

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
		Affiliate:     affiliateService,
		Company:       companyService,
		Currency:      currencyService,
		News:          newsService,
		Downloadable:  downloadService,
		Domain:        domainService,
		License:       licenseService,
		APIKey:        apiKeyService,
		KB:            kbService,
		Product:       productService,
		Server:        serverService,
		System:        systemService,
		Page:          pageService,
		Activity:      activityService,
		Notification:  notificationService,
		AdminNotif:    adminNotifService,
		Antispam:      antispamService,
		Fraud:         fraudChecker,
		Tax:           taxCalculator,
		Health:        healthService,
		Formbuilder:   formbuilderService,
		Extension:     extensionUsecase.NewExtensionService(repos.Extension),
		Redirect:      redirectUsecase.NewRedirectService(repos.Redirect),
		CookieConsent: cookieconsentUsecase.NewCookieConsentService(repos.Extension),
		Theme:         themeUsecase.NewThemeService(repos.Theme),
		SEO:           seoUsecase.NewSEOService(repos.Page, repos.News, repos.Product),
		Widget:        widgetUsecase.NewWidgetService(),
		Gateways:      gatewayRegistry,
		Hooks:         hookManager,
		Cache:         appCache,
		WSHub:         wsHub,
		MassMail:      massMailService,
		EventBus:      eventBus,
	}
}
