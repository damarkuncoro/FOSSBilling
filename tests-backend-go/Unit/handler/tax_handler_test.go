package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestTDD_AdminTaxEndpointsAndAliases(t *testing.T) {
	ts, promoRepo, staffRepo, productRepo := setupTestServer()
	defer ts.Close()
	ctx := context.Background()
	setupTestAdminsAndPromos(ctx, promoRepo, staffRepo, productRepo)

	// 1. Unauthenticated requests should return 401 Unauthorized
	t.Run("Unauthenticated access to tax endpoints", func(t *testing.T) {
		reqTaxes, _ := http.NewRequest(http.MethodGet, ts.URL+"/api/v1/admin/taxes", nil)
		respTaxes, err := http.DefaultClient.Do(reqTaxes)
		if err != nil || respTaxes.StatusCode != http.StatusUnauthorized {
			t.Fatalf("Expected 401 Unauthorized for unauth /api/v1/admin/taxes, got %v", respTaxes.StatusCode)
		}

		reqTaxRules, _ := http.NewRequest(http.MethodGet, ts.URL+"/api/v1/admin/tax-rules", nil)
		respTaxRules, err := http.DefaultClient.Do(reqTaxRules)
		if err != nil || respTaxRules.StatusCode != http.StatusUnauthorized {
			t.Fatalf("Expected 401 Unauthorized for unauth /api/v1/admin/tax-rules, got %v", respTaxRules.StatusCode)
		}
	})

	// 2. Admin Login
	loginBody, _ := json.Marshal(map[string]string{
		"email":    "admin@fossbilling.org",
		"password": "SuperSecretAdmin123!",
	})
	loginResp, err := http.Post(ts.URL+"/api/v1/admin/auth/login", "application/json", bytes.NewBuffer(loginBody))
	if err != nil || loginResp.StatusCode != http.StatusOK {
		t.Fatalf("Admin login failed: %v", err)
	}
	var adminData struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	_ = json.NewDecoder(loginResp.Body).Decode(&adminData)
	adminToken := adminData.Data.Token

	// 3. Create Tax Rule via POST /api/v1/admin/tax-rules
	t.Run("Create Tax Rule via /admin/tax-rules", func(t *testing.T) {
		taxPayload := map[string]interface{}{
			"name":        "VAT Standard",
			"country":     "GB",
			"state":       "",
			"rate":        20.0,
			"tax_exempt":  false,
		}
		body, _ := json.Marshal(taxPayload)
		req, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/v1/admin/tax-rules", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+adminToken)
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != http.StatusCreated {
			t.Fatalf("Expected 201 Created for creating tax rule, got %v", resp.StatusCode)
		}

		var created struct {
			Data struct {
				ID   int64   `json:"id"`
				Name string  `json:"name"`
				Rate float64 `json:"rate"`
			} `json:"data"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&created)
		if created.Data.Name != "VAT Standard" || created.Data.Rate != 20.0 {
			t.Errorf("Unexpected created tax rule data: %+v", created.Data)
		}
	})

	// 4. Create Tax Rule via POST /api/v1/admin/taxes (Alias)
	t.Run("Create Tax Rule via alias /admin/taxes", func(t *testing.T) {
		taxPayload := map[string]interface{}{
			"name":        "PPN Indonesia",
			"country":     "ID",
			"state":       "",
			"rate":        11.0,
			"tax_exempt":  false,
		}
		body, _ := json.Marshal(taxPayload)
		req, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/v1/admin/taxes", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+adminToken)
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != http.StatusCreated {
			t.Fatalf("Expected 201 Created for creating tax rule via alias, got %v", resp.StatusCode)
		}
	})

	// 5. List Tax Rules via GET /api/v1/admin/tax-rules and GET /api/v1/admin/taxes
	t.Run("List Tax Rules via both endpoints", func(t *testing.T) {
		// Test /api/v1/admin/tax-rules
		reqRules, _ := http.NewRequest(http.MethodGet, ts.URL+"/api/v1/admin/tax-rules", nil)
		reqRules.Header.Set("Authorization", "Bearer "+adminToken)
		respRules, err := http.DefaultClient.Do(reqRules)
		if err != nil || respRules.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200 OK for /admin/tax-rules, got %v", respRules.StatusCode)
		}

		var rulesData struct {
			Data []struct {
				ID   int64  `json:"id"`
				Name string `json:"name"`
			} `json:"data"`
		}
		_ = json.NewDecoder(respRules.Body).Decode(&rulesData)
		if len(rulesData.Data) < 2 {
			t.Errorf("Expected at least 2 tax rules, got %d", len(rulesData.Data))
		}

		// Test /api/v1/admin/taxes
		reqTaxes, _ := http.NewRequest(http.MethodGet, ts.URL+"/api/v1/admin/taxes", nil)
		reqTaxes.Header.Set("Authorization", "Bearer "+adminToken)
		respTaxes, err := http.DefaultClient.Do(reqTaxes)
		if err != nil || respTaxes.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200 OK for /admin/taxes, got %v", respTaxes.StatusCode)
		}

		var taxesData struct {
			Data []struct {
				ID   int64  `json:"id"`
				Name string `json:"name"`
			} `json:"data"`
		}
		_ = json.NewDecoder(respTaxes.Body).Decode(&taxesData)
		if len(taxesData.Data) != len(rulesData.Data) {
			t.Errorf("Expected taxes and tax-rules to return identical count, got %d vs %d", len(taxesData.Data), len(rulesData.Data))
		}
	})

	// 6. Update Tax Rule via PUT /api/v1/admin/tax-rules/{id} and PUT /api/v1/admin/taxes/{id}
	t.Run("Update Tax Rule", func(t *testing.T) {
		updatePayload := map[string]interface{}{
			"name": "VAT Standard Updated",
			"rate": 21.0,
		}
		body, _ := json.Marshal(updatePayload)
		req, _ := http.NewRequest(http.MethodPut, ts.URL+"/api/v1/admin/tax-rules/1", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+adminToken)
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200 OK for update tax rule, got %v", resp.StatusCode)
		}
	})

	// 7. Delete Tax Rule via DELETE /api/v1/admin/tax-rules/{id} or /api/v1/admin/taxes/{id}
	t.Run("Delete Tax Rule", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/api/v1/admin/taxes/1", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200 OK for delete tax rule, got %v", resp.StatusCode)
		}
	})
}
