package domain

import (
	"context"
	"time"
)

type ProductCategory struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Description string    `json:"description,omitempty"`
	Icon        string    `json:"icon,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Server struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Hostname    string    `json:"hostname"`
	IP          string    `json:"ip"`
	Manager     string    `json:"manager"` // cpanel, plesk, directadmin, etc.
	Status      string    `json:"status"`  // online, offline
	AccessKey   string    `json:"-"`
	IsDefault   bool      `json:"is_default"`
	MaxAccounts int       `json:"max_accounts"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TLD struct {
	ID                int64   `json:"id"`
	Tld               string  `json:"tld"`
	RegistrarID       string  `json:"registrar_id"`
	PriceRegistration float64 `json:"price_registration"`
	PriceRenewal      float64 `json:"price_renewal"`
	PriceTransfer     float64 `json:"price_transfer"`
	MinYears          int     `json:"min_years"`
	IsActive          bool    `json:"is_active"`
}

type CatalogRepository interface {
	ListCategories(ctx context.Context) ([]*ProductCategory, error)
	GetCategoryByID(ctx context.Context, id int64) (*ProductCategory, error)

	ListServers(ctx context.Context) ([]*Server, error)
	GetServerByID(ctx context.Context, id int64) (*Server, error)

	ListTlds(ctx context.Context) ([]*TLD, error)
}
