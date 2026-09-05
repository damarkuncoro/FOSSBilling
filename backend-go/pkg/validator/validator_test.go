package validator

import "testing"

func TestValidator_CheckEmail(t *testing.T) {
	v := New()
	v.CheckEmail("email", "valid.user@example.com")
	if !v.IsValid() {
		t.Errorf("Expected valid email, got errors: %v", v)
	}

	v2 := New()
	v2.CheckEmail("email", "invalid-email-format")
	if v2.IsValid() {
		t.Error("Expected error for invalid email format, got valid")
	}

	v3 := New()
	v3.CheckEmail("email", "")
	if v3.IsValid() {
		t.Error("Expected error for empty email, got valid")
	}
}

func TestValidator_CheckRequired(t *testing.T) {
	v := New()
	v.CheckRequired("name", "John")
	if !v.IsValid() {
		t.Errorf("Expected valid, got errors: %v", v)
	}

	v2 := New()
	v2.CheckRequired("name", "   ")
	if v2.IsValid() {
		t.Error("Expected error for whitespace string, got valid")
	}
}

func TestValidator_CheckMinLength(t *testing.T) {
	v := New()
	v.CheckMinLength("password", "123456", 6)
	if !v.IsValid() {
		t.Errorf("Expected valid length, got errors: %v", v)
	}

	v2 := New()
	v2.CheckMinLength("password", "123", 6)
	if v2.IsValid() {
		t.Error("Expected error for short password, got valid")
	}
}
