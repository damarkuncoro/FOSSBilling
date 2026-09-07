package admin

import (
	"net/http"

	seoUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/seo"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type SEOHandler struct {
	svc *seoUsecase.SEOService
}

func NewSEOHandler(svc *seoUsecase.SEOService) *SEOHandler {
	return &SEOHandler{svc: svc}
}

func (h *SEOHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
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
