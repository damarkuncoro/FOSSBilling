package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"net"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

// TestRegression_AuditLogsRoute mematikan rute audit logs tetap konsisten di /system/audit-logs
func TestRegression_AuditLogsRoute(t *testing.T) {
	mux := setupTestMux() // Fungsi helper yang memuat setupRoutes kita

	req, _ := http.NewRequest("GET", "/api/v1/admin/system/audit-logs", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Kita mengharapkan 401 Unauthorized (karena tanpa token), BUKAN 404 Not Found
	assert.NotEqual(t, http.StatusNotFound, rr.Code, "Route /api/v1/admin/system/audit-logs must exist")
}

// TestRegression_WebSocketHijacker mematikan middleware tidak merusak interface Hijacker yang dibutuhkan WebSocket
func TestRegression_WebSocketHijacker(t *testing.T) {
	server := httptest.NewServer(setupTestMux())
	defer server.Close()

	// Ubah http:// ke ws://
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/v1/admin/system/ws"

	dialer := websocket.Dialer{}
	_, resp, err := dialer.Dial(wsURL, nil)

	// Jika error adalah 401, berarti Hijacking berhasil (sampai ke middleware auth)
	// Jika error 500 dengan pesan "does not implement http.Hijacker", berarti regresi terulang
	if err != nil {
		assert.NotContains(t, err.Error(), "500", "WebSocket should not fail with 500 Hijacker error")
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "WebSocket should reach Auth middleware")
	}
}
