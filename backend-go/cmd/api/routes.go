package main

import (
	"net/http"
	"os"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/config"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/http/admin"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/http/client"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/http/guest"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type AppHandlers struct {
	GuestAuth      *guest.AuthHandler
	GuestCart      *guest.CartHandler
	GuestWebhook   *guest.WebhookHandler
	GuestCurrency  *guest.CurrencyHandler
	GuestNews      *guest.NewsHandler
	GuestCompany   *guest.CompanyHandler
	ClientProfile  *client.ProfileHandler
	ClientOrder    *client.OrderHandler
	ClientInvoice  *client.InvoiceHandler
	ClientDeposit  *client.DepositHandler
	ClientSupport  *client.SupportHandler
	ClientDownload *client.DownloadHandler
	ClientAPIKey   *client.APIKeyHandler
	AdminAuth      *admin.StaffAuthHandler
	AdminStaff     *admin.StaffManagementHandler
	AdminClient    *admin.ClientManagementHandler
	AdminInvoice   *admin.InvoiceManagementHandler
	AdminStats     *admin.StatsHandler
	AdminCurrency  *admin.CurrencyHandler
	AdminNews      *admin.NewsHandler
	AdminMassMail  *admin.MassMailHandler
	AdminCompany   *admin.CompanyHandler
	AdminCatalog   *admin.CatalogHandler
	AdminBilling   *admin.BillingModuleHandler
	AdminSystem    *admin.SystemModuleHandler
}

func setupRoutes(cfg *config.Config, h *AppHandlers, rateLimiter *middleware.RateLimiter) http.Handler {
	mux := http.NewServeMux()

	// System Base
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]interface{}{"status": "ok", "environment": cfg.AppEnv}, nil)
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

	// Public Guest
	mux.Handle("POST /api/v1/guest/auth/register", rateLimiter.RateLimit(http.HandlerFunc(h.GuestAuth.Register)))
	mux.Handle("POST /api/v1/guest/auth/login", rateLimiter.RateLimit(http.HandlerFunc(h.GuestAuth.Login)))
	mux.Handle("POST /api/v1/guest/cart/calculate", rateLimiter.RateLimit(http.HandlerFunc(h.GuestCart.Calculate)))
	mux.Handle("POST /api/v1/guest/cart/checkout", rateLimiter.RateLimit(http.HandlerFunc(h.GuestCart.Checkout)))
	mux.Handle("POST /api/v1/guest/gateways/{gateway}/webhook", http.HandlerFunc(h.GuestWebhook.HandleGatewayWebhook))
	mux.Handle("GET /api/v1/guest/currencies", rateLimiter.RateLimit(http.HandlerFunc(h.GuestCurrency.List)))
	mux.Handle("GET /api/v1/guest/news", rateLimiter.RateLimit(http.HandlerFunc(h.GuestNews.List)))
	mux.Handle("GET /api/v1/guest/news/{slug}", rateLimiter.RateLimit(http.HandlerFunc(h.GuestNews.Get)))
	mux.Handle("GET /api/v1/guest/company", rateLimiter.RateLimit(http.HandlerFunc(h.GuestCompany.GetCompany)))

	// Client Protected
	cAuth := middleware.RequireAuth(cfg.JWTSecret, "client", "admin", "superadmin")
	mux.Handle("GET /api/v1/client/profile", cAuth(http.HandlerFunc(h.ClientProfile.GetProfile)))
	mux.Handle("PUT /api/v1/client/profile", cAuth(http.HandlerFunc(h.ClientProfile.UpdateProfile)))
	mux.Handle("POST /api/v1/client/profile/change-password", cAuth(http.HandlerFunc(h.ClientProfile.ChangePassword)))
	mux.Handle("GET /api/v1/client/orders", cAuth(http.HandlerFunc(h.ClientOrder.ListOrders)))
	mux.Handle("GET /api/v1/client/orders/{id}", cAuth(http.HandlerFunc(h.ClientOrder.GetOrder)))
	mux.Handle("GET /api/v1/client/invoices", cAuth(http.HandlerFunc(h.ClientInvoice.ListInvoices)))
	mux.Handle("GET /api/v1/client/invoices/{id}", cAuth(http.HandlerFunc(h.ClientInvoice.GetInvoice)))
	mux.Handle("GET /api/v1/client/invoices/{id}/pdf", cAuth(http.HandlerFunc(h.ClientInvoice.DownloadPDF)))
	mux.Handle("POST /api/v1/client/invoices/{id}/pay-balance", cAuth(http.HandlerFunc(h.ClientInvoice.PayWithBalance)))
	mux.Handle("POST /api/v1/client/funds/deposit", cAuth(http.HandlerFunc(h.ClientDeposit.DepositFunds)))
	mux.Handle("POST /api/v1/client/support/tickets", cAuth(http.HandlerFunc(h.ClientSupport.OpenTicket)))
	mux.Handle("GET /api/v1/client/support/tickets", cAuth(http.HandlerFunc(h.ClientSupport.ListTickets)))
	mux.Handle("GET /api/v1/client/support/tickets/{id}", cAuth(http.HandlerFunc(h.ClientSupport.GetTicket)))
	mux.Handle("POST /api/v1/client/support/tickets/{id}/reply", cAuth(http.HandlerFunc(h.ClientSupport.ReplyTicket)))
	mux.Handle("POST /api/v1/client/support/tickets/{id}/close", cAuth(http.HandlerFunc(h.ClientSupport.CloseTicket)))
	mux.Handle("GET /api/v1/client/downloads/{id}/link", cAuth(http.HandlerFunc(h.ClientDownload.GenerateLink)))
	mux.Handle("GET /api/v1/client/downloads/{id}/file", http.HandlerFunc(h.ClientDownload.StreamFile))
	mux.Handle("GET /api/v1/client/api-keys", cAuth(http.HandlerFunc(h.ClientAPIKey.List)))
	mux.Handle("POST /api/v1/client/api-keys", cAuth(http.HandlerFunc(h.ClientAPIKey.Generate)))
	mux.Handle("DELETE /api/v1/client/api-keys/{id}", cAuth(http.HandlerFunc(h.ClientAPIKey.Revoke)))

	// Admin Auth & Protected
	mux.Handle("POST /api/v1/admin/auth/login", rateLimiter.RateLimit(http.HandlerFunc(h.AdminAuth.Login)))
	aAuth := middleware.RequireAuth(cfg.JWTSecret, "admin", "superadmin", "support", "billing")
	mux.Handle("GET /api/v1/admin/stats/dashboard", aAuth(http.HandlerFunc(h.AdminStats.GetDashboard)))
	mux.Handle("GET /api/v1/admin/clients", aAuth(http.HandlerFunc(h.AdminClient.ListClients)))
	mux.Handle("POST /api/v1/admin/clients", aAuth(http.HandlerFunc(h.AdminClient.CreateClient)))
	mux.Handle("GET /api/v1/admin/clients/{id}", aAuth(http.HandlerFunc(h.AdminClient.GetClient)))
	mux.Handle("PUT /api/v1/admin/clients/{id}", aAuth(http.HandlerFunc(h.AdminClient.UpdateClient)))
	mux.Handle("DELETE /api/v1/admin/clients/{id}", aAuth(http.HandlerFunc(h.AdminClient.DeleteClient)))
	mux.Handle("GET /api/v1/admin/invoices", aAuth(http.HandlerFunc(h.AdminInvoice.ListInvoices)))
	mux.Handle("POST /api/v1/admin/invoices", aAuth(http.HandlerFunc(h.AdminInvoice.CreateInvoice)))
	mux.Handle("GET /api/v1/admin/orders", aAuth(http.HandlerFunc(h.AdminStaff.ListOrders)))
	mux.Handle("POST /api/v1/admin/orders/{id}/suspend", aAuth(http.HandlerFunc(h.AdminStaff.SuspendOrder)))
	mux.Handle("POST /api/v1/admin/orders/{id}/unsuspend", aAuth(http.HandlerFunc(h.AdminStaff.UnsuspendOrder)))
	mux.Handle("POST /api/v1/admin/orders/{id}/activate", aAuth(http.HandlerFunc(h.AdminStaff.ActivateOrder)))
	mux.Handle("GET /api/v1/admin/support/tickets", aAuth(http.HandlerFunc(h.AdminStaff.ListTickets)))
	mux.Handle("POST /api/v1/admin/support/tickets/{id}/reply", aAuth(http.HandlerFunc(h.AdminStaff.ReplyTicket)))
	mux.Handle("GET /api/v1/admin/audit-logs", aAuth(http.HandlerFunc(h.AdminAuth.GetAuditLogs)))

	// Admin Currencies, News, MassMail, Company
	mux.Handle("GET /api/v1/admin/currencies", aAuth(http.HandlerFunc(h.AdminCurrency.List)))
	mux.Handle("POST /api/v1/admin/currencies", aAuth(http.HandlerFunc(h.AdminCurrency.Create)))
	mux.Handle("PUT /api/v1/admin/currencies/{code}", aAuth(http.HandlerFunc(h.AdminCurrency.Update)))
	mux.Handle("DELETE /api/v1/admin/currencies/{code}", aAuth(http.HandlerFunc(h.AdminCurrency.Delete)))
	mux.Handle("POST /api/v1/admin/currencies/{code}/default", aAuth(http.HandlerFunc(h.AdminCurrency.SetDefault)))
	mux.Handle("GET /api/v1/admin/news", aAuth(http.HandlerFunc(h.AdminNews.List)))
	mux.Handle("POST /api/v1/admin/news", aAuth(http.HandlerFunc(h.AdminNews.Create)))
	mux.Handle("PUT /api/v1/admin/news/{id}", aAuth(http.HandlerFunc(h.AdminNews.Update)))
	mux.Handle("DELETE /api/v1/admin/news/{id}", aAuth(http.HandlerFunc(h.AdminNews.Delete)))
	mux.Handle("GET /api/v1/admin/mass-mail", aAuth(http.HandlerFunc(h.AdminMassMail.List)))
	mux.Handle("POST /api/v1/admin/mass-mail", aAuth(http.HandlerFunc(h.AdminMassMail.Create)))
	mux.Handle("POST /api/v1/admin/mass-mail/{id}/send", aAuth(http.HandlerFunc(h.AdminMassMail.Send)))
	mux.Handle("GET /api/v1/admin/company", aAuth(http.HandlerFunc(h.AdminCompany.GetCompany)))
	mux.Handle("PUT /api/v1/admin/company", aAuth(http.HandlerFunc(h.AdminCompany.UpdateCompany)))

	// Admin Catalog (Products, TLDs, Servers)
	mux.Handle("GET /api/v1/admin/products", aAuth(http.HandlerFunc(h.AdminCatalog.ListProducts)))
	mux.Handle("POST /api/v1/admin/products", aAuth(http.HandlerFunc(h.AdminCatalog.CreateProduct)))
	mux.Handle("PUT /api/v1/admin/products/{id}", aAuth(http.HandlerFunc(h.AdminCatalog.UpdateProduct)))
	mux.Handle("DELETE /api/v1/admin/products/{id}", aAuth(http.HandlerFunc(h.AdminCatalog.DeleteProduct)))
	mux.Handle("GET /api/v1/admin/product-categories", aAuth(http.HandlerFunc(h.AdminCatalog.ListProductCategories)))
	mux.Handle("GET /api/v1/admin/domains/tlds", aAuth(http.HandlerFunc(h.AdminCatalog.ListTlds)))
	mux.Handle("POST /api/v1/admin/domains/tlds", aAuth(http.HandlerFunc(h.AdminCatalog.CreateTld)))
	mux.Handle("DELETE /api/v1/admin/domains/tlds/{id}", aAuth(http.HandlerFunc(h.AdminCatalog.DeleteTld)))
	mux.Handle("GET /api/v1/admin/domains/registrars", aAuth(http.HandlerFunc(h.AdminCatalog.ListRegistrars)))
	mux.Handle("GET /api/v1/admin/servers", aAuth(http.HandlerFunc(h.AdminCatalog.ListServers)))
	mux.Handle("POST /api/v1/admin/servers", aAuth(http.HandlerFunc(h.AdminCatalog.CreateServer)))
	mux.Handle("POST /api/v1/admin/servers/{id}/test", aAuth(http.HandlerFunc(h.AdminCatalog.TestServer)))
	mux.Handle("DELETE /api/v1/admin/servers/{id}", aAuth(http.HandlerFunc(h.AdminCatalog.DeleteServer)))

	// Admin Billing (Gateways, Tax, Coupons, Email, Reports)
	mux.Handle("GET /api/v1/admin/gateways", aAuth(http.HandlerFunc(h.AdminBilling.ListGateways)))
	mux.Handle("GET /api/v1/admin/tax-rules", aAuth(http.HandlerFunc(h.AdminBilling.ListTaxRules)))
	mux.Handle("POST /api/v1/admin/tax-rules", aAuth(http.HandlerFunc(h.AdminBilling.CreateTaxRule)))
	mux.Handle("DELETE /api/v1/admin/tax-rules/{id}", aAuth(http.HandlerFunc(h.AdminBilling.DeleteTaxRule)))
	mux.Handle("GET /api/v1/admin/coupons", aAuth(http.HandlerFunc(h.AdminBilling.ListCoupons)))
	mux.Handle("POST /api/v1/admin/coupons", aAuth(http.HandlerFunc(h.AdminBilling.CreateCoupon)))
	mux.Handle("DELETE /api/v1/admin/coupons/{id}", aAuth(http.HandlerFunc(h.AdminBilling.DeleteCoupon)))
	mux.Handle("GET /api/v1/admin/email-templates", aAuth(http.HandlerFunc(h.AdminBilling.ListEmailTemplates)))
	mux.Handle("GET /api/v1/admin/settings/mail", aAuth(http.HandlerFunc(h.AdminBilling.GetMailConfig)))
	mux.Handle("POST /api/v1/admin/settings/mail/test", aAuth(http.HandlerFunc(h.AdminBilling.SendTestEmail)))
	mux.Handle("GET /api/v1/admin/reports/financial", aAuth(http.HandlerFunc(h.AdminBilling.GetFinancialReports)))

	// Admin System (Security, Health, Pages, KB, Extensions)
	mux.Handle("GET /api/v1/admin/settings/security", aAuth(http.HandlerFunc(h.AdminSystem.GetSecuritySettings)))
	mux.Handle("GET /api/v1/admin/system/status", aAuth(http.HandlerFunc(h.AdminSystem.GetSystemStatus)))
	mux.Handle("POST /api/v1/admin/system/cron/run", aAuth(http.HandlerFunc(h.AdminSystem.TriggerCron)))
	mux.Handle("POST /api/v1/admin/system/cache/clear", aAuth(http.HandlerFunc(h.AdminSystem.ClearCache)))
	mux.Handle("GET /api/v1/admin/pages", aAuth(http.HandlerFunc(h.AdminSystem.ListPages)))
	mux.Handle("GET /api/v1/admin/knowledgebase", aAuth(http.HandlerFunc(h.AdminSystem.ListKnowledgebase)))
	mux.Handle("GET /api/v1/admin/extensions", aAuth(http.HandlerFunc(h.AdminSystem.ListExtensions)))

	return middleware.Logger(middleware.CORS(mux))
}
