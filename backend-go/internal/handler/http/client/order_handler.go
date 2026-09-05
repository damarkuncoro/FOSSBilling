package client

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/fossbilling/backend-go/internal/domain"
	"github.com/fossbilling/backend-go/internal/handler/middleware"
	appErrors "github.com/fossbilling/backend-go/pkg/errors"
	"github.com/fossbilling/backend-go/pkg/response"
)

type OrderHandler struct {
	orderRepo domain.OrderRepository
}

func NewOrderHandler(orderRepo domain.OrderRepository) *OrderHandler {
	return &OrderHandler{orderRepo: orderRepo}
}

func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	if clientID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}

	limit := 20
	offset := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}

	orders, total, err := h.orderRepo.ListByClientID(r.Context(), clientID, limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve orders", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, orders, &response.Meta{
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	if clientID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	orderIDStr := parts[len(parts)-1]
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid order ID", nil)
		return
	}

	order, err := h.orderRepo.GetByID(r.Context(), orderID)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Order not found", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get order", err.Error())
		return
	}

	if order.ClientID != clientID {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "You do not have access to this order", nil)
		return
	}

	response.JSON(w, http.StatusOK, order, nil)
}
