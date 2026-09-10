package guest

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/catalog"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/request"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type ProductHandler struct{ svc *catalog.ProductService }

func NewProductHandler(s *catalog.ProductService) *ProductHandler { return &ProductHandler{s} }

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	l, o := request.GetLimitOffset(r)
	ps, tot, err := h.svc.ListProducts(r.Context(), l, o)
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, ps, &response.Meta{Total: tot, Limit: l, Offset: o})
}
