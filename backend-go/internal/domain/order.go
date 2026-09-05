package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/fossbilling/backend-go/pkg/decimal"
)

type OrderStatus string

const (
	OrderStatusPendingSetup OrderStatus = "pending_setup"
	OrderStatusActive       OrderStatus = "active"
	OrderStatusSuspended    OrderStatus = "suspended"
	OrderStatusCanceled     OrderStatus = "canceled"
	OrderStatusTerminated   OrderStatus = "terminated"
)

type Order struct {
	ID               int64           `json:"id"`
	ClientID         int64           `json:"client_id"`
	ProductID        int64           `json:"product_id"`
	InvoiceID        *int64          `json:"invoice_id,omitempty"`
	Status           OrderStatus     `json:"status"`
	Title            string          `json:"title"`
	Period           string          `json:"period"` // 1M, 3M, 1Y, etc.
	Price            decimal.Money   `json:"price"`
	Currency         string          `json:"currency"`
	Config           json.RawMessage `json:"config,omitempty"`
	ActivatedAt      *time.Time      `json:"activated_at,omitempty"`
	ExpiresAt        *time.Time      `json:"expires_at,omitempty"`
	NextDueDate      *time.Time      `json:"next_due_date,omitempty"`
	SuspendedAt      *time.Time      `json:"suspended_at,omitempty"`
	SuspensionReason *string         `json:"suspension_reason,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

type OrderRepository interface {
	GetByID(ctx context.Context, id int64) (*Order, error)
	ListByClientID(ctx context.Context, clientID int64, limit, offset int) ([]*Order, int, error)
	List(ctx context.Context, limit, offset int) ([]*Order, int, error)
	ListDueOrders(ctx context.Context, dueBefore time.Time) ([]*Order, error)
	ListOverdueSuspensions(ctx context.Context, overdueDays int) ([]*Order, error)
	Create(ctx context.Context, order *Order) error
	Update(ctx context.Context, order *Order) error
	UpdateStatus(ctx context.Context, id int64, status OrderStatus, reason *string) error
}

