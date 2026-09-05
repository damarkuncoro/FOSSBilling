package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

type TransactionType string

const (
	TransactionTypePayment TransactionType = "payment"
	TransactionTypeRefund  TransactionType = "refund"
)

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusComplete  TransactionStatus = "complete"
	TransactionStatusFailed    TransactionStatus = "failed"
	TransactionStatusRefunded  TransactionStatus = "refunded"
)

type Transaction struct {
	ID         int64             `json:"id"`
	InvoiceID  *int64            `json:"invoice_id,omitempty"`
	GatewayID  string            `json:"gateway_id"`
	TxnID      string            `json:"txn_id"`
	Type       TransactionType   `json:"type"`
	Amount     decimal.Money     `json:"amount"`
	Currency   string            `json:"currency"`
	Status     TransactionStatus `json:"status"`
	RawPayload json.RawMessage   `json:"raw_payload,omitempty"`
	CreatedAt  time.Time         `json:"created_at"`
}

type PaymentGateway interface {
	ID() string
	Name() string
	GetPaymentURL(ctx context.Context, invoice *Invoice) (string, error)
	ProcessWebhook(ctx context.Context, payload []byte, headers map[string]string) (*Transaction, error)
	ProcessRefund(ctx context.Context, txnID string, amount decimal.Money) error
}

type TransactionRepository interface {
	GetByID(ctx context.Context, id int64) (*Transaction, error)
	GetByTxnID(ctx context.Context, gatewayID, txnID string) (*Transaction, error)
	Create(ctx context.Context, txn *Transaction) error
	UpdateStatus(ctx context.Context, id int64, status TransactionStatus) error
}
