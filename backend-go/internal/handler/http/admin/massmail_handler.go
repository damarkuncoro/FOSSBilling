package admin

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/fossbilling/backend-go/internal/handler/middleware"
	"github.com/fossbilling/backend-go/internal/usecase/massmail"
	"github.com/fossbilling/backend-go/internal/usecase/staff"
	"github.com/fossbilling/backend-go/pkg/response"
)

type MassMailHandler struct {
	massMailService *massmail.MassMailService
	staffService    *staff.StaffService
}

func NewMassMailHandler(massMailService *massmail.MassMailService, staffService *staff.StaffService) *MassMailHandler {
	return &MassMailHandler{
		massMailService: massMailService,
		staffService:    staffService,
	}
}

type CreateCampaignRequest struct {
	Subject string `json:"subject"`
	Content string `json:"content"`
}

func (h *MassMailHandler) List(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "system", "read")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for mass mail", nil)
		return
	}

	limit := 20
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

	campaigns, total, err := h.massMailService.List(r.Context(), limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve campaigns", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, campaigns, &response.Meta{
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

func (h *MassMailHandler) Create(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "system", "write")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for mass mail", nil)
		return
	}

	var req CreateCampaignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", nil)
		return
	}

	campaign, err := h.massMailService.Create(r.Context(), staffID, req.Subject, req.Content)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "CREATE_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusCreated, campaign, nil)
}

func (h *MassMailHandler) Send(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "system", "write")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for mass mail", nil)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	// e.g. api/v1/admin/mass-mail/1/send -> id is parts[len(parts)-2]
	var idStr string
	if len(parts) >= 2 && parts[len(parts)-1] == "send" {
		idStr = parts[len(parts)-2]
	} else {
		idStr = parts[len(parts)-1]
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid campaign ID", nil)
		return
	}

	campaign, err := h.massMailService.Send(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "SEND_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, campaign, nil)
}
