package client

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	license "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/license"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/request"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type LicenseHandler struct{ svc *license.LicenseService }

func NewLicenseHandler(s *license.LicenseService) *LicenseHandler { return &LicenseHandler{s} }

func (h *LicenseHandler) ListLicenses(w http.ResponseWriter, r *http.Request) {
	cid := middleware.GetClientID(r.Context()); if cid == 0 { response.Error(w, 401, "UNAUTHORIZED", "Login required", nil); return }
	ls, err := h.svc.ListClientLicenses(r.Context(), cid)
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, ls, nil)
}

func (h *LicenseHandler) ResetLicenseLock(w http.ResponseWriter, r *http.Request) {
	cid := middleware.GetClientID(r.Context()); if cid == 0 { response.Error(w, 401, "UNAUTHORIZED", "Login required", nil); return }
	if err := h.svc.ResetLicenseLock(r.Context(), cid, request.GetID(r)); err != nil { response.Error(w, 403, "FORBIDDEN", err.Error(), nil); return }
	response.JSON(w, 200, map[string]any{"success": true}, nil)
}
