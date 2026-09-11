package main

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
)

func registerAdminRoutes(mux *http.ServeMux, h *AppHandlers, aAuth func(http.Handler) http.Handler, rateLimiter, authRateLimiter *middleware.RateLimiter) {
	authRoutes(mux, h, authRateLimiter, aAuth)
	clientRoutes(mux, h, aAuth)
	billingRoutes(mux, h, aAuth)
	catalogRoutes(mux, h, aAuth)
	systemRoutes(mux, h, aAuth)
	extensionRoutes(mux, h, aAuth)
	contentRoutes(mux, h, aAuth)
}

func authRoutes(mux *http.ServeMux, h *AppHandlers, arl *middleware.RateLimiter, aAuth func(http.Handler) http.Handler) {
	mux.Handle("POST /api/v1/admin/auth/login", arl.RateLimit(http.HandlerFunc(h.AdminAuth.Login)))
	mux.Handle("POST /api/v1/admin/auth/verify-2fa", arl.RateLimit(http.HandlerFunc(h.AdminAuth.VerifyTwoFactor)))
	mux.Handle("POST /api/v1/admin/auth/2fa/setup", aAuth(http.HandlerFunc(h.AdminAuth.SetupTwoFactor)))
	mux.Handle("POST /api/v1/admin/auth/2fa/enable", aAuth(http.HandlerFunc(h.AdminAuth.EnableTwoFactor)))
	mux.Handle("POST /api/v1/admin/auth/2fa/disable", aAuth(http.HandlerFunc(h.AdminAuth.DisableTwoFactor)))
	mux.Handle("GET /api/v1/admin/staff", aAuth(http.HandlerFunc(h.AdminStaff.ListStaff)))
	mux.Handle("GET /api/v1/admin/system/audit-logs", aAuth(http.HandlerFunc(h.AdminAuth.GetAuditLogs)))
}

func clientRoutes(mux *http.ServeMux, h *AppHandlers, aAuth func(http.Handler) http.Handler) {
	mux.Handle("GET /api/v1/admin/clients", aAuth(http.HandlerFunc(h.AdminClient.ListClients)))
	mux.Handle("POST /api/v1/admin/clients", aAuth(http.HandlerFunc(h.AdminClient.CreateClient)))
	mux.Handle("GET /api/v1/admin/clients/{id}", aAuth(http.HandlerFunc(h.AdminClient.GetClient)))
	mux.Handle("PUT /api/v1/admin/clients/{id}", aAuth(http.HandlerFunc(h.AdminClient.UpdateClient)))
	mux.Handle("DELETE /api/v1/admin/clients/{id}", aAuth(http.HandlerFunc(h.AdminClient.DeleteClient)))
	mux.Handle("POST /api/v1/admin/clients/{id}/impersonate", aAuth(http.HandlerFunc(h.AdminClient.ImpersonateClient)))
}

