package gateways

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/payment"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

type StripeGateway struct {
	secretKey      string
	publishableKey string
	webhookSecret  string
	client         *http.Client
}

func NewStripeGateway(secretKey, publishableKey, webhookSecret string) *StripeGateway {
	return &StripeGateway{
		secretKey:      secretKey,
		publishableKey: publishableKey,
		webhookSecret:  webhookSecret,
		client:         &http.Client{Timeout: 15 * time.Second},
	}
}

func (g *StripeGateway) ID() string   { return "stripe" }
func (g *StripeGateway) Name() string { return "Stripe (Credit / Debit Card, Apple Pay, Google Pay)" }
func (g *StripeGateway) Type() string { return "cc" }

func (g *StripeGateway) InitiatePayment(ctx context.Context, req payment.PaymentRequest) (*payment.PaymentResponse, error) {
	// Fallback for dev
	if g.secretKey == "" || g.secretKey == "sk_test" {
		sessionID := fmt.Sprintf("MOCK-CS-%d", req.InvoiceID)
		return &payment.PaymentResponse{
			GatewayID:     g.ID(),
			TransactionID: fmt.Sprintf("MOCK-STRIPE-%d", req.InvoiceID),
			RedirectURL:   "https://checkout.stripe.com/c/pay/" + sessionID,
			Token:         sessionID,
		}, nil
	}

	apiURL := "https://api.stripe.com/v1/checkout/sessions"

	// Create form values for Stripe API
	data := url.Values{}
	data.Set("success_url", req.ReturnURL)
	data.Set("cancel_url", req.CancelURL)
	data.Set("mode", "payment")
	data.Set("customer_email", req.ClientEmail)
	data.Set("client_reference_id", fmt.Sprintf("%d", req.InvoiceID))
	data.Set("line_items[0][price_data][currency]", strings.ToLower(req.Currency))
	data.Set("line_items[0][price_data][product_data][name]", req.InvoiceNr)
	data.Set("line_items[0][price_data][unit_amount]", fmt.Sprintf("%d", int64(req.Amount.ToFloat()*100)))
	data.Set("line_items[0][quantity]", "1")
	data.Set("metadata[invoice_id]", fmt.Sprintf("%d", req.InvoiceID))

	apiReq, err := http.NewRequestWithContext(ctx, "POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	apiReq.Header.Set("Authorization", "Bearer "+g.secretKey)
	apiReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := g.client.Do(apiReq)
	if err != nil {
		return nil, fmt.Errorf("stripe api error: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("stripe error (%d): %s", resp.StatusCode, string(respBody))
	}

	var stripeResp struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	}
	if err := json.Unmarshal(respBody, &stripeResp); err != nil {
		return nil, err
	}

	return &payment.PaymentResponse{
		GatewayID:     g.ID(),
		TransactionID: stripeResp.ID,
		RedirectURL:   stripeResp.URL,
		Token:         stripeResp.ID,
	}, nil
}

func (g *StripeGateway) ParseWebhook(r *http.Request) (*payment.WebhookResult, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var payload struct {
		ID   string `json:"id"`
		Type string `json:"type"`
		Data struct {
			Object struct {
				ID          string `json:"id"`
				AmountTotal int64  `json:"amount_total"`
				Currency    string `json:"currency"`
				Metadata    struct {
					InvoiceID string `json:"invoice_id"`
				} `json:"metadata"`
				PaymentStatus string `json:"payment_status"`
			} `json:"object"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}

	invID, _ := strconv.ParseInt(payload.Data.Object.Metadata.InvoiceID, 10, 64)
	if invID == 0 {
		// Fallback check from transaction ID or custom parsing
		parts := strings.Split(payload.Data.Object.ID, "_")
		if len(parts) >= 3 {
			invID, _ = strconv.ParseInt(parts[2], 10, 64)
		}
	}

	// Stripe amounts are in cents/smallest currency unit
	amtFloat := float64(payload.Data.Object.AmountTotal) / 100.0
	isPaid := payload.Type == "checkout.session.completed" || payload.Data.Object.PaymentStatus == "paid"

	return &payment.WebhookResult{
		GatewayID:     g.ID(),
		TransactionID: payload.Data.Object.ID,
		InvoiceID:     invID,
		Amount:        decimal.FromFloat(amtFloat),
		Currency:      strings.ToUpper(payload.Data.Object.Currency),
		IsPaid:        isPaid,
		RawPayload:    body,
	}, nil
}
