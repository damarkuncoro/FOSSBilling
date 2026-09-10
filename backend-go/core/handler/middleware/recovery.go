package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic: %v\n%s", err, debug.Stack())
				response.Error(w, 500, "ERR", "Internal error", nil)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
