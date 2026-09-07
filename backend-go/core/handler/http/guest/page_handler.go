package guest

import (
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/page"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type PageHandler struct {
	pageService *page.PageService
}

func NewPageHandler(pageService *page.PageService) *PageHandler {
	return &PageHandler{pageService: pageService}
}

func (h *PageHandler) GetPage(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	p, err := h.pageService.GetPage(r.Context(), slug)
	if err != nil {
		response.Error(w, http.StatusNotFound, "NOT_FOUND", "page not found", nil)
		return
	}
	if !p.Published {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "page is not published", nil)
		return
	}
	response.JSON(w, http.StatusOK, p, nil)
}
