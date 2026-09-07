package main

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/i18n"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

// registerGuestRoutes configures all public guest API endpoints
func registerGuestRoutes(mux *http.ServeMux, h *AppHandlers, rateLimiter *middleware.RateLimiter) {
	mux.Handle("POST /api/v1/guest/auth/register", rateLimiter.RateLimit(http.HandlerFunc(h.GuestAuth.Register)))
	mux.Handle("POST /api/v1/guest/auth/login", rateLimiter.RateLimit(http.HandlerFunc(h.GuestAuth.Login)))
	mux.Handle("POST /api/v1/guest/cart/calculate", rateLimiter.RateLimit(http.HandlerFunc(h.GuestCart.Calculate)))
	mux.Handle("POST /api/v1/guest/cart/checkout", http.HandlerFunc(h.GuestCart.Checkout))
	mux.Handle("POST /api/v1/guest/webhook/custom", http.HandlerFunc(h.GuestWebhook.HandleGatewayWebhook))
	mux.Handle("POST /api/v1/guest/gateways/{gateway}/webhook", http.HandlerFunc(h.GuestWebhook.HandleGatewayWebhook))
	mux.Handle("GET /api/v1/guest/currencies", rateLimiter.RateLimit(http.HandlerFunc(h.GuestCurrency.List)))
	mux.Handle("GET /api/v1/guest/news", rateLimiter.RateLimit(http.HandlerFunc(h.GuestNews.List)))
	mux.Handle("GET /api/v1/guest/news/{slug}", rateLimiter.RateLimit(http.HandlerFunc(h.GuestNews.Get)))
	mux.Handle("GET /api/v1/guest/pages/{slug}", rateLimiter.RateLimit(http.HandlerFunc(h.GuestPage.GetPage)))
	mux.Handle("GET /api/v1/guest/company", rateLimiter.RateLimit(http.HandlerFunc(h.GuestCompany.GetCompany)))
	mux.Handle("GET /api/v1/guest/domains/check", rateLimiter.RateLimit(http.HandlerFunc(h.GuestDomain.CheckAvailability)))
	mux.Handle("GET /api/v1/guest/forms/{id}", rateLimiter.RateLimit(http.HandlerFunc(h.GuestFormbuilder.GetForm)))
	mux.Handle("GET /api/v1/guest/redirects/lookup", rateLimiter.RateLimit(http.HandlerFunc(h.GuestRedirect.Lookup)))
	mux.Handle("GET /api/v1/guest/cookie-consent", rateLimiter.RateLimit(http.HandlerFunc(h.GuestCookieConsent.GetConfig)))
	mux.Handle("GET /api/v1/guest/theme/active", rateLimiter.RateLimit(http.HandlerFunc(h.GuestTheme.GetActiveTheme)))
	mux.Handle("GET /api/v1/guest/widgets", rateLimiter.RateLimit(http.HandlerFunc(h.GuestWidget.GetSlotWidgets)))
	mux.Handle("GET /api/v1/guest/system/geoip", rateLimiter.RateLimit(http.HandlerFunc(h.AdminSystem.ResolveGeoIP)))
	mux.Handle("GET /api/v1/guest/locales", rateLimiter.RateLimit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, i18n.SupportedLocales, nil)
	})))
}
