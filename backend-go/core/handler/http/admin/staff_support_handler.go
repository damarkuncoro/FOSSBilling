package admin

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/request"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

func (h *StaffManagementHandler) GetTicket(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "support", "read") { return }
	t, err := h.supportService.GetTicket(r.Context(), request.GetID(r), 0)
	if err != nil { response.Error(w, 404, "NOT_FOUND", "Ticket not found", nil); return }
	response.JSON(w, 200, t, nil)
}

func (h *StaffManagementHandler) ListTickets(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "support", "read") { return }
	l, o := request.GetLimitOffset(r)
	ts, tot, err := h.supportService.ListAllTickets(r.Context(), l, o)
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, ts, &response.Meta{Total: tot, Limit: l, Offset: o})
}

func (h *StaffManagementHandler) ReplyTicket(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "support", "write") { return }
	var req struct{ Message string }; if request.Decode(r, &req) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	m, err := h.supportService.StaffReply(r.Context(), request.GetID(r), middleware.GetClientID(r.Context()), req.Message)
	if err != nil { response.Error(w, 400, "ERR", err.Error(), nil); return }
	response.JSON(w, 201, m, nil)
}
