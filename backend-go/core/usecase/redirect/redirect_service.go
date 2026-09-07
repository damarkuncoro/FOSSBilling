package redirect

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
)

type RedirectService struct {
	repo domain.RedirectRepository
}

func NewRedirectService(repo domain.RedirectRepository) *RedirectService {
	return &RedirectService{repo: repo}
}

func (s *RedirectService) SanitizePath(p string) string {
	p = strings.TrimSpace(p)
	p = strings.TrimPrefix(p, "/")
	return "/" + p
}

func (s *RedirectService) ValidateTarget(target string) error {
	target = strings.TrimSpace(target)
	if target == "" {
		return errors.New("target redirect URL cannot be empty")
	}

	if strings.HasPrefix(target, "/") {
		return nil // internal relative redirect
	}

	u, err := url.Parse(target)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		return errors.New("target must be a valid relative path (starting with '/') or an absolute HTTP/HTTPS URL")
	}
	return nil
}

func (s *RedirectService) ListRedirects(ctx context.Context, limit, offset int) ([]*domain.Redirect, int, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *RedirectService) GetRedirect(ctx context.Context, id int64) (*domain.Redirect, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *RedirectService) MatchRedirect(ctx context.Context, reqPath string) (*domain.Redirect, error) {
	normalized := s.SanitizePath(reqPath)

	r, err := s.repo.GetByPath(ctx, normalized)
	if err != nil {
		return nil, err
	}

	if !r.IsEnabled {
		return nil, appErrors.ErrNotFound
	}

	// Increment hit count asynchronously or inline
	_ = s.repo.IncrementHitCount(ctx, r.ID)
	return r, nil
}

type CreateRedirectDTO struct {
	Path       string `json:"path"`
	Target     string `json:"target"`
	StatusCode int    `json:"status_code"`
	IsEnabled  bool   `json:"is_enabled"`
}

func (s *RedirectService) CreateRedirect(ctx context.Context, dto CreateRedirectDTO) (*domain.Redirect, error) {
	path := s.SanitizePath(dto.Path)
	if path == "/" {
		return nil, errors.New("cannot redirect root '/' path")
	}

	if err := s.ValidateTarget(dto.Target); err != nil {
		return nil, err
	}

	// Check if already exists
	existing, err := s.repo.GetByPath(ctx, path)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("redirect rule for path '%s' already exists", path)
	}

	code := dto.StatusCode
	if code != 301 && code != 302 && code != 307 && code != 308 {
		code = 301
	}

	item := &domain.Redirect{
		Path:       path,
		Target:     strings.TrimSpace(dto.Target),
		StatusCode: code,
		IsEnabled:  dto.IsEnabled,
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

type UpdateRedirectDTO struct {
	ID         int64   `json:"id"`
	Path       *string `json:"path"`
	Target     *string `json:"target"`
	StatusCode *int    `json:"status_code"`
	IsEnabled  *bool   `json:"is_enabled"`
}

func (s *RedirectService) UpdateRedirect(ctx context.Context, dto UpdateRedirectDTO) (*domain.Redirect, error) {
	item, err := s.repo.GetByID(ctx, dto.ID)
	if err != nil {
		return nil, err
	}

	if dto.Path != nil {
		p := s.SanitizePath(*dto.Path)
		if p == "/" {
			return nil, errors.New("cannot redirect root '/' path")
		}
		item.Path = p
	}

	if dto.Target != nil {
		if err := s.ValidateTarget(*dto.Target); err != nil {
			return nil, err
		}
		item.Target = strings.TrimSpace(*dto.Target)
	}

	if dto.StatusCode != nil {
		code := *dto.StatusCode
		if code == 301 || code == 302 || code == 307 || code == 308 {
			item.StatusCode = code
		}
	}

	if dto.IsEnabled != nil {
		item.IsEnabled = *dto.IsEnabled
	}

	if err := s.repo.Update(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *RedirectService) DeleteRedirect(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
