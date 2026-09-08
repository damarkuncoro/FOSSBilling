package guest

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strings"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/antispam"
	authUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/auth"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/i18n"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type AuthHandler struct {
	authUsecase     *authUsecase.AuthUsecase
	antispamService *antispam.AntispamService
}

func NewAuthHandler(authUsecase *authUsecase.AuthUsecase, antispamService ...*antispam.AntispamService) *AuthHandler {
	var as *antispam.AntispamService
	if len(antispamService) > 0 {
		as = antispamService[0]
	}
	return &AuthHandler{
		authUsecase:     authUsecase,
		antispamService: as,
	}
}

func (h *AuthHandler) SetAntispam(as *antispam.AntispamService) {
	h.antispamService = as
}

func getClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req authUsecase.RegisterDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", nil)
		return
	}

	clientIP := getClientIP(r)
	res, validationErrs, err := h.authUsecase.Register(r.Context(), req, clientIP)
	if err != nil {
		if errors.Is(err, appErrors.ErrInvalidInput) {
			// Check if it's an antispam error mapped to validationErrs
			response.Error(w, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "Validation errors occurred", validationErrs)
			return
		}
		if errors.Is(err, appErrors.ErrDuplicateEntry) {
			response.Error(w, http.StatusConflict, "EMAIL_EXISTS", "Email is already registered", validationErrs)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to register client", err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, res, nil)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req authUsecase.LoginDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", nil)
		return
	}

	res, err := h.authUsecase.Login(r.Context(), req)
	if err != nil {
		locale := i18n.LocaleFromContext(r.Context())
		if errors.Is(err, appErrors.ErrUnauthorized) || errors.Is(err, appErrors.ErrNotFound) {
			msg := i18n.T(locale, "invalid_credentials")
			response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", msg, nil)
			return
		}
		response.Error(w, http.StatusBadRequest, "LOGIN_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, res, nil)
}

func (h *AuthHandler) VerifyTwoFactor(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", nil)
		return
	}

	res, err := h.authUsecase.VerifyTwoFactor(r.Context(), req.Email, req.Code)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, res, nil)
}
