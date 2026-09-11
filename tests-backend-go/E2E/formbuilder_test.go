package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestE2E_FormbuilderAndCheckoutFlow(t *testing.T) {
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
		"password": "admin123",
	}
	adminBody, _ := json.Marshal(adminLoginPayload)
	adminLoginResp, err := http.Post(baseURL+"/api/v1/admin/auth/login", "application/json", bytes.NewBuffer(adminBody))
	if err != nil {
		t.Fatalf("Admin login request failed: %v", err)
	}
	if adminLoginResp.StatusCode != http.StatusOK {
		adminLoginPayload["password"] = "SuperSecretAdmin123!"
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

	// 2. Create Custom Form
	formPayload := map[string]interface{}{
		"name": "VPS Configuration Form E2E",
		"style": map[string]interface{}{
			"type":       "vertical",
			"show_title": true,
		},
	}
	formBody, _ := json.Marshal(formPayload)
	formReq, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/admin/forms", bytes.NewBuffer(formBody))
	formReq.Header.Set("Authorization", "Bearer "+adminToken)
	formReq.Header.Set("Content-Type", "application/json")
	formResp, err := http.DefaultClient.Do(formReq)
	if err != nil || formResp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to create form: %v", err)
	}
	var formData struct {
		Data struct {
			ID int64 `json:"id"`
		} `json:"data"`
	}
	_ = json.NewDecoder(formResp.Body).Decode(&formData)
	formID := formData.Data.ID

	// 3. Add Field to Form (Hostname)
	fieldPayload := map[string]interface{}{
		"name":     "hostname",
		"label":    "Server Hostname",
		"type":     "text",
		"required": true,
	}
	fieldBody, _ := json.Marshal(fieldPayload)
	fieldReq, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/admin/forms/%d/fields", baseURL, formID), bytes.NewBuffer(fieldBody))
	fieldReq.Header.Set("Authorization", "Bearer "+adminToken)
	fieldReq.Header.Set("Content-Type", "application/json")
	fieldResp, err := http.DefaultClient.Do(fieldReq)
	if err != nil || fieldResp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to add field to form: %v", err)
	}

	// 4. Create Product and Link Form
	productPayload := map[string]interface{}{
		"type":        "hosting",
		"name":        "Cloud VPS with Form E2E",
		"slug":        fmt.Sprintf("vps-form-e2e-%d", time.Now().Unix()),
		"status":        "enabled",
		"form_id":       formID,
		"setup_type":    "recurring",
		"price_monthly": 15.00,
		"stock":         100,
	}
	prodBody, _ := json.Marshal(productPayload)
	prodReq, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/admin/products", bytes.NewBuffer(prodBody))
	prodReq.Header.Set("Authorization", "Bearer "+adminToken)
	prodReq.Header.Set("Content-Type", "application/json")
	prodResp, err := http.DefaultClient.Do(prodReq)
	if err != nil || prodResp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to create product: %v", err)
	}
	var createdProduct struct {
		Data struct {
			ID int64 `json:"id"`
		} `json:"data"`
	}
	_ = json.NewDecoder(prodResp.Body).Decode(&createdProduct)
	productID := createdProduct.Data.ID

	// 5. Register Client
	uniqueEmail := fmt.Sprintf("form.user.%d@example.com", time.Now().UnixNano())
	regPayload := map[string]string{
		"email":      uniqueEmail,
		"password":   "SecurePassword123!",
		"first_name": "Form",
		"last_name":  "Tester",
		"country":    "ID",
		"currency":   "USD",
	}
	regBody, _ := json.Marshal(regPayload)
	regResp, err := http.Post(baseURL+"/api/v1/guest/auth/register", "application/json", bytes.NewBuffer(regBody))
	if err != nil || regResp.StatusCode != http.StatusCreated {
		t.Fatalf("Register failed: %v", err)
	}
	var regData struct {
		Data struct {
			Token  string `json:"token"`
			Client struct {
				ID int64 `json:"id"`
			} `json:"client"`
		} `json:"data"`
	}
	_ = json.NewDecoder(regResp.Body).Decode(&regData)
	token := regData.Data.Token
	clientID := regData.Data.Client.ID

	// 6. Checkout with Custom Config
	configPayload := map[string]string{
		"hostname": "vps.myserver.com",
	}

	checkoutPayload := map[string]interface{}{
		"client_id": clientID,
		"items": []map[string]interface{}{
			{
				"product_id": productID,
				"title":      "Cloud VPS Starter",
				"period":     "1M",
				"price":      15.00,
				"quantity":   1,
				"config":     configPayload,
			},
		},
	}
	checkoutBody, _ := json.Marshal(checkoutPayload)
	req, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/guest/cart/checkout", bytes.NewBuffer(checkoutBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	checkResp, err := http.DefaultClient.Do(req)
	if err != nil || checkResp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(checkResp.Body)
		t.Fatalf("Checkout failed with status %d: %s", checkResp.StatusCode, string(body))
	}

	var checkoutData struct {
		Data struct {
			Orders []struct {
				ID int64 `json:"id"`
			} `json:"orders"`
		} `json:"data"`
	}
	_ = json.NewDecoder(checkResp.Body).Decode(&checkoutData)
	orderID := checkoutData.Data.Orders[0].ID

	// 7. Verify Order Configuration
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/client/orders/%d", baseURL, orderID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	orderResp, _ := http.DefaultClient.Do(req)
	var orderData struct {
		Data struct {
			Config json.RawMessage `json:"config"`
		} `json:"data"`
	}
	_ = json.NewDecoder(orderResp.Body).Decode(&orderData)

	var finalConfig map[string]interface{}
	_ = json.Unmarshal(orderData.Data.Config, &finalConfig)

	if finalConfig["hostname"] != "vps.myserver.com" {
		t.Errorf("Order hostname configuration mismatch, got: %v", finalConfig["hostname"])
	}
}
