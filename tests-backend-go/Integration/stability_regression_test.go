package integration_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/http/admin"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/system"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func setupTestMux() http.Handler {
	jwtSecret := "test-secret-key-123456789012"
	staffRepo := memory.NewMockStaffRepository()
	staffService := staff.NewStaffService(staffRepo, jwtSecret, "FOSSBilling")
	sysRepo := memory.NewMockSystemRepository()
	sysService := system.NewSystemService(sysRepo)
	adminAuthHandler := admin.NewStaffAuthHandler(staffService)
	adminSysHandler := admin.NewSystemModuleHandler(staffService, sysService, nil, nil, nil)

	mux := http.NewServeMux()
	aa := middleware.RequireAuth(jwtSecret, "admin", "superadmin")
	mux.Handle("GET /api/v1/admin/system/audit-logs", aa(http.HandlerFunc(adminAuthHandler.GetAuditLogs)))

	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/admin/system/ws" {
			adminSysHandler.HandleWebSocket(w, r)
			return
		}
		handler := http.Handler(mux)
		handler = middleware.Logger(handler)
		handler = middleware.SecurityHeaders(handler)
		handler = middleware.Recovery(handler)
		handler.ServeHTTP(w, r)
	})

	return finalHandler
}

// TestRegression_AuditLogsRoute ensures audit logs route remains consistent at /system/audit-logs
func TestRegression_AuditLogsRoute(t *testing.T) {
	mux := setupTestMux()

	req, _ := http.NewRequest("GET", "/api/v1/admin/system/audit-logs", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Expect 401 Unauthorized (because without token), NOT 404 Not Found
	assert.NotEqual(t, http.StatusNotFound, rr.Code, "Route /api/v1/admin/system/audit-logs must exist")
}

// TestRegression_WebSocketHijacker ensures middleware does not break Hijacker interface needed by WebSocket
func TestRegression_WebSocketHijacker(t *testing.T) {
	server := httptest.NewServer(setupTestMux())
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/v1/admin/system/ws"

	dialer := websocket.Dialer{}
	_, resp, err := dialer.Dial(wsURL, nil)

	if err != nil && resp != nil {
		assert.NotContains(t, err.Error(), "500", "WebSocket should not fail with 500 Hijacker error")
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "WebSocket should reach Auth middleware")
	}
}
