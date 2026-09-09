package main

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
)

// registerAdminRoutes configures all staff and administrator protected endpoints
func registerAdminRoutes(mux *http.ServeMux, h *AppHandlers, aAuth func(http.Handler) http.Handler, rateLimiter *middleware.RateLimiter) {
	// Admin Auth
	mux.Handle("POST /api/v1/admin/auth/login", rateLimiter.RateLimit(http.HandlerFunc(h.AdminAuth.Login)))
	mux.Handle("POST /api/v1/admin/auth/verify-2fa", rateLimiter.RateLimit(http.HandlerFunc(h.AdminAuth.VerifyTwoFactor)))
	mux.Handle("POST /api/v1/admin/auth/2fa/setup", aAuth(http.HandlerFunc(h.AdminAuth.SetupTwoFactor)))
	mux.Handle("POST /api/v1/admin/auth/2fa/enable", aAuth(http.HandlerFunc(h.AdminAuth.EnableTwoFactor)))
	mux.Handle("POST /api/v1/admin/auth/2fa/disable", aAuth(http.HandlerFunc(h.AdminAuth.DisableTwoFactor)))

	// Core Operations & Clients
	mux.Handle("GET /api/v1/admin/stats/dashboard", aAuth(http.HandlerFunc(h.AdminStats.GetDashboard)))
	mux.Handle("GET /api/v1/admin/clients", aAuth(http.HandlerFunc(h.AdminClient.ListClients)))
	mux.Handle("POST /api/v1/admin/clients", aAuth(http.HandlerFunc(h.AdminClient.CreateClient)))
	mux.Handle("GET /api/v1/admin/clients/{id}", aAuth(http.HandlerFunc(h.AdminClient.GetClient)))
	mux.Handle("PUT /api/v1/admin/clients/{id}", aAuth(http.HandlerFunc(h.AdminClient.UpdateClient)))
	mux.Handle("DELETE /api/v1/admin/clients/{id}", aAuth(http.HandlerFunc(h.AdminClient.DeleteClient)))
	mux.Handle("GET /api/v1/admin/invoices", aAuth(http.HandlerFunc(h.AdminInvoice.ListInvoices)))
	mux.Handle("POST /api/v1/admin/invoices", aAuth(http.HandlerFunc(h.AdminInvoice.CreateInvoice)))
	mux.Handle("GET /api/v1/admin/invoices/{id}", aAuth(http.HandlerFunc(h.AdminInvoice.GetInvoice)))
	mux.Handle("POST /api/v1/admin/invoices/{id}/refund", aAuth(http.HandlerFunc(h.AdminInvoice.RefundInvoice)))
	mux.Handle("DELETE /api/v1/admin/invoices/{id}", aAuth(http.HandlerFunc(h.AdminInvoice.DeleteInvoice)))
	mux.Handle("GET /api/v1/admin/orders", aAuth(http.HandlerFunc(h.AdminStaff.ListOrders)))
	mux.Handle("POST /api/v1/admin/orders/{id}/suspend", aAuth(http.HandlerFunc(h.AdminStaff.SuspendOrder)))
	mux.Handle("POST /api/v1/admin/orders/{id}/unsuspend", aAuth(http.HandlerFunc(h.AdminStaff.UnsuspendOrder)))
	mux.Handle("POST /api/v1/admin/orders/{id}/activate", aAuth(http.HandlerFunc(h.AdminStaff.ActivateOrder)))
	mux.Handle("POST /api/v1/admin/orders/{id}/sync", aAuth(http.HandlerFunc(h.AdminStaff.SyncOrder)))
	mux.Handle("POST /api/v1/admin/orders/{id}/change-password", aAuth(http.HandlerFunc(h.AdminStaff.ChangeOrderPassword)))
	mux.Handle("GET /api/v1/admin/support/tickets", aAuth(http.HandlerFunc(h.AdminStaff.ListTickets)))
	mux.Handle("POST /api/v1/admin/support/tickets/{id}/reply", aAuth(http.HandlerFunc(h.AdminStaff.ReplyTicket)))
	mux.Handle("GET /api/v1/admin/activity", aAuth(http.HandlerFunc(h.AdminActivity.ListLogs)))

	// Currencies, News, MassMail, Company
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

	// Catalog (Products, TLDs, Servers)
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

	// Billing (Gateways, Tax, Coupons, Email, Reports)
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
	mux.Handle("GET /api/v1/admin/reports/invoices/csv", aAuth(http.HandlerFunc(h.AdminBilling.ExportInvoicesCSV)))

	// System (Security, Health, Pages, KB)
	mux.Handle("GET /api/v1/admin/settings/security", aAuth(http.HandlerFunc(h.AdminSystem.GetSecuritySettings)))
	mux.Handle("PUT /api/v1/admin/settings/security", aAuth(http.HandlerFunc(h.AdminSystem.UpdateSecuritySettings)))
	mux.Handle("GET /api/v1/admin/notifications", aAuth(http.HandlerFunc(h.AdminNotification.List)))
	mux.Handle("PUT /api/v1/admin/notifications/{id}/read", aAuth(http.HandlerFunc(h.AdminNotification.MarkRead)))
	mux.Handle("POST /api/v1/admin/notifications/mark-all-read", aAuth(http.HandlerFunc(h.AdminNotification.MarkAllRead)))
	mux.Handle("GET /api/v1/admin/system/status", aAuth(http.HandlerFunc(h.AdminSystem.GetSystemStatus)))
	mux.Handle("POST /api/v1/admin/system/cron/run", aAuth(http.HandlerFunc(h.AdminSystem.TriggerCron)))
	mux.Handle("POST /api/v1/admin/system/cache/clear", aAuth(http.HandlerFunc(h.AdminSystem.ClearCache)))
	mux.Handle("GET /api/v1/admin/pages", aAuth(http.HandlerFunc(h.AdminSystem.ListPages)))
	mux.Handle("GET /api/v1/admin/knowledgebase/articles", aAuth(http.HandlerFunc(h.AdminKB.ListArticles)))
	mux.Handle("POST /api/v1/admin/knowledgebase/articles", aAuth(http.HandlerFunc(h.AdminKB.CreateArticle)))
	mux.Handle("PUT /api/v1/admin/knowledgebase/articles/{id}", aAuth(http.HandlerFunc(h.AdminKB.UpdateArticle)))
	mux.Handle("DELETE /api/v1/admin/knowledgebase/articles/{id}", aAuth(http.HandlerFunc(h.AdminKB.DeleteArticle)))
	mux.Handle("GET /api/v1/admin/knowledgebase/categories", aAuth(http.HandlerFunc(h.AdminKB.ListCategories)))
	mux.Handle("POST /api/v1/admin/knowledgebase/categories", aAuth(http.HandlerFunc(h.AdminKB.CreateCategory)))

	// Extensions & Marketplace Hub
	mux.Handle("GET /api/v1/admin/extensions", aAuth(http.HandlerFunc(h.AdminExtension.ListExtensions)))
	mux.Handle("GET /api/v1/admin/extensions/marketplace", aAuth(http.HandlerFunc(h.AdminExtension.ListMarketplace)))
	mux.Handle("GET /api/v1/admin/extensions/marketplace/{id}/readme", aAuth(http.HandlerFunc(h.AdminExtension.GetMarketplaceReadme)))
	mux.Handle("GET /api/v1/admin/extensions/{id}", aAuth(http.HandlerFunc(h.AdminExtension.GetExtension)))
	mux.Handle("POST /api/v1/admin/extensions/{id}/activate", aAuth(http.HandlerFunc(h.AdminExtension.Activate)))
	mux.Handle("POST /api/v1/admin/extensions/{id}/deactivate", aAuth(http.HandlerFunc(h.AdminExtension.Deactivate)))
	mux.Handle("POST /api/v1/admin/extensions/{id}/install", aAuth(http.HandlerFunc(h.AdminExtension.Install)))
	mux.Handle("POST /api/v1/admin/extensions/{id}/uninstall", aAuth(http.HandlerFunc(h.AdminExtension.Uninstall)))
	mux.Handle("GET /api/v1/admin/extensions/{id}/config", aAuth(http.HandlerFunc(h.AdminExtension.GetConfig)))
	mux.Handle("PUT /api/v1/admin/extensions/{id}/config", aAuth(http.HandlerFunc(h.AdminExtension.UpdateConfig)))

	// Antispam & Abuse Protection
	mux.Handle("GET /api/v1/admin/antispam/config", aAuth(http.HandlerFunc(h.AdminAntispam.GetConfig)))
	mux.Handle("PUT /api/v1/admin/antispam/config", aAuth(http.HandlerFunc(h.AdminAntispam.UpdateConfig)))
	mux.Handle("GET /api/v1/admin/antispam/blocked-ips", aAuth(http.HandlerFunc(h.AdminAntispam.ListBlockedIPs)))
	mux.Handle("POST /api/v1/admin/antispam/blocked-ips", aAuth(http.HandlerFunc(h.AdminAntispam.AddBlockedIP)))
	mux.Handle("DELETE /api/v1/admin/antispam/blocked-ips/{ip}", aAuth(http.HandlerFunc(h.AdminAntispam.DeleteBlockedIP)))

	// Formbuilder (Custom Order Forms & Dynamic Fields)
	mux.Handle("GET /api/v1/admin/forms", aAuth(http.HandlerFunc(h.AdminFormbuilder.ListForms)))
	mux.Handle("POST /api/v1/admin/forms", aAuth(http.HandlerFunc(h.AdminFormbuilder.CreateForm)))
	mux.Handle("GET /api/v1/admin/forms/{id}", aAuth(http.HandlerFunc(h.AdminFormbuilder.GetForm)))
	mux.Handle("PUT /api/v1/admin/forms/{id}", aAuth(http.HandlerFunc(h.AdminFormbuilder.UpdateForm)))
	mux.Handle("DELETE /api/v1/admin/forms/{id}", aAuth(http.HandlerFunc(h.AdminFormbuilder.DeleteForm)))
	mux.Handle("POST /api/v1/admin/forms/{id}/fields", aAuth(http.HandlerFunc(h.AdminFormbuilder.AddField)))
	mux.Handle("PUT /api/v1/admin/forms/fields/{field_id}", aAuth(http.HandlerFunc(h.AdminFormbuilder.UpdateField)))
	mux.Handle("DELETE /api/v1/admin/forms/fields/{field_id}", aAuth(http.HandlerFunc(h.AdminFormbuilder.DeleteField)))

	// Redirects Management (Vanity Links & HTTP Forwarding Rules)
	mux.Handle("GET /api/v1/admin/redirects", aAuth(http.HandlerFunc(h.AdminRedirect.ListRedirects)))
	mux.Handle("POST /api/v1/admin/redirects", aAuth(http.HandlerFunc(h.AdminRedirect.CreateRedirect)))
	mux.Handle("GET /api/v1/admin/redirects/{id}", aAuth(http.HandlerFunc(h.AdminRedirect.GetRedirect)))
	mux.Handle("PUT /api/v1/admin/redirects/{id}", aAuth(http.HandlerFunc(h.AdminRedirect.UpdateRedirect)))
	mux.Handle("DELETE /api/v1/admin/redirects/{id}", aAuth(http.HandlerFunc(h.AdminRedirect.DeleteRedirect)))

	// Cookie Consent Configuration
	mux.Handle("GET /api/v1/admin/cookie-consent", aAuth(http.HandlerFunc(h.AdminCookieConsent.GetConfig)))
	mux.Handle("PUT /api/v1/admin/cookie-consent", aAuth(http.HandlerFunc(h.AdminCookieConsent.UpdateConfig)))

	// Theme Management
	mux.Handle("GET /api/v1/admin/themes", aAuth(http.HandlerFunc(h.AdminTheme.ListThemes)))
	mux.Handle("GET /api/v1/admin/themes/current", aAuth(http.HandlerFunc(h.AdminTheme.GetCurrentTheme)))
	mux.Handle("POST /api/v1/admin/themes/select", aAuth(http.HandlerFunc(h.AdminTheme.SelectTheme)))
	mux.Handle("GET /api/v1/admin/themes/{code}/config", aAuth(http.HandlerFunc(h.AdminTheme.GetConfig)))
	mux.Handle("PUT /api/v1/admin/themes/{code}/config", aAuth(http.HandlerFunc(h.AdminTheme.UpdateConfig)))

	// SEO & Search Engine Management
	mux.Handle("GET /api/v1/admin/seo/info", aAuth(http.HandlerFunc(h.AdminSEO.GetInfo)))
	mux.Handle("POST /api/v1/admin/seo/ping", aAuth(http.HandlerFunc(h.AdminSEO.PingSearchEngines)))

	// Widget Registry
	mux.Handle("GET /api/v1/admin/widgets", aAuth(http.HandlerFunc(h.AdminWidget.GetRegistry)))

	// System Tools & GeoIP Utilities
	mux.Handle("GET /api/v1/admin/system/tools/password", aAuth(http.HandlerFunc(h.AdminSystem.GeneratePassword)))
	mux.Handle("GET /api/v1/admin/system/geoip", aAuth(http.HandlerFunc(h.AdminSystem.ResolveGeoIP)))
}
