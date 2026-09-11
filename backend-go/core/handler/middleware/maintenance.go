package middleware

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/system"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

func MaintenanceMode(svc *system.SystemService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip check for admin routes or specific guest paths
			if isPathExempt(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			isMaintenance := svc.GetIntSetting(r.Context(), "system", "maintenance_mode", 0) == 1
			if isMaintenance {
				response.Error(w, http.StatusServiceUnavailable, "MAINTENANCE", "System is currently under maintenance. Please try again later.", nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func isPathExempt(path string) bool {
	exemptPrefixes := []string{"/api/v1/admin/", "/health", "/docs"}
	for _, p := range exemptPrefixes {
		if len(path) >= len(p) && path[:len(p)] == p {
			return true
		}
	}
	return false
}
