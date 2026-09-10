package admin

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/request"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type StaffAuthHandler struct{ svc *staff.StaffService }

func NewStaffAuthHandler(s *staff.StaffService) *StaffAuthHandler { return &StaffAuthHandler{s} }

func (h *StaffAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req staff.StaffLoginDTO; if request.Decode(r, &req) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	res, err := h.svc.Login(r.Context(), req, r.RemoteAddr)
	if err != nil { response.Error(w, 401, "UNAUTHORIZED", err.Error(), nil); return }
	response.JSON(w, 200, res, nil)
}

func (h *StaffAuthHandler) VerifyTwoFactor(w http.ResponseWriter, r *http.Request) {
	var req struct{ Email, Code string }; _ = request.Decode(r, &req)
	res, err := h.svc.VerifyTwoFactor(r.Context(), req.Email, req.Code, r.RemoteAddr)
	if err != nil { response.Error(w, 401, "UNAUTHORIZED", err.Error(), nil); return }
	response.JSON(w, 200, res, nil)
}

func (h *StaffAuthHandler) SetupTwoFactor(w http.ResponseWriter, r *http.Request) {
	s, q, err := h.svc.SetupTwoFactor(r.Context(), middleware.GetClientID(r.Context()))
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, map[string]string{"secret": s, "qr_url": q}, nil)
}

func (h *StaffAuthHandler) EnableTwoFactor(w http.ResponseWriter, r *http.Request) {
	var req struct{ Code string }; _ = request.Decode(r, &req)
	if err := h.svc.EnableTwoFactor(r.Context(), middleware.GetClientID(r.Context()), req.Code); err != nil {
		response.Error(w, 400, "ERR", err.Error(), nil); return
	}
	response.JSON(w, 200, map[string]bool{"success": true}, nil)
}

func (h *StaffAuthHandler) DisableTwoFactor(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.DisableTwoFactor(r.Context(), middleware.GetClientID(r.Context())); err != nil {
		response.Error(w, 500, "ERR", err.Error(), nil); return
	}
	response.JSON(w, 200, map[string]bool{"success": true}, nil)
}

func (h *StaffAuthHandler) GetAuditLogs(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.svc, "system", "read") { return }
	l, o := request.GetLimitOffset(r)
	ls, tot, err := h.svc.ListAuditLogs(r.Context(), l, o)
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, ls, &response.Meta{Total: tot, Limit: l, Offset: o})
}
