package admin

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/stats"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type StatsHandler struct {
	svc *stats.StatsService; staffSvc *staff.StaffService
}

func NewStatsHandler(s *stats.StatsService, ss *staff.StaffService) *StatsHandler { return &StatsHandler{s, ss} }

func (h *StatsHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffSvc, "system", "read") { return }
	d, err := h.svc.CalculateDashboard(r.Context())
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, d, nil)
}
