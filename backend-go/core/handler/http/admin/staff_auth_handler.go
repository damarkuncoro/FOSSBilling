package admin

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type StaffAuthHandler struct {
	staffService *staff.StaffService
}

func NewStaffAuthHandler(staffService *staff.StaffService) *StaffAuthHandler {
	return &StaffAuthHandler{staffService: staffService}
}

func (h *StaffAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req staff.StaffLoginDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", nil)
		return
	}

	res, err := h.staffService.Login(r.Context(), req)
	if err != nil {
		if errors.Is(err, appErrors.ErrUnauthorized) {
			response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid staff credentials", nil)
			return
		}
		response.Error(w, http.StatusBadRequest, "LOGIN_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, res, nil)
}

func (h *StaffAuthHandler) VerifyTwoFactor(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	res, err := h.staffService.VerifyTwoFactor(r.Context(), req.Email, req.Code)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
		return
	}
	response.JSON(w, http.StatusOK, res, nil)
}

func (h *StaffAuthHandler) SetupTwoFactor(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	secret, qrURL, err := h.staffService.SetupTwoFactor(r.Context(), staffID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"secret": secret, "qr_url": qrURL}, nil)
}

func (h *StaffAuthHandler) EnableTwoFactor(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	var req struct {
		Code string `json:"code"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	if err := h.staffService.EnableTwoFactor(r.Context(), staffID, req.Code); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_CODE", err.Error(), nil)
		return
	}
	response.JSON(w, http.StatusOK, map[string]bool{"success": true}, nil)
}

func (h *StaffAuthHandler) DisableTwoFactor(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	if err := h.staffService.DisableTwoFactor(r.Context(), staffID); err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}
	response.JSON(w, http.StatusOK, map[string]bool{"success": true}, nil)
}

func (h *StaffAuthHandler) GetAuditLogs(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "system", "read")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: system", nil)
		return
	}

	limit := 50
	offset := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}

	logs, total, err := h.staffService.ListAuditLogs(r.Context(), limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve audit logs", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, logs, &response.Meta{
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}
