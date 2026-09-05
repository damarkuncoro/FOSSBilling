package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
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
