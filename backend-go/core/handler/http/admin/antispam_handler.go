package admin

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/antispam"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type AntispamHandler struct {
	staffService *staff.StaffService
	service      *antispam.AntispamService
}

func NewAntispamHandler(staffService *staff.StaffService, service *antispam.AntispamService) *AntispamHandler {
	return &AntispamHandler{staffService: staffService, service: service}
}

func (h *AntispamHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "antispam", "read")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: antispam", nil)
		return
	}

	config, err := h.service.GetConfig(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}
	response.JSON(w, http.StatusOK, config, nil)
}

func (h *AntispamHandler) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "antispam", "write")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: antispam", nil)
		return
	}

	var cfg domain.AntispamConfig
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", nil)
		return
	}

	if err := h.service.UpdateConfig(r.Context(), &cfg); err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, cfg, nil)
}

func (h *AntispamHandler) ListBlockedIPs(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "antispam", "read")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: antispam", nil)
		return
	}

	ips, err := h.service.ListBlockedIPs(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}
	response.JSON(w, http.StatusOK, ips, nil)
}

type AddBlockedIPDTO struct {
	IP     string `json:"ip"`
	Reason string `json:"reason"`
}

func (h *AntispamHandler) AddBlockedIP(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "antispam", "write")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: antispam", nil)
		return
	}

	var dto AddBlockedIPDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", nil)
		return
	}

	if dto.IP == "" {
		response.Error(w, http.StatusBadRequest, "VALIDATION_FAILED", "IP is required", nil)
		return
	}

	blocked, err := h.service.BlockIP(r.Context(), dto.IP, dto.Reason)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_IP", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusCreated, blocked, nil)
}

func (h *AntispamHandler) DeleteBlockedIP(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "antispam", "delete")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: antispam", nil)
		return
	}

	ip := r.PathValue("ip")
	if ip == "" {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "IP path parameter is required", nil)
		return
	}

	if err := h.service.UnblockIP(r.Context(), ip); err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Blocked IP not found", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "IP successfully unblocked"}, nil)
}
