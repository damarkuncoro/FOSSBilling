package domain

import (
	"context"
	"time"
)

type CompanySettings struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	Address1       string    `json:"address_1"`
	Address2       string    `json:"address_2"`
	City           string    `json:"city"`
	State          string    `json:"state"`
	Postcode       string    `json:"postcode"`
	Country        string    `json:"country"`
	VatNumber      string    `json:"vat_number"`
	LogoURL        string    `json:"logo_url"`
	LogoDarkURL    string    `json:"logo_dark_url"`
	FaviconURL     string    `json:"favicon_url"`
	TermsURL       string    `json:"terms_url"`
	EmailSignature string    `json:"email_signature"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CompanyRepository interface {
	Get(ctx context.Context) (*CompanySettings, error)
	Update(ctx context.Context, settings *CompanySettings) error
}
