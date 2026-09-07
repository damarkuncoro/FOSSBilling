package theme

import (
	"context"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type ThemeService struct {
	repo domain.ThemeRepository
}

func NewThemeService(repo domain.ThemeRepository) *ThemeService {
	return &ThemeService{repo: repo}
}

func (s *ThemeService) ListThemes(ctx context.Context, target domain.ThemeTarget) ([]*domain.Theme, error) {
	return s.repo.List(ctx, target)
}

func (s *ThemeService) GetTheme(ctx context.Context, code string) (*domain.Theme, error) {
	return s.repo.GetByCode(ctx, code)
}

func (s *ThemeService) GetCurrentTheme(ctx context.Context, target domain.ThemeTarget) (*domain.Theme, error) {
	return s.repo.GetCurrent(ctx, target)
}

func (s *ThemeService) SelectTheme(ctx context.Context, code string, target domain.ThemeTarget) error {
	return s.repo.SetCurrent(ctx, code, target)
}

func (s *ThemeService) GetConfig(ctx context.Context, code string) (map[string]interface{}, error) {
	return s.repo.GetConfig(ctx, code)
}

func (s *ThemeService) UpdateConfig(ctx context.Context, code string, config map[string]interface{}) error {
	return s.repo.UpdateConfig(ctx, code, config)
}
