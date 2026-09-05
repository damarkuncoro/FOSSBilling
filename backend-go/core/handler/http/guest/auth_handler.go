package guest

import (
	"encoding/json"
	"errors"
	"net/http"

	authUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/auth"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type AuthHandler struct {
	authUsecase *authUsecase.AuthUsecase
}

func NewAuthHandler(authUsecase *authUsecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req authUsecase.RegisterDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", nil)
		return
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
