package guest

import (
	"net/http"

	widgetUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/widget"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type WidgetHandler struct {
	svc *widgetUsecase.WidgetService
}

func NewWidgetHandler(svc *widgetUsecase.WidgetService) *WidgetHandler {
	return &WidgetHandler{svc: svc}
}

func (h *WidgetHandler) GetSlotWidgets(w http.ResponseWriter, r *http.Request) {
	slot := r.URL.Query().Get("slot")
	if slot == "" {
		response.Error(w, http.StatusBadRequest, "MISSING_SLOT", "Slot parameter is required", nil)
		return
	}

	widgets := h.svc.GetWidgetsForSlot(r.Context(), slot)
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"slot":    slot,
		"widgets": widgets,
	}, nil)
}
