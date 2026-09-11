package client

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/affiliate"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type AffiliateHandler struct {
	svc *affiliate.AffiliateService
}

func NewAffiliateHandler(s *affiliate.AffiliateService) *AffiliateHandler {
	return &AffiliateHandler{svc: s}
}

func (h *AffiliateHandler) GetMyAffiliate(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	aff, err := h.svc.GetAffiliate(r.Context(), clientID)
	if err != nil {
		response.Error(w, 404, "NOT_FOUND", "Affiliate program not joined", nil)
		return
	}
	response.JSON(w, 200, aff, nil)
}

func (h *AffiliateHandler) Join(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	aff, err := h.svc.ActivateAffiliate(r.Context(), clientID)
	if err != nil {
		response.Error(w, 500, "ERR", err.Error(), nil)
		return
	}
	response.JSON(w, 201, aff, nil)
}
