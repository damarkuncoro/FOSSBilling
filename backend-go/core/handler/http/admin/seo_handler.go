package admin

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	seoUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/seo"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type SEOHandler struct {
	staffService *staff.StaffService
	svc          *seoUsecase.SEOService
}

func NewSEOHandler(staffService *staff.StaffService, svc *seoUsecase.SEOService) *SEOHandler {
	return &SEOHandler{staffService: staffService, svc: svc}
}

func (h *SEOHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "system", "read")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: system", nil)
		return
	}

	baseURL := r.Header.Get("X-Forwarded-Host")
	if baseURL == "" {
		baseURL = r.Host
	}
	scheme := "https"
	if r.TLS == nil && !r.URL.IsAbs() {
		scheme = "http"
	}
	fullURL := scheme + "://" + baseURL

	info := h.svc.GetInfo(fullURL)
	response.JSON(w, http.StatusOK, info, nil)
}

func (h *SEOHandler) PingSearchEngines(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "system", "write")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: system", nil)
		return
	}

	baseURL := r.Header.Get("X-Forwarded-Host")
	if baseURL == "" {
		baseURL = r.Host
	}
	scheme := "https"
	if r.TLS == nil && !r.URL.IsAbs() {
		scheme = "http"
	}
	fullURL := scheme + "://" + baseURL

	results, err := h.svc.PingSearchEngines(r.Context(), fullURL)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to ping search engines", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Search engine ping dispatched successfully",
		"results": results,
	}, nil)
}
