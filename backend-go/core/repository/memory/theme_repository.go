package memory

import (
	"context"
	"sync"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
)

type ThemeRepository struct {
	mu     sync.RWMutex
	themes map[string]*domain.Theme
}

func NewThemeRepository() *ThemeRepository {
	now := time.Now().UTC()
	themes := map[string]*domain.Theme{
		"huraga": {
			Code:        "huraga",
			Name:        "Huraga Modern Client Portal",
			Description: "Clean, responsive client dashboard with dark/light mode and Tailwind CSS.",
			Version:     "2.0.0",
			Author:      "FOSSBilling Core Team",
			Target:      domain.ThemeTargetClient,
			IsCurrent:   true,
			Config: map[string]interface{}{
				"primary_color": "#0284c7",
				"dark_mode":     true,
				"show_hero":     true,
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		"default_admin": {
			Code:        "default_admin",
			Name:        "Tabler Admin Dashboard",
			Description: "Unified administrative portal built with React, Lucide Icons, and shadcn/ui.",
			Version:     "2.0.0",
			Author:      "FOSSBilling Core Team",
			Target:      domain.ThemeTargetAdmin,
			IsCurrent:   true,
			Config: map[string]interface{}{
				"compact_sidebar": false,
				"dense_tables":    false,
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		"antigravity_sleek": {
			Code:        "antigravity_sleek",
			Name:        "Antigravity Sleek Theme",
			Description: "Futuristic dark-glass theme with micro-animations and glowing accents.",
			Version:     "1.0.0",
			Author:      "Nusantara Developers",
			Target:      domain.ThemeTargetClient,
			IsCurrent:   false,
			Config: map[string]interface{}{
				"glassmorphism": true,
				"glow_accent":   "#6366f1",
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	return &ThemeRepository{themes: themes}
}

func (r *ThemeRepository) List(ctx context.Context, target domain.ThemeTarget) ([]*domain.Theme, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var list []*domain.Theme
	for _, t := range r.themes {
		if target == "" || t.Target == target {
			cpy := *t
			list = append(list, &cpy)
		}
	}
	return list, nil
}

func (r *ThemeRepository) GetByCode(ctx context.Context, code string) (*domain.Theme, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.themes[code]
	if !ok {
		return nil, appErrors.ErrNotFound
	}
	cpy := *t
	return &cpy, nil
}

func (r *ThemeRepository) GetCurrent(ctx context.Context, target domain.ThemeTarget) (*domain.Theme, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, t := range r.themes {
		if t.Target == target && t.IsCurrent {
			cpy := *t
			return &cpy, nil
		}
	}
	return nil, appErrors.ErrNotFound
}

func (r *ThemeRepository) SetCurrent(ctx context.Context, code string, target domain.ThemeTarget) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	targetTheme, ok := r.themes[code]
	if !ok {
		return appErrors.ErrNotFound
	}

	if targetTheme.Target != target {
		return appErrors.ErrNotFound
	}

	for _, t := range r.themes {
		if t.Target == target {
			t.IsCurrent = (t.Code == code)
		}
	}
	return nil
}

func (r *ThemeRepository) GetConfig(ctx context.Context, code string) (map[string]interface{}, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.themes[code]
	if !ok {
		return nil, appErrors.ErrNotFound
	}
	return t.Config, nil
}

func (r *ThemeRepository) UpdateConfig(ctx context.Context, code string, config map[string]interface{}) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	t, ok := r.themes[code]
	if !ok {
		return appErrors.ErrNotFound
	}
	t.Config = config
	t.UpdatedAt = time.Now().UTC()
	return nil
}
