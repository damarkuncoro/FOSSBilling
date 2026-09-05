package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/fossbilling/backend-go/internal/domain"
	"github.com/fossbilling/backend-go/internal/repository/memory"
	appErrors "github.com/fossbilling/backend-go/pkg/errors"
	"github.com/fossbilling/backend-go/pkg/decimal"
)

func setupAuthUsecase() (*AuthUsecase, *memory.MockClientRepository) {
	repo := memory.NewMockClientRepository()
	uc := NewAuthUsecase(repo, "test-jwt-secret-key-32-bytes-long")
	return uc, repo
}

func TestAuthUsecase_Register(t *testing.T) {
	ctx := context.Background()
	uc, _ := setupAuthUsecase()

	req := RegisterDTO{
		Email:     "john.doe@example.com",
		Password:  "StrongPassword123!",
		FirstName: "John",
		LastName:  "Doe",
		Company:   "Acme Corp",
		Country:   "US",
		Currency:  "USD",
	}

	// 1. Successful Registration
	res, vErrs, err := uc.Register(ctx, req)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	if vErrs != nil && !vErrs.IsValid() {
		t.Fatalf("Expected valid, got validation errors: %v", vErrs)
	}
	if res.Token == "" {
		t.Error("Expected JWT token, got empty string")
	}
	if res.Client.Email != req.Email {
		t.Errorf("Client.Email = %s; want %s", res.Client.Email, req.Email)
	}
	if res.Client.FirstName != req.FirstName {
		t.Errorf("Client.FirstName = %s; want %s", res.Client.FirstName, req.FirstName)
	}

	// 2. Duplicate Registration should fail
	_, vErrs, err = uc.Register(ctx, req)
	if err == nil {
		t.Error("Expected error on duplicate email registration, got nil")
	}
	if !errors.Is(err, appErrors.ErrDuplicateEntry) {
		t.Errorf("Expected ErrDuplicateEntry, got %v", err)
	}

	// 3. Validation Failure (invalid email, short password)
	invalidReq := RegisterDTO{
		Email:    "not-an-email",
		Password: "123",
	}
	_, vErrs, err = uc.Register(ctx, invalidReq)
	if err == nil {
		t.Error("Expected validation error, got nil")
	}
	if vErrs == nil || vErrs.IsValid() {
		t.Error("Expected validation errors map to be populated")
	}
}

func TestAuthUsecase_Login(t *testing.T) {
	ctx := context.Background()
	uc, repo := setupAuthUsecase()

	// Register user first
	regReq := RegisterDTO{
		Email:     "jane@example.com",
		Password:  "CorrectPassword123",
		FirstName: "Jane",
		LastName:  "Smith",
	}
	_, _, err := uc.Register(ctx, regReq)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// 1. Successful Login
	loginRes, err := uc.Login(ctx, LoginDTO{
		Email:    "jane@example.com",
		Password: "CorrectPassword123",
	})
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if loginRes.Token == "" {
		t.Error("Expected JWT token on successful login")
	}
	if loginRes.Client.Email != "jane@example.com" {
		t.Errorf("Login Email = %s; want jane@example.com", loginRes.Client.Email)
	}

	// 2. Wrong Password
	_, err = uc.Login(ctx, LoginDTO{
		Email:    "jane@example.com",
		Password: "WrongPassword!",
	})
	if !errors.Is(err, appErrors.ErrUnauthorized) {
		t.Errorf("Expected ErrUnauthorized for wrong password, got: %v", err)
	}

	// 3. Non-existent User
	_, err = uc.Login(ctx, LoginDTO{
		Email:    "nobody@example.com",
		Password: "SomePassword",
	})
	if !errors.Is(err, appErrors.ErrUnauthorized) {
		t.Errorf("Expected ErrUnauthorized for missing user, got: %v", err)
	}

	// 4. Inactive / Suspended User
	client, _ := repo.GetByEmail(ctx, "jane@example.com")
	client.Status = domain.ClientStatusSuspended
	_ = repo.Update(ctx, client)

	_, err = uc.Login(ctx, LoginDTO{
		Email:    "jane@example.com",
		Password: "CorrectPassword123",
	})
	if err == nil {
		t.Error("Expected error when logging in with suspended account, got nil")
	}
}

func TestAuthUsecase_Profile(t *testing.T) {
	ctx := context.Background()
	uc, repo := setupAuthUsecase()

	regRes, _, err := uc.Register(ctx, RegisterDTO{
		Email:     "profile.test@example.com",
		Password:  "password123",
		FirstName: "Alice",
		LastName:  "Wonder",
	})
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Add some balance to client
	_ = repo.AddBalanceTransaction(ctx, &domain.ClientBalance{
		ClientID: regRes.Client.ID,
		Type:     domain.BalanceTypeCredit,
		Amount:   decimal.FromFloat(50.00),
	})

	// 1. Get Profile
	profile, err := uc.GetProfile(ctx, regRes.Client.ID)
	if err != nil {
		t.Fatalf("GetProfile failed: %v", err)
	}
	if profile.Client.FirstName != "Alice" {
		t.Errorf("FirstName = %s; want Alice", profile.Client.FirstName)
	}
	if profile.Balance.String() != "50.00" {
		t.Errorf("Balance = %s; want 50.00", profile.Balance.String())
	}

	// 2. Update Profile
	updateRes, err := uc.UpdateProfile(ctx, regRes.Client.ID, UpdateProfileDTO{
		FirstName: "Alicia",
		LastName:  "Keys",
		City:      "Jakarta",
		Country:   "ID",
		Phone:     "08123456789",
	})
	if err != nil {
		t.Fatalf("UpdateProfile failed: %v", err)
	}
	if updateRes.Client.FirstName != "Alicia" {
		t.Errorf("Updated FirstName = %s; want Alicia", updateRes.Client.FirstName)
	}
	if updateRes.Client.City != "Jakarta" {
		t.Errorf("Updated City = %s; want Jakarta", updateRes.Client.City)
	}
}
