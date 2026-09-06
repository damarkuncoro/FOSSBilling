package domain

import (
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

// Strongly-typed event payloads for reactive event-driven architecture

type ClientRegisteredPayload struct {
	ClientID  int64     `json:"client_id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

type InvoicePaidPayload struct {
	InvoiceID int64         `json:"invoice_id"`
	ClientID  int64         `json:"client_id"`
	Amount    decimal.Money `json:"amount"`
	Currency  string        `json:"currency"`
	GatewayID string        `json:"gateway_id"`
	TxnID     string        `json:"txn_id"`
	PaidAt    time.Time     `json:"paid_at"`
}

type OrderActivatedPayload struct {
	OrderID     int64     `json:"order_id"`
	ClientID    int64     `json:"client_id"`
	ProductID   int64     `json:"product_id"`
	Title       string    `json:"title"`
	ActivatedAt time.Time `json:"activated_at"`
}

type OrderSuspendedPayload struct {
	OrderID     int64     `json:"order_id"`
	ClientID    int64     `json:"client_id"`
	Reason      string    `json:"reason"`
	SuspendedAt time.Time `json:"suspended_at"`
}

type TicketOpenedPayload struct {
	TicketID  int64     `json:"ticket_id"`
	ClientID  int64     `json:"client_id"`
	Subject   string    `json:"subject"`
	Priority  string    `json:"priority"`
	CreatedAt time.Time `json:"created_at"`
}
