package client

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	authUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/auth"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/request"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type ProfileHandler struct {
	auth *authUsecase.AuthUsecase; pwd *authUsecase.PasswordUsecase
}

func NewProfileHandler(u *authUsecase.AuthUsecase, p *authUsecase.PasswordUsecase) *ProfileHandler { return &ProfileHandler{u, p} }

func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	cid := middleware.GetClientID(r.Context()); if cid == 0 { response.Error(w, 401, "UNAUTHORIZED", "Login required", nil); return }
	p, err := h.auth.GetProfile(r.Context(), cid)
	if err != nil { response.Error(w, 404, "NOT_FOUND", "Not found", nil); return }
	response.JSON(w, 200, p, nil)
}

func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	cid := middleware.GetClientID(r.Context()); if cid == 0 { response.Error(w, 401, "UNAUTHORIZED", "Login required", nil); return }
	var req authUsecase.UpdateProfileDTO; if request.Decode(r, &req) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	p, err := h.auth.UpdateProfile(r.Context(), cid, req)
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, p, nil)
}

func (h *ProfileHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	cid := middleware.GetClientID(r.Context()); if cid == 0 { response.Error(w, 401, "UNAUTHORIZED", "Login required", nil); return }
	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if request.Decode(r, &req) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	if err := h.pwd.ChangePassword(r.Context(), cid, req.CurrentPassword, req.NewPassword); err != nil {
		response.Error(w, 401, "UNAUTHORIZED", err.Error(), nil); return
	}
	response.JSON(w, 200, map[string]any{"success": true}, nil)
}

func (h *ProfileHandler) SetupTwoFactor(w http.ResponseWriter, r *http.Request) {
	cid := middleware.GetClientID(r.Context()); if cid == 0 { response.Error(w, 401, "UNAUTHORIZED", "Login required", nil); return }
	res, err := h.auth.SetupTwoFactor(r.Context(), cid)
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, res, nil)
}

func (h *ProfileHandler) EnableTwoFactor(w http.ResponseWriter, r *http.Request) {
	cid := middleware.GetClientID(r.Context()); if cid == 0 { response.Error(w, 401, "UNAUTHORIZED", "Login required", nil); return }
	var req struct{ Code string `json:"code"` }; _ = request.Decode(r, &req)
	if err := h.auth.EnableTwoFactor(r.Context(), cid, req.Code); err != nil { response.Error(w, 400, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, map[string]bool{"success": true}, nil)
}

func (h *ProfileHandler) DisableTwoFactor(w http.ResponseWriter, r *http.Request) {
	cid := middleware.GetClientID(r.Context()); if cid == 0 { response.Error(w, 401, "UNAUTHORIZED", "Login required", nil); return }
	if err := h.auth.DisableTwoFactor(r.Context(), cid); err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, map[string]bool{"success": true}, nil)
}
