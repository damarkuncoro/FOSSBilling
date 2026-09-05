package apikey_test

import (
	"context"
	"strings"
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/apikey"
)

func TestAPIKeyService_Lifecycle(t *testing.T) {
	repo := memory.NewMockAPIKeyRepository()
	service := apikey.NewAPIKeyService(repo)
	ctx := context.Background()

	// 1. Generate Key
	k, err := service.GenerateKey(ctx, 1, "Production Automation Script", 30)
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}
	if !strings.HasPrefix(k.Key, "fb_") {
		t.Errorf("expected key to start with 'fb_', got %s", k.Key)
	}
	if k.ExpiresAt == nil {
		t.Errorf("expected expiration date to be set")
	}

	// 2. List Keys
	list, err := service.ListKeys(ctx, 1)
	if err != nil || len(list) != 1 {
		t.Fatalf("expected 1 key, got %d", len(list))
	}

	// 3. Revoke Key
	err = service.RevokeKey(ctx, k.ID, 1)
	if err != nil {
		t.Fatalf("RevokeKey failed: %v", err)
	}

	// 4. Verify Deleted
	listAfter, _ := service.ListKeys(ctx, 1)
	if len(listAfter) != 0 {
		t.Errorf("expected 0 keys after revoke, got %d", len(listAfter))
	}
}
