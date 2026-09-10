package guest

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/system"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type SystemHandler struct {
	healthUsecase *system.HealthUsecase
}

func NewSystemHandler(healthUsecase *system.HealthUsecase) *SystemHandler {
	return &SystemHandler{healthUsecase: healthUsecase}
}

func (h *SystemHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	status := h.healthUsecase.Check(r.Context())

	code := http.StatusOK
	if status.Status == "error" {
		code = http.StatusServiceUnavailable
	}

	response.JSON(w, code, status, nil)
}
