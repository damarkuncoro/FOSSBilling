package guest

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/antispam"
	authUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/auth"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/request"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type AuthHandler struct {
	authUsecase     *authUsecase.AuthUsecase
	antispamService *antispam.AntispamService
}

func NewAuthHandler(u *authUsecase.AuthUsecase, as ...*antispam.AntispamService) *AuthHandler {
	var a *antispam.AntispamService
	if len(as) > 0 {
		a = as[0]
	}
	return &AuthHandler{u, a}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req authUsecase.RegisterDTO; if request.Decode(r, &req) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	res, v, err := h.authUsecase.Register(r.Context(), req, r.RemoteAddr)
	if err != nil { response.Error(w, 400, "ERR", err.Error(), v); return }
	response.JSON(w, 201, res, nil)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req authUsecase.LoginDTO; if request.Decode(r, &req) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	res, err := h.authUsecase.Login(r.Context(), req)
	if err != nil { response.Error(w, 401, "UNAUTHORIZED", err.Error(), nil); return }
	response.JSON(w, 200, res, nil)
}

func (h *AuthHandler) VerifyTwoFactor(w http.ResponseWriter, r *http.Request) {
	var req struct{ Email, Code string }; if request.Decode(r, &req) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	res, err := h.authUsecase.VerifyTwoFactor(r.Context(), req.Email, req.Code)
	if err != nil { response.Error(w, 401, "UNAUTHORIZED", err.Error(), nil); return }
	response.JSON(w, 200, res, nil)
}
