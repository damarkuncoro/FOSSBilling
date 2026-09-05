package domain

import (
	"context"
	"encoding/json"
	"time"
)

type ProductType string

const (
	ProductTypeHosting      ProductType = "hosting"
	ProductTypeDomain       ProductType = "domain"
	ProductTypeLicense      ProductType = "license"
	ProductTypeDownloadable ProductType = "downloadable"
	ProductTypeCustom       ProductType = "custom"
)

type SetupType string

const (
	SetupTypeFree      SetupType = "free"
	SetupTypeOnetime   SetupType = "onetime"
	SetupTypeRecurring SetupType = "recurring"
)

type Product struct {
	ID          int64           `json:"id"`
	CategoryID  *int64          `json:"category_id,omitempty"`
	Type        ProductType     `json:"type"`
	Name        string          `json:"name"`
	Slug        string          `json:"slug"`
	Description string          `json:"description"`
	Status      string          `json:"status"` // enabled, disabled, hidden
	SetupType   SetupType       `json:"setup_type"`
	Config      json.RawMessage `json:"config,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type ProductRepository interface {
	GetByID(ctx context.Context, id int64) (*Product, error)
	GetBySlug(ctx context.Context, slug string) (*Product, error)
	List(ctx context.Context, limit, offset int) ([]*Product, int, error)
	Create(ctx context.Context, product *Product) error
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id int64) error
}
