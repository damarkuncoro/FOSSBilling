package admin

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/stats"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type StatsHandler struct {
	statsService *stats.StatsService
	staffService *staff.StaffService
}

func NewStatsHandler(statsService *stats.StatsService, staffService *staff.StaffService) *StatsHandler {
	return &StatsHandler{
		statsService: statsService,
		staffService: staffService,
	}
}

func (h *StatsHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "system", "read")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: system", nil)
		return
	}

	dashboard, err := h.statsService.CalculateDashboard(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to calculate stats", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, dashboard, nil)
}
