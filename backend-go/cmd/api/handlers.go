package main

import (
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/http/admin"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/http/client"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/http/guest"
)

// InitHandlers wires domain services and repositories to HTTP endpoint handlers
func InitHandlers(services *Services, repos *Repositories) *AppHandlers {
	return &AppHandlers{
		GuestAuth:          guest.NewAuthHandler(services.Auth, services.Antispam),
		GuestCart:          guest.NewCartHandler(services.Cart),
		GuestWebhook:       guest.NewWebhookHandler(services.Webhook),
		GuestCurrency:      guest.NewCurrencyHandler(services.Currency),
		GuestNews:          guest.NewNewsHandler(services.News),
		GuestKB:            guest.NewKBHandler(services.KB),
		GuestPage:          guest.NewPageHandler(services.Page),
		GuestCompany:       guest.NewCompanyHandler(services.Company),
		GuestDomain:        guest.NewDomainHandler(services.Domain),
		GuestFormbuilder:   guest.NewFormbuilderHandler(services.Formbuilder),
		GuestRedirect:      guest.NewRedirectHandler(services.Redirect),
		GuestCookieConsent: guest.NewCookieConsentHandler(services.CookieConsent),
		GuestTheme:         guest.NewThemeHandler(services.Theme),
		GuestSEO:           guest.NewSEOHandler(services.SEO),
		GuestWidget:        guest.NewWidgetHandler(services.Widget),
		ClientProfile:      client.NewProfileHandler(services.Auth, services.Password),
		ClientOrder:        client.NewOrderHandler(services.Order),
		ClientDomain:       client.NewDomainHandler(services.Domain),
		ClientInvoice:      client.NewInvoiceHandler(repos.Invoice, repos.Client, services.Invoice, services.Payment),
		ClientDeposit:      client.NewDepositHandler(services.Invoice, services.Payment, repos.Client),
		ClientSupport:      client.NewSupportHandler(services.Support),
		ClientActivity:     client.NewActivityHandler(services.Activity),
		ClientNotification: client.NewNotificationHandler(services.Notification),
		ClientDownload:     client.NewDownloadHandler(services.Downloadable),
		ClientLicense:      client.NewLicenseHandler(services.License),
		ClientAPIKey:       client.NewAPIKeyHandler(services.APIKey),
		AdminAuth:          admin.NewStaffAuthHandler(services.Staff),
		AdminStaff:         admin.NewStaffManagementHandler(services.Staff, repos.Client, repos.Order, services.Order, services.Support),
		AdminClient:        admin.NewClientManagementHandler(services.Staff, repos.Client),
		AdminInvoice:       admin.NewInvoiceManagementHandler(services.Staff, repos.Invoice, repos.Client, services.Invoice),
		AdminStats:         admin.NewStatsHandler(services.Stats, services.Staff),
		AdminCurrency:      admin.NewCurrencyHandler(services.Currency, services.Staff),
		AdminNews:          admin.NewNewsHandler(services.News, services.Staff),
		AdminKB:            admin.NewKBHandler(services.Staff, services.KB),
		AdminMassMail:      admin.NewMassMailHandler(services.MassMail, services.Staff),
		AdminCompany:       admin.NewCompanyHandler(services.Company, services.Staff),
		AdminCatalog:       admin.NewCatalogHandler(services.Staff, services.Product, services.Server, repos.Catalog),
		AdminBilling:       admin.NewBillingModuleHandler(services.Staff, services.Stats, services.Tax, repos.Promo),
		AdminSystem:        admin.NewSystemModuleHandler(services.Staff, services.System, services.Page),
		AdminActivity:      admin.NewActivityHandler(services.Staff, services.Activity),
		AdminAntispam:      admin.NewAntispamHandler(services.Staff, services.Antispam),
		AdminFormbuilder:   admin.NewFormbuilderHandler(services.Staff, services.Formbuilder),
		AdminExtension:     admin.NewExtensionHandler(services.Staff, services.Extension),
		AdminRedirect:      admin.NewRedirectHandler(services.Staff, services.Redirect),
		AdminCookieConsent: admin.NewCookieConsentHandler(services.Staff, services.CookieConsent),
		AdminNotification:  admin.NewAdminNotificationHandler(services.Staff, services.AdminNotif),
		AdminTheme:         admin.NewThemeHandler(services.Staff, services.Theme),
		AdminSEO:           admin.NewSEOHandler(services.Staff, services.SEO),
		AdminWidget:        admin.NewWidgetHandler(services.Staff, services.Widget),
	}
}
