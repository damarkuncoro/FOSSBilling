package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/metrics"
)

func Metrics() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "/ws") {
				next.ServeHTTP(w, r)
				return
			}
			start := time.Now()

			rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rw, r)

			duration := time.Since(start).Seconds()
			status := strconv.Itoa(rw.status)

			metrics.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
			metrics.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
		})
	}
}
