package domain

import (
	"context"
	"time"
)

type ExtensionType string

const (
	ExtensionTypeMod       ExtensionType = "mod"
	ExtensionTypeTheme     ExtensionType = "theme"
	ExtensionTypeGateway   ExtensionType = "gateway"
	ExtensionTypeRegistrar ExtensionType = "registrar"
	ExtensionTypePlugin    ExtensionType = "plugin"
	ExtensionTypeService   ExtensionType = "service"
)

type ExtensionStatus string

const (
	ExtensionStatusActive   ExtensionStatus = "active"
	ExtensionStatusInactive ExtensionStatus = "inactive"
	ExtensionStatusCore     ExtensionStatus = "core"
)

type Extension struct {
	ID          string                 `json:"id"` // Unique slug/id, e.g. "antispam", "formbuilder"
	Name        string                 `json:"name"`
	Type        ExtensionType          `json:"type"`
	Version     string                 `json:"version"`
	Description string                 `json:"description"`
	Author      string                 `json:"author"`
	AuthorURL   string                 `json:"author_url,omitempty"`
	Icon        string                 `json:"icon,omitempty"`
	Status      ExtensionStatus        `json:"status"`
	HasSettings bool                   `json:"has_settings"`
	Config      map[string]interface{} `json:"config,omitempty"`
	Manifest    map[string]interface{} `json:"manifest,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type ExtensionFilter struct {
	Type        string
	Status      string
	Installed   *bool
	Active      *bool
	HasSettings *bool
	Search      string
}

type MarketplaceExtension struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Type        ExtensionType `json:"type"`
	Version     string        `json:"version"`
	Description string        `json:"description"`
	Author      string        `json:"author"`
	IconURL     string        `json:"icon_url"`
	DownloadURL string        `json:"download_url"`
	Rating      float64       `json:"rating"`
	Downloads   int           `json:"downloads"`
	Readme      string        `json:"readme,omitempty"`
}

type ExtensionRepository interface {
	List(ctx context.Context, filter ExtensionFilter) ([]*Extension, error)
	GetByID(ctx context.Context, id string) (*Extension, error)
	Create(ctx context.Context, ext *Extension) error
	Update(ctx context.Context, ext *Extension) error
	Delete(ctx context.Context, id string) error
	GetConfig(ctx context.Context, extID string) (map[string]interface{}, error)
	UpdateConfig(ctx context.Context, extID string, config map[string]interface{}) error
}
