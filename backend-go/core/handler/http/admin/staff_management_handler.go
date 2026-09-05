package admin

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/order"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/support"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type StaffManagementHandler struct {
	staffService   *staff.StaffService
	clientRepo     domain.ClientRepository
	orderRepo      domain.OrderRepository
	orderService   *order.OrderService
	supportService *support.SupportService
}

func NewStaffManagementHandler(
	staffService *staff.StaffService,
	clientRepo domain.ClientRepository,
	orderRepo domain.OrderRepository,
	orderService *order.OrderService,
	supportService *support.SupportService,
) *StaffManagementHandler {
	return &StaffManagementHandler{
		staffService:   staffService,
		clientRepo:     clientRepo,
		orderRepo:      orderRepo,
		orderService:   orderService,
		supportService: supportService,
	}
}

func (h *StaffManagementHandler) ListClients(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "clients", "read")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: clients", nil)
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

	clients, total, err := h.clientRepo.List(r.Context(), limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve clients", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, clients, &response.Meta{
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

func (h *StaffManagementHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "orders", "read")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: orders", nil)
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

	orders, total, err := h.orderRepo.List(r.Context(), limit, offset)
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

type suspendRequest struct {
	Reason string `json:"reason"`
}

func (h *StaffManagementHandler) SuspendOrder(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "orders", "write")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: orders", nil)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid path", nil)
		return
	}
	orderID, err := strconv.ParseInt(parts[len(parts)-2], 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid order ID", nil)
		return
	}

	var req suspendRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	if req.Reason == "" {
		req.Reason = "Suspended by admin"
	}

	res, err := h.orderService.Suspend(r.Context(), orderID, req.Reason)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "SUSPEND_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, res, nil)
}

func (h *StaffManagementHandler) UnsuspendOrder(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid path", nil)
		return
	}
	orderID, err := strconv.ParseInt(parts[len(parts)-2], 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid order ID", nil)
		return
	}

	res, err := h.orderService.Unsuspend(r.Context(), orderID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "UNSUSPEND_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, res, nil)
}

func (h *StaffManagementHandler) ActivateOrder(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid path", nil)
		return
	}
	orderID, err := strconv.ParseInt(parts[len(parts)-2], 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid order ID", nil)
		return
	}

	res, err := h.orderService.Activate(r.Context(), orderID, time.Now().UTC())
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ACTIVATE_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, res, nil)
}
