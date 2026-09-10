package main

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/i18n"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

func registerGuestRoutes(mux *http.ServeMux, h *AppHandlers, rl, arl *middleware.RateLimiter) {
	guestAuthRoutes(mux, h, arl)
	guestCatalogRoutes(mux, h, rl)
	guestContentRoutes(mux, h, rl)
}

func guestAuthRoutes(mux *http.ServeMux, h *AppHandlers, arl *middleware.RateLimiter) {
	mux.Handle("POST /api/v1/guest/auth/register", arl.RateLimit(http.HandlerFunc(h.GuestAuth.Register)))
	mux.Handle("POST /api/v1/guest/auth/login", arl.RateLimit(http.HandlerFunc(h.GuestAuth.Login)))
	mux.Handle("POST /api/v1/guest/auth/verify-2fa", arl.RateLimit(http.HandlerFunc(h.GuestAuth.VerifyTwoFactor)))
}

func guestCatalogRoutes(mux *http.ServeMux, h *AppHandlers, rl *middleware.RateLimiter) {
	mux.Handle("GET /api/v1/guest/products", rl.RateLimit(http.HandlerFunc(h.GuestProduct.List)))
	mux.Handle("POST /api/v1/guest/cart/calculate", rl.RateLimit(http.HandlerFunc(h.GuestCart.Calculate)))
	mux.Handle("POST /api/v1/guest/cart/checkout", http.HandlerFunc(h.GuestCart.Checkout))
	mux.Handle("GET /api/v1/guest/currencies", rl.RateLimit(http.HandlerFunc(h.GuestCurrency.List)))
	mux.Handle("GET /api/v1/guest/domains/check", rl.RateLimit(http.HandlerFunc(h.GuestDomain.CheckAvailability)))
	mux.Handle("GET /api/v1/guest/forms/{id}", rl.RateLimit(http.HandlerFunc(h.GuestFormbuilder.GetForm)))
}

func guestContentRoutes(mux *http.ServeMux, h *AppHandlers, rl *middleware.RateLimiter) {
	mux.Handle("GET /api/v1/guest/news", rl.RateLimit(http.HandlerFunc(h.GuestNews.List)))
	mux.Handle("GET /api/v1/guest/news/{slug}", rl.RateLimit(http.HandlerFunc(h.GuestNews.Get)))
	mux.Handle("GET /api/v1/guest/kb/categories", rl.RateLimit(http.HandlerFunc(h.GuestKB.ListCategories)))
	mux.Handle("GET /api/v1/guest/kb/articles", rl.RateLimit(http.HandlerFunc(h.GuestKB.ListArticles)))
	mux.Handle("GET /api/v1/guest/kb/articles/detail", rl.RateLimit(http.HandlerFunc(h.GuestKB.GetArticle)))
	mux.Handle("GET /api/v1/guest/pages/{slug}", rl.RateLimit(http.HandlerFunc(h.GuestPage.GetPage)))
	mux.Handle("GET /api/v1/guest/company", rl.RateLimit(http.HandlerFunc(h.GuestCompany.GetCompany)))
	mux.Handle("GET /api/v1/guest/redirects/lookup", rl.RateLimit(http.HandlerFunc(h.GuestRedirect.Lookup)))
	mux.Handle("GET /api/v1/guest/cookie-consent", rl.RateLimit(http.HandlerFunc(h.GuestCookieConsent.GetConfig)))
	mux.Handle("GET /api/v1/guest/theme/active", rl.RateLimit(http.HandlerFunc(h.GuestTheme.GetActiveTheme)))
	mux.Handle("GET /api/v1/guest/widgets", rl.RateLimit(http.HandlerFunc(h.GuestWidget.GetSlotWidgets)))
	mux.Handle("GET /api/v1/guest/system/health", rl.RateLimit(http.HandlerFunc(h.GuestSystem.HealthCheck)))
	mux.Handle("GET /api/v1/guest/system/geoip", rl.RateLimit(http.HandlerFunc(h.AdminSystem.ResolveGeoIP)))
	mux.Handle("POST /api/v1/guest/webhook/custom", http.HandlerFunc(h.GuestWebhook.HandleGatewayWebhook))
	mux.Handle("POST /api/v1/guest/gateways/{gateway}/webhook", http.HandlerFunc(h.GuestWebhook.HandleGatewayWebhook))
	mux.Handle("GET /api/v1/guest/locales", rl.RateLimit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, i18n.SupportedLocales, nil)
	})))
}
