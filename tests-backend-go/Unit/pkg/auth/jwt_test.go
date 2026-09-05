package auth_test

import (
	"testing"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/auth"
)

func TestGenerateAndValidateToken(t *testing.T) {
	secret := "test-secret-key-32-chars-long-abc"
	clientID := int64(42)
	email := "client@example.com"
	role := "client"

	tokenStr, err := auth.GenerateToken(secret, clientID, email, role, 1*time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	claims, err := auth.ValidateToken(secret, tokenStr)
	if err != nil {
		t.Fatalf("ValidateToken returned error: %v", err)
	}

	if claims.ClientID != clientID {
		t.Errorf("claims.ClientID = %d; want %d", claims.ClientID, clientID)
	}
	if claims.Email != email {
		t.Errorf("claims.Email = %s; want %s", claims.Email, email)
	}
	if claims.Role != role {
		t.Errorf("claims.Role = %s; want %s", claims.Role, role)
	}
}

func TestValidateToken_InvalidSecret(t *testing.T) {
	tokenStr, _ := auth.GenerateToken("secret-1", 1, "a@b.com", "client", 1*time.Hour)

	_, err := auth.ValidateToken("wrong-secret-key", tokenStr)
	if err == nil {
		t.Error("ValidateToken with wrong secret should fail, but passed")
	}
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	secret := "secret-key"
	// Expired 10 seconds ago
	tokenStr, err := auth.GenerateToken(secret, 1, "a@b.com", "client", -10*time.Second)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	_, err = auth.ValidateToken(secret, tokenStr)
	if err == nil {
		t.Error("ValidateToken with expired token should fail, but passed")
	}
}
