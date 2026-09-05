package guest

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fossbilling/backend-go/internal/usecase/cart"
	"github.com/fossbilling/backend-go/pkg/response"
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
