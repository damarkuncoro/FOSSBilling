package domain

import (
	"context"
	"time"
)

type Currency struct {
	ID                 int64     `json:"id"`
	Code               string    `json:"code"` // e.g. "USD", "IDR", "EUR"
	Title              string    `json:"title"`
	ConversionRate     float64   `json:"conversion_rate"`
	Format             string    `json:"format"` // e.g. "$ {{price}}" or "Rp {{price}}"
	PriceFormat        string    `json:"price_format"` // 1 | 2 (decimal places)
	IsDefault          bool      `json:"is_default"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type CurrencyRepository interface {
	GetByCode(ctx context.Context, code string) (*Currency, error)
	GetDefault(ctx context.Context) (*Currency, error)
	List(ctx context.Context) ([]*Currency, error)
	Create(ctx context.Context, currency *Currency) error
	Update(ctx context.Context, currency *Currency) error
	SetDefault(ctx context.Context, code string) error
	Delete(ctx context.Context, code string) error
}
