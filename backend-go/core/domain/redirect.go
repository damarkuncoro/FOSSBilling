package domain

import (
	"context"
	"time"
)

// Redirect represents a URL redirection rule.
type Redirect struct {
	ID         int64     `json:"id" db:"id"`
	Path       string    `json:"path" db:"path"`               // e.g. "/promo/summer" or "summer-sale"
	Target     string    `json:"target" db:"target"`           // e.g. "/order?product=3" or "https://fossbilling.org"
	StatusCode int       `json:"status_code" db:"status_code"` // 301, 302, 307, 308 (default 301)
	IsEnabled  bool      `json:"is_enabled" db:"is_enabled"`
	HitCount   int64     `json:"hit_count" db:"hit_count"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// RedirectRepository defines the data access contract for Redirects.
type RedirectRepository interface {
	List(ctx context.Context, limit, offset int) ([]*Redirect, int, error)
	GetByID(ctx context.Context, id int64) (*Redirect, error)
	GetByPath(ctx context.Context, path string) (*Redirect, error)
	Create(ctx context.Context, r *Redirect) error
	Update(ctx context.Context, r *Redirect) error
	Delete(ctx context.Context, id int64) error
	IncrementHitCount(ctx context.Context, id int64) error
}
