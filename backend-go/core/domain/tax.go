package domain

import (
	"context"
	"time"
)

type TaxRule struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Country   *string   `json:"country,omitempty"`
	State     *string   `json:"state,omitempty"`
	Rate      float64   `json:"rate"` // percentage e.g. 11.00
	IsActive  bool      `json:"is_active"`
	TaxExempt bool      `json:"tax_exempt"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TaxRepository interface {
	GetByID(ctx context.Context, id int64) (*TaxRule, error)
	List(ctx context.Context) ([]*TaxRule, error)
	GetByLocation(ctx context.Context, country, state string) (*TaxRule, error)
	Create(ctx context.Context, rule *TaxRule) error
	Update(ctx context.Context, rule *TaxRule) error
	Delete(ctx context.Context, id int64) error
}
