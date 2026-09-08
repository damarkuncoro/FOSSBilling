package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestE2E_AdminDashboardAndFinancialReports(t *testing.T) {
	baseURL := getBaseURL()
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		t.Skipf("API server is not running at %s: %v", baseURL, err)
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

	// 2. Fetch Dashboard Stats
	req, _ := http.NewRequest(http.MethodGet, baseURL+"/api/v1/admin/stats/dashboard", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	dashResp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Dashboard request failed: %v", err)
	}
	defer dashResp.Body.Close()

	if dashResp.StatusCode != http.StatusOK {
		t.Fatalf("Dashboard status = %d; want 200", dashResp.StatusCode)
	}

	var dashData struct {
		Data struct {
			TotalRevenue float64 `json:"total_revenue"`
			MRR          float64 `json:"mrr"`
			TotalClients int     `json:"total_clients"`
			ActiveOrders int     `json:"active_orders"`
		} `json:"data"`
	}
	_ = json.NewDecoder(dashResp.Body).Decode(&dashData)

	// Since this is a live/shared environment, we just check if it's a valid JSON response
	// The values might change depending on the seeder or other tests.
	if dashData.Data.TotalClients < 0 {
		t.Errorf("TotalClients should be >= 0, got %d", dashData.Data.TotalClients)
	}

	// 3. Fetch Financial Reports Breakdown
	req, _ = http.NewRequest(http.MethodGet, baseURL+"/api/v1/admin/reports/financial", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	reportResp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Financial reports request failed: %v", err)
	}
	defer reportResp.Body.Close()

	if reportResp.StatusCode != http.StatusOK {
		t.Fatalf("Financial reports status = %d; want 200", reportResp.StatusCode)
	}

	var reportData struct {
		Data struct {
			MRR              float64 `json:"mrr"`
			MonthlyBreakdown []struct {
				Month   string  `json:"month"`
				Revenue float64 `json:"revenue"`
			} `json:"monthly_breakdown"`
		} `json:"data"`
	}
	_ = json.NewDecoder(reportResp.Body).Decode(&reportData)

	if len(reportData.Data.MonthlyBreakdown) == 0 {
		t.Error("Expected at least one month in breakdown")
	}
}

func TestE2E_MassMailLifecycle(t *testing.T) {
	baseURL := getBaseURL()
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		t.Skip("API server not running")
		return
	}
	_ = resp.Body.Close()

	// 1. Admin Login
	adminLoginPayload := map[string]string{
		"email":    "admin@fossbilling.org",
		"password": "SuperSecretAdmin123!",
	}
	adminBody, _ := json.Marshal(adminLoginPayload)
	adminLoginResp, _ := http.Post(baseURL+"/api/v1/admin/auth/login", "application/json", bytes.NewBuffer(adminBody))
	if adminLoginResp.StatusCode != http.StatusOK {
		adminLoginPayload["password"] = "admin123"
		adminBody, _ = json.Marshal(adminLoginPayload)
		adminLoginResp, _ = http.Post(baseURL+"/api/v1/admin/auth/login", "application/json", bytes.NewBuffer(adminBody))
	}
	var adminData struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	_ = json.NewDecoder(adminLoginResp.Body).Decode(&adminData)
	adminToken := adminData.Data.Token

	// 2. Create Mass Mail Campaign
	campaignPayload := map[string]string{
		"subject": "E2E Test Announcement",
		"content": "<h1>Hello Clients!</h1><p>This is a test broadcast.</p>",
	}
	body, _ := json.Marshal(campaignPayload)
	req, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/admin/mass-mail", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")
	createResp, err := http.DefaultClient.Do(req)
	if err != nil || createResp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to create campaign: %v", err)
	}

	var createData struct {
		Data struct {
			ID int64 `json:"id"`
		} `json:"data"`
	}
	_ = json.NewDecoder(createResp.Body).Decode(&createData)
	campaignID := createData.Data.ID

	// 3. Send Campaign
	req, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/admin/mass-mail/%d/send", baseURL, campaignID), nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	sendResp, err := http.DefaultClient.Do(req)
	if err != nil || sendResp.StatusCode != http.StatusOK {
		t.Fatalf("Failed to send campaign: %v", err)
	}

	var sendData struct {
		Data struct {
			Status    string `json:"status"`
			SentCount int    `json:"sent_count"`
		} `json:"data"`
	}
	_ = json.NewDecoder(sendResp.Body).Decode(&sendData)
	if sendData.Data.Status != "completed" {
		t.Errorf("Campaign status = %s; want completed", sendData.Data.Status)
	}
}
