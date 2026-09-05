package company_test

import (
	"context"
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/company"
)

func TestCompanyService_GetAndUpdate(t *testing.T) {
	repo := memory.NewMockCompanyRepository()
	service := company.NewCompanyService(repo)
	ctx := context.Background()

	// 1. Get Company Settings
	settings, err := service.GetCompany(ctx)
	if err != nil {
		t.Fatalf("GetCompany failed: %v", err)
	}
	if settings.Name != "FOSSBilling Cloud Solutions" {
		t.Errorf("expected default name 'FOSSBilling Cloud Solutions', got %s", settings.Name)
	}

	// 2. Get Public Company Info
	pub, err := service.GetPublicCompany(ctx)
	if err != nil {
		t.Fatalf("GetPublicCompany failed: %v", err)
	}
	if pub["name"] != "FOSSBilling Cloud Solutions" {
		t.Errorf("expected public name 'FOSSBilling Cloud Solutions', got %v", pub["name"])
	}
	if _, exists := pub["id"]; exists {
		t.Errorf("public company info should not expose internal id")
	}

	// 3. Update Company Settings
	updated, err := service.UpdateCompany(ctx, &domain.CompanySettings{
		Name:           "Acme Cloud Hosting",
		Email:          "contact@acmehosting.com",
		Phone:          "+1 800-555-0199",
		Address1:       "100 Tech Blvd",
		City:           "San Francisco",
		State:          "CA",
		Postcode:       "94107",
		Country:        "US",
		VatNumber:      "US123456789",
		LogoURL:        "/branding/acme-light.svg",
		LogoDarkURL:    "/branding/acme-dark.svg",
		FaviconURL:     "/branding/favicon.ico",
		TermsURL:       "https://acmehosting.com/terms",
		EmailSignature: "--\nAcme Hosting Support Team",
	})
	if err != nil {
		t.Fatalf("UpdateCompany failed: %v", err)
	}
	if updated.Name != "Acme Cloud Hosting" {
		t.Errorf("expected updated name 'Acme Cloud Hosting', got %s", updated.Name)
	}

	// 4. Update with missing name or email should fail validation
	_, err = service.UpdateCompany(ctx, &domain.CompanySettings{
		Name:  "",
		Email: "test@example.com",
	})
	if err == nil {
		t.Errorf("expected error when updating with empty name, got nil")
	}

	_, err = service.UpdateCompany(ctx, &domain.CompanySettings{
		Name:  "Valid Name",
		Email: "",
	})
	if err == nil {
		t.Errorf("expected error when updating with empty email, got nil")
	}
}
