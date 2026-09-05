package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/auth"
)

func setupTestAdminsAndPromos(ctx context.Context, promoRepo domain.PromoRepository, staffRepo domain.StaffRepository) {
	_ = promoRepo.Create(ctx, &domain.Promo{
		Code:        "MERDEKA20",
		Description: "20% discount",
		Type:        domain.PromoTypePercentage,
		Value:       200000,
		Active:      true,
	})

	group := &domain.AdminGroup{
		Name: "Administrators",
		Permissions: map[string][]string{
			"clients": {"*"}, "orders": {"*"}, "support": {"*"}, "system": {"*"},
		},
	}
	_ = staffRepo.CreateGroup(ctx, group)
	passHash, _ := auth.HashPassword("SuperSecretAdmin123!")
	_ = staffRepo.Create(ctx, &domain.Staff{
		GroupID:      group.ID,
		Email:        "admin@fossbilling.org",
		PasswordHash: passHash,
		Name:         "Super Admin",
		Role:         domain.StaffRoleSuperAdmin,
		Status:       "active",
	})
}

func TestHTTP_RegistrationAndCheckoutLifecycle(t *testing.T) {
	ts, promoRepo, staffRepo := setupTestServer()
	defer ts.Close()
	ctx := context.Background()
	setupTestAdminsAndPromos(ctx, promoRepo, staffRepo)

	res, err := http.Get(ts.URL + "/health")
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("Health check failed: %v, status: %d", err, res.StatusCode)
	}

	regBody, _ := json.Marshal(map[string]interface{}{
		"email": "ahmad.dhani@example.com", "password": "Password123!",
		"first_name": "Ahmad", "last_name": "Dhani", "country": "ID", "currency": "USD",
	})
	regResp, err := http.Post(ts.URL+"/api/v1/guest/auth/register", "application/json", bytes.NewBuffer(regBody))
	if err != nil || regResp.StatusCode != http.StatusCreated {
		t.Fatalf("Register failed: %v, status: %d", err, regResp.StatusCode)
	}

	var regData struct {
		Data struct {
			Token  string `json:"token"`
			Client struct{ ID int64 `json:"id"` } `json:"client"`
		} `json:"data"`
	}
	_ = json.NewDecoder(regResp.Body).Decode(&regData)
	clientToken, clientID := regData.Data.Token, regData.Data.Client.ID

	cartBody, _ := json.Marshal(map[string]interface{}{
		"client_id": clientID, "promo_code": "MERDEKA20",
		"items": []map[string]interface{}{
			{"product_id": 1, "title": "Cloud Hosting Pro", "period": "1M", "price": 1000000, "quantity": 1},
		},
	})
	checkoutResp, err := http.Post(ts.URL+"/api/v1/guest/cart/checkout", "application/json", bytes.NewBuffer(cartBody))
	if err != nil || checkoutResp.StatusCode != http.StatusCreated {
		t.Fatalf("Checkout failed: %v, status: %d", err, checkoutResp.StatusCode)
	}

	var checkoutData struct {
		Data struct {
			Orders  []domain.Order `json:"orders"`
			Invoice domain.Invoice `json:"invoice"`
		} `json:"data"`
	}
	_ = json.NewDecoder(checkoutResp.Body).Decode(&checkoutData)
	invoiceID := checkoutData.Data.Invoice.ID
	orderID := checkoutData.Data.Orders[0].ID

	req, _ := http.NewRequest("GET", ts.URL+"/api/v1/client/orders", nil)
	req.Header.Set("Authorization", "Bearer "+clientToken)
	ordersResp, err := http.DefaultClient.Do(req)
	if err != nil || ordersResp.StatusCode != http.StatusOK {
		t.Fatalf("Get client orders failed: %v, status: %d", err, ordersResp.StatusCode)
	}

	webhookBody, _ := json.Marshal(map[string]interface{}{
		"gateway_id": "midtrans",
		"txn_id":     fmt.Sprintf("TRX-%d", time.Now().UnixNano()),
		"invoice_id": invoiceID,
		"amount":     checkoutData.Data.Invoice.Total,
		"currency":   "USD",
	})
	webhookResp, err := http.Post(ts.URL+"/api/v1/guest/gateways/midtrans/webhook", "application/json", bytes.NewBuffer(webhookBody))
	if err != nil || webhookResp.StatusCode != http.StatusOK {
		t.Fatalf("Webhook failed: %v, status: %d", err, webhookResp.StatusCode)
	}

	req, _ = http.NewRequest("GET", fmt.Sprintf("%s/api/v1/client/orders/%d", ts.URL, orderID), nil)
	req.Header.Set("Authorization", "Bearer "+clientToken)
	getOrdResp, _ := http.DefaultClient.Do(req)
	var ordData struct{ Data domain.Order `json:"data"` }
	_ = json.NewDecoder(getOrdResp.Body).Decode(&ordData)
	if ordData.Data.Status != domain.OrderStatusActive {
		t.Errorf("Expected order status active after payment, got: %s", ordData.Data.Status)
	}
}
