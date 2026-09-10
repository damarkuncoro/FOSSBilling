package memory

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
)

type MockExtensionRepository struct {
	mu         sync.RWMutex
	extensions map[string]*domain.Extension
}

func NewMockExtensionRepository() *MockExtensionRepository {
	repo := &MockExtensionRepository{
		extensions: make(map[string]*domain.Extension),
	}

	// Seed built-in core extensions
	defaultExts := []*domain.Extension{
		{
			ID:          "antispam",
			Name:        "Anti-Spam & Abuse Shield",
			Type:        domain.ExtensionTypePlugin,
			Version:     "2.0.0",
			Description: "StopForumSpam, Cloudflare Turnstile, and temporary email protection.",
			Author:      "FOSSBilling Core Team",
			Status:      domain.ExtensionStatusActive,
			HasSettings: true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "formbuilder",
			Name:        "Custom Order Formbuilder",
			Type:        domain.ExtensionTypeMod,
			Version:     "2.0.0",
			Description: "Dynamic custom checkout fields for products and provisioning options.",
			Author:      "FOSSBilling Core Team",
			Status:      domain.ExtensionStatusActive,
			HasSettings: true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "servicehosting",
			Name:        "cPanel & DirectAdmin Hosting Provisioner",
			Type:        domain.ExtensionTypeService,
			Version:     "2.1.0",
			Description: "Automated provisioning for shared and reseller hosting servers.",
			Author:      "FOSSBilling Core Team",
			Status:      domain.ExtensionStatusActive,
			HasSettings: true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "midtrans",
			Name:        "Midtrans Payment Gateway",
			Type:        domain.ExtensionTypeGateway,
			Version:     "1.3.0",
			Description: "SNAP, QRIS, Virtual Account, and Credit Card payments via Midtrans.",
			Author:      "Nusantara Developers",
			Status:      domain.ExtensionStatusActive,
			HasSettings: true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "stripe",
			Name:        "Stripe Checkout & Elements",
			Type:        domain.ExtensionTypeGateway,
			Version:     "2.0.0",
			Description: "Global credit card and debit payment processing.",
			Author:      "FOSSBilling Core Team",
			Status:      domain.ExtensionStatusActive,
			HasSettings: true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, ext := range defaultExts {
		repo.extensions[ext.ID] = ext
	}

	return repo
}

func (r *MockExtensionRepository) List(ctx context.Context, filter domain.ExtensionFilter) ([]*domain.Extension, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var list []*domain.Extension
	for _, ext := range r.extensions {
		if filter.Type != "" && string(ext.Type) != filter.Type {
			continue
		}
		if filter.Status != "" && string(ext.Status) != filter.Status {
			continue
		}
		if filter.Active != nil {
			isActive := ext.Status == domain.ExtensionStatusActive || ext.Status == domain.ExtensionStatusCore
			if *filter.Active != isActive {
				continue
			}
		}
		if filter.HasSettings != nil && ext.HasSettings != *filter.HasSettings {
			continue
		}
		if filter.Search != "" {
			term := strings.ToLower(filter.Search)
			if !strings.Contains(strings.ToLower(ext.Name), term) &&
				!strings.Contains(strings.ToLower(ext.Description), term) &&
				!strings.Contains(strings.ToLower(ext.ID), term) {
				continue
			}
		}

		copyExt := *ext
		list = append(list, &copyExt)
	}

	return list, nil
}

func (r *MockExtensionRepository) GetByID(ctx context.Context, id string) (*domain.Extension, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ext, exists := r.extensions[id]
	if !exists {
		return nil, appErrors.ErrNotFound
	}
	copyExt := *ext
	return &copyExt, nil
}

func (r *MockExtensionRepository) Create(ctx context.Context, ext *domain.Extension) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.extensions[ext.ID]; exists {
		return appErrors.ErrDuplicate
	}

	if ext.CreatedAt.IsZero() {
		ext.CreatedAt = time.Now()
	}
	ext.UpdatedAt = time.Now()
	r.extensions[ext.ID] = ext
	return nil
}

func (r *MockExtensionRepository) Update(ctx context.Context, ext *domain.Extension) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.extensions[ext.ID]; !exists {
		return appErrors.ErrNotFound
	}
	ext.UpdatedAt = time.Now()
	r.extensions[ext.ID] = ext
	return nil
}

func (r *MockExtensionRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	ext, exists := r.extensions[id]
	if !exists {
		return appErrors.ErrNotFound
	}
	if ext.Status == domain.ExtensionStatusCore {
		return appErrors.ErrForbidden
	}

	delete(r.extensions, id)
	return nil
}

func (r *MockExtensionRepository) GetConfig(ctx context.Context, extID string) (map[string]interface{}, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ext, exists := r.extensions[extID]
	if !exists {
		return nil, appErrors.ErrNotFound
	}
	if ext.Config == nil {
		return make(map[string]interface{}), nil
	}
	return ext.Config, nil
}

func (r *MockExtensionRepository) UpdateConfig(ctx context.Context, extID string, config map[string]interface{}) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	ext, exists := r.extensions[extID]
	if !exists {
		return appErrors.ErrNotFound
	}
	ext.Config = config
	ext.UpdatedAt = time.Now()
	return nil
}
