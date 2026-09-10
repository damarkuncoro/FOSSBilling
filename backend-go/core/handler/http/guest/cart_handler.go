package guest

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/cart"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/request"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type CartHandler struct { svc *cart.CartService }

func NewCartHandler(s *cart.CartService) *CartHandler { return &CartHandler{s} }

func (h *CartHandler) Calculate(w http.ResponseWriter, r *http.Request) {
	var c cart.Cart; if request.Decode(r, &c) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	if id := middleware.GetClientID(r.Context()); id > 0 { c.ClientID = id }
	if err := h.svc.CalculateTotals(r.Context(), &c); err != nil { response.Error(w, 400, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, map[string]any{"items": c.Items, "promo_code": c.PromoCode, "subtotal": c.Subtotal, "discount": c.Discount, "tax": c.Tax, "total": c.Total}, nil)
}

func (h *CartHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	var c cart.Cart; if request.Decode(r, &c) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	if id := middleware.GetClientID(r.Context()); id > 0 { c.ClientID = id }
	if c.ClientID == 0 { response.Error(w, 401, "UNAUTHORIZED", "Login required", nil); return }
	res, err := h.svc.Checkout(r.Context(), &c, r.RemoteAddr)
	if err != nil { response.Error(w, 400, "ERR", err.Error(), nil); return }
	response.JSON(w, 201, res, nil)
}
