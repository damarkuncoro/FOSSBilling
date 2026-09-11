package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestE2E_AdminClientsAndSupportTicketDetail(t *testing.T) {
	baseURL := getBaseURL()
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		t.Skipf("API server is not running at %s (skipping live e2e test): %v", baseURL, err)
		return
	}
	_ = resp.Body.Close()

	// 1. Admin Login
	adminLoginPayload := map[string]string{
		"email":    "admin@fossbilling.org",
		"password": "SuperSecretAdmin123!",
	}
	adminBody, _ := json.Marshal(adminLoginPayload)
	adminLoginResp, err := http.Post(baseURL+"/api/v1/admin/auth/login", "application/json", bytes.NewBuffer(adminBody))
	if err != nil {
		t.Fatalf("Admin login request failed: %v", err)
	}
	if adminLoginResp.StatusCode != http.StatusOK {
		// Fallback to simpler password if seeder used it
		adminLoginPayload["password"] = "admin123"
		adminBody, _ = json.Marshal(adminLoginPayload)
		adminLoginResp, _ = http.Post(baseURL+"/api/v1/admin/auth/login", "application/json", bytes.NewBuffer(adminBody))
	}

	if adminLoginResp.StatusCode != http.StatusOK {
		t.Fatalf("Admin login failed with status %d", adminLoginResp.StatusCode)
	}
	defer adminLoginResp.Body.Close()
	var adminData struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	_ = json.NewDecoder(adminLoginResp.Body).Decode(&adminData)
	adminToken := adminData.Data.Token

	// 2. Test Admin List Clients (Verify Fix for 500 error)
	req, _ := http.NewRequest(http.MethodGet, baseURL+"/api/v1/admin/clients", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	clientsResp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Admin list clients failed: %v", err)
	}
	defer clientsResp.Body.Close()

	if clientsResp.StatusCode != http.StatusOK {
		t.Errorf("Admin list clients status = %d; want 200 (possible regression or scanning error)", clientsResp.StatusCode)
	}

	// 3. Setup: Register a client and open a ticket to test ticket detail
	uniqueEmail := fmt.Sprintf("ticket.admin.tester.%d@example.com", time.Now().UnixNano())
	regPayload := map[string]string{
		"email":      uniqueEmail,
		"password":   "SecurePassword123!",
		"first_name": "Support",
		"last_name":  "Tester",
		"country":    "ID",
		"currency":   "USD",
	}
	regBody, _ := json.Marshal(regPayload)
	regResp, _ := http.Post(baseURL+"/api/v1/guest/auth/register", "application/json", bytes.NewBuffer(regBody))
	var regData struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	_ = json.NewDecoder(regResp.Body).Decode(&regData)
	clientToken := regData.Data.Token

	ticketBody, _ := json.Marshal(map[string]interface{}{
		"subject":  "E2E Admin Ticket Test",
		"message":  "Testing admin ticket detail view",
		"priority": "low",
	})
	tReq, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/client/support/tickets", bytes.NewBuffer(ticketBody))
	tReq.Header.Set("Authorization", "Bearer "+clientToken)
	tReq.Header.Set("Content-Type", "application/json")
	tResp, _ := http.DefaultClient.Do(tReq)

	// The response from OpenTicket is the ticket object itself wrapped in Data
	var wrapData struct { Data struct { ID int64 `json:"id" `} `json:"data"`}
	_ = json.NewDecoder(tResp.Body).Decode(&wrapData)
	ticketID := wrapData.Data.ID

	if ticketID == 0 {
		t.Fatal("Failed to create ticket for testing admin view")
	}

	// 4. Test Admin Get Support Ticket Detail (Verify Fix for 404 error)
	detailURL := fmt.Sprintf("%s/api/v1/admin/support/tickets/%d", baseURL, ticketID)
	detailReq, _ := http.NewRequest(http.MethodGet, detailURL, nil)
	detailReq.Header.Set("Authorization", "Bearer "+adminToken)
	detailResp, err := http.DefaultClient.Do(detailReq)
	if err != nil {
		t.Fatalf("Admin get ticket detail failed: %v", err)
	}
	defer detailResp.Body.Close()

	if detailResp.StatusCode != http.StatusOK {
		t.Errorf("Admin get ticket detail status = %d; want 200 (possible missing route/handler)", detailResp.StatusCode)
	}
}
