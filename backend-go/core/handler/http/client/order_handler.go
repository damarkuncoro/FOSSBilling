package client

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	orderUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/order"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type OrderHandler struct {
	orderService *orderUsecase.OrderService
}

func NewOrderHandler(orderService *orderUsecase.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
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

	orders, total, err := h.orderService.ListByClientID(r.Context(), clientID, limit, offset)
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

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid order ID", nil)
		return
	}

	order, err := h.orderService.GetByIDForClient(r.Context(), clientID, id)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Order not found", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get order", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, order, nil)
}

func (h *OrderHandler) SyncStatus(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	if clientID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid order ID", nil)
		return
	}

	status, err := h.orderService.SyncServiceStatus(r.Context(), clientID, id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "SYNC_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, status, nil)
}

func (h *OrderHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	if clientID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid order ID", nil)
		return
	}

	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_BODY", "Failed to parse JSON body", nil)
		return
	}

	if err := h.orderService.ChangeServicePassword(r.Context(), clientID, id, req.Password); err != nil {
		response.Error(w, http.StatusInternalServerError, "PASSWORD_CHANGE_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Service password changed successfully"}, nil)
}
