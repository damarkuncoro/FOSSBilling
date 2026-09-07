package client

import (
	"net/http"
	"strconv"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/activity"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type ActivityHandler struct {
	activityService *activity.ActivityService
}

func NewActivityHandler(activityService *activity.ActivityService) *ActivityHandler {
	return &ActivityHandler{
		activityService: activityService,
	}
}

func (h *ActivityHandler) ListMyLogs(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	logs, total, err := h.activityService.ListClientLogs(r.Context(), clientID, limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	meta := map[string]interface{}{
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}
	response.JSON(w, http.StatusOK, logs, meta)
}
