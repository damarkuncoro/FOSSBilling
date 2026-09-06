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

// AppHandlers bundles all presentation layer HTTP handlers
type AppHandlers struct {
	GuestAuth      *guest.AuthHandler
	GuestCart      *guest.CartHandler
	GuestWebhook   *guest.WebhookHandler
	GuestCurrency  *guest.CurrencyHandler
	GuestNews      *guest.NewsHandler
	GuestCompany   *guest.CompanyHandler
	GuestDomain    *guest.DomainHandler
	ClientProfile  *client.ProfileHandler
	ClientOrder    *client.OrderHandler
	ClientDomain   *client.DomainHandler
	ClientInvoice  *client.InvoiceHandler
	ClientDeposit  *client.DepositHandler
	ClientSupport  *client.SupportHandler
	ClientDownload *client.DownloadHandler
	ClientLicense  *client.LicenseHandler
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

// setupRoutes initializes system routes and dispatches to role-scoped routers
func setupRoutes(cfg *config.Config, h *AppHandlers, rateLimiter *middleware.RateLimiter) http.Handler {
	mux := http.NewServeMux()

	// 1. System Base (Health check & OpenAPI specs)
	startTime := os.Getenv("BOOT_TIME")
	if startTime == "" {
		startTime = "active"
	}

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]interface{}{
			"status":      "healthy",
			"version":     "2.0.0-golang",
			"environment": cfg.AppEnv,
			"subsystems": map[string]string{
				"api":          "healthy",
				"scheduler":    "ready",
				"security":     "active",
				"provisioning": "ready",
				"payment":      "ready",
			},
		}, nil)
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

	// 2. Auth Middlewares
	clientAuth := middleware.RequireAuth(cfg.JWTSecret, "client", "admin", "superadmin")
	adminAuth := middleware.RequireAuth(cfg.JWTSecret, "admin", "superadmin", "support", "billing")

	// 3. Register Role-Scoped Routes
	registerGuestRoutes(mux, h, rateLimiter)
	registerClientRoutes(mux, h, clientAuth)
	registerAdminRoutes(mux, h, adminAuth, rateLimiter)

	return middleware.Recovery(middleware.SecurityHeaders(middleware.Logger(middleware.CORS(mux))))
}
