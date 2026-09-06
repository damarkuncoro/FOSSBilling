package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

// Recovery handles panics gracefully and returns a standardized JSON 500 error response
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("🔥 [Panic Recovered] %v\nStacktrace:\n%s", err, debug.Stack())
				response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected internal server error occurred", nil)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
