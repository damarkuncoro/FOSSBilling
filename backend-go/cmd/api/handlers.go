package main

import (
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/http/admin"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/http/client"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/http/guest"
)

// InitHandlers wires domain services and repositories to HTTP endpoint handlers
func InitHandlers(services *Services, repos *Repositories) *AppHandlers {
	return &AppHandlers{
		GuestAuth:      guest.NewAuthHandler(services.Auth),
		GuestCart:      guest.NewCartHandler(services.Cart),
		GuestWebhook:   guest.NewWebhookHandler(services.Webhook),
		GuestCurrency:  guest.NewCurrencyHandler(services.Currency),
		GuestNews:      guest.NewNewsHandler(services.News),
		GuestCompany:   guest.NewCompanyHandler(services.Company),
		ClientProfile:  client.NewProfileHandler(services.Auth, services.Password),
		ClientOrder:    client.NewOrderHandler(repos.Order),
		ClientInvoice:  client.NewInvoiceHandler(repos.Invoice, repos.Client, services.Invoice),
		ClientDeposit:  client.NewDepositHandler(services.Invoice),
		ClientSupport:  client.NewSupportHandler(services.Support),
		ClientDownload: client.NewDownloadHandler(services.Downloadable),
		ClientAPIKey:   client.NewAPIKeyHandler(services.APIKey),
		AdminAuth:      admin.NewStaffAuthHandler(services.Staff),
		AdminStaff:     admin.NewStaffManagementHandler(services.Staff, repos.Client, repos.Order, services.Order, services.Support),
		AdminClient:    admin.NewClientManagementHandler(services.Staff, repos.Client),
		AdminInvoice:   admin.NewInvoiceManagementHandler(repos.Invoice, repos.Client, services.Invoice),
		AdminStats:     admin.NewStatsHandler(services.Stats, services.Staff),
		AdminCurrency:  admin.NewCurrencyHandler(services.Currency, services.Staff),
		AdminNews:      admin.NewNewsHandler(services.News, services.Staff),
		AdminMassMail:  admin.NewMassMailHandler(services.MassMail, services.Staff),
		AdminCompany:   admin.NewCompanyHandler(services.Company, services.Staff),
		AdminCatalog:   admin.NewCatalogHandler(),
		AdminBilling:   admin.NewBillingModuleHandler(services.Stats),
		AdminSystem:    admin.NewSystemModuleHandler(),
	}
}
