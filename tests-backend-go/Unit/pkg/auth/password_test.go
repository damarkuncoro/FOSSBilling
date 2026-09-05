package auth_test

import (
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/auth"
)

func TestHashAndCheckPassword(t *testing.T) {
	password := "SecretP@ssw0rd!123"

	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hash == password {
		t.Fatal("HashPassword returned raw password instead of hash")
	}

	// Correct password check
	if !auth.CheckPassword(password, hash) {
		t.Errorf("CheckPassword(%q, hash) = false; want true", password)
	}

	// Wrong password check
	if auth.CheckPassword("WrongPassword", hash) {
		t.Errorf("CheckPassword(wrong, hash) = true; want false")
	}

	// Empty password check
	if auth.CheckPassword("", hash) {
		t.Errorf("CheckPassword('', hash) = true; want false")
	}
}
