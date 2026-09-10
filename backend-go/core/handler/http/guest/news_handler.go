package guest

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/news"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/request"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type NewsHandler struct{ svc *news.NewsService }

func NewNewsHandler(s *news.NewsService) *NewsHandler { return &NewsHandler{s} }

func (h *NewsHandler) List(w http.ResponseWriter, r *http.Request) {
	l, o := request.GetLimitOffset(r)
	ps, tot, err := h.svc.ListPublished(r.Context(), l, o)
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, ps, &response.Meta{Total: tot, Limit: l, Offset: o})
}

func (h *NewsHandler) Get(w http.ResponseWriter, r *http.Request) {
	p, err := h.svc.GetBySlug(r.Context(), r.PathValue("slug"))
	if err != nil { response.Error(w, 404, "NOT_FOUND", "Not found", nil); return }
	response.JSON(w, 200, p, nil)
}
