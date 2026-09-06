package guest

import (
	"net/http"
	"strings"

	domainUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type DomainHandler struct {
	domainService *domainUsecase.DomainService
}

func NewDomainHandler(domainService *domainUsecase.DomainService) *DomainHandler {
	return &DomainHandler{domainService: domainService}
}

// CheckAvailability handles GET /api/v1/guest/domains/check?domain=example.com
func (h *DomainHandler) CheckAvailability(w http.ResponseWriter, r *http.Request) {
	domainName := strings.TrimSpace(r.URL.Query().Get("domain"))
	if domainName == "" {
		response.Error(w, http.StatusBadRequest, "VALIDATION_FAILED", "Domain query parameter is required", nil)
		return
	}

	result, err := h.domainService.CheckAvailability(r.Context(), domainName)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_DOMAIN", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, result, nil)
}

