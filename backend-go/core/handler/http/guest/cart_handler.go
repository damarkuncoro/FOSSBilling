package guest

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/cart"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type CartHandler struct {
	cartService *cart.CartService
}

func NewCartHandler(cartService *cart.CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

func (h *CartHandler) Calculate(w http.ResponseWriter, r *http.Request) {
	var cartReq cart.Cart
	if err := json.NewDecoder(r.Body).Decode(&cartReq); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", nil)
		return
	}

	// If user is logged in, use their ID for tax/promo calculations
	if clientID := middleware.GetClientID(r.Context()); clientID > 0 {
		cartReq.ClientID = clientID
	} else {
		cartReq.ClientID = 0
	}

	if err := h.cartService.CalculateTotals(r.Context(), &cartReq); err != nil {
		response.Error(w, http.StatusBadRequest, "CALCULATION_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, cartReq, nil)
}

func (h *CartHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	var cartReq cart.Cart
	if err := json.NewDecoder(r.Body).Decode(&cartReq); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", nil)
		return
	}

	// Use authenticated client ID if available
	clientID := middleware.GetClientID(r.Context())
	if clientID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Please login to complete checkout", nil)
		return
	}
	cartReq.ClientID = clientID

	res, err := h.cartService.Checkout(r.Context(), &cartReq)
	if err != nil {
		if errors.Is(err, cart.ErrEmptyCart) {
			response.Error(w, http.StatusBadRequest, "EMPTY_CART", "Cart is empty", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "CHECKOUT_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusCreated, res, nil)
}
