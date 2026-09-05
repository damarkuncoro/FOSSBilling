package domain

import (
	"context"
	"time"

	"github.com/fossbilling/backend-go/pkg/decimal"
)

type PromoType string

const (
	PromoTypePercentage PromoType = "percentage"
	PromoTypeAbsolute   PromoType = "absolute"
)

type Promo struct {
	ID             int64         `json:"id"`
	Code           string        `json:"code"`
	Description    string        `json:"description"`
	Type           PromoType     `json:"type"` // percentage or absolute
	Value          decimal.Money `json:"value"` // percentage (e.g. 200000 = 20.00%) or fixed money amount
	MaxUses        int           `json:"max_uses"`
	UsedCount      int           `json:"used_count"`
	OncePerClient  bool          `json:"once_per_client"`
	StartDate      *time.Time    `json:"start_date,omitempty"`
	EndDate        *time.Time    `json:"end_date,omitempty"`
	Active         bool          `json:"active"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

type PromoRedemption struct {
	ID        int64     `json:"id"`
	PromoID   int64     `json:"promo_id"`
	ClientID  int64     `json:"client_id"`
	OrderID   *int64    `json:"order_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type PromoRepository interface {
	GetByCode(ctx context.Context, code string) (*Promo, error)
	GetRedemptionCount(ctx context.Context, promoID int64, clientID int64) (int, error)
	IncrementUsed(ctx context.Context, promoID int64, clientID int64, orderID *int64) error
	Create(ctx context.Context, promo *Promo) error
}
