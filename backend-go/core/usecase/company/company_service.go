package company

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type CompanyService interface {
	GetCompany(ctx context.Context) (*domain.CompanySettings, error)
	GetPublicCompany(ctx context.Context) (map[string]interface{}, error)
	UpdateCompany(ctx context.Context, settings *domain.CompanySettings) (*domain.CompanySettings, error)
}

type companyService struct {
	repo       domain.CompanyRepository
	systemRepo domain.SystemRepository
}

func NewCompanyService(repo domain.CompanyRepository, systemRepo domain.SystemRepository) CompanyService {
	return &companyService{repo: repo, systemRepo: systemRepo}
}

func (s *companyService) GetCompany(ctx context.Context) (*domain.CompanySettings, error) {
	return s.repo.Get(ctx)
}

func (s *companyService) GetPublicCompany(ctx context.Context) (map[string]interface{}, error) {
	c, err := s.repo.Get(ctx)
	if err != nil {
		return nil, err
	}

	// Fetch branding
	brandingSettings, _ := s.systemRepo.ListSettings(ctx, "branding")
	branding := make(map[string]interface{})
	for _, setting := range brandingSettings {
		var val interface{}
		_ = json.Unmarshal(setting.Value, &val)
		branding[setting.Key] = val
	}

	return map[string]interface{}{
		"name":          c.Name,
		"email":         c.Email,
		"phone":         c.Phone,
		"address_1":     c.Address1,
		"city":          c.City,
		"state":         c.State,
		"country":       c.Country,
		"vat_number":    c.VatNumber,
		"logo_url":      c.LogoURL,
		"branding":      branding,
	}, nil
}

func (s *companyService) UpdateCompany(ctx context.Context, settings *domain.CompanySettings) (*domain.CompanySettings, error) {
	if settings.Name == "" {
		return nil, errors.New("company name is required")
	}
	if settings.Email == "" {
		return nil, errors.New("company email is required")
	}

	if err := s.repo.Update(ctx, settings); err != nil {
		return nil, err
	}

	return s.repo.Get(ctx)
}
