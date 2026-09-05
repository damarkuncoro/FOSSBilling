package domain

import (
	"context"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

type InvoiceStatus string

const (
	InvoiceStatusUnpaid   InvoiceStatus = "unpaid"
	InvoiceStatusPaid     InvoiceStatus = "paid"
	InvoiceStatusRefunded InvoiceStatus = "refunded"
	InvoiceStatusCanceled InvoiceStatus = "canceled"
)

type Invoice struct {
	ID           int64         `json:"id"`
	Serie        string        `json:"serie"`
	Nr           string        `json:"nr"`
	ClientID     int64         `json:"client_id"`
	Status       InvoiceStatus `json:"status"`
	Currency     string        `json:"currency"`
	CurrencyRate float64       `json:"currency_rate"`
	Subtotal     decimal.Money `json:"subtotal"`
	Tax          decimal.Money `json:"tax"`
	Total        decimal.Money `json:"total"`
	TaxRate      float64       `json:"tax_rate"`
	DueAt        time.Time     `json:"due_at"`
	PaidAt       *time.Time    `json:"paid_at,omitempty"`
	Items        []InvoiceItem `json:"items,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

type InvoiceItem struct {
	ID        int64         `json:"id"`
	InvoiceID int64         `json:"invoice_id"`
	OrderID   *int64        `json:"order_id,omitempty"`
	Title     string        `json:"title"`
	Period    *string       `json:"period,omitempty"`
	Price     decimal.Money `json:"price"`
	Quantity  int           `json:"quantity"`
	Unit      string        `json:"unit"`
	Taxable   bool          `json:"taxable"`
	CreatedAt time.Time     `json:"created_at"`
}

type InvoiceRepository interface {
	GetByID(ctx context.Context, id int64) (*Invoice, error)
	ListByClientID(ctx context.Context, clientID int64, limit, offset int) ([]*Invoice, int, error)
	List(ctx context.Context, limit, offset int) ([]*Invoice, int, error)
	Create(ctx context.Context, invoice *Invoice, items []InvoiceItem) error
	MarkAsPaid(ctx context.Context, id int64, paidAt time.Time) error
	Update(ctx context.Context, invoice *Invoice) error
}

