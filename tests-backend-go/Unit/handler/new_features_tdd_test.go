package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestTDD_NewFeaturesBackendEndpoints(t *testing.T) {
	ts, promoRepo, staffRepo := setupTestServer()
	defer ts.Close()
	ctx := context.Background()
	setupTestAdminsAndPromos(ctx, promoRepo, staffRepo)

	// 1. Register Client for Testing
	regBody, _ := json.Marshal(map[string]interface{}{
		"email": "feature.tester@example.com", "password": "OldPassword123!",
		"first_name": "Feature", "last_name": "Tester", "country": "ID", "currency": "USD",
	})
	regResp, err := http.Post(ts.URL+"/api/v1/guest/auth/register", "application/json", bytes.NewBuffer(regBody))
	if err != nil || regResp.StatusCode != http.StatusCreated {
		t.Fatalf("Register failed: %v", err)
	}

	var regData struct {
		Data struct {
			Token string `json:"token"`
			Client struct{ ID int64 `json:"id"` } `json:"client"`
		} `json:"data"`
	}
	_ = json.NewDecoder(regResp.Body).Decode(&regData)
	clientToken, clientID := regData.Data.Token, regData.Data.Client.ID

	// 2. TDD: Client Change Password
	t.Run("Change Password - Wrong Current Password", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"current_password": "WrongPassword!",
			"new_password":     "NewPassword123!",
		})
		req, _ := http.NewRequest("POST", ts.URL+"/api/v1/client/profile/change-password", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+clientToken)
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("Expected 401 Unauthorized for wrong password, got %d", resp.StatusCode)
		}
	})

	t.Run("Change Password - Success", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"current_password": "OldPassword123!",
			"new_password":     "NewSecurePassword456!",
		})
		req, _ := http.NewRequest("POST", ts.URL+"/api/v1/client/profile/change-password", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+clientToken)
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200 OK for valid password change, got %d", resp.StatusCode)
		}
	})

	// 3. TDD: Client Deposit Funds
	t.Run("Deposit Funds - Create Top-Up Invoice", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"amount":   75.50,
			"currency": "USD",
		})
		req, _ := http.NewRequest("POST", ts.URL+"/api/v1/client/funds/deposit", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+clientToken)
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != http.StatusCreated {
			t.Fatalf("Expected 201 Created for deposit funds, got %d", resp.StatusCode)
		}

		var depData struct {
			Data struct {
				InvoiceID int64   `json:"invoice_id"`
				Total     float64 `json:"total"`
				Status    string  `json:"status"`
			} `json:"data"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&depData)
		if depData.Data.InvoiceID == 0 {
			t.Errorf("Expected deposit invoice ID > 0, got %v", depData.Data.InvoiceID)
		}
	})

	// 4. TDD: Admin Login & Manual Invoice Creation
	loginBody, _ := json.Marshal(map[string]string{"email": "admin@fossbilling.org", "password": "SuperSecretAdmin123!"})
	loginResp, _ := http.Post(ts.URL+"/api/v1/admin/auth/login", "application/json", bytes.NewBuffer(loginBody))
	var adminData struct {
		Data struct { Token string `json:"token"` } `json:"data"`
	}
	_ = json.NewDecoder(loginResp.Body).Decode(&adminData)
	adminToken := adminData.Data.Token

	t.Run("Admin - Create Custom Invoice & List Invoices", func(t *testing.T) {
		invBody, _ := json.Marshal(map[string]interface{}{
			"client_id": clientID,
			"currency":  "USD",
			"due_days":  7,
			"items": []map[string]interface{}{
				{"title": "Custom Cloud VPS - 8 vCPU", "price": 49.99, "quantity": 1, "taxable": false},
				{"title": "Dedicated IPv4 Addon", "price": 4.00, "quantity": 2, "taxable": false},
			},
		})
		req, _ := http.NewRequest("POST", ts.URL+"/api/v1/admin/invoices", bytes.NewBuffer(invBody))
		req.Header.Set("Authorization", "Bearer "+adminToken)
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != http.StatusCreated {
			t.Fatalf("Expected 201 Created for admin create invoice, got %d", resp.StatusCode)
		}

		// List Admin Invoices
		listReq, _ := http.NewRequest("GET", ts.URL+"/api/v1/admin/invoices", nil)
		listReq.Header.Set("Authorization", "Bearer "+adminToken)
		listResp, err := http.DefaultClient.Do(listReq)
		if err != nil || listResp.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200 OK for admin list invoices, got %d", listResp.StatusCode)
		}

		var listData struct {
			Data []struct {
				ID int64 `json:"id"`
				Total int64 `json:"total"`
			} `json:"data"`
		}
		_ = json.NewDecoder(listResp.Body).Decode(&listData)
		if len(listData.Data) == 0 {
			t.Errorf("Expected at least 1 invoice in admin list, got 0")
		}
	})
}
