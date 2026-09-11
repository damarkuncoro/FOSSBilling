package main

import (
	"net/http"
)

func registerClientRoutes(mux *http.ServeMux, h *AppHandlers, cAuth func(http.Handler) http.Handler) {
	clientAccountRoutes(mux, h, cAuth)
	clientBillingRoutes(mux, h, cAuth)
	clientSupportRoutes(mux, h, cAuth)
	clientServiceRoutes(mux, h, cAuth)
}

func clientAccountRoutes(mux *http.ServeMux, h *AppHandlers, cAuth func(http.Handler) http.Handler) {
	mux.Handle("GET /api/v1/client/profile", cAuth(http.HandlerFunc(h.ClientProfile.GetProfile)))
	mux.Handle("PUT /api/v1/client/profile", cAuth(http.HandlerFunc(h.ClientProfile.UpdateProfile)))
	mux.Handle("POST /api/v1/client/profile/change-password", cAuth(http.HandlerFunc(h.ClientProfile.ChangePassword)))
	mux.Handle("POST /api/v1/client/profile/2fa/setup", cAuth(http.HandlerFunc(h.ClientProfile.SetupTwoFactor)))
	mux.Handle("POST /api/v1/client/profile/2fa/enable", cAuth(http.HandlerFunc(h.ClientProfile.EnableTwoFactor)))
	mux.Handle("POST /api/v1/client/profile/2fa/disable", cAuth(http.HandlerFunc(h.ClientProfile.DisableTwoFactor)))
	mux.Handle("GET /api/v1/client/activity", cAuth(http.HandlerFunc(h.ClientActivity.ListMyLogs)))
	mux.Handle("GET /api/v1/client/notifications", cAuth(http.HandlerFunc(h.ClientNotification.List)))
	mux.Handle("PUT /api/v1/client/notifications/{id}/read", cAuth(http.HandlerFunc(h.ClientNotification.MarkAsRead)))
	mux.Handle("GET /api/v1/client/api-keys", cAuth(http.HandlerFunc(h.ClientAPIKey.List)))
	mux.Handle("POST /api/v1/client/api-keys", cAuth(http.HandlerFunc(h.ClientAPIKey.Generate)))
	mux.Handle("DELETE /api/v1/client/api-keys/{id}", cAuth(http.HandlerFunc(h.ClientAPIKey.Revoke)))
	mux.Handle("GET /api/v1/client/affiliate", cAuth(http.HandlerFunc(h.ClientAffiliate.GetMyAffiliate)))
	mux.Handle("POST /api/v1/client/affiliate/join", cAuth(http.HandlerFunc(h.ClientAffiliate.Join)))
}

func clientBillingRoutes(mux *http.ServeMux, h *AppHandlers, cAuth func(http.Handler) http.Handler) {
	mux.Handle("GET /api/v1/client/invoices", cAuth(http.HandlerFunc(h.ClientInvoice.ListInvoices)))
	mux.Handle("GET /api/v1/client/invoices/{id}", cAuth(http.HandlerFunc(h.ClientInvoice.GetInvoice)))
	mux.Handle("GET /api/v1/client/invoices/{id}/pdf", cAuth(http.HandlerFunc(h.ClientInvoice.DownloadPDF)))
	mux.Handle("POST /api/v1/client/invoices/{id}/pay-balance", cAuth(http.HandlerFunc(h.ClientInvoice.PayWithBalance)))
	mux.Handle("POST /api/v1/client/invoices/{id}/pay-gateway", cAuth(http.HandlerFunc(h.ClientInvoice.PayWithGateway)))
	mux.Handle("POST /api/v1/client/funds/deposit", cAuth(http.HandlerFunc(h.ClientDeposit.DepositFunds)))
}

func clientSupportRoutes(mux *http.ServeMux, h *AppHandlers, cAuth func(http.Handler) http.Handler) {
	mux.Handle("POST /api/v1/client/support/tickets", cAuth(http.HandlerFunc(h.ClientSupport.OpenTicket)))
	mux.Handle("GET /api/v1/client/support/tickets", cAuth(http.HandlerFunc(h.ClientSupport.ListTickets)))
	mux.Handle("GET /api/v1/client/support/tickets/{id}", cAuth(http.HandlerFunc(h.ClientSupport.GetTicket)))
	mux.Handle("POST /api/v1/client/support/tickets/{id}/reply", cAuth(http.HandlerFunc(h.ClientSupport.ReplyTicket)))
	mux.Handle("POST /api/v1/client/support/tickets/{id}/close", cAuth(http.HandlerFunc(h.ClientSupport.CloseTicket)))
}

func clientServiceRoutes(mux *http.ServeMux, h *AppHandlers, cAuth func(http.Handler) http.Handler) {
	mux.Handle("GET /api/v1/client/orders", cAuth(http.HandlerFunc(h.ClientOrder.ListOrders)))
	mux.Handle("GET /api/v1/client/orders/{id}", cAuth(http.HandlerFunc(h.ClientOrder.GetOrder)))
	mux.Handle("POST /api/v1/client/orders/{id}/sync", cAuth(http.HandlerFunc(h.ClientOrder.SyncStatus)))
	mux.Handle("POST /api/v1/client/orders/{id}/change-password", cAuth(http.HandlerFunc(h.ClientOrder.ChangePassword)))
	mux.Handle("GET /api/v1/client/domains", cAuth(http.HandlerFunc(h.ClientDomain.ListDomains)))
	mux.Handle("GET /api/v1/client/domains/{id}/dns", cAuth(http.HandlerFunc(h.ClientDomain.ListDNSRecords)))
	mux.Handle("POST /api/v1/client/domains/{id}/dns", cAuth(http.HandlerFunc(h.ClientDomain.AddDNSRecord)))
	mux.Handle("DELETE /api/v1/client/domains/{id}/dns/{record_id}", cAuth(http.HandlerFunc(h.ClientDomain.DeleteDNSRecord)))
	mux.Handle("PUT /api/v1/client/domains/{id}/nameservers", cAuth(http.HandlerFunc(h.ClientDomain.UpdateNameservers)))
	mux.Handle("POST /api/v1/client/domains/{id}/toggle-autorenew", cAuth(http.HandlerFunc(h.ClientDomain.ToggleAutoRenew)))
	mux.Handle("GET /api/v1/client/licenses", cAuth(http.HandlerFunc(h.ClientLicense.ListLicenses)))
	mux.Handle("POST /api/v1/client/licenses/{id}/reset", cAuth(http.HandlerFunc(h.ClientLicense.ResetLicenseLock)))
	mux.Handle("GET /api/v1/client/downloads", cAuth(http.HandlerFunc(h.ClientDownload.ListDownloads)))
	mux.Handle("GET /api/v1/client/downloads/{id}/link", cAuth(http.HandlerFunc(h.ClientDownload.GenerateLink)))
	mux.Handle("GET /api/v1/client/downloads/{id}/file", http.HandlerFunc(h.ClientDownload.StreamFile))
}
