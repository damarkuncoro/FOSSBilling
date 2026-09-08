package gateways

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/payment"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

type MidtransGateway struct {
	serverKey string
	clientKey string
	isProd    bool
	client    *http.Client
}

func NewMidtransGateway(serverKey, clientKey string, isProd bool) *MidtransGateway {
	return &MidtransGateway{
		serverKey: serverKey,
		clientKey: clientKey,
		isProd:    isProd,
		client:    &http.Client{Timeout: 15 * time.Second},
	}
}

func (g *MidtransGateway) ID() string   { return "midtrans" }
func (g *MidtransGateway) Name() string { return "Midtrans Payment Gateway (QRIS, GoPay, VA)" }
func (g *MidtransGateway) Type() string { return "wallet" }

func (g *MidtransGateway) getBaseURL() string {
	if g.isProd {
		return "https://app.midtrans.com/snap/v1"
	}
	return "https://app.sandbox.midtrans.com/snap/v1"
}

func (g *MidtransGateway) InitiatePayment(ctx context.Context, req payment.PaymentRequest) (*payment.PaymentResponse, error) {
	// If keys are empty, fallback to mock (for dev environment)
	if g.serverKey == "" || g.serverKey == "server_key" {
		snapToken := fmt.Sprintf("MOCK-SNAP-%d", req.InvoiceID)
		return &payment.PaymentResponse{
			GatewayID:     g.ID(),
			TransactionID: fmt.Sprintf("MOCK-MID-%d", req.InvoiceID),
			RedirectURL:   "https://app.sandbox.midtrans.com/snap/v2/vtweb/" + snapToken,
			Token:         snapToken,
		}, nil
	}

	orderID := fmt.Sprintf("INV-%d-%d", req.InvoiceID, time.Now().Unix())
	payload := map[string]interface{}{
		"transaction_details": map[string]interface{}{
			"order_id":     orderID,
			"gross_amount": int64(req.Amount.ToFloat()), // Midtrans expects integer for IDR
		},
		"customer_details": map[string]interface{}{
			"first_name": req.ClientName,
			"email":      req.ClientEmail,
		},
		"expiry": map[string]interface{}{
			"duration": 24,
			"unit":     "hours",
		},
	}

	body, _ := json.Marshal(payload)
	apiReq, err := http.NewRequestWithContext(ctx, "POST", g.getBaseURL()+"/transactions", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(g.serverKey + ":"))
	apiReq.Header.Set("Authorization", "Basic "+auth)
	apiReq.Header.Set("Content-Type", "application/json")
	apiReq.Header.Set("Accept", "application/json")

	resp, err := g.client.Do(apiReq)
	if err != nil {
		return nil, fmt.Errorf("midtrans api error: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("midtrans error (%d): %s", resp.StatusCode, string(respBody))
	}

	var midtransResp struct {
		Token       string `json:"token"`
		RedirectURL string `json:"redirect_url"`
	}
	if err := json.Unmarshal(respBody, &midtransResp); err != nil {
		return nil, err
	}

	return &payment.PaymentResponse{
		GatewayID:     g.ID(),
		TransactionID: orderID,
		RedirectURL:   midtransResp.RedirectURL,
		Token:         midtransResp.Token,
	}, nil
}

func (g *MidtransGateway) ParseWebhook(r *http.Request) (*payment.WebhookResult, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var payload struct {
		TransactionStatus string `json:"transaction_status"`
		OrderID           string `json:"order_id"`
		GrossAmount       string `json:"gross_amount"`
		TransactionID     string `json:"transaction_id"`
		StatusCode        string `json:"status_code"`
		Currency          string `json:"currency"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}

	// Extract invoice ID from order_id (e.g., "INV-101" or "101")
	invIDStr := strings.TrimPrefix(payload.OrderID, "INV-")
	invID, _ := strconv.ParseInt(invIDStr, 10, 64)

	amtFloat, _ := strconv.ParseFloat(payload.GrossAmount, 64)
	isPaid := payload.TransactionStatus == "settlement" || payload.TransactionStatus == "capture"

	return &payment.WebhookResult{
		GatewayID:     g.ID(),
		TransactionID: payload.TransactionID,
		InvoiceID:     invID,
		Amount:        decimal.FromFloat(amtFloat),
		Currency:      payload.Currency,
		IsPaid:        isPaid,
		RawPayload:    body,
	}, nil
}
