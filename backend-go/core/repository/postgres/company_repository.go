package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CompanyRepository struct {
	pool *pgxpool.Pool
}

func NewCompanyRepository(pool *pgxpool.Pool) *CompanyRepository {
	return &CompanyRepository{pool: pool}
}

func (r *CompanyRepository) Get(ctx context.Context) (*domain.CompanySettings, error) {
	query := `
		SELECT id, name, email, phone, address_1, address_2, city, state, postcode, country, vat_number,
		       logo_url, logo_dark_url, favicon_url, terms_url, email_signature, updated_at
		FROM company_settings
		WHERE id = 1
		LIMIT 1
	`
	var c domain.CompanySettings
	err := r.pool.QueryRow(ctx, query).Scan(
		&c.ID, &c.Name, &c.Email, &c.Phone, &c.Address1, &c.Address2, &c.City, &c.State,
		&c.Postcode, &c.Country, &c.VatNumber, &c.LogoURL, &c.LogoDarkURL, &c.FaviconURL,
		&c.TermsURL, &c.EmailSignature, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Return default settings if table is empty
			return &domain.CompanySettings{
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
			}, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *CompanyRepository) Update(ctx context.Context, s *domain.CompanySettings) error {
	query := `
		INSERT INTO company_settings (id, name, email, phone, address_1, address_2, city, state, postcode, country, vat_number,
		                              logo_url, logo_dark_url, favicon_url, terms_url, email_signature, updated_at)
		VALUES (1, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, NOW())
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			email = EXCLUDED.email,
			phone = EXCLUDED.phone,
			address_1 = EXCLUDED.address_1,
			address_2 = EXCLUDED.address_2,
			city = EXCLUDED.city,
			state = EXCLUDED.state,
			postcode = EXCLUDED.postcode,
			country = EXCLUDED.country,
			vat_number = EXCLUDED.vat_number,
			logo_url = EXCLUDED.logo_url,
			logo_dark_url = EXCLUDED.logo_dark_url,
			favicon_url = EXCLUDED.favicon_url,
			terms_url = EXCLUDED.terms_url,
			email_signature = EXCLUDED.email_signature,
			updated_at = NOW()
	`
	_, err := r.pool.Exec(ctx, query,
		s.Name, s.Email, s.Phone, s.Address1, s.Address2, s.City, s.State, s.Postcode, s.Country, s.VatNumber,
		s.LogoURL, s.LogoDarkURL, s.FaviconURL, s.TermsURL, s.EmailSignature,
	)
	if err != nil {
		return appErrors.ErrInternal
	}
	return nil
}
