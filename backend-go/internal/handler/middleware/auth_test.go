package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fossbilling/backend-go/pkg/auth"
)

func TestRequireAuth_ValidToken(t *testing.T) {
	jwtSecret := "test-secret-32-characters-key-12"
	token, _ := auth.GenerateToken(jwtSecret, 99, "user@test.com", "client", 1*time.Hour)

	var capturedID int64
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedID = GetClientID(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	mw := RequireAuth(jwtSecret, "client")
	ts := httptest.NewServer(mw(dummyHandler))
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d; want %d", resp.StatusCode, http.StatusOK)
	}
	if capturedID != 99 {
		t.Errorf("Captured ClientID = %d; want 99", capturedID)
	}
}

func TestRequireAuth_MissingOrInvalidToken(t *testing.T) {
	jwtSecret := "test-secret-key"
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mw := RequireAuth(jwtSecret, "client")
	ts := httptest.NewServer(mw(dummyHandler))
	defer ts.Close()

	// 1. Missing header
	req1, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
	resp1, _ := http.DefaultClient.Do(req1)
	if resp1.StatusCode != http.StatusUnauthorized {
		t.Errorf("Missing token StatusCode = %d; want 401", resp1.StatusCode)
	}

	// 2. Invalid Token
	req2, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
	req2.Header.Set("Authorization", "Bearer invalid-junk-token")
	resp2, _ := http.DefaultClient.Do(req2)
	if resp2.StatusCode != http.StatusUnauthorized {
		t.Errorf("Invalid token StatusCode = %d; want 401", resp2.StatusCode)
	}
}
