package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

func getBaseURL() string {
	url := os.Getenv("TEST_API_URL")
	if url == "" {
		url = "http://localhost:8080"
	}
	return url
}

func TestE2E_HealthCheck(t *testing.T) {
	baseURL := getBaseURL()
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		t.Skipf("API server is not running at %s (skipping live e2e test): %v", baseURL, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Health status = %d; want 200", resp.StatusCode)
	}
}

func TestE2E_GuestCurrenciesAndNews(t *testing.T) {
	baseURL := getBaseURL()
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		t.Skipf("API server is not running at %s (skipping live e2e test): %v", baseURL, err)
		return
	}
	_ = resp.Body.Close()

	// 1. Get Guest Currencies
	currResp, err := http.Get(baseURL + "/api/v1/guest/currencies")
	if err != nil {
		t.Fatalf("Guest currencies failed: %v", err)
	}
	defer currResp.Body.Close()
	if currResp.StatusCode != http.StatusOK {
		t.Errorf("Currencies status = %d; want 200", currResp.StatusCode)
	}

	// 2. Get Guest News
	newsResp, err := http.Get(baseURL + "/api/v1/guest/news")
	if err != nil {
		t.Fatalf("Guest news failed: %v", err)
	}
	defer newsResp.Body.Close()
	if newsResp.StatusCode != http.StatusOK {
		t.Errorf("News status = %d; want 200", newsResp.StatusCode)
	}
}

