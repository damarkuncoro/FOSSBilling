package cookieconsent

import (
	"context"
	"encoding/json"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type CookieConsentService struct{ repo domain.ExtensionRepository }

func NewCookieConsentService(r domain.ExtensionRepository) *CookieConsentService { return &CookieConsentService{r} }

func (s *CookieConsentService) GetConfig(ctx context.Context) (*domain.CookieConsentConfig, error) {
	def := &domain.CookieConsentConfig{Enabled: true, Message: "This site uses cookies.", ButtonText: "Accept", Position: "bottom", PrivacyPolicyURL: "/privacy-policy", Theme: "classic"}
	ext, err := s.repo.GetByID(ctx, "cookieconsent"); if err != nil || ext == nil || ext.Config == nil { return def, nil }
	b, _ := json.Marshal(ext.Config); var cfg domain.CookieConsentConfig
	if err := json.Unmarshal(b, &cfg); err != nil { return def, nil }; return &cfg, nil
}

func (s *CookieConsentService) UpdateConfig(ctx context.Context, cfg domain.CookieConsentConfig) error {
	m := make(map[string]any); b, _ := json.Marshal(cfg); _ = json.Unmarshal(b, &m)
	ext, err := s.repo.GetByID(ctx, "cookieconsent")
	if err != nil || ext == nil {
		return s.repo.Create(ctx, &domain.Extension{ID: "cookieconsent", Name: "Cookie Consent", Type: domain.ExtensionTypeMod, Version: "1.0.0", Status: domain.ExtensionStatusActive, HasSettings: true, Config: m})
	}
	return s.repo.UpdateConfig(ctx, "cookieconsent", m)
}
