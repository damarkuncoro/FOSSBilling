package company

import (
	"context"
	"errors"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type CompanyService interface {
	GetCompany(ctx context.Context) (*domain.CompanySettings, error)
	GetPublicCompany(ctx context.Context) (map[string]interface{}, error)
	UpdateCompany(ctx context.Context, settings *domain.CompanySettings) (*domain.CompanySettings, error)
}

type companyService struct {
	repo domain.CompanyRepository
}

func NewCompanyService(repo domain.CompanyRepository) CompanyService {
	return &companyService{repo: repo}
}

func (s *companyService) GetCompany(ctx context.Context) (*domain.CompanySettings, error) {
	return s.repo.Get(ctx)
}

func (s *companyService) GetPublicCompany(ctx context.Context) (map[string]interface{}, error) {
	c, err := s.repo.Get(ctx)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"name":          c.Name,
		"email":         c.Email,
		"phone":         c.Phone,
		"address_1":     c.Address1,
		"address_2":     c.Address2,
		"city":          c.City,
		"state":         c.State,
		"postcode":      c.Postcode,
		"country":       c.Country,
		"vat_number":    c.VatNumber,
		"logo_url":      c.LogoURL,
		"logo_dark_url": c.LogoDarkURL,
		"favicon_url":   c.FaviconURL,
		"terms_url":     c.TermsURL,
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
