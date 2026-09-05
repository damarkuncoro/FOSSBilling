package memory

import (
	"context"
	"sync"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type MockCompanyRepository struct {
	mu       sync.RWMutex
	settings *domain.CompanySettings
}

func NewMockCompanyRepository() *MockCompanyRepository {
	return &MockCompanyRepository{
		settings: &domain.CompanySettings{
			ID:             1,
			Name:           "FOSSBilling Cloud Solutions",
			Email:          "billing@fossbilling.org",
			Phone:          "+62 812-3456-7890",
			Address1:       "Jl. Sudirman No. 88, SCBD Area",
			Address2:       "Tower 2, Floor 18",
			City:           "Jakarta Selatan",
			State:          "DKI Jakarta",
			Postcode:       "12190",
			Country:        "ID",
			VatNumber:      "01.234.567.8-901.000",
			LogoURL:        "/branding/logo-light.svg",
			LogoDarkURL:    "/branding/logo-dark.svg",
			FaviconURL:     "/branding/favicon.svg",
			TermsURL:       "https://fossbilling.org/terms",
			EmailSignature: "--\nBest regards,\nThe FOSSBilling Cloud Team\nhttps://fossbilling.org",
			UpdatedAt:      time.Now().UTC(),
		},
	}
}

func (r *MockCompanyRepository) Get(ctx context.Context) (*domain.CompanySettings, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	// Return copy
	clone := *r.settings
	return &clone, nil
}

func (r *MockCompanyRepository) Update(ctx context.Context, s *domain.CompanySettings) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	s.ID = 1
	s.UpdatedAt = time.Now().UTC()
	r.settings = s
	return nil
}
