package main

import (
	"net/http"
)

// registerClientRoutes configures all client portal protected endpoints
func registerClientRoutes(mux *http.ServeMux, h *AppHandlers, cAuth func(http.Handler) http.Handler) {
	mux.Handle("GET /api/v1/client/profile", cAuth(http.HandlerFunc(h.ClientProfile.GetProfile)))
	mux.Handle("PUT /api/v1/client/profile", cAuth(http.HandlerFunc(h.ClientProfile.UpdateProfile)))
	mux.Handle("POST /api/v1/client/profile/change-password", cAuth(http.HandlerFunc(h.ClientProfile.ChangePassword)))
	mux.Handle("GET /api/v1/client/orders", cAuth(http.HandlerFunc(h.ClientOrder.ListOrders)))
	mux.Handle("GET /api/v1/client/orders/{id}", cAuth(http.HandlerFunc(h.ClientOrder.GetOrder)))
	mux.Handle("GET /api/v1/client/domains", cAuth(http.HandlerFunc(h.ClientDomain.ListDomains)))
	mux.Handle("PUT /api/v1/client/domains/{id}/nameservers", cAuth(http.HandlerFunc(h.ClientDomain.UpdateNameservers)))
	mux.Handle("POST /api/v1/client/domains/{id}/toggle-autorenew", cAuth(http.HandlerFunc(h.ClientDomain.ToggleAutoRenew)))
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
	mux.Handle("GET /api/v1/client/licenses", cAuth(http.HandlerFunc(h.ClientLicense.ListLicenses)))
	mux.Handle("POST /api/v1/client/licenses/{id}/reset", cAuth(http.HandlerFunc(h.ClientLicense.ResetLicenseLock)))
	mux.Handle("GET /api/v1/client/downloads", cAuth(http.HandlerFunc(h.ClientDownload.ListDownloads)))
	mux.Handle("GET /api/v1/client/downloads/{id}/link", cAuth(http.HandlerFunc(h.ClientDownload.GenerateLink)))
	mux.Handle("GET /api/v1/client/downloads/{id}/file", http.HandlerFunc(h.ClientDownload.StreamFile))
	mux.Handle("GET /api/v1/client/api-keys", cAuth(http.HandlerFunc(h.ClientAPIKey.List)))
	mux.Handle("POST /api/v1/client/api-keys", cAuth(http.HandlerFunc(h.ClientAPIKey.Generate)))
	mux.Handle("DELETE /api/v1/client/api-keys/{id}", cAuth(http.HandlerFunc(h.ClientAPIKey.Revoke)))
}
