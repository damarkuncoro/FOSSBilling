package cookieconsent

import (
	"context"
	"encoding/json"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type CookieConsentService struct {
	extensionRepo domain.ExtensionRepository
}

func NewCookieConsentService(extensionRepo domain.ExtensionRepository) *CookieConsentService {
	return &CookieConsentService{
		extensionRepo: extensionRepo,
	}
}

func (s *CookieConsentService) GetConfig(ctx context.Context) (*domain.CookieConsentConfig, error) {
	ext, err := s.extensionRepo.GetByID(ctx, "cookieconsent")
	if err != nil || ext == nil || ext.Config == nil {
		// Return default
		return &domain.CookieConsentConfig{
			Enabled:          true,
			Message:          "This website uses cookies. By continuing to use this website, you consent to our use of these cookies.",
			ButtonText:       "Accept",
			Position:         "bottom",
			PrivacyPolicyURL: "/privacy-policy",
			Theme:            "classic",
		}, nil
	}

	cfgBytes, _ := json.Marshal(ext.Config)
	var cfg domain.CookieConsentConfig
	if err := json.Unmarshal(cfgBytes, &cfg); err != nil {
		return &domain.CookieConsentConfig{
			Enabled:          true,
			Message:          "This website uses cookies. By continuing to use this website, you consent to our use of these cookies.",
			ButtonText:       "Accept",
			Position:         "bottom",
			PrivacyPolicyURL: "/privacy-policy",
			Theme:            "classic",
		}, nil
	}

	return &cfg, nil
}

func (s *CookieConsentService) UpdateConfig(ctx context.Context, cfg domain.CookieConsentConfig) error {
	rawMap := make(map[string]interface{})
	cfgBytes, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	_ = json.Unmarshal(cfgBytes, &rawMap)

	ext, err := s.extensionRepo.GetByID(ctx, "cookieconsent")
	if err != nil || ext == nil {
		ext = &domain.Extension{
			ID:          "cookieconsent",
			Name:        "Cookie Consent Banner",
			Type:        domain.ExtensionTypeMod,
			Version:     "1.0.0",
			Description: "GDPR & ePrivacy compliant cookie consent banner",
			Author:      "FOSSBilling Team",
			Status:      domain.ExtensionStatusActive,
			HasSettings: true,
			Config:      rawMap,
		}
		return s.extensionRepo.Create(ctx, ext)
	}

	ext.Config = rawMap
	return s.extensionRepo.UpdateConfig(ctx, "cookieconsent", rawMap)
}
