package main

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
)

// registerAdminRoutes configures all staff and administrator protected endpoints
func registerAdminRoutes(mux *http.ServeMux, h *AppHandlers, aAuth func(http.Handler) http.Handler, rateLimiter *middleware.RateLimiter) {
	// Admin Auth
	mux.Handle("POST /api/v1/admin/auth/login", rateLimiter.RateLimit(http.HandlerFunc(h.AdminAuth.Login)))

	// Core Operations & Clients
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

	// System (Security, Health, Pages, KB, Extensions)
	mux.Handle("GET /api/v1/admin/settings/security", aAuth(http.HandlerFunc(h.AdminSystem.GetSecuritySettings)))
	mux.Handle("GET /api/v1/admin/system/status", aAuth(http.HandlerFunc(h.AdminSystem.GetSystemStatus)))
	mux.Handle("POST /api/v1/admin/system/cron/run", aAuth(http.HandlerFunc(h.AdminSystem.TriggerCron)))
	mux.Handle("POST /api/v1/admin/system/cache/clear", aAuth(http.HandlerFunc(h.AdminSystem.ClearCache)))
	mux.Handle("GET /api/v1/admin/pages", aAuth(http.HandlerFunc(h.AdminSystem.ListPages)))
	mux.Handle("GET /api/v1/admin/knowledgebase", aAuth(http.HandlerFunc(h.AdminSystem.ListKnowledgebase)))
	mux.Handle("GET /api/v1/admin/extensions", aAuth(http.HandlerFunc(h.AdminSystem.ListExtensions)))
}
