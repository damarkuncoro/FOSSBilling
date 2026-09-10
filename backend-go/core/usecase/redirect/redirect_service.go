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

type RedirectService struct{ repo domain.RedirectRepository }

func NewRedirectService(r domain.RedirectRepository) *RedirectService { return &RedirectService{r} }

func (s *RedirectService) san(p string) string { return "/" + strings.Trim(strings.TrimSpace(p), "/") }

func (s *RedirectService) val(t string) error {
	t = strings.TrimSpace(t); if t == "" { return errors.New("empty target") }
	if strings.HasPrefix(t, "/") { return nil }
	u, err := url.Parse(t); if err != nil || (u.Scheme != "http" && u.Scheme != "https") { return errors.New("invalid target") }
	return nil
}

func (s *RedirectService) ListRedirects(ctx context.Context, l, o int) ([]*domain.Redirect, int, error) { return s.repo.List(ctx, l, o) }
func (s *RedirectService) GetRedirect(ctx context.Context, id int64) (*domain.Redirect, error) { return s.repo.GetByID(ctx, id) }

func (s *RedirectService) MatchRedirect(ctx context.Context, p string) (*domain.Redirect, error) {
	r, err := s.repo.GetByPath(ctx, s.san(p)); if err != nil || !r.IsEnabled { return nil, appErrors.ErrNotFound }
	_ = s.repo.IncrementHitCount(ctx, r.ID); return r, nil
}

type CreateRedirectDTO struct { Path, Target string; StatusCode int; IsEnabled bool }

func (s *RedirectService) CreateRedirect(ctx context.Context, d CreateRedirectDTO) (*domain.Redirect, error) {
	p := s.san(d.Path); if p == "/" { return nil, errors.New("cannot redirect /") }
	if err := s.val(d.Target); err != nil { return nil, err }
	if ex, _ := s.repo.GetByPath(ctx, p); ex != nil { return nil, fmt.Errorf("exists: %s", p) }
	c := d.StatusCode; if c != 301 && c != 302 && c != 307 && c != 308 { c = 301 }
	r := &domain.Redirect{Path: p, Target: strings.TrimSpace(d.Target), StatusCode: c, IsEnabled: d.IsEnabled}
	return r, s.repo.Create(ctx, r)
}

type UpdateRedirectDTO struct { ID int64; Path, Target *string; StatusCode *int; IsEnabled *bool }

func (s *RedirectService) UpdateRedirect(ctx context.Context, d UpdateRedirectDTO) (*domain.Redirect, error) {
	r, err := s.repo.GetByID(ctx, d.ID); if err != nil { return nil, err }
	if d.Path != nil { p := s.san(*d.Path); if p == "/" { return nil, errors.New("cannot redirect /") }; r.Path = p }
	if d.Target != nil { if err := s.val(*d.Target); err != nil { return nil, err }; r.Target = strings.TrimSpace(*d.Target) }
	if d.StatusCode != nil { c := *d.StatusCode; if c == 301 || c == 302 || c == 307 || c == 308 { r.StatusCode = c } }
	if d.IsEnabled != nil { r.IsEnabled = *d.IsEnabled }
	return r, s.repo.Update(ctx, r)
}

func (s *RedirectService) DeleteRedirect(ctx context.Context, id int64) error { return s.repo.Delete(ctx, id) }