func TestE2E_LiveAuthProfileAndAPIKeysFlow(t *testing.T) {
	baseURL := getBaseURL()
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		t.Skipf("API server is not running at %s (skipping live e2e test): %v", baseURL, err)
		return
	}
	_ = resp.Body.Close()

	uniqueEmail := fmt.Sprintf("test.user.%d@example.com", time.Now().UnixNano())

	// 1. Register
	regPayload := map[string]string{
		"email":      uniqueEmail,
		"password":   "SecurePassword123!",
		"first_name": "Testing",
		"last_name":  "User",
		"country":    "ID",
		"currency":   "IDR",
	}
	body, _ := json.Marshal(regPayload)
	regResp, err := http.Post(baseURL+"/api/v1/guest/auth/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	defer regResp.Body.Close()

	if regResp.StatusCode != http.StatusCreated {
		t.Fatalf("Register status = %d; want 201", regResp.StatusCode)
	}

	var regData struct {
		Success bool `json:"success"`
		Data    struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	_ = json.NewDecoder(regResp.Body).Decode(&regData)
	token := regData.Data.Token
	if token == "" {
		t.Fatal("Expected JWT token from register response")
	}

	// 2. Login
	loginPayload := map[string]string{
		"email":    uniqueEmail,
		"password": "SecurePassword123!",
	}
	body, _ = json.Marshal(loginPayload)
	loginResp, err := http.Post(baseURL+"/api/v1/guest/auth/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	defer loginResp.Body.Close()

	if loginResp.StatusCode != http.StatusOK {
		t.Fatalf("Login status = %d; want 200", loginResp.StatusCode)
	}

	// 3. Get Profile
	req, _ := http.NewRequest(http.MethodGet, baseURL+"/api/v1/client/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	profileResp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Profile request failed: %v", err)
	}
	defer profileResp.Body.Close()

	if profileResp.StatusCode != http.StatusOK {
		t.Fatalf("Profile status = %d; want 200", profileResp.StatusCode)
	}

	// 4. Generate API Key
	keyReqBody, _ := json.Marshal(map[string]interface{}{
		"name":        "Test Integration Key",
		"expire_days": 30,
	})
	keyReq, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/client/api-keys", bytes.NewBuffer(keyReqBody))
	keyReq.Header.Set("Authorization", "Bearer "+token)
	keyReq.Header.Set("Content-Type", "application/json")
	keyResp, err := http.DefaultClient.Do(keyReq)
	if err != nil {
		t.Fatalf("API Key generation failed: %v", err)
	}
	defer keyResp.Body.Close()

	if keyResp.StatusCode != http.StatusCreated {
		t.Fatalf("API Key status = %d; want 201", keyResp.StatusCode)
	}

	// 5. List API Keys
	listReq, _ := http.NewRequest(http.MethodGet, baseURL+"/api/v1/client/api-keys", nil)
	listReq.Header.Set("Authorization", "Bearer "+token)
	listResp, err := http.DefaultClient.Do(listReq)
	if err != nil {
		t.Fatalf("List API Keys failed: %v", err)
	}
	defer listResp.Body.Close()

	if listResp.StatusCode != http.StatusOK {
		t.Fatalf("List API Keys status = %d; want 200", listResp.StatusCode)
	}
}

func TestE2E_LiveSupportTicketFlow(t *testing.T) {
	baseURL := getBaseURL()
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		t.Skipf("API server is not running at %s (skipping live e2e test): %v", baseURL, err)
		return
	}
	_ = resp.Body.Close()

	uniqueEmail := fmt.Sprintf("ticket.user.%d@example.com", time.Now().UnixNano())

	// 1. Register
	regPayload := map[string]string{
		"email":      uniqueEmail,
		"password":   "SecurePassword123!",
		"first_name": "Support",
		"last_name":  "Tester",
		"country":    "ID",
		"currency":   "USD",
	}
	body, _ := json.Marshal(regPayload)
	regResp, err := http.Post(baseURL+"/api/v1/guest/auth/register", "application/json", bytes.NewBuffer(body))
	if err != nil || regResp.StatusCode != http.StatusCreated {
		t.Fatalf("Register failed: %v", err)
	}
	defer regResp.Body.Close()

	var regData struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	_ = json.NewDecoder(regResp.Body).Decode(&regData)
	token := regData.Data.Token

	// 2. Open Support Ticket
	ticketBody, _ := json.Marshal(map[string]interface{}{
		"subject":  "VPS High CPU Utilization Inquiry",
		"message":  "Hello, we notice intermittent CPU spikes on our VPS instance.",
		"priority": "high",
	})
	req, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/client/support/tickets", bytes.NewBuffer(ticketBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	tResp, err := http.DefaultClient.Do(req)
	if err != nil || tResp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to open ticket: %v", err)
	}
	defer tResp.Body.Close()

	// 3. List Tickets
	req, _ = http.NewRequest(http.MethodGet, baseURL+"/api/v1/client/support/tickets", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	listResp, err := http.DefaultClient.Do(req)
	if err != nil || listResp.StatusCode != http.StatusOK {
		t.Fatalf("Failed to list tickets: %v", err)
	}
	defer listResp.Body.Close()

	var listData struct {
		Data []struct {
			ID     int64  `json:"id"`
			Status string `json:"status"`
		} `json:"data"`
	}
	_ = json.NewDecoder(listResp.Body).Decode(&listData)
	if len(listData.Data) == 0 {
		t.Fatal("Expected at least one ticket in list")
	}
	ticketID := listData.Data[0].ID

	// 4. Client Reply to Ticket
	replyPayload := map[string]string{
		"message": "Update: The issue is still occurring even after reboot.",
	}
	body, _ = json.Marshal(replyPayload)
	req, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/client/support/tickets/%d/reply", baseURL, ticketID), bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	replyResp, err := http.DefaultClient.Do(req)
	if err != nil || replyResp.StatusCode != http.StatusCreated {
		t.Fatalf("Client failed to reply to ticket: %v", err)
	}
	defer replyResp.Body.Close()

	// 5. Verify Ticket Status is now "awaiting_staff"
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/client/support/tickets/%d", baseURL, ticketID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	detailResp, _ := http.DefaultClient.Do(req)
	var detailData struct {
		Data struct {
			Ticket struct {
				Status string `json:"status"`
			} `json:"ticket"`
		} `json:"data"`
	}
	_ = json.NewDecoder(detailResp.Body).Decode(&detailData)
	if detailData.Data.Ticket.Status != "awaiting_staff" {
		t.Errorf("Expected ticket status awaiting_staff, got %s", detailData.Data.Ticket.Status)
	}
}

