package domain

import (
	"context"
	"time"
)

type APIKey struct {
	ID        int64      `json:"id"`
	ClientID  int64      `json:"client_id"`
	Name      string     `json:"name"`
	Key       string     `json:"key"`
	Secret    string     `json:"-"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type APIKeyRepository interface {
	GetByID(ctx context.Context, id int64) (*APIKey, error)
	GetByKey(ctx context.Context, key string) (*APIKey, error)
	ListByClientID(ctx context.Context, clientID int64) ([]*APIKey, error)
	Create(ctx context.Context, apiKey *APIKey) error
	Delete(ctx context.Context, id, clientID int64) error
}