func billingRoutes(mux *http.ServeMux, h *AppHandlers, aAuth func(http.Handler) http.Handler) {
	mux.Handle("GET /api/v1/admin/invoices", aAuth(http.HandlerFunc(h.AdminInvoice.ListInvoices)))
	mux.Handle("POST /api/v1/admin/invoices", aAuth(http.HandlerFunc(h.AdminInvoice.CreateInvoice)))
	mux.Handle("GET /api/v1/admin/invoices/{id}", aAuth(http.HandlerFunc(h.AdminInvoice.GetInvoice)))
	mux.Handle("GET /api/v1/admin/invoices/{id}/pdf", aAuth(http.HandlerFunc(h.AdminInvoice.DownloadPDF)))
	mux.Handle("POST /api/v1/admin/invoices/{id}/refund", aAuth(http.HandlerFunc(h.AdminInvoice.RefundInvoice)))
	mux.Handle("DELETE /api/v1/admin/invoices/{id}", aAuth(http.HandlerFunc(h.AdminInvoice.DeleteInvoice)))
	mux.Handle("GET /api/v1/admin/gateways", aAuth(http.HandlerFunc(h.AdminBilling.ListGateways)))
	mux.Handle("GET /api/v1/admin/tax-rules", aAuth(http.HandlerFunc(h.AdminBilling.ListTaxRules)))
	mux.Handle("POST /api/v1/admin/tax-rules", aAuth(http.HandlerFunc(h.AdminBilling.CreateTaxRule)))
	mux.Handle("PUT /api/v1/admin/tax-rules/{id}", aAuth(http.HandlerFunc(h.AdminBilling.UpdateTaxRule)))
	mux.Handle("DELETE /api/v1/admin/tax-rules/{id}", aAuth(http.HandlerFunc(h.AdminBilling.DeleteTaxRule)))
	mux.Handle("GET /api/v1/admin/taxes", aAuth(http.HandlerFunc(h.AdminBilling.ListTaxRules)))
	mux.Handle("POST /api/v1/admin/taxes", aAuth(http.HandlerFunc(h.AdminBilling.CreateTaxRule)))
	mux.Handle("PUT /api/v1/admin/taxes/{id}", aAuth(http.HandlerFunc(h.AdminBilling.UpdateTaxRule)))
	mux.Handle("DELETE /api/v1/admin/taxes/{id}", aAuth(http.HandlerFunc(h.AdminBilling.DeleteTaxRule)))
	mux.Handle("GET /api/v1/admin/coupons", aAuth(http.HandlerFunc(h.AdminBilling.ListCoupons)))
	mux.Handle("POST /api/v1/admin/coupons", aAuth(http.HandlerFunc(h.AdminBilling.CreateCoupon)))
	mux.Handle("DELETE /api/v1/admin/coupons/{id}", aAuth(http.HandlerFunc(h.AdminBilling.DeleteCoupon)))
	mux.Handle("GET /api/v1/admin/reports/financial", aAuth(http.HandlerFunc(h.AdminBilling.GetFinancialReports)))
	mux.Handle("GET /api/v1/admin/reports/invoices/csv", aAuth(http.HandlerFunc(h.AdminBilling.ExportInvoicesCSV)))
	mux.Handle("GET /api/v1/admin/affiliates/{id}", aAuth(http.HandlerFunc(h.AdminAffiliate.GetAffiliate)))
}

func catalogRoutes(mux *http.ServeMux, h *AppHandlers, aAuth func(http.Handler) http.Handler) {
	mux.Handle("GET /api/v1/admin/products", aAuth(http.HandlerFunc(h.AdminCatalog.ListProducts)))
	mux.Handle("POST /api/v1/admin/products", aAuth(http.HandlerFunc(h.AdminCatalog.CreateProduct)))
	mux.Handle("PUT /api/v1/admin/products/{id}", aAuth(http.HandlerFunc(h.AdminCatalog.UpdateProduct)))
	mux.Handle("DELETE /api/v1/admin/products/{id}", aAuth(http.HandlerFunc(h.AdminCatalog.DeleteProduct)))
	mux.Handle("GET /api/v1/admin/products/categories", aAuth(http.HandlerFunc(h.AdminCatalog.ListProductCategories)))
	mux.Handle("GET /api/v1/admin/domains/tlds", aAuth(http.HandlerFunc(h.AdminCatalog.ListTlds)))
	mux.Handle("POST /api/v1/admin/domains/tlds", aAuth(http.HandlerFunc(h.AdminCatalog.CreateTld)))
	mux.Handle("DELETE /api/v1/admin/domains/tlds/{id}", aAuth(http.HandlerFunc(h.AdminCatalog.DeleteTld)))
	mux.Handle("GET /api/v1/admin/domains/registrars", aAuth(http.HandlerFunc(h.AdminCatalog.ListRegistrars)))
	mux.Handle("GET /api/v1/admin/servers", aAuth(http.HandlerFunc(h.AdminCatalog.ListServers)))
	mux.Handle("POST /api/v1/admin/servers", aAuth(http.HandlerFunc(h.AdminCatalog.CreateServer)))
	mux.Handle("POST /api/v1/admin/servers/{id}/test", aAuth(http.HandlerFunc(h.AdminCatalog.TestServer)))
	mux.Handle("DELETE /api/v1/admin/servers/{id}", aAuth(http.HandlerFunc(h.AdminCatalog.DeleteServer)))
}

