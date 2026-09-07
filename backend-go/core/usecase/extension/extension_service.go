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
	ErrCoreExtensionImmutable = errors.New("core extensions cannot be deactivated or uninstalled")
	ErrExtensionNotFound      = errors.New("extension not found in system or marketplace")
)

type ExtensionService struct {
	repo domain.ExtensionRepository
}

func NewExtensionService(repo domain.ExtensionRepository) *ExtensionService {
	return &ExtensionService{repo: repo}
}

func (s *ExtensionService) ListExtensions(ctx context.Context, filter domain.ExtensionFilter) ([]*domain.Extension, error) {
	return s.repo.List(ctx, filter)
}

func (s *ExtensionService) GetExtension(ctx context.Context, id string) (*domain.Extension, error) {
	return s.repo.GetByID(ctx, strings.TrimSpace(id))
}

func (s *ExtensionService) ActivateExtension(ctx context.Context, id string) (*domain.Extension, error) {
	ext, err := s.repo.GetByID(ctx, strings.TrimSpace(id))
	if err != nil {
		return nil, err
	}

	if ext.Status == domain.ExtensionStatusCore {
		return ext, nil
	}

	ext.Status = domain.ExtensionStatusActive
	if err := s.repo.Update(ctx, ext); err != nil {
		return nil, err
	}
	return ext, nil
}

func (s *ExtensionService) DeactivateExtension(ctx context.Context, id string) (*domain.Extension, error) {
	ext, err := s.repo.GetByID(ctx, strings.TrimSpace(id))
	if err != nil {
		return nil, err
	}

	if ext.Status == domain.ExtensionStatusCore {
		return nil, ErrCoreExtensionImmutable
	}

	ext.Status = domain.ExtensionStatusInactive
	if err := s.repo.Update(ctx, ext); err != nil {
		return nil, err
	}
	return ext, nil
}

func (s *ExtensionService) InstallExtension(ctx context.Context, marketplaceID string) (*domain.Extension, error) {
	// Look up marketplace item
	marketplaceList := s.getMarketplaceCatalog()
	var found *domain.MarketplaceExtension
	for _, m := range marketplaceList {
		if m.ID == marketplaceID {
			found = m
			break
		}
	}

	if found == nil {
		return nil, ErrExtensionNotFound
	}

	// Check if already installed
	existing, _ := s.repo.GetByID(ctx, marketplaceID)
	if existing != nil {
		existing.Status = domain.ExtensionStatusActive
		_ = s.repo.Update(ctx, existing)
		return existing, nil
	}

	newExt := &domain.Extension{
		ID:          found.ID,
		Name:        found.Name,
		Type:        found.Type,
		Version:     found.Version,
		Description: found.Description,
		Author:      found.Author,
		Icon:        found.IconURL,
		Status:      domain.ExtensionStatusActive,
		HasSettings: true,
		Config:      make(map[string]interface{}),
		Manifest: map[string]interface{}{
			"download_url": found.DownloadURL,
			"rating":       found.Rating,
		},
	}

	if err := s.repo.Create(ctx, newExt); err != nil {
		return nil, err
	}
	return newExt, nil
}

func (s *ExtensionService) UninstallExtension(ctx context.Context, id string) error {
	ext, err := s.repo.GetByID(ctx, strings.TrimSpace(id))
	if err != nil {
		return err
	}

	if ext.Status == domain.ExtensionStatusCore {
		return ErrCoreExtensionImmutable
	}

	return s.repo.Delete(ctx, strings.TrimSpace(id))
}

func (s *ExtensionService) GetExtensionConfig(ctx context.Context, id string) (map[string]interface{}, error) {
	return s.repo.GetConfig(ctx, strings.TrimSpace(id))
}

func (s *ExtensionService) UpdateExtensionConfig(ctx context.Context, id string, config map[string]interface{}) error {
	return s.repo.UpdateConfig(ctx, strings.TrimSpace(id), config)
}

