package gateways

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/payment"
)

type CustomGateway struct{}

func NewCustomGateway() *CustomGateway {
	return &CustomGateway{}
}

func (g *CustomGateway) ID() string   { return "custom" }
func (g *CustomGateway) Name() string { return "Custom / Test Gateway (Always Success)" }
func (g *CustomGateway) Type() string { return "test" }

func (g *CustomGateway) InitiatePayment(ctx context.Context, req payment.PaymentRequest) (*payment.PaymentResponse, error) {
	txnID := fmt.Sprintf("CUSTOM-TXN-%d", time.Now().UnixNano())

	return &payment.PaymentResponse{
		GatewayID:     g.ID(),
		TransactionID: txnID,
		RedirectURL:   fmt.Sprintf("/api/v1/guest/webhook/custom?invoice_id=%d&amount=%s&currency=%s&txn_id=%s", req.InvoiceID, req.Amount.String(), req.Currency, txnID),
		Token:         txnID,
	}
}

func (g *CustomGateway) ParseWebhook(r *http.Request) (*payment.WebhookResult, error) {
	// In this mock, we get data from query params for easy testing
	q := r.URL.Query()

	var invID int64
	fmt.Sscanf(q.Get("invoice_id"), "%d", &invID)

	return &payment.WebhookResult{
		GatewayID:     g.ID(),
		TransactionID: q.Get("txn_id"),
		InvoiceID:     invID,
		IsPaid:        true,
		Currency:      q.Get("currency"),
	}, nil
}
