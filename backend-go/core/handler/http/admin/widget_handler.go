package admin

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

func (h *WidgetHandler) GetRegistry(w http.ResponseWriter, r *http.Request) {
	registry := h.svc.GetRegistry(r.Context())
	response.JSON(w, http.StatusOK, registry, nil)
}
