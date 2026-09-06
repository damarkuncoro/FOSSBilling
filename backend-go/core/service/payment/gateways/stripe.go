package gateways

import (
	"context"
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

type StripeGateway struct {
	secretKey     string
	publishableKey string
	webhookSecret string
}

func NewStripeGateway(secretKey, publishableKey, webhookSecret string) *StripeGateway {
	return &StripeGateway{
		secretKey:      secretKey,
		publishableKey: publishableKey,
		webhookSecret: webhookSecret,
	}
}

func (g *StripeGateway) ID() string   { return "stripe" }
func (g *StripeGateway) Name() string { return "Stripe (Credit / Debit Card, Apple Pay, Google Pay)" }
func (g *StripeGateway) Type() string { return "cc" }

func (g *StripeGateway) InitiatePayment(ctx context.Context, req payment.PaymentRequest) (*payment.PaymentResponse, error) {
	sessionID := fmt.Sprintf("cs_test_%d_%d", req.InvoiceID, time.Now().Unix())
	checkoutURL := fmt.Sprintf("https://checkout.stripe.com/c/pay/%s", sessionID)

	return &payment.PaymentResponse{
		GatewayID:     g.ID(),
		TransactionID: fmt.Sprintf("txn_stripe_%d", req.InvoiceID),
		RedirectURL:   checkoutURL,
		Token:         sessionID,
		Metadata: map[string]interface{}{
			"invoice_nr": req.InvoiceNr,
			"pub_key":    g.publishableKey,
		},
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
