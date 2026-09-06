package main

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
)

// registerGuestRoutes configures all public guest API endpoints
func registerGuestRoutes(mux *http.ServeMux, h *AppHandlers, rateLimiter *middleware.RateLimiter) {
	mux.Handle("POST /api/v1/guest/auth/register", rateLimiter.RateLimit(http.HandlerFunc(h.GuestAuth.Register)))
	mux.Handle("POST /api/v1/guest/auth/login", rateLimiter.RateLimit(http.HandlerFunc(h.GuestAuth.Login)))
	mux.Handle("POST /api/v1/guest/cart/calculate", rateLimiter.RateLimit(http.HandlerFunc(h.GuestCart.Calculate)))
	mux.Handle("POST /api/v1/guest/cart/checkout", rateLimiter.RateLimit(http.HandlerFunc(h.GuestCart.Checkout)))
	mux.Handle("POST /api/v1/guest/gateways/{gateway}/webhook", http.HandlerFunc(h.GuestWebhook.HandleGatewayWebhook))
	mux.Handle("GET /api/v1/guest/currencies", rateLimiter.RateLimit(http.HandlerFunc(h.GuestCurrency.List)))
	mux.Handle("GET /api/v1/guest/news", rateLimiter.RateLimit(http.HandlerFunc(h.GuestNews.List)))
	mux.Handle("GET /api/v1/guest/news/{slug}", rateLimiter.RateLimit(http.HandlerFunc(h.GuestNews.Get)))
	mux.Handle("GET /api/v1/guest/company", rateLimiter.RateLimit(http.HandlerFunc(h.GuestCompany.GetCompany)))
}
