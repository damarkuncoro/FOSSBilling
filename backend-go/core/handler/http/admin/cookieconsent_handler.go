package admin

import (
	"encoding/json"
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	cookieconsentUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/cookieconsent"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type CookieConsentHandler struct {
	staffService *staff.StaffService
	svc          *cookieconsentUsecase.CookieConsentService
}

func NewCookieConsentHandler(staffService *staff.StaffService, svc *cookieconsentUsecase.CookieConsentService) *CookieConsentHandler {
	return &CookieConsentHandler{staffService: staffService, svc: svc}
}

func (h *CookieConsentHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "system", "read")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: system", nil)
		return
	}

	cfg, err := h.svc.GetConfig(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get cookie consent configuration", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, cfg, nil)
}

func (h *CookieConsentHandler) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "system", "write")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: system", nil)
		return
	}

	var req domain.CookieConsentConfig
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_BODY", "Invalid JSON payload", err.Error())
		return
	}

	if err := h.svc.UpdateConfig(r.Context(), req); err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update cookie consent configuration", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Cookie consent configuration updated successfully",
		"config":  req,
	}, nil)
}
