package company

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type CompanyService interface {
	GetCompany(ctx context.Context) (*domain.CompanySettings, error)
	GetPublicCompany(ctx context.Context) (map[string]any, error)
	UpdateCompany(ctx context.Context, s *domain.CompanySettings) (*domain.CompanySettings, error)
}

type companyService struct {
	repo domain.CompanyRepository; sysRepo domain.SystemRepository
}

func NewCompanyService(r domain.CompanyRepository, sr domain.SystemRepository) CompanyService {
	return &companyService{r, sr}
}

func (s *companyService) GetCompany(ctx context.Context) (*domain.CompanySettings, error) { return s.repo.Get(ctx) }

func (s *companyService) GetPublicCompany(ctx context.Context) (map[string]any, error) {
	c, err := s.repo.Get(ctx); if err != nil { return nil, err }
	bs, _ := s.sysRepo.ListSettings(ctx, "branding")
	b := make(map[string]any); for _, st := range bs { var v any; _ = json.Unmarshal(st.Value, &v); b[st.Key] = v }
	return map[string]any{"name": c.Name, "email": c.Email, "phone": c.Phone, "address_1": c.Address1, "city": c.City, "state": c.State, "country": c.Country, "vat_number": c.VatNumber, "logo_url": c.LogoURL, "branding": b}, nil
}

func (s *companyService) UpdateCompany(ctx context.Context, c *domain.CompanySettings) (*domain.CompanySettings, error) {
	if c.Name == "" || c.Email == "" { return nil, errors.New("name/email required") }
	if err := s.repo.Update(ctx, c); err != nil { return nil, err }; return s.repo.Get(ctx)
}
