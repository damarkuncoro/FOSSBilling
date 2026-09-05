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

func TestTDD_AllInteractiveButtonActions(t *testing.T) {
	ts, promoRepo, staffRepo := setupTestServer()
	defer ts.Close()
	ctx := context.Background()
	setupTestAdminsAndPromos(ctx, promoRepo, staffRepo)

	// 1. Client Register Action
	regBody, _ := json.Marshal(map[string]interface{}{
		"email": "tester.buttons@example.com", "password": "Password123!",
		"first_name": "Button", "last_name": "Tester", "country": "ID", "currency": "USD",
	})
	regResp, err := http.Post(ts.URL+"/api/v1/guest/auth/register", "application/json", bytes.NewBuffer(regBody))
	if err != nil || regResp.StatusCode != http.StatusCreated {
		t.Fatalf("Register button failed: %v", err)
	}

	var regData struct {
		Data struct {
			Token string `json:"token"`
			Client struct{ ID int64 `json:"id"` } `json:"client"`
		} `json:"data"`
	}
	_ = json.NewDecoder(regResp.Body).Decode(&regData)
	clientToken, clientID := regData.Data.Token, regData.Data.Client.ID

	// 2. Cart "Apply Promo Code" Button
	promoBody, _ := json.Marshal(map[string]interface{}{
		"client_id": clientID, "promo_code": "MERDEKA20",
		"items": []map[string]interface{}{
			{"product_id": 1, "title": "NVMe Cloud", "period": "1M", "price": 1000000, "quantity": 1},
		},
	})
	calcResp, err := http.Post(ts.URL+"/api/v1/guest/cart/calculate", "application/json", bytes.NewBuffer(promoBody))
	if err != nil || calcResp.StatusCode != http.StatusOK {
		t.Fatalf("Apply promo button failed: %v", err)
	}
	var calcData struct {
		Data struct {
			Subtotal int64 `json:"subtotal"`
			Discount int64 `json:"discount"`
			Total    int64 `json:"total"`
		} `json:"data"`
	}
	_ = json.NewDecoder(calcResp.Body).Decode(&calcData)
	if calcData.Data.Discount != 200000 {
		t.Errorf("Expected 20%% discount (200000), got %d", calcData.Data.Discount)
	}

	// 3. Cart "Checkout" Button
	checkoutResp, err := http.Post(ts.URL+"/api/v1/guest/cart/checkout", "application/json", bytes.NewBuffer(promoBody))
	if err != nil || checkoutResp.StatusCode != http.StatusCreated {
		t.Fatalf("Checkout button failed: %v", err)
	}
	var coData struct {
		Data struct {
			Orders  []domain.Order `json:"orders"`
			Invoice domain.Invoice `json:"invoice"`
		} `json:"data"`
	}
	_ = json.NewDecoder(checkoutResp.Body).Decode(&coData)
	if len(coData.Data.Orders) == 0 {
		t.Fatal("Expected orders created")
	}
	orderID := coData.Data.Orders[0].ID
	invoiceID := coData.Data.Invoice.ID

	// 4. Invoice "Pay With Balance" Button (Failed initially without balance)
	payReq, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/client/invoices/%d/pay-balance", ts.URL, invoiceID), nil)
	payReq.Header.Set("Authorization", "Bearer "+clientToken)
	payResp, _ := http.DefaultClient.Do(payReq)
	if payResp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400 when paying with insufficient balance, got %d", payResp.StatusCode)
	}

	// 5. Invoice "Download PDF" Button
	pdfReq, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/client/invoices/%d/pdf", ts.URL, invoiceID), nil)
	pdfReq.Header.Set("Authorization", "Bearer "+clientToken)
	pdfResp, err := http.DefaultClient.Do(pdfReq)
	if err != nil || pdfResp.StatusCode != http.StatusOK {
		t.Fatalf("Download PDF button failed: %v, code: %d", err, pdfResp.StatusCode)
	}

	// 6. Admin Login Button
	adminLoginBody, _ := json.Marshal(map[string]string{
		"email": "admin@fossbilling.org", "password": "SuperSecretAdmin123!",
	})
	adminResp, _ := http.Post(ts.URL+"/api/v1/admin/auth/login", "application/json", bytes.NewBuffer(adminLoginBody))
	var adminData struct{ Data struct{ Token string `json:"token"` } `json:"data"` }
	_ = json.NewDecoder(adminResp.Body).Decode(&adminData)
	adminToken := adminData.Data.Token

	// 7. Admin "Activate Order" Button
	actReq, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/admin/orders/%d/activate", ts.URL, orderID), nil)
	actReq.Header.Set("Authorization", "Bearer "+adminToken)
	actResp, _ := http.DefaultClient.Do(actReq)
	if actResp.StatusCode != http.StatusOK {
		t.Fatalf("Admin Activate Order button failed, code: %d", actResp.StatusCode)
	}

	// 8. Admin "Suspend Order" Button
	suspBody, _ := json.Marshal(map[string]string{"reason": "TDD testing suspension"})
	suspReq, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/admin/orders/%d/suspend", ts.URL, orderID), bytes.NewBuffer(suspBody))
	suspReq.Header.Set("Authorization", "Bearer "+adminToken)
	suspResp, _ := http.DefaultClient.Do(suspReq)
	if suspResp.StatusCode != http.StatusOK {
		t.Fatalf("Admin Suspend Order button failed, code: %d", suspResp.StatusCode)
	}

	// 9. Admin "Unsuspend Order" Button
	unsuspReq, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/admin/orders/%d/unsuspend", ts.URL, orderID), nil)
	unsuspReq.Header.Set("Authorization", "Bearer "+adminToken)
	unsuspResp, _ := http.DefaultClient.Do(unsuspReq)
	if unsuspResp.StatusCode != http.StatusOK {
		t.Fatalf("Admin Unsuspend Order button failed, code: %d", unsuspResp.StatusCode)
	}

	// 10. Support "Submit Ticket", "Staff Reply", and "Close Ticket" Buttons
	tReqBody, _ := json.Marshal(map[string]interface{}{
		"helpdesk_id": 1, "subject": "Button Action Test", "message": "Testing interactive replies", "priority": "medium",
	})
	openReq, _ := http.NewRequest("POST", ts.URL+"/api/v1/client/support/tickets", bytes.NewBuffer(tReqBody))
	openReq.Header.Set("Authorization", "Bearer "+clientToken)
	openResp, _ := http.DefaultClient.Do(openReq)
	var tData struct{ Data domain.Ticket `json:"data"` }
	_ = json.NewDecoder(openResp.Body).Decode(&tData)
	ticketID := tData.Data.ID

	// Staff reply button
	repBody, _ := json.Marshal(map[string]string{"message": "Acknowledged by admin"})
	repReq, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/admin/support/tickets/%d/reply", ts.URL, ticketID), bytes.NewBuffer(repBody))
	repReq.Header.Set("Authorization", "Bearer "+adminToken)
	repResp, _ := http.DefaultClient.Do(repReq)
	if repResp.StatusCode != http.StatusCreated {
		t.Fatalf("Staff Reply button failed, code: %d", repResp.StatusCode)
	}

	// Close ticket button
	closeReq, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/client/support/tickets/%d/close", ts.URL, ticketID), nil)
	closeReq.Header.Set("Authorization", "Bearer "+clientToken)
	closeResp, _ := http.DefaultClient.Do(closeReq)
	if closeResp.StatusCode != http.StatusOK {
		t.Fatalf("Client Close Ticket button failed, code: %d", closeResp.StatusCode)
	}
}
