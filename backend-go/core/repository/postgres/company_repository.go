package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CompanyRepository struct{ pool *pgxpool.Pool }

func NewCompanyRepository(p *pgxpool.Pool) *CompanyRepository { return &CompanyRepository{p} }

func (r *CompanyRepository) Get(ctx context.Context) (*domain.CompanySettings, error) {
	q := `SELECT id, name, email, phone, address_1, address_2, city, state, postcode, country, vat_number, logo_url, logo_dark_url, favicon_url, terms_url, email_signature, updated_at FROM company_settings WHERE id = 1 LIMIT 1`
	var c domain.CompanySettings
	err := r.pool.QueryRow(ctx, q).Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Address1, &c.Address2, &c.City, &c.State, &c.Postcode, &c.Country, &c.VatNumber, &c.LogoURL, &c.LogoDarkURL, &c.FaviconURL, &c.TermsURL, &c.EmailSignature, &c.UpdatedAt)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return &domain.CompanySettings{ID: 1, Name: "FOSSBilling Enterprise", Email: "billing@example.com", Phone: "", Address1: "Corporate Headquarters", City: "", State: "", Postcode: "", Country: "US", VatNumber: "", LogoURL: "/branding/logo-light.svg", LogoDarkURL: "/branding/logo-dark.svg", FaviconURL: "/branding/favicon.svg", TermsURL: "", EmailSignature: "--\nBest regards,\nThe Support Team", UpdatedAt: time.Now().UTC()}, nil
	}
	return &c, err
}

func (r *CompanyRepository) Update(ctx context.Context, s *domain.CompanySettings) error {
	q := `INSERT INTO company_settings (id, name, email, phone, address_1, address_2, city, state, postcode, country, vat_number, logo_url, logo_dark_url, favicon_url, terms_url, email_signature, updated_at) VALUES (1, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, NOW()) ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, email = EXCLUDED.email, phone = EXCLUDED.phone, address_1 = EXCLUDED.address_1, address_2 = EXCLUDED.address_2, city = EXCLUDED.city, state = EXCLUDED.state, postcode = EXCLUDED.postcode, country = EXCLUDED.country, vat_number = EXCLUDED.vat_number, logo_url = EXCLUDED.logo_url, logo_dark_url = EXCLUDED.logo_dark_url, favicon_url = EXCLUDED.favicon_url, terms_url = EXCLUDED.terms_url, email_signature = EXCLUDED.email_signature, updated_at = NOW()`
	_, err := r.pool.Exec(ctx, q, s.Name, s.Email, s.Phone, s.Address1, s.Address2, s.City, s.State, s.Postcode, s.Country, s.VatNumber, s.LogoURL, s.LogoDarkURL, s.FaviconURL, s.TermsURL, s.EmailSignature); return err
}
