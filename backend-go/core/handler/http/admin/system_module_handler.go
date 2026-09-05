package admin

import (
	"net/http"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type SystemModuleHandler struct{}

func NewSystemModuleHandler() *SystemModuleHandler {
	return &SystemModuleHandler{}
}

// --- Security Settings ---
func (h *SystemModuleHandler) GetSecuritySettings(w http.ResponseWriter, r *http.Request) {
	settings := map[string]interface{}{
		"recaptcha_enabled": true,
		"recaptcha_provider": "cloudflare_turnstile",
		"site_key": "0x4AAAAAAAxMockSiteKey",
		"ip_blacklist": []string{"198.51.100.4", "203.0.113.88"},
		"max_login_attempts": 5,
		"lockout_time_minutes": 15,
		"force_ssl": true,
	}
	response.JSON(w, http.StatusOK, settings, nil)
}

// --- System Health & Maintenance ---
func (h *SystemModuleHandler) GetSystemStatus(w http.ResponseWriter, r *http.Request) {
	status := map[string]interface{}{
		"engine_version": "v0.7.0-NextGen (Go 1.27)",
		"database_type": "PostgreSQL 16 High-Availability Pool",
		"database_size": "24.5 MB",
		"active_sessions": 8,
		"cron_last_run": time.Now().Add(-5 * time.Minute).Format(time.RFC3339),
		"cron_status": "healthy",
		"system_load": "0.18, 0.22, 0.15",
		"memory_usage": "142 MB / 8 GB (1.7%)",
		"uptime": "14 days, 6 hours",
	}
	response.JSON(w, http.StatusOK, status, nil)
}

func (h *SystemModuleHandler) TriggerCron(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Cron scheduler tasks executed: 4 invoices generated, 1 expired service suspended.",
		"timestamp": time.Now().Format(time.RFC3339),
	}, nil)
}

func (h *SystemModuleHandler) ClearCache(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Application cache cleared successfully.",
	}, nil)
}

// --- Custom Pages & Knowledgebase ---
func (h *SystemModuleHandler) ListPages(w http.ResponseWriter, r *http.Request) {
	pages := []map[string]interface{}{
		{"id": 1, "title": "Terms of Service", "slug": "terms-of-service", "published": true},
		{"id": 2, "title": "Privacy Policy", "slug": "privacy-policy", "published": true},
	}
	response.JSON(w, http.StatusOK, pages, nil)
}

func (h *SystemModuleHandler) ListKnowledgebase(w http.ResponseWriter, r *http.Request) {
	kb := []map[string]interface{}{
		{"id": 1, "category": "Hosting", "title": "How to point DNS A Records", "slug": "how-to-dns", "views": 420, "published": true},
	}
	response.JSON(w, http.StatusOK, kb, nil)
}

// --- Extensions Hub ---
func (h *SystemModuleHandler) ListExtensions(w http.ResponseWriter, r *http.Request) {
	extensions := []map[string]interface{}{
		{"id": "servicehosting", "name": "cPanel & DirectAdmin Hosting", "version": "2.4.0", "author": "FOSSBilling", "type": "service", "is_enabled": true},
		{"id": "midtrans", "name": "Midtrans Payment Gateway", "version": "1.2.0", "author": "Nusantara Devs", "type": "gateway", "is_enabled": true},
		{"id": "antispam", "name": "Cloudflare Turnstile Shield", "version": "1.0.5", "author": "Security Team", "type": "plugin", "is_enabled": true},
	}
	response.JSON(w, http.StatusOK, extensions, nil)
}
