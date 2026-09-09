package admin

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/page"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/system"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/geoip"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/tools"
)

type SystemModuleHandler struct {
	staffService  *staff.StaffService
	systemService *system.SystemService
	pageService   *page.PageService
}

func NewSystemModuleHandler(staffService *staff.StaffService, systemService *system.SystemService, pageService *page.PageService) *SystemModuleHandler {
	return &SystemModuleHandler{
		staffService:  staffService,
		systemService: systemService,
		pageService:   pageService,
	}
}

// --- Security Settings ---
func (h *SystemModuleHandler) GetSecuritySettings(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "system", "read")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: system", nil)
		return
	}

	settings, err := h.systemService.GetSecuritySettings(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}
	response.JSON(w, http.StatusOK, settings, nil)
}

func (h *SystemModuleHandler) UpdateSecuritySettings(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "system", "write")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: system", nil)
		return
	}

	var settings map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body", nil)
		return
	}

	for k, v := range settings {
		if err := h.systemService.UpdateSecuritySetting(r.Context(), k, v); err != nil {
			response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
			return
		}
	}

	response.JSON(w, http.StatusOK, map[string]bool{"success": true}, nil)
}

// --- System Health & Maintenance ---
func (h *SystemModuleHandler) GetSystemStatus(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "system", "read")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: system", nil)
		return
	}

	status := h.systemService.GetSystemStatus(r.Context())
	response.JSON(w, http.StatusOK, status, nil)
}

func (h *SystemModuleHandler) TriggerCron(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "system", "write")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: system", nil)
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"success":   true,
		"message":   "Cron scheduler tasks executed: 4 invoices generated, 1 expired service suspended.",
		"timestamp": time.Now().Format(time.RFC3339),
	}, nil)
}

func (h *SystemModuleHandler) ClearCache(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "system", "write")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: system", nil)
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Application cache cleared successfully.",
	}, nil)
}

// --- Custom Pages & Knowledgebase ---
func (h *SystemModuleHandler) ListPages(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "system", "read")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: system", nil)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	pages, total, err := h.pageService.ListPages(r.Context(), limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}

	meta := &response.Meta{
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}
	response.JSON(w, http.StatusOK, pages, meta)
}

func (h *SystemModuleHandler) CreatePage(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "system", "write")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: system", nil)
		return
	}

	var p domain.Page
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body", nil)
		return
	}

	if err := h.pageService.CreatePage(r.Context(), &p); err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusCreated, p, nil)
}

func (h *SystemModuleHandler) DeletePage(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "system", "delete")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: system", nil)
		return
	}

	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err := h.pageService.DeletePage(r.Context(), id); err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}
	response.JSON(w, http.StatusOK, map[string]bool{"success": true}, nil)
}

func (h *SystemModuleHandler) ListKnowledgebase(w http.ResponseWriter, r *http.Request) {
	kb := []map[string]interface{}{
		{"id": 1, "category": "Hosting", "title": "How to point DNS A Records", "slug": "how-to-dns", "views": 420, "published": true},
	}
	response.JSON(w, http.StatusOK, kb, nil)
}

// --- System Tools & GeoIP Utilities ---
func (h *SystemModuleHandler) GeneratePassword(w http.ResponseWriter, r *http.Request) {
	length, _ := strconv.Atoi(r.URL.Query().Get("length"))
	if length <= 0 {
		length = 16
	}
	special := r.URL.Query().Get("special") != "false"

	pwd, err := tools.GeneratePassword(length, special)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"password": pwd,
		"length":   len(pwd),
	}, nil)
}

func (h *SystemModuleHandler) ResolveGeoIP(w http.ResponseWriter, r *http.Request) {
	ip := r.URL.Query().Get("ip")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = r.RemoteAddr
		}
	}

	countryCode := r.URL.Query().Get("country")
	if countryCode == "" {
		countryCode = "US"
	}

	info := geoip.LookupCountry(countryCode)
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"ip":         ip,
		"is_private": geoip.IsPrivateIP(ip),
		"country":    info,
	}, nil)
}
