package theme

import (
	"context"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type ThemeService struct{ repo domain.ThemeRepository }

func NewThemeService(r domain.ThemeRepository) *ThemeService { return &ThemeService{r} }

func (s *ThemeService) ListThemes(ctx context.Context, t domain.ThemeTarget) ([]*domain.Theme, error) { return s.repo.List(ctx, t) }
func (s *ThemeService) GetTheme(ctx context.Context, c string) (*domain.Theme, error) { return s.repo.GetByCode(ctx, c) }
func (s *ThemeService) GetCurrentTheme(ctx context.Context, t domain.ThemeTarget) (*domain.Theme, error) { return s.repo.GetCurrent(ctx, t) }
func (s *ThemeService) SelectTheme(ctx context.Context, c string, t domain.ThemeTarget) error { return s.repo.SetCurrent(ctx, c, t) }
func (s *ThemeService) GetConfig(ctx context.Context, c string) (map[string]any, error) { return s.repo.GetConfig(ctx, c) }
func (s *ThemeService) UpdateConfig(ctx context.Context, c string, cfg map[string]any) error { return s.repo.UpdateConfig(ctx, c, cfg) }
