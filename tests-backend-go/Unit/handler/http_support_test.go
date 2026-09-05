package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

func TestHTTP_SupportAndAdminLifecycle(t *testing.T) {
	ts, promoRepo, staffRepo := setupTestServer()
	defer ts.Close()
	ctx := context.Background()
	setupTestAdminsAndPromos(ctx, promoRepo, staffRepo)

	regBody, _ := json.Marshal(map[string]interface{}{
		"email": "supporter@example.com", "password": "Password123!",
		"first_name": "Support", "last_name": "Tester", "country": "ID", "currency": "USD",
	})
	regResp, err := http.Post(ts.URL+"/api/v1/guest/auth/register", "application/json", bytes.NewBuffer(regBody))
	if err != nil || regResp.StatusCode != http.StatusCreated {
		t.Fatalf("Register failed: %v, status: %d", err, regResp.StatusCode)
	}

	var regData struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	_ = json.NewDecoder(regResp.Body).Decode(&regData)
	clientToken := regData.Data.Token

	// 1. Open Ticket
	ticketBody, _ := json.Marshal(map[string]interface{}{
		"helpdesk_id": 1, "subject": "Need help with SSH keys",
		"message": "How do I add my public key?", "priority": "high",
	})
	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/client/support/tickets", bytes.NewBuffer(ticketBody))
	req.Header.Set("Authorization", "Bearer "+clientToken)
	ticketResp, err := http.DefaultClient.Do(req)
	if err != nil || ticketResp.StatusCode != http.StatusCreated {
		t.Fatalf("Open ticket failed: %v, status: %d", err, ticketResp.StatusCode)
	}

	var ticketData struct{ Data domain.Ticket `json:"data"` }
	_ = json.NewDecoder(ticketResp.Body).Decode(&ticketData)
	ticketID := ticketData.Data.ID

	// 2. Admin Login
	adminLoginBody, _ := json.Marshal(map[string]string{
		"email": "admin@fossbilling.org", "password": "SuperSecretAdmin123!",
	})
	adminLoginResp, err := http.Post(ts.URL+"/api/v1/admin/auth/login", "application/json", bytes.NewBuffer(adminLoginBody))
	if err != nil || adminLoginResp.StatusCode != http.StatusOK {
		t.Fatalf("Admin login failed: %v, status: %d", err, adminLoginResp.StatusCode)
	}
	var adminAuthData struct{ Data struct{ Token string `json:"token"` } `json:"data"` }
	_ = json.NewDecoder(adminLoginResp.Body).Decode(&adminAuthData)
	adminToken := adminAuthData.Data.Token

	// 3. Staff Reply
	staffReplyBody, _ := json.Marshal(map[string]string{
		"message": "You can paste your SSH public key in the cPanel security section.",
	})
	req, _ = http.NewRequest("POST", fmt.Sprintf("%s/api/v1/admin/support/tickets/%d/reply", ts.URL, ticketID), bytes.NewBuffer(staffReplyBody))
	req.Header.Set("Authorization", "Bearer "+adminToken)
	staffReplyResp, err := http.DefaultClient.Do(req)
	if err != nil || staffReplyResp.StatusCode != http.StatusCreated {
		t.Fatalf("Staff reply failed: %v, status: %d", err, staffReplyResp.StatusCode)
	}

	// 4. Client Closes Ticket
	req, _ = http.NewRequest("POST", fmt.Sprintf("%s/api/v1/client/support/tickets/%d/close", ts.URL, ticketID), nil)
	req.Header.Set("Authorization", "Bearer "+clientToken)
	closeResp, err := http.DefaultClient.Do(req)
	if err != nil || closeResp.StatusCode != http.StatusOK {
		t.Fatalf("Close ticket failed: %v, status: %d", err, closeResp.StatusCode)
	}

	// 5. Admin Checks Audit Logs
	req, _ = http.NewRequest("GET", ts.URL+"/api/v1/admin/audit-logs", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	auditResp, err := http.DefaultClient.Do(req)
	if err != nil || auditResp.StatusCode != http.StatusOK {
		t.Fatalf("Get audit logs failed: %v, status: %d", err, auditResp.StatusCode)
	}

	var auditData struct{ Data []domain.AuditLog `json:"data"` }
	_ = json.NewDecoder(auditResp.Body).Decode(&auditData)
	if len(auditData.Data) == 0 {
		t.Error("Expected at least 1 audit log entry")
	}
}
