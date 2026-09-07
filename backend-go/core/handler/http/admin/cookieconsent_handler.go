package admin

import (
	"encoding/json"
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	cookieconsentUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/cookieconsent"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type CookieConsentHandler struct {
	svc *cookieconsentUsecase.CookieConsentService
}

func NewCookieConsentHandler(svc *cookieconsentUsecase.CookieConsentService) *CookieConsentHandler {
	return &CookieConsentHandler{svc: svc}
}

func (h *CookieConsentHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	cfg, err := h.svc.GetConfig(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get cookie consent configuration", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, cfg, nil)
}

func (h *CookieConsentHandler) UpdateConfig(w http.ResponseWriter, r *http.Request) {
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
