package guest

import (
	"net/http"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/page"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
	"github.com/go-chi/chi/v5"
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
		response.Error(w, http.StatusNotFound, "page not found")
		return
	}
	if !p.Published {
		response.Error(w, http.StatusForbidden, "page is not published")
		return
	}
	response.JSON(w, http.StatusOK, p, nil)
}
