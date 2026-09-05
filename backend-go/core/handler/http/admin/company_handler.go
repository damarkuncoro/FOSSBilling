package admin

import (
	"encoding/json"
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/company"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type CompanyHandler struct {
	companyService company.CompanyService
	staffService   *staff.StaffService
}

func NewCompanyHandler(companyService company.CompanyService, staffService *staff.StaffService) *CompanyHandler {
	return &CompanyHandler{
		companyService: companyService,
		staffService:   staffService,
	}
}

func (h *CompanyHandler) GetCompany(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "system", "read")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: system", nil)
		return
	}

	data, err := h.companyService.GetCompany(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve company settings", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, data, nil)
}

func (h *CompanyHandler) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "system", "manage")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: system", nil)
		return
	}

	var req domain.CompanySettings
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON payload", err.Error())
		return
	}

	updated, err := h.companyService.UpdateCompany(r.Context(), &req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, updated, nil)
}
