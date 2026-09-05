package client

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fossbilling/backend-go/internal/handler/middleware"
	authUsecase "github.com/fossbilling/backend-go/internal/usecase/auth"
	appErrors "github.com/fossbilling/backend-go/pkg/errors"
	"github.com/fossbilling/backend-go/pkg/response"
)

type ProfileHandler struct {
	authUsecase *authUsecase.AuthUsecase
}

func NewProfileHandler(authUsecase *authUsecase.AuthUsecase) *ProfileHandler {
	return &ProfileHandler{authUsecase: authUsecase}
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