func TestE2E_FullCheckoutAndPaymentFlow(t *testing.T) {
	baseURL := getBaseURL()
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		t.Skipf("API server is not running at %s (skipping live e2e test): %v", baseURL, err)
		return
	}
	_ = resp.Body.Close()

	// 1. Admin Login & Create Product (to avoid FK violation)
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
		body, _ := io.ReadAll(adminLoginResp.Body)
		t.Fatalf("Admin login failed with status %d: %s", adminLoginResp.StatusCode, string(body))
	}
	defer adminLoginResp.Body.Close()
	var adminData struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	_ = json.NewDecoder(adminLoginResp.Body).Decode(&adminData)
	adminToken := adminData.Data.Token

	productPayload := map[string]interface{}{
		"type":        "hosting",
		"name":        "Cloud VPS Starter E2E",
		"slug":        fmt.Sprintf("cloud-vps-starter-e2e-%d", time.Now().UnixNano()),
		"description": "Auto-generated for E2E testing",
		"status":      "enabled",
		"setup_type":  "recurring",
		"stock":       100,
	}
	prodBody, _ := json.Marshal(productPayload)
	prodReq, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/admin/products", bytes.NewBuffer(prodBody))
	prodReq.Header.Set("Authorization", "Bearer "+adminToken)
	prodReq.Header.Set("Content-Type", "application/json")
	prodResp, err := http.DefaultClient.Do(prodReq)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}
	// We don't check for 201 strictly here because it might already exist or ID 101 might be taken
	// But we try to get a valid product ID.
	var createdProduct struct {
		Data struct {
			ID int64 `json:"id"`
		} `json:"data"`
	}
	if prodResp.StatusCode == http.StatusCreated {
		_ = json.NewDecoder(prodResp.Body).Decode(&createdProduct)
	} else {
		// Fallback: list products and take the first one
		listReq, _ := http.NewRequest(http.MethodGet, baseURL+"/api/v1/admin/products", nil)
		listReq.Header.Set("Authorization", "Bearer "+adminToken)
		listResp, _ := http.DefaultClient.Do(listReq)
		var listData struct {
			Data []struct {
				ID int64 `json:"id"`
			} `json:"data"`
		}
		_ = json.NewDecoder(listResp.Body).Decode(&listData)
		if len(listData.Data) > 0 {
			createdProduct.Data.ID = listData.Data[0].ID
		} else {
			t.Fatal("No products available for checkout test")
		}
	}
	productID := createdProduct.Data.ID

	uniqueEmail := fmt.Sprintf("checkout.user.%d@example.com", time.Now().UnixNano())

	// 2. Register Client
	regPayload := map[string]string{
		"email":      uniqueEmail,
		"password":   "SecurePassword123!",
		"first_name": "E2E",
		"last_name":  "Buyer",
		"country":    "ID",
		"currency":   "USD",
	}
	body, _ := json.Marshal(regPayload)
	regResp, err := http.Post(baseURL+"/api/v1/guest/auth/register", "application/json", bytes.NewBuffer(body))
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

	// 3. Checkout (Creating Order & Invoice)
	checkoutPayload := map[string]interface{}{
		"client_id": clientID,
		"items": []map[string]interface{}{
			{
				"product_id": productID,
				"title":      "Cloud VPS Starter",
				"period":     "1M",
				"price":      9.99,
				"quantity":   1,
			},
		},
	}
	body, _ = json.Marshal(checkoutPayload)
	req, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/guest/cart/checkout", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	checkResp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Checkout failed: %v", err)
	}
	defer checkResp.Body.Close()
	if checkResp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(checkResp.Body)
		t.Fatalf("Checkout failed with status %d: %s", checkResp.StatusCode, string(body))
	}

	var checkoutData struct {
		Data struct {
			Invoice struct {
				ID    int64   `json:"id"`
				Total float64 `json:"total"`
			} `json:"invoice"`
		} `json:"data"`
	}
	_ = json.NewDecoder(checkResp.Body).Decode(&checkoutData)
	invID := checkoutData.Data.Invoice.ID

	// 4. Simulate Webhook Payment (Custom Gateway)
	webhookURL := fmt.Sprintf("%s/api/v1/guest/webhook/custom?invoice_id=%d&txn_id=E2E-TEST-TXN-%d&currency=USD", baseURL, invID, time.Now().UnixNano())
	webResp, err := http.Post(webhookURL, "application/json", nil)
	if err != nil {
		t.Fatalf("Webhook simulation failed with err: %v", err)
	}
	defer webResp.Body.Close()
	if webResp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(webResp.Body)
		t.Fatalf("Webhook simulation failed with status %d: %s", webResp.StatusCode, string(b))
	}

	// 5. Verify Invoice is Paid
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/client/invoices/%d", baseURL, invID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	invResp, err := http.DefaultClient.Do(req)
	if err != nil || invResp.StatusCode != http.StatusOK {
		t.Fatalf("Failed to fetch invoice: %v", err)
	}
	var invData struct {
		Data struct {
			Status string `json:"status"`
		} `json:"data"`
	}
	_ = json.NewDecoder(invResp.Body).Decode(&invData)
	if invData.Data.Status != "paid" {
		t.Errorf("Invoice status = %s; want paid", invData.Data.Status)
	}
}

