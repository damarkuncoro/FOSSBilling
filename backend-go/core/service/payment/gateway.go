package payment

import (
	"context"
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

// PaymentRequest contains information required to initiate a checkout transaction
type PaymentRequest struct {
	InvoiceID   int64         `json:"invoice_id"`
	InvoiceNr   string        `json:"invoice_nr"`
	Amount      decimal.Money `json:"amount"`
	Currency    string        `json:"currency"`
	ClientEmail string        `json:"client_email"`
	ClientName  string        `json:"client_name"`
	ReturnURL   string        `json:"return_url"`
	CancelURL   string        `json:"cancel_url"`
}

// PaymentResponse contains gateway checkout redirect URL or payment token
type PaymentResponse struct {
	GatewayID     string                 `json:"gateway_id"`
	TransactionID string                 `json:"transaction_id"`
	RedirectURL   string                 `json:"redirect_url,omitempty"`
	Token         string                 `json:"token,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// WebhookResult represents the normalized result of an inbound webhook event
type WebhookResult struct {
	GatewayID     string                 `json:"gateway_id"`
	TransactionID string                 `json:"transaction_id"`
	InvoiceID     int64                  `json:"invoice_id"`
	Amount        decimal.Money          `json:"amount"`
	Currency      string                 `json:"currency"`
	IsPaid        bool                   `json:"is_paid"`
	RawPayload    []byte                 `json:"raw_payload"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// PaymentGateway interface standardizes operations across different payment service providers
type PaymentGateway interface {
	ID() string
	Name() string
	Type() string
	InitiatePayment(ctx context.Context, req PaymentRequest) (*PaymentResponse, error)
	ParseWebhook(r *http.Request) (*WebhookResult, error)
}
