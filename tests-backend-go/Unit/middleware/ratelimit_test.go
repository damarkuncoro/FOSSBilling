package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
)

func TestRateLimiter_AllowAndBlock(t *testing.T) {
	// Limit 3 requests, 1 token refill per 100ms
	limiter := middleware.NewRateLimiter(3, 100*time.Millisecond)

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrapped := limiter.RateLimit(dummyHandler)

	// 1st request -> allow
	req1 := httptest.NewRequest("GET", "/test", nil)
	req1.RemoteAddr = "192.168.1.1:1234"
	rec1 := httptest.NewRecorder()
	wrapped.ServeHTTP(rec1, req1)
	if rec1.Code != http.StatusOK {
		t.Errorf("req1 expected 200, got %d", rec1.Code)
	}

	// 2nd request -> allow
	rec2 := httptest.NewRecorder()
	wrapped.ServeHTTP(rec2, req1)
	if rec2.Code != http.StatusOK {
		t.Errorf("req2 expected 200, got %d", rec2.Code)
	}

	// 3rd request -> allow
	rec3 := httptest.NewRecorder()
	wrapped.ServeHTTP(rec3, req1)
	if rec3.Code != http.StatusOK {
		t.Errorf("req3 expected 200, got %d", rec3.Code)
	}

	// 4th request -> blocked (429)
	rec4 := httptest.NewRecorder()
	wrapped.ServeHTTP(rec4, req1)
	if rec4.Code != http.StatusTooManyRequests {
		t.Errorf("req4 expected 429, got %d", rec4.Code)
	}

	// Wait 250ms for tokens to replenish
	time.Sleep(250 * time.Millisecond)

	// 5th request -> allowed again
	rec5 := httptest.NewRecorder()
	wrapped.ServeHTTP(rec5, req1)
	if rec5.Code != http.StatusOK {
		t.Errorf("req5 expected 200 after refill, got %d", rec5.Code)
	}
}
