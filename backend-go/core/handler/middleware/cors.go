package middleware

import (
	"net/http"
	"strings"
)

func CORS(allowed string) func(http.Handler) http.Handler {
	if allowed == "" { allowed = "*" }; os := strings.Split(allowed, ",")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			org, ok := r.Header.Get("Origin"), false
			if allowed == "*" { ok = true } else { for _, o := range os { if o == org { ok = true; break } } }
			if ok && org != "" { w.Header().Set("Access-Control-Allow-Origin", org) } else if allowed == "*" { w.Header().Set("Access-Control-Allow-Origin", "*") }
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			if r.Method == "OPTIONS" { w.WriteHeader(200); return }
			next.ServeHTTP(w, r)
		})
	}
}
