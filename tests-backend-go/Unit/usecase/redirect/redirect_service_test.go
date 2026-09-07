package redirect_test

import (
	"context"
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/redirect"
)

func TestRedirectService_CRUD(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewRedirectRepository()
	svc := redirect.NewRedirectService(repo)

	// 1. Create redirect rule
	dto := redirect.CreateRedirectDTO{
		Path:       "promo-vps",
		Target:     "/order?product=3",
		StatusCode: 302,
		IsEnabled:  true,
	}

	created, err := svc.CreateRedirect(ctx, dto)
	if err != nil {
		t.Fatalf("Failed to create redirect: %v", err)
	}

	if created.Path != "/promo-vps" {
		t.Errorf("Expected sanitized path '/promo-vps', got '%s'", created.Path)
	}
	if created.Target != "/order?product=3" {
		t.Errorf("Expected target '/order?product=3', got '%s'", created.Target)
	}

	// 2. Duplicate path validation
	_, err = svc.CreateRedirect(ctx, dto)
	if err == nil {
		t.Fatal("Expected error creating duplicate redirect path, got nil")
	}

	// 3. Match redirect
	matched, err := svc.MatchRedirect(ctx, "promo-vps")
	if err != nil {
		t.Fatalf("Failed to match redirect: %v", err)
	}
	if matched.ID != created.ID {
		t.Errorf("Expected matched ID %d, got %d", created.ID, matched.ID)
	}

	// 4. Update redirect
	newTarget := "https://fossbilling.org/pricing"
	newCode := 301
	updated, err := svc.UpdateRedirect(ctx, redirect.UpdateRedirectDTO{
		ID:         created.ID,
		Target:     &newTarget,
		StatusCode: &newCode,
	})
	if err != nil {
		t.Fatalf("Failed to update redirect: %v", err)
	}
	if updated.Target != newTarget || updated.StatusCode != 301 {
		t.Errorf("Updated redirect mismatch: %+v", updated)
	}

	// 5. List redirects
	list, total, err := svc.ListRedirects(ctx, 10, 0)
	if err != nil {
		t.Fatalf("Failed to list redirects: %v", err)
	}
	if total != 1 || len(list) != 1 {
		t.Errorf("Expected total=1, len=1, got total=%d, len=%d", total, len(list))
	}

	// 6. Delete redirect
	err = svc.DeleteRedirect(ctx, created.ID)
	if err != nil {
		t.Fatalf("Failed to delete redirect: %v", err)
	}

	_, err = svc.GetRedirect(ctx, created.ID)
	if err == nil {
		t.Fatal("Expected error getting deleted redirect, got nil")
	}
}
