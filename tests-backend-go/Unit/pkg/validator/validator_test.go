package validator_test

import (
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/validator"
)

func TestValidator_CheckEmail(t *testing.T) {
	v := validator.New()
	v.CheckEmail("email", "valid.user@example.com")
	if !v.IsValid() {
		t.Errorf("Expected valid email, got errors: %v", v)
	}

	v2 := validator.New()
	v2.CheckEmail("email", "invalid-email-format")
	if v2.IsValid() {
		t.Error("Expected error for invalid email format, got valid")
	}

	v3 := validator.New()
	v3.CheckEmail("email", "")
	if v3.IsValid() {
		t.Error("Expected error for empty email, got valid")
	}
}

func TestValidator_CheckRequired(t *testing.T) {
	v := validator.New()
	v.CheckRequired("name", "John")
	if !v.IsValid() {
		t.Errorf("Expected valid, got errors: %v", v)
	}

	v2 := validator.New()
	v2.CheckRequired("name", "   ")
	if v2.IsValid() {
		t.Error("Expected error for whitespace string, got valid")
	}
}

func TestValidator_CheckMinLength(t *testing.T) {
	v := validator.New()
	v.CheckMinLength("password", "123456", 6)
	if !v.IsValid() {
		t.Errorf("Expected valid length, got errors: %v", v)
	}

	v2 := validator.New()
	v2.CheckMinLength("password", "123", 6)
	if v2.IsValid() {
		t.Error("Expected error for short password, got valid")
	}
}

func TestValidator_DomainAndURL(t *testing.T) {
	v := validator.New()
	v.CheckURL("website", "https://fossbilling.org")
	v.CheckIP("server_ip", "192.168.1.1")
	v.CheckSLD("domain_name", "my-hosting-store")
	v.CheckTLD("domain_tld", "com")
	v.CheckSlug("slug", "cloud-vps-pro")
	v.CheckPasswordStrength("password", "P@ssw0rd123!", 8)

	if !v.IsValid() {
		t.Errorf("Expected all to be valid, got: %v", v)
	}

	invalidV := validator.New()
	invalidV.CheckURL("website", "ftp://invalid-url")
	invalidV.CheckIP("server_ip", "999.999.999.999")
	invalidV.CheckSLD("domain_name", "invalid.domain.with.dots")
	invalidV.CheckTLD("domain_tld", "123")
	invalidV.CheckSlug("slug", "Invalid Slug With Spaces!")
	invalidV.CheckPasswordStrength("password", "simple", 8)

	if invalidV.IsValid() {
		t.Error("Expected validation errors for invalid inputs, but passed")
	}
}