func TestE2E_LocalesGeoIPAndDocs(t *testing.T) {
	baseURL := getBaseURL()
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		t.Skipf("API server is not running at %s (skipping live e2e test): %v", baseURL, err)
		return
	}
	_ = resp.Body.Close()

	// 1. Check Locales
	locResp, err := http.Get(baseURL + "/api/v1/guest/locales")
	if err != nil {
		t.Fatalf("Locales request failed: %v", err)
	}
	defer locResp.Body.Close()
	if locResp.StatusCode != http.StatusOK {
		t.Errorf("Locales status = %d; want 200", locResp.StatusCode)
	}

	// 2. Check i18n switching via header
	req, _ := http.NewRequest(http.MethodGet, baseURL+"/api/v1/guest/currencies", nil)
	req.Header.Set("Accept-Language", "id-ID")
	idResp, err := http.DefaultClient.Do(req)
	if err != nil || idResp.StatusCode != http.StatusOK {
		t.Errorf("Indonesian locale request failed")
	}

	// 2b. Check translated error message
	loginPayload := map[string]string{"email": "wrong@email.com", "password": "wrong"}
	body, _ := json.Marshal(loginPayload)
	req, _ = http.NewRequest(http.MethodPost, baseURL+"/api/v1/guest/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Accept-Language", "id-ID")
	req.Header.Set("Content-Type", "application/json")
	loginResp, _ := http.DefaultClient.Do(req)
	var loginErr struct {
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	_ = json.NewDecoder(loginResp.Body).Decode(&loginErr)
	// From i18n.go: "id_ID" -> "invalid_credentials": "Email atau kata sandi tidak valid"
	if loginErr.Error.Message != "Email atau kata sandi tidak valid" {
		t.Errorf("Expected Indonesian error message, got: %s", loginErr.Error.Message)
	}

	// 3. Check GeoIP lookup
	geoResp, err := http.Get(baseURL + "/api/v1/guest/system/geoip?ip=8.8.8.8")
	if err != nil {
		t.Fatalf("GeoIP request failed: %v", err)
	}
	defer geoResp.Body.Close()
	if geoResp.StatusCode != http.StatusOK {
		t.Errorf("GeoIP status = %d; want 200", geoResp.StatusCode)
	}

	// 4. Check OpenAPI specification endpoint
	specResp, err := http.Get(baseURL + "/openapi.json")
	if err != nil {
		t.Fatalf("OpenAPI spec failed: %v", err)
	}
	defer specResp.Body.Close()
	if specResp.StatusCode != http.StatusOK {
		t.Errorf("OpenAPI spec status = %d; want 200", specResp.StatusCode)
	}

	// 5. Check Scalar Docs endpoint
	docsResp, err := http.Get(baseURL + "/docs")
	if err != nil {
		t.Fatalf("Scalar docs failed: %v", err)
	}
	defer docsResp.Body.Close()
	if docsResp.StatusCode != http.StatusOK {
		t.Errorf("Docs status = %d; want 200", docsResp.StatusCode)
	}
}

