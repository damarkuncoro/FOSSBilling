package auth

import "testing"

func TestHashAndCheckPassword(t *testing.T) {
	password := "SecretP@ssw0rd!123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hash == password {
		t.Fatal("HashPassword returned raw password instead of hash")
	}

	// Correct password check
	if !CheckPassword(password, hash) {
		t.Errorf("CheckPassword(%q, hash) = false; want true", password)
	}

	// Wrong password check
	if CheckPassword("WrongPassword", hash) {
		t.Errorf("CheckPassword(wrong, hash) = true; want false")
	}

	// Empty password check
	if CheckPassword("", hash) {
		t.Errorf("CheckPassword('', hash) = true; want false")
	}
}
