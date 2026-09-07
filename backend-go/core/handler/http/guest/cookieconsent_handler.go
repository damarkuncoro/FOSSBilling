package guest

import (
	"net/http"

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
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get cookie consent banner settings", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, cfg, nil)
}
