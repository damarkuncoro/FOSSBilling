package client

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	orderUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/order"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/request"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type OrderHandler struct{ svc *orderUsecase.OrderService }

func NewOrderHandler(s *orderUsecase.OrderService) *OrderHandler { return &OrderHandler{s} }

func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	cid := middleware.GetClientID(r.Context()); if cid == 0 { response.Error(w, 401, "UNAUTHORIZED", "Login required", nil); return }
	l, o := request.GetLimitOffset(r)
	os, tot, err := h.svc.ListByClientID(r.Context(), cid, l, o)
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, os, &response.Meta{Total: tot, Limit: l, Offset: o})
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	cid := middleware.GetClientID(r.Context()); if cid == 0 { response.Error(w, 401, "UNAUTHORIZED", "Login required", nil); return }
	o, err := h.svc.GetByIDForClient(r.Context(), cid, request.GetID(r))
	if err != nil { response.Error(w, 404, "NOT_FOUND", "Not found", nil); return }
	response.JSON(w, 200, o, nil)
}

func (h *OrderHandler) SyncStatus(w http.ResponseWriter, r *http.Request) {
	cid := middleware.GetClientID(r.Context()); if cid == 0 { response.Error(w, 401, "UNAUTHORIZED", "Login required", nil); return }
	o, err := h.svc.GetByIDForClient(r.Context(), cid, request.GetID(r))
	if err != nil { response.Error(w, 404, "NOT_FOUND", "Not found", nil); return }
	st, err := h.svc.SyncRemote(r.Context(), o)
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, st, nil)
}

func (h *OrderHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	cid := middleware.GetClientID(r.Context()); if cid == 0 { response.Error(w, 401, "UNAUTHORIZED", "Login required", nil); return }
	o, err := h.svc.GetByIDForClient(r.Context(), cid, request.GetID(r))
	if err != nil { response.Error(w, 404, "NOT_FOUND", "Not found", nil); return }
	var req struct{ Password string }; if request.Decode(r, &req) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	if err := h.svc.ChangePasswordRemote(r.Context(), o, req.Password); err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, map[string]string{"message": "Changed"}, nil)
}
