package middleware

import (
	"net/http"
	"strings"

	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/geoip"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

// AdminGeofence restricts access to admin routes based on IP country
func AdminGeofence(allowedCountries string) func(http.Handler) http.Handler {
	if allowedCountries == "" {
		return func(next http.Handler) http.Handler { return next }
	}

	allowed := strings.Split(strings.ToUpper(allowedCountries), ",")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only apply to admin routes
			if !strings.HasPrefix(r.URL.Path, "/api/v1/admin/") {
				next.ServeHTTP(w, r)
				return
			}

			ip := r.RemoteAddr
			if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
				ip = strings.Split(xff, ",")[0]
			}

			res := geoip.LookupCountry(ip)

			isAllowed := false
			for _, c := range allowed {
				if c == res.ISOCode {
					isAllowed = true
					break
				}
			}

			if !isAllowed && res.ISOCode != "XX" { // Skip if unknown/internal
				response.Error(w, http.StatusForbidden, "GEOFENCE", "Administrative access is restricted from your country ("+res.Name+")", nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
