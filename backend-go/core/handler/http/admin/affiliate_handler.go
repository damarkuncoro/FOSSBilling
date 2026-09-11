package admin

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/affiliate"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type AffiliateHandler struct {
	svc *affiliate.AffiliateService
	staffSvc *staff.StaffService
}

func NewAffiliateHandler(s *affiliate.AffiliateService, ss *staff.StaffService) *AffiliateHandler {
	return &AffiliateHandler{svc: s, staffSvc: ss}
}

func (h *AffiliateHandler) GetAffiliate(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffSvc, "billing", "read") { return }
	// Implementation for admin to view affiliate details
	response.JSON(w, 200, map[string]string{"status": "ok"}, nil)
}
