package admin

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/massmail"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/request"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type MassMailHandler struct {
	svc *massmail.MassMailService; staffSvc *staff.StaffService
}

func NewMassMailHandler(s *massmail.MassMailService, ss *staff.StaffService) *MassMailHandler { return &MassMailHandler{s, ss} }

func (h *MassMailHandler) List(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffSvc, "system", "read") { return }
	l, o := request.GetLimitOffset(r)
	cs, tot, err := h.svc.List(r.Context(), l, o)
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, cs, &response.Meta{Total: tot, Limit: l, Offset: o})
}

func (h *MassMailHandler) Create(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffSvc, "system", "write") { return }
	var req struct{ Subject, Content string }; if request.Decode(r, &req) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	c, err := h.svc.Create(r.Context(), middleware.GetClientID(r.Context()), req.Subject, req.Content)
	if err != nil { response.Error(w, 400, "ERR", err.Error(), nil); return }
	response.JSON(w, 201, c, nil)
}

func (h *MassMailHandler) Send(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffSvc, "system", "write") { return }
	c, err := h.svc.Send(r.Context(), request.GetID(r))
	if err != nil { response.Error(w, 400, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, c, nil)
}
