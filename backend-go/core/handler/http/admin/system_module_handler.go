package admin

import (
	"net/http"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/page"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/system"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/cache"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/notifications"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/request"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type SystemModuleHandler struct {
	staffService *staff.StaffService; systemService *system.SystemService; pageService *page.PageService; cache cache.Cache; wsHub *notifications.WSHub
}

func NewSystemModuleHandler(s *staff.StaffService, sys *system.SystemService, p *page.PageService, c cache.Cache, ws *notifications.WSHub) *SystemModuleHandler {
	return &SystemModuleHandler{s, sys, p, c, ws}
}

func (h *SystemModuleHandler) GetSecuritySettings(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "system", "read") { return }
	s, err := h.systemService.GetSecuritySettings(r.Context())
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, s, nil)
}

func (h *SystemModuleHandler) UpdateSecuritySettings(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "system", "write") { return }
	var s map[string]any; if request.Decode(r, &s) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	for k, v := range s { _ = h.systemService.UpdateSetting(r.Context(), "security", k, v) }
	response.JSON(w, 200, map[string]bool{"success": true}, nil)
}

func (h *SystemModuleHandler) GetBrandingSettings(w http.ResponseWriter, r *http.Request) {
	s, err := h.systemService.GetBrandingSettings(r.Context())
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, s, nil)
}

func (h *SystemModuleHandler) UpdateBrandingSettings(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "system", "write") { return }
	var s map[string]any; if request.Decode(r, &s) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	for k, v := range s { _ = h.systemService.UpdateSetting(r.Context(), "branding", k, v) }
	response.JSON(w, 200, map[string]bool{"success": true}, nil)
}

func (h *SystemModuleHandler) ExportBackup(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "system", "write") { return }
	b, err := h.systemService.CreateBackup(r.Context())
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, b, nil)
}

func (h *SystemModuleHandler) GetSystemStatus(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "system", "read") { return }
	response.JSON(w, 200, h.systemService.GetSystemStatus(r.Context()), nil)
}

func (h *SystemModuleHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	if h.wsHub != nil {
		h.wsHub.HandleWebSocket(w, r)
	}
}

func (h *SystemModuleHandler) TriggerCron(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, 200, map[string]any{"success": true, "timestamp": time.Now()}, nil)
}
func (h *SystemModuleHandler) GeneratePassword(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, 200, map[string]string{"password": "demo-password-123"}, nil)
}
func (h *SystemModuleHandler) ResolveGeoIP(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, 200, map[string]string{"country": "ID"}, nil)
}

func (h *SystemModuleHandler) ClearCache(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "system", "write") { return }
	_ = h.systemService.PurgeAllCache(r.Context(), h.cache)
	response.JSON(w, 200, map[string]bool{"success": true}, nil)
}

func (h *SystemModuleHandler) ListPages(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "system", "read") { return }
	l, o := request.GetLimitOffset(r)
	ps, tot, err := h.pageService.ListPages(r.Context(), l, o)
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, ps, &response.Meta{Total: tot, Limit: l, Offset: o})
}

func (h *SystemModuleHandler) CreatePage(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "system", "write") { return }
	var p domain.Page; if request.Decode(r, &p) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	if err := h.pageService.CreatePage(r.Context(), &p); err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 201, p, nil)
}

func (h *SystemModuleHandler) GetSystemService() *system.SystemService { return h.systemService }

func (h *SystemModuleHandler) DeletePage(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "system", "delete") { return }
	_ = h.pageService.DeletePage(r.Context(), request.GetID(r))
	response.JSON(w, 200, map[string]bool{"success": true}, nil)
}