func (s *ExtensionService) FetchMarketplace(ctx context.Context, extType string) ([]*domain.MarketplaceExtension, error) {
	catalog := s.getMarketplaceCatalog()
	if extType == "" || extType == "all" {
		return catalog, nil
	}

	var filtered []*domain.MarketplaceExtension
	for _, item := range catalog {
		if string(item.Type) == extType {
			filtered = append(filtered, item)
		}
	}
	return filtered, nil
}

func (s *ExtensionService) GetMarketplaceReadme(ctx context.Context, id string) (string, error) {
	catalog := s.getMarketplaceCatalog()
	for _, item := range catalog {
		if item.ID == id {
			if item.Readme != "" {
				return item.Readme, nil
			}
			return fmt.Sprintf("# %s\n\n%s\n\n**Author:** %s\n**Version:** %s", item.Name, item.Description, item.Author, item.Version), nil
		}
	}
	return "", appErrors.ErrNotFound
}

func (s *ExtensionService) getMarketplaceCatalog() []*domain.MarketplaceExtension {
	return []*domain.MarketplaceExtension{
		{
			ID:          "tripay",
			Name:        "TriPay Indonesia Payment Gateway",
			Type:        domain.ExtensionTypeGateway,
			Version:     "1.4.0",
			Description: "Automated payment channels including QRIS, BCA, Mandiri, BRI, Alfamart, Indomaret.",
			Author:      "TriPay Community",
			IconURL:     "https://raw.githubusercontent.com/FOSSBilling/Extensions/main/tripay/icon.svg",
			DownloadURL: "https://extensions.fossbilling.org/download/tripay-1.4.0.zip",
			Rating:      4.9,
			Downloads:   1420,
			Readme:      "# TriPay Payment Gateway\nOfficial TriPay integration for FOSSBilling providing instantaneous QRIS & VA reconciliation.",
		},
		{
			ID:          "xendit",
			Name:        "Xendit Global & SEA Payments",
			Type:        domain.ExtensionTypeGateway,
			Version:     "2.0.1",
			Description: "Accept credit cards, e-wallets (OVO, DANA, ShopeePay), and recurring subscriptions.",
			Author:      "Xendit Integrators",
			IconURL:     "https://raw.githubusercontent.com/FOSSBilling/Extensions/main/xendit/icon.svg",
			DownloadURL: "https://extensions.fossbilling.org/download/xendit-2.0.1.zip",
			Rating:      4.8,
			Downloads:   2150,
		},
		{
			ID:          "cyberpanel",
			Name:        "CyberPanel Server Provisioning",
			Type:        domain.ExtensionTypeService,
			Version:     "1.1.0",
			Description: "Direct automated provisioning and management for LiteSpeed CyberPanel web servers.",
			Author:      "CyberPanel Devs",
			IconURL:     "https://raw.githubusercontent.com/FOSSBilling/Extensions/main/cyberpanel/icon.svg",
			DownloadURL: "https://extensions.fossbilling.org/download/cyberpanel-1.1.0.zip",
			Rating:      4.7,
			Downloads:   980,
		},
		{
			ID:          "telegram_notif",
			Name:        "Telegram Bot Staff & Client Alerts",
			Type:        domain.ExtensionTypePlugin,
			Version:     "1.5.0",
			Description: "Instant notifications for paid invoices, new support tickets, and service provisioning events.",
			Author:      "FOSSBilling Community",
			IconURL:     "https://raw.githubusercontent.com/FOSSBilling/Extensions/main/telegram/icon.svg",
			DownloadURL: "https://extensions.fossbilling.org/download/telegram-1.5.0.zip",
			Rating:      5.0,
			Downloads:   3890,
		},
		{
			ID:          "dark_theme",
			Name:        "Onyx Dark Client Theme",
			Type:        domain.ExtensionTypeTheme,
			Version:     "1.0.2",
			Description: "Ultra-sleek OLED dark mode theme with glassmorphic accents.",
			Author:      "ThemeMasters",
			IconURL:     "https://raw.githubusercontent.com/FOSSBilling/Extensions/main/onyx/icon.svg",
			DownloadURL: "https://extensions.fossbilling.org/download/onyx-1.0.2.zip",
			Rating:      4.9,
			Downloads:   4200,
		},
	}
}
