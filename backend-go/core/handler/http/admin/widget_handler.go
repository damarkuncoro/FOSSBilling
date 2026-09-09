package admin

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	widgetUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/widget"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type WidgetHandler struct {
	staffService *staff.StaffService
	svc          *widgetUsecase.WidgetService
}

func NewWidgetHandler(staffService *staff.StaffService, svc *widgetUsecase.WidgetService) *WidgetHandler {
	return &WidgetHandler{staffService: staffService, svc: svc}
}

func (h *WidgetHandler) GetRegistry(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "system", "read")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: system", nil)
		return
	}

	registry := h.svc.GetRegistry(r.Context())
	response.JSON(w, http.StatusOK, registry, nil)
}