func systemRoutes(mux *http.ServeMux, h *AppHandlers, aAuth func(http.Handler) http.Handler) {
	mux.Handle("GET /api/v1/admin/system/ws", aAuth(http.HandlerFunc(h.AdminSystem.HandleWebSocket)))
	mux.Handle("GET /api/v1/admin/stats/dashboard", aAuth(http.HandlerFunc(h.AdminStats.GetDashboard)))
	mux.Handle("GET /api/v1/admin/stats/projections", aAuth(http.HandlerFunc(h.AdminStats.GetRevenueProjection)))
	mux.Handle("GET /api/v1/admin/system/status", aAuth(http.HandlerFunc(h.AdminSystem.GetSystemStatus)))
	mux.Handle("POST /api/v1/admin/system/cron/run", aAuth(http.HandlerFunc(h.AdminSystem.TriggerCron)))
	mux.Handle("POST /api/v1/admin/system/cache/clear", aAuth(http.HandlerFunc(h.AdminSystem.ClearCache)))
	mux.Handle("POST /api/v1/admin/system/backup/export", aAuth(http.HandlerFunc(h.AdminSystem.ExportBackup)))
	mux.Handle("GET /api/v1/admin/settings/security", aAuth(http.HandlerFunc(h.AdminSystem.GetSecuritySettings)))
	mux.Handle("PUT /api/v1/admin/settings/security", aAuth(http.HandlerFunc(h.AdminSystem.UpdateSecuritySettings)))
	mux.Handle("GET /api/v1/admin/settings/branding", aAuth(http.HandlerFunc(h.AdminSystem.GetBrandingSettings)))
	mux.Handle("PUT /api/v1/admin/settings/branding", aAuth(http.HandlerFunc(h.AdminSystem.UpdateBrandingSettings)))
	mux.Handle("GET /api/v1/admin/settings/email-templates", aAuth(http.HandlerFunc(h.AdminEmailTemplate.List)))
	mux.Handle("GET /api/v1/admin/settings/email-templates/{code}", aAuth(http.HandlerFunc(h.AdminEmailTemplate.Get)))
	mux.Handle("PUT /api/v1/admin/settings/email-templates", aAuth(http.HandlerFunc(h.AdminEmailTemplate.Update)))
	mux.Handle("GET /api/v1/admin/currencies", aAuth(http.HandlerFunc(h.AdminCurrency.List)))
	mux.Handle("POST /api/v1/admin/currencies", aAuth(http.HandlerFunc(h.AdminCurrency.Create)))
	mux.Handle("PUT /api/v1/admin/currencies/{code}", aAuth(http.HandlerFunc(h.AdminCurrency.Update)))
	mux.Handle("DELETE /api/v1/admin/currencies/{code}", aAuth(http.HandlerFunc(h.AdminCurrency.Delete)))
	mux.Handle("POST /api/v1/admin/currencies/sync", aAuth(http.HandlerFunc(h.AdminCurrency.SyncRates)))
	mux.Handle("POST /api/v1/admin/currencies/{code}/default", aAuth(http.HandlerFunc(h.AdminCurrency.SetDefault)))
	mux.Handle("GET /api/v1/admin/company", aAuth(http.HandlerFunc(h.AdminCompany.GetCompany)))
	mux.Handle("PUT /api/v1/admin/company", aAuth(http.HandlerFunc(h.AdminCompany.UpdateCompany)))
	mux.Handle("GET /api/v1/admin/activity", aAuth(http.HandlerFunc(h.AdminActivity.ListLogs)))
	mux.Handle("GET /api/v1/admin/activity/trend", aAuth(http.HandlerFunc(h.AdminActivity.GetActivityTrend)))
}

