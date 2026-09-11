package admin

import (
	"net/http"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/order"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/support"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/request"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type StaffManagementHandler struct {
	staffService *staff.StaffService; clientRepo domain.ClientRepository; orderRepo domain.OrderRepository; orderService *order.OrderService; supportService *support.SupportService
}

func NewStaffManagementHandler(s *staff.StaffService, cr domain.ClientRepository, or domain.OrderRepository, os *order.OrderService, ss *support.SupportService) *StaffManagementHandler {
	return &StaffManagementHandler{s, cr, or, os, ss}
}

func (h *StaffManagementHandler) ListClients(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "clients", "read") { return }
	l, o := request.GetLimitOffset(r)
	cs, tot, err := h.clientRepo.List(r.Context(), l, o)
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, cs, &response.Meta{Total: tot, Limit: l, Offset: o})
}

func (h *StaffManagementHandler) ListStaff(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "staff", "read") { return }
	l, o := request.GetLimitOffset(r)
	sl, tot, err := h.staffService.ListStaff(r.Context(), l, o)
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, sl, &response.Meta{Total: tot, Limit: l, Offset: o})
}

func (h *StaffManagementHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "orders", "read") { return }
	l, o := request.GetLimitOffset(r)
	os, tot, err := h.orderRepo.List(r.Context(), l, o)
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, os, &response.Meta{Total: tot, Limit: l, Offset: o})
}

func (h *StaffManagementHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "orders", "read") { return }
	o, err := h.orderRepo.GetByID(r.Context(), request.GetID(r))
	if err != nil { response.Error(w, 404, "NOT_FOUND", "Order not found", nil); return }
	response.JSON(w, 200, o, nil)
}

func (h *StaffManagementHandler) SuspendOrder(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "orders", "write") { return }
	var req struct{ Reason string }; _ = request.Decode(r, &req)
	if req.Reason == "" { req.Reason = "Suspended by admin" }
	res, err := h.orderService.Suspend(r.Context(), request.GetID(r), req.Reason)
	if err != nil { response.Error(w, 400, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, res, nil)
}

func (h *StaffManagementHandler) UnsuspendOrder(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "orders", "write") { return }
	res, err := h.orderService.Unsuspend(r.Context(), request.GetID(r))
	if err != nil { response.Error(w, 400, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, res, nil)
}

func (h *StaffManagementHandler) ActivateOrder(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "orders", "write") { return }
	res, err := h.orderService.Activate(r.Context(), request.GetID(r), time.Now().UTC())
	if err != nil { response.Error(w, 400, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, res, nil)
}

func (h *StaffManagementHandler) SyncOrder(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "orders", "write") { return }
	o, err := h.orderRepo.GetByID(r.Context(), request.GetID(r))
	if err != nil { response.Error(w, 404, "NOT_FOUND", "Order not found", nil); return }
	st, err := h.orderService.SyncRemote(r.Context(), o)
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, st, nil)
}

func (h *StaffManagementHandler) ChangeOrderPassword(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "orders", "write") { return }
	o, err := h.orderRepo.GetByID(r.Context(), request.GetID(r))
	if err != nil { response.Error(w, 404, "NOT_FOUND", "Order not found", nil); return }
	var req struct{ Password string }; if request.Decode(r, &req) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	if err := h.orderService.ChangePasswordRemote(r.Context(), o, req.Password); err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, map[string]string{"message": "Changed"}, nil)
}
