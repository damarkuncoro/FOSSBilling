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

	// 1. Antispam validation if enabled
	if h.antispamService != nil {
		clientIP := getClientIP(r)
		if err := h.antispamService.ValidateSignup(r.Context(), req.Email, clientIP, req.Honeypot, req.CaptchaToken); err != nil {
			if errors.Is(err, antispam.ErrIPBlocked) {
				response.Error(w, http.StatusForbidden, "IP_BLOCKED", err.Error(), nil)
				return
			}
			if errors.Is(err, antispam.ErrHoneypotTriggered) {
				response.Error(w, http.StatusBadRequest, "BOT_DETECTED", err.Error(), nil)
				return
			}
			if errors.Is(err, antispam.ErrCaptchaFailed) {
				response.Error(w, http.StatusBadRequest, "CAPTCHA_FAILED", err.Error(), nil)
				return
			}
			if errors.Is(err, antispam.ErrDisposableEmail) {
				response.Error(w, http.StatusBadRequest, "DISPOSABLE_EMAIL", err.Error(), nil)
				return
			}
			response.Error(w, http.StatusBadRequest, "SPAM_DETECTED", err.Error(), nil)
			return
		}
	}

	res, validationErrs, err := h.authUsecase.Register(r.Context(), req)
	if err != nil {
		if errors.Is(err, appErrors.ErrInvalidInput) {
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
		if errors.Is(err, appErrors.ErrUnauthorized) || errors.Is(err, appErrors.ErrNotFound) {
			response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid email or password", nil)
			return
		}
		response.Error(w, http.StatusBadRequest, "LOGIN_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, res, nil)
}
