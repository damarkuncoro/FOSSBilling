package client

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/support"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/request"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type SupportHandler struct{ svc *support.SupportService }

func NewSupportHandler(s *support.SupportService) *SupportHandler { return &SupportHandler{s} }

func (h *SupportHandler) OpenTicket(w http.ResponseWriter, r *http.Request) {
	cid := middleware.GetClientID(r.Context()); if cid == 0 { response.Error(w, 401, "UNAUTHORIZED", "Login required", nil); return }
	var req struct { HelpdeskID int64; Subject, Message string; RelType *string; RelID *int64 }
	if request.Decode(r, &req) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	t, err := h.svc.OpenTicket(r.Context(), support.CreateTicketDTO{ClientID: cid, HelpdeskID: req.HelpdeskID, Subject: req.Subject, Message: req.Message, Priority: "medium", RelType: req.RelType, RelID: req.RelID, IPAddress: r.RemoteAddr})
	if err != nil { response.Error(w, 400, "ERR", err.Error(), nil); return }
	response.JSON(w, 201, t, nil)
}

func (h *SupportHandler) ListTickets(w http.ResponseWriter, r *http.Request) {
	cid := middleware.GetClientID(r.Context()); if cid == 0 { response.Error(w, 401, "UNAUTHORIZED", "Login required", nil); return }
	l, o := request.GetLimitOffset(r)
	ts, tot, err := h.svc.ListClientTickets(r.Context(), cid, l, o)
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, ts, &response.Meta{Total: tot, Limit: l, Offset: o})
}

func (h *SupportHandler) GetTicket(w http.ResponseWriter, r *http.Request) {
	cid := middleware.GetClientID(r.Context()); if cid == 0 { response.Error(w, 401, "UNAUTHORIZED", "Login required", nil); return }
	d, err := h.svc.GetTicket(r.Context(), request.GetID(r), cid)
	if err != nil { response.Error(w, 403, "FORBIDDEN", err.Error(), nil); return }
	response.JSON(w, 200, d, nil)
}

func (h *SupportHandler) ReplyTicket(w http.ResponseWriter, r *http.Request) {
	cid := middleware.GetClientID(r.Context()); if cid == 0 { response.Error(w, 401, "UNAUTHORIZED", "Login required", nil); return }
	var req struct{ Message string }; if request.Decode(r, &req) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	m, err := h.svc.ClientReply(r.Context(), request.GetID(r), cid, req.Message, r.RemoteAddr)
	if err != nil { response.Error(w, 400, "ERR", err.Error(), nil); return }
	response.JSON(w, 201, m, nil)
}

func (h *SupportHandler) CloseTicket(w http.ResponseWriter, r *http.Request) {
	cid := middleware.GetClientID(r.Context()); if cid == 0 { response.Error(w, 401, "UNAUTHORIZED", "Login required", nil); return }
	if err := h.svc.CloseTicket(r.Context(), request.GetID(r), cid); err != nil { response.Error(w, 400, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, map[string]string{"status": "closed"}, nil)
}
