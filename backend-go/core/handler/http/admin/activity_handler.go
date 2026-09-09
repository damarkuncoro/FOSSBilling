package admin

import (
	"net/http"
	"strconv"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/activity"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type ActivityHandler struct {
	staffService    *staff.StaffService
	activityService *activity.ActivityService
}

func NewActivityHandler(staffService *staff.StaffService, activityService *activity.ActivityService) *ActivityHandler {
	return &ActivityHandler{
		staffService:    staffService,
		activityService: activityService,
	}
}

func (h *ActivityHandler) ListLogs(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "activity", "read")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: activity", nil)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	logs, total, err := h.activityService.ListLogs(r.Context(), limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}

	meta := &response.Meta{
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}
	response.JSON(w, http.StatusOK, logs, meta)
}
