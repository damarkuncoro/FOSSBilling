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
