package guest

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/company"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type CompanyHandler struct {
	companyService company.CompanyService
}

func NewCompanyHandler(companyService company.CompanyService) *CompanyHandler {
	return &CompanyHandler{companyService: companyService}
}

func (h *CompanyHandler) GetCompany(w http.ResponseWriter, r *http.Request) {
	data, err := h.companyService.GetPublicCompany(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, data, nil)
}
