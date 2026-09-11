package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestE2E_AdminOrdersManagement(t *testing.T) {
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

	// 2. Setup: Ensure an order exists (Register -> Checkout)
	uniqueEmail := fmt.Sprintf("order.admin.tester.%d@example.com", time.Now().UnixNano())
	regPayload := map[string]string{
		"email":      uniqueEmail,
		"password":   "SecurePassword123!",
		"first_name": "Order",
		"last_name":  "Tester",
		"country":    "ID",
		"currency":   "USD",
	}
	regBody, _ := json.Marshal(regPayload)
	regResp, _ := http.Post(baseURL+"/api/v1/guest/auth/register", "application/json", bytes.NewBuffer(regBody))
	var regData struct {
		Data struct {
			Token  string `json:"token"`
			Client struct { ID int64 `json:"id"` } `json:"client"`
		} `json:"data"`
	}
	_ = json.NewDecoder(regResp.Body).Decode(&regData)
	clientToken := regData.Data.Token
	clientID := regData.Data.Client.ID

	// Create a new product with stock to avoid "insufficient stock" error
	uniqueSlug := fmt.Sprintf("test-prod-%d", time.Now().UnixNano())
	productPayload := map[string]interface{}{
		"type":        "hosting",
		"name":        "Test Product E2E",
		"slug":        uniqueSlug,
		"status":      "enabled",
		"setup_type":  "recurring",
		"stock":       100,
	}
	prodBody, _ := json.Marshal(productPayload)
	createProdReq, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/admin/products", bytes.NewBuffer(prodBody))
	createProdReq.Header.Set("Authorization", "Bearer "+adminToken)
	createProdReq.Header.Set("Content-Type", "application/json")
	createProdResp, err := http.DefaultClient.Do(createProdReq)
	if err != nil { t.Fatalf("Failed to create product: %v", err) }

	var createdProdData struct { Data struct { ID int64 `json:"id"` } `json:"data"` }
	_ = json.NewDecoder(createProdResp.Body).Decode(&createdProdData)
	productID := createdProdData.Data.ID
	if productID == 0 { t.Fatal("Failed to get ID of created product") }

	checkoutPayload := map[string]interface{}{
		"client_id": clientID,
		"items": []map[string]interface{}{
			{ "product_id": productID, "title": "Test Product", "period": "1M", "price": 10.0, "quantity": 1 },
		},
	}
	checkBody, _ := json.Marshal(checkoutPayload)
	checkReq, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/guest/cart/checkout", bytes.NewBuffer(checkBody))
	checkReq.Header.Set("Authorization", "Bearer "+clientToken)
	checkReq.Header.Set("Content-Type", "application/json")
	checkResp, _ := http.DefaultClient.Do(checkReq)
	if checkResp.StatusCode != http.StatusCreated { t.Fatalf("Checkout failed with status %d", checkResp.StatusCode) }

	// 3. Test Admin List Orders
	listReq, _ := http.NewRequest(http.MethodGet, baseURL+"/api/v1/admin/orders", nil)
	listReq.Header.Set("Authorization", "Bearer "+adminToken)
	listResp, err := http.DefaultClient.Do(listReq)
	if err != nil { t.Fatalf("Admin list orders failed: %v", err) }
	defer listResp.Body.Close()

	if listResp.StatusCode != http.StatusOK {
		t.Errorf("Admin list orders status = %d; want 200 (possible missing route)", listResp.StatusCode)
	}

	var orderListData struct { Data []struct { ID int64 `json:"id"` } `json:"data"` }
	_ = json.NewDecoder(listResp.Body).Decode(&orderListData)
	if len(orderListData.Data) == 0 { t.Fatal("Expected at least one order in list") }
	orderID := orderListData.Data[0].ID

	// 4. Test Admin Get Order Detail
	detailURL := fmt.Sprintf("%s/api/v1/admin/orders/%d", baseURL, orderID)
	detailReq, _ := http.NewRequest(http.MethodGet, detailURL, nil)
	detailReq.Header.Set("Authorization", "Bearer "+adminToken)
	detailResp, err := http.DefaultClient.Do(detailReq)
	if err != nil { t.Fatalf("Admin get order detail failed: %v", err) }
	defer detailResp.Body.Close()

	if detailResp.StatusCode != http.StatusOK {
		t.Errorf("Admin get order detail status = %d; want 200 (possible missing route/handler)", detailResp.StatusCode)
	}
}
