package guest

import (
	"net/http"

	seoUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/seo"
)

type SEOHandler struct {
	svc *seoUsecase.SEOService
}

func NewSEOHandler(svc *seoUsecase.SEOService) *SEOHandler {
	return &SEOHandler{svc: svc}
}

func (h *SEOHandler) GetSitemap(w http.ResponseWriter, r *http.Request) {
	baseURL := r.Header.Get("X-Forwarded-Host")
	if baseURL == "" {
		baseURL = r.Host
	}
	scheme := "https"
	if r.TLS == nil && !r.URL.IsAbs() {
		scheme = "http"
	}
	fullURL := scheme + "://" + baseURL

	xmlStr, err := h.svc.GenerateSitemapXML(r.Context(), fullURL)
	if err != nil {
		http.Error(w, "Failed to generate sitemap", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(xmlStr))
}
