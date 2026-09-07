package domain

import (
	"context"
	"time"
)

type ThemeTarget string

const (
	ThemeTargetClient ThemeTarget = "client"
	ThemeTargetAdmin  ThemeTarget = "admin"
)

type Theme struct {
	Code        string                 `json:"code"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Version     string                 `json:"version"`
	Author      string                 `json:"author"`
	AuthorURL   string                 `json:"author_url,omitempty"`
	PreviewURL  string                 `json:"preview_url,omitempty"`
	Target      ThemeTarget            `json:"target"`
	IsCurrent   bool                   `json:"is_current"`
	Config      map[string]interface{} `json:"config,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type ThemeRepository interface {
	List(ctx context.Context, target ThemeTarget) ([]*Theme, error)
	GetByCode(ctx context.Context, code string) (*Theme, error)
	GetCurrent(ctx context.Context, target ThemeTarget) (*Theme, error)
	SetCurrent(ctx context.Context, code string, target ThemeTarget) error
	GetConfig(ctx context.Context, code string) (map[string]interface{}, error)
	UpdateConfig(ctx context.Context, code string, config map[string]interface{}) error
}
