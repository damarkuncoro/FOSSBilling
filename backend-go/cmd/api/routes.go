package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/config"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/http/admin"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/http/client"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/http/guest"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/i18n"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

// AppHandlers bundles all presentation layer HTTP handlers
type AppHandlers struct {
	GuestAuth          *guest.AuthHandler
	GuestCart          *guest.CartHandler
	GuestProduct       *guest.ProductHandler
	GuestWebhook       *guest.WebhookHandler
	GuestCurrency      *guest.CurrencyHandler
	GuestNews          *guest.NewsHandler
	GuestPage          *guest.PageHandler
	GuestCompany       *guest.CompanyHandler
	GuestDomain        *guest.DomainHandler
	GuestKB            *guest.KBHandler
	GuestFormbuilder   *guest.FormbuilderHandler
	GuestRedirect      *guest.RedirectHandler
	GuestCookieConsent *guest.CookieConsentHandler
	GuestTheme         *guest.ThemeHandler
	GuestSEO           *guest.SEOHandler
	GuestWidget        *guest.WidgetHandler
	GuestSystem        *guest.SystemHandler
	ClientProfile      *client.ProfileHandler
	ClientOrder        *client.OrderHandler
	ClientDomain       *client.DomainHandler
	ClientInvoice      *client.InvoiceHandler
	ClientDeposit      *client.DepositHandler
	ClientSupport      *client.SupportHandler
	ClientActivity     *client.ActivityHandler
	ClientNotification *client.NotificationHandler
	ClientDownload     *client.DownloadHandler
	ClientLicense      *client.LicenseHandler
	ClientAPIKey       *client.APIKeyHandler
	ClientAffiliate    *client.AffiliateHandler
	AdminAuth          *admin.StaffAuthHandler
	AdminStaff         *admin.StaffManagementHandler
	AdminClient        *admin.ClientManagementHandler
	AdminInvoice       *admin.InvoiceManagementHandler
	AdminStats         *admin.StatsHandler
	AdminCurrency      *admin.CurrencyHandler
	AdminNews          *admin.NewsHandler
	AdminKB            *admin.KBHandler
	AdminMassMail      *admin.MassMailHandler
	AdminCompany       *admin.CompanyHandler
	AdminCatalog       *admin.CatalogHandler
	AdminBilling       *admin.BillingModuleHandler
	AdminSystem        *admin.SystemModuleHandler
	AdminActivity      *admin.ActivityHandler
	AdminAntispam      *admin.AntispamHandler
	AdminFormbuilder   *admin.FormbuilderHandler
	AdminExtension     *admin.ExtensionHandler
	AdminRedirect      *admin.RedirectHandler
	AdminCookieConsent *admin.CookieConsentHandler
	AdminNotification  *admin.AdminNotificationHandler
	AdminTheme         *admin.ThemeHandler
	AdminSEO           *admin.SEOHandler
	AdminWidget        *admin.WidgetHandler
	AdminEmailTemplate *admin.EmailTemplateHandler
	AdminAffiliate     *admin.AffiliateHandler
}

// setupRoutes initializes system routes and dispatches to role-scoped routers
func setupRoutes(cfg *config.Config, h *AppHandlers, rateLimiter, authRateLimiter *middleware.RateLimiter) http.Handler {
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
	mux.Handle("GET /metrics", promhttp.Handler())
	mux.HandleFunc("GET /openapi.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		data, err := os.ReadFile("docs/openapi.json")
		if err != nil {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "OpenAPI spec not found", nil)
			return
		}
		_, _ = w.Write(data)
	})
	mux.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		html := fmt.Sprintf(`<!doctype html>
<html>
  <head>
    <title>%s API Documentation</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <script
      id="api-reference"
      data-url="/openapi.json"
      src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`, cfg.CompanyName)
		_, _ = w.Write([]byte(html))
	})
	mux.HandleFunc("GET /sitemap.xml", h.GuestSEO.GetSitemap)

	// 2. Auth Middlewares
	clientAuth := middleware.RequireAuth(cfg.JWTSecret, "client", "admin", "superadmin")
	adminAuth := middleware.RequireAuth(cfg.JWTSecret, "admin", "superadmin", "support", "billing")

	// 3. Register Role-Scoped Routes
	registerGuestRoutes(mux, h, rateLimiter, authRateLimiter)
	registerClientRoutes(mux, h, clientAuth)
	registerAdminRoutes(mux, h, adminAuth, rateLimiter, authRateLimiter)

	// 4. Dev/Simulation Routes (Only in non-prod)
	if cfg.AppEnv != "production" {
		registerDevRoutes(mux, h)
	}

	return middleware.Recovery(middleware.SecurityHeaders(middleware.AdminGeofence(cfg.AllowedCountries)(middleware.Metrics()(middleware.Logger(middleware.MaintenanceMode(h.AdminSystem.GetSystemService())(middleware.CORS(cfg.AllowedOrigins)(i18n.LocaleMiddleware(mux))))))))
}
