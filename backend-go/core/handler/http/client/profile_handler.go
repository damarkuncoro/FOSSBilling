package client

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	authUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/auth"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type ProfileHandler struct {
	authUsecase     *authUsecase.AuthUsecase
	passwordUsecase *authUsecase.PasswordUsecase
}

func NewProfileHandler(authUsecase *authUsecase.AuthUsecase, passwordUsecase *authUsecase.PasswordUsecase) *ProfileHandler {
	return &ProfileHandler{
		authUsecase:     authUsecase,
		passwordUsecase: passwordUsecase,
	}
}

func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	if clientID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}

	profile, err := h.authUsecase.GetProfile(r.Context(), clientID)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Client profile not found", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve profile", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, profile, nil)
}

func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	if clientID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}

	var req authUsecase.UpdateProfileDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", nil)
		return
	}

	profile, err := h.authUsecase.UpdateProfile(r.Context(), clientID, req)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Client profile not found", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update profile", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, profile, nil)
}

type changePasswordReq struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

func (h *ProfileHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	if clientID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}

	var req changePasswordReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request payload", nil)
		return
	}

	if h.passwordUsecase == nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Password service uninitialized", nil)
		return
	}

	if err := h.passwordUsecase.ChangePassword(r.Context(), clientID, req.CurrentPassword, req.NewPassword); err != nil {
		if errors.Is(err, appErrors.ErrUnauthorized) {
			response.Error(w, http.StatusUnauthorized, "INVALID_PASSWORD", "Current password is incorrect", nil)
			return
		}
		response.Error(w, http.StatusBadRequest, "PASSWORD_UPDATE_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Password changed successfully",
	}, nil)
}
