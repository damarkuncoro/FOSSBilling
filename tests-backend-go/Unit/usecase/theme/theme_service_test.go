package theme_test

import (
	"context"
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/theme"
)

func TestThemeService_ListAndSelect(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewThemeRepository()
	svc := theme.NewThemeService(repo)

	// 1. List client themes
	list, err := svc.ListThemes(ctx, domain.ThemeTargetClient)
	if err != nil {
		t.Fatalf("Failed to list themes: %v", err)
	}
	if len(list) < 2 {
		t.Errorf("Expected at least 2 client themes, got %d", len(list))
	}

	// 2. Get current client theme
	current, err := svc.GetCurrentTheme(ctx, domain.ThemeTargetClient)
	if err != nil {
		t.Fatalf("Failed to get current theme: %v", err)
	}
	if current.Code != "huraga" {
		t.Errorf("Expected current client theme 'huraga', got '%s'", current.Code)
	}

	// 3. Select alternative theme
	err = svc.SelectTheme(ctx, "antigravity_sleek", domain.ThemeTargetClient)
	if err != nil {
		t.Fatalf("Failed to switch theme: %v", err)
	}

	newCurrent, err := svc.GetCurrentTheme(ctx, domain.ThemeTargetClient)
	if err != nil || newCurrent.Code != "antigravity_sleek" {
		t.Errorf("Expected new current theme 'antigravity_sleek', got '%s'", newCurrent.Code)
	}

	// 4. Update and Get config
	cfg, err := svc.GetConfig(ctx, "antigravity_sleek")
	if err != nil {
		t.Fatalf("Failed to get theme config: %v", err)
	}
	cfg["primary_accent"] = "#a855f7"
	err = svc.UpdateConfig(ctx, "antigravity_sleek", cfg)
	if err != nil {
		t.Fatalf("Failed to update theme config: %v", err)
	}
}