func extensionRoutes(mux *http.ServeMux, h *AppHandlers, aAuth func(http.Handler) http.Handler) {
	mux.Handle("GET /api/v1/admin/extensions", aAuth(http.HandlerFunc(h.AdminExtension.ListExtensions)))
	mux.Handle("POST /api/v1/admin/extensions/{id}/activate", aAuth(http.HandlerFunc(h.AdminExtension.Activate)))
	mux.Handle("POST /api/v1/admin/extensions/{id}/deactivate", aAuth(http.HandlerFunc(h.AdminExtension.Deactivate)))
	mux.Handle("GET /api/v1/admin/forms", aAuth(http.HandlerFunc(h.AdminFormbuilder.ListForms)))
	mux.Handle("POST /api/v1/admin/forms", aAuth(http.HandlerFunc(h.AdminFormbuilder.CreateForm)))
	mux.Handle("GET /api/v1/admin/forms/{id}", aAuth(http.HandlerFunc(h.AdminFormbuilder.GetForm)))
	mux.Handle("PUT /api/v1/admin/forms/{id}", aAuth(http.HandlerFunc(h.AdminFormbuilder.UpdateForm)))
	mux.Handle("DELETE /api/v1/admin/forms/{id}", aAuth(http.HandlerFunc(h.AdminFormbuilder.DeleteForm)))
	mux.Handle("POST /api/v1/admin/forms/{id}/fields", aAuth(http.HandlerFunc(h.AdminFormbuilder.AddField)))
	mux.Handle("PUT /api/v1/admin/forms/fields/{field_id}", aAuth(http.HandlerFunc(h.AdminFormbuilder.UpdateField)))
	mux.Handle("DELETE /api/v1/admin/forms/fields/{field_id}", aAuth(http.HandlerFunc(h.AdminFormbuilder.DeleteField)))
}

func contentRoutes(mux *http.ServeMux, h *AppHandlers, aAuth func(http.Handler) http.Handler) {
	mux.Handle("GET /api/v1/admin/news", aAuth(http.HandlerFunc(h.AdminNews.List)))
	mux.Handle("POST /api/v1/admin/news", aAuth(http.HandlerFunc(h.AdminNews.Create)))
	mux.Handle("PUT /api/v1/admin/news/{id}", aAuth(http.HandlerFunc(h.AdminNews.Update)))
	mux.Handle("DELETE /api/v1/admin/news/{id}", aAuth(http.HandlerFunc(h.AdminNews.Delete)))
	mux.Handle("GET /api/v1/admin/pages", aAuth(http.HandlerFunc(h.AdminSystem.ListPages)))
	mux.Handle("GET /api/v1/admin/knowledgebase/articles", aAuth(http.HandlerFunc(h.AdminKB.ListArticles)))
	mux.Handle("POST /api/v1/admin/knowledgebase/articles", aAuth(http.HandlerFunc(h.AdminKB.CreateArticle)))
	mux.Handle("PUT /api/v1/admin/knowledgebase/articles/{id}", aAuth(http.HandlerFunc(h.AdminKB.UpdateArticle)))
	mux.Handle("DELETE /api/v1/admin/knowledgebase/articles/{id}", aAuth(http.HandlerFunc(h.AdminKB.DeleteArticle)))
	mux.Handle("GET /api/v1/admin/mass-mail", aAuth(http.HandlerFunc(h.AdminMassMail.List)))
	mux.Handle("POST /api/v1/admin/mass-mail", aAuth(http.HandlerFunc(h.AdminMassMail.Create)))
	mux.Handle("POST /api/v1/admin/mass-mail/{id}/send", aAuth(http.HandlerFunc(h.AdminMassMail.Send)))
	mux.Handle("GET /api/v1/admin/notifications", aAuth(http.HandlerFunc(h.AdminNotification.List)))
	mux.Handle("GET /api/v1/admin/orders", aAuth(http.HandlerFunc(h.AdminStaff.ListOrders)))
	mux.Handle("GET /api/v1/admin/orders/{id}", aAuth(http.HandlerFunc(h.AdminStaff.GetOrder)))
	mux.Handle("POST /api/v1/admin/orders/{id}/suspend", aAuth(http.HandlerFunc(h.AdminStaff.SuspendOrder)))
	mux.Handle("POST /api/v1/admin/orders/{id}/unsuspend", aAuth(http.HandlerFunc(h.AdminStaff.UnsuspendOrder)))
	mux.Handle("POST /api/v1/admin/orders/{id}/activate", aAuth(http.HandlerFunc(h.AdminStaff.ActivateOrder)))
	mux.Handle("GET /api/v1/admin/support/tickets", aAuth(http.HandlerFunc(h.AdminStaff.ListTickets)))
	mux.Handle("GET /api/v1/admin/support/tickets/{id}", aAuth(http.HandlerFunc(h.AdminStaff.GetTicket)))
	mux.Handle("POST /api/v1/admin/support/tickets/{id}/reply", aAuth(http.HandlerFunc(h.AdminStaff.ReplyTicket)))
}
