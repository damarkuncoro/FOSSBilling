package extension

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
)

var (
	ErrImmutable = errors.New("core immutable")
	ErrNotFound  = errors.New("not found")
)

type ExtensionService struct{ repo domain.ExtensionRepository }

func NewExtensionService(r domain.ExtensionRepository) *ExtensionService { return &ExtensionService{r} }

func (s *ExtensionService) ListExtensions(ctx context.Context, f domain.ExtensionFilter) ([]*domain.Extension, error) { return s.repo.List(ctx, f) }
func (s *ExtensionService) GetExtension(ctx context.Context, id string) (*domain.Extension, error) { return s.repo.GetByID(ctx, strings.TrimSpace(id)) }

func (s *ExtensionService) setStatus(ctx context.Context, id string, st domain.ExtensionStatus, checkCore bool) (*domain.Extension, error) {
	e, err := s.repo.GetByID(ctx, strings.TrimSpace(id)); if err != nil { return nil, err }
	if checkCore && e.Status == domain.ExtensionStatusCore { return nil, ErrImmutable }
	e.Status = st; return e, s.repo.Update(ctx, e)
}

func (s *ExtensionService) ActivateExtension(ctx context.Context, id string) (*domain.Extension, error) { return s.setStatus(ctx, id, domain.ExtensionStatusActive, false) }
func (s *ExtensionService) DeactivateExtension(ctx context.Context, id string) (*domain.Extension, error) { return s.setStatus(ctx, id, domain.ExtensionStatusInactive, true) }

func (s *ExtensionService) InstallExtension(ctx context.Context, mid string) (*domain.Extension, error) {
	var f *domain.MarketplaceExtension; for _, m := range s.getMarketplaceCatalog() { if m.ID == mid { f = m; break } }
	if f == nil { return nil, ErrNotFound }
	if ex, _ := s.repo.GetByID(ctx, mid); ex != nil { ex.Status = domain.ExtensionStatusActive; _ = s.repo.Update(ctx, ex); return ex, nil }
	ne := &domain.Extension{ID: f.ID, Name: f.Name, Type: f.Type, Version: f.Version, Description: f.Description, Author: f.Author, Icon: f.IconURL, Status: domain.ExtensionStatusActive, HasSettings: true, Config: map[string]any{}, Manifest: map[string]any{"download_url": f.DownloadURL, "rating": f.Rating}}
	return ne, s.repo.Create(ctx, ne)
}

func (s *ExtensionService) UninstallExtension(ctx context.Context, id string) error {
	e, err := s.repo.GetByID(ctx, strings.TrimSpace(id)); if err != nil { return err }
	if e.Status == domain.ExtensionStatusCore { return ErrImmutable }; return s.repo.Delete(ctx, strings.TrimSpace(id))
}

func (s *ExtensionService) GetExtensionConfig(ctx context.Context, id string) (map[string]any, error) { return s.repo.GetConfig(ctx, strings.TrimSpace(id)) }
func (s *ExtensionService) UpdateExtensionConfig(ctx context.Context, id string, cfg map[string]any) error { return s.repo.UpdateConfig(ctx, strings.TrimSpace(id), cfg) }

func (s *ExtensionService) FetchMarketplace(ctx context.Context, t string) ([]*domain.MarketplaceExtension, error) {
	cat := s.getMarketplaceCatalog(); if t == "" || t == "all" { return cat, nil }
	var res []*domain.MarketplaceExtension; for _, i := range cat { if string(i.Type) == t { res = append(res, i) } }; return res, nil
}

func (s *ExtensionService) GetMarketplaceReadme(ctx context.Context, id string) (string, error) {
	for _, i := range s.getMarketplaceCatalog() { if i.ID == id { if i.Readme != "" { return i.Readme, nil }; return fmt.Sprintf("# %s\n\n%s", i.Name, i.Description), nil } }
	return "", appErrors.ErrNotFound
}

func (s *ExtensionService) getMarketplaceCatalog() []*domain.MarketplaceExtension {
	return []*domain.MarketplaceExtension{
		{ID: "tripay", Name: "TriPay", Type: domain.ExtensionTypeGateway, Version: "1.4.0", Description: "QRIS, VA payments.", Author: "TriPay", DownloadURL: "https://x.org/dl.zip", Rating: 4.9},
		{ID: "stripe_pro", Name: "Stripe Pro", Type: domain.ExtensionTypeGateway, Version: "2.1.0", Description: "Enhanced Stripe.", Author: "FOSSBilling", DownloadURL: "https://x.org/sp.zip", Rating: 4.8},
		{ID: "telegram_notif", Name: "Telegram Notif", Type: domain.ExtensionTypePlugin, Version: "1.5.0", Description: "Staff alerts.", Author: "FOSSBilling", DownloadURL: "https://x.org/tg.zip", Rating: 5.0},
	}
}
