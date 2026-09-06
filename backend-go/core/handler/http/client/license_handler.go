package client

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	licenseUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/license"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type LicenseHandler struct {
	licenseService *licenseUsecase.LicenseService
}

func NewLicenseHandler(licenseService *licenseUsecase.LicenseService) *LicenseHandler {
	return &LicenseHandler{licenseService: licenseService}
}

// ListLicenses handles GET /api/v1/client/licenses
func (h *LicenseHandler) ListLicenses(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	if clientID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}

	list, err := h.licenseService.ListClientLicenses(r.Context(), clientID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, list, nil)
}

// ResetLicenseLock handles POST /api/v1/client/licenses/{id}/reset
func (h *LicenseHandler) ResetLicenseLock(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	if clientID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid license ID", nil)
		return
	}

	if err := h.licenseService.ResetLicenseLock(r.Context(), clientID, id); err != nil {
		if errors.Is(err, licenseUsecase.ErrUnauthorizedLicense) {
			response.Error(w, http.StatusForbidden, "FORBIDDEN", "You do not own this license", nil)
			return
		}
		if errors.Is(err, licenseUsecase.ErrLicenseNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "License not found", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{"id": id, "message": "Domain and IP lock reset successfully"}, nil)
}
