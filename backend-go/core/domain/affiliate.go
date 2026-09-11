package domain

import (
	"context"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

type AffiliateStatus string

const (
	AffiliateStatusActive   AffiliateStatus = "active"
	AffiliateStatusPending  AffiliateStatus = "pending"
	AffiliateStatusSuspended AffiliateStatus = "suspended"
)

type Affiliate struct {
	ID             int64           `json:"id"`
	ClientID       int64           `json:"client_id"`
	CommissionRate float64         `json:"commission_rate"` // e.g. 10.0 for 10%
	Status         AffiliateStatus `json:"status"`
	Balance        decimal.Money   `json:"balance"`
	TotalEarned    decimal.Money   `json:"total_earned"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type AffiliateReferral struct {
	ID          int64         `json:"id"`
	AffiliateID int64         `json:"affiliate_id"`
	ClientID    int64         `json:"client_id"` // Referred client
	OrderID     int64         `json:"order_id"`
	Amount      decimal.Money `json:"amount"` // Commission amount
	Status      string        `json:"status"` // pending, approved, canceled
	CreatedAt   time.Time     `json:"created_at"`
}

type AffiliateRepository interface {
	GetByClientID(ctx context.Context, clientID int64) (*Affiliate, error)
	Create(ctx context.Context, aff *Affiliate) error
	Update(ctx context.Context, aff *Affiliate) error

	AddReferral(ctx context.Context, ref *AffiliateReferral) error
	ListReferrals(ctx context.Context, affiliateID int64) ([]*AffiliateReferral, error)
}
