package widget_test

import (
	"context"
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/widget"
)

func TestWidgetService_SlotsAndRegistry(t *testing.T) {
	ctx := context.Background()
	svc := widget.NewWidgetService()

	// 1. Get widgets for client.footer slot
	footerWidgets := svc.GetWidgetsForSlot(ctx, "client.footer")
	if len(footerWidgets) == 0 {
		t.Fatal("Expected at least 1 widget for client.footer slot")
	}

	// 2. Register custom widget with higher priority (lower number)
	svc.RegisterWidget(&domain.Widget{
		ID:        "custom_footer_banner",
		Module:    "custompages",
		Slot:      "client.footer",
		Title:     "Custom Promotion Banner",
		Component: "PromoBannerWidget",
		Priority:  5,
	})

	updatedFooter := svc.GetWidgetsForSlot(ctx, "client.footer")
	if len(updatedFooter) != 2 {
		t.Fatalf("Expected 2 widgets for client.footer, got %d", len(updatedFooter))
	}

	// Verify priority ordering (Priority 5 should come before Priority 100)
	if updatedFooter[0].ID != "custom_footer_banner" {
		t.Errorf("Expected 'custom_footer_banner' first due to priority, got '%s'", updatedFooter[0].ID)
	}

	// 3. Get full registry
	reg := svc.GetRegistry(ctx)
	if len(reg) < 2 {
		t.Errorf("Expected at least 2 slots in registry, got %d", len(reg))
	}

	// 4. Unregister widget
	svc.UnregisterWidget("custom_footer_banner")
	afterRemove := svc.GetWidgetsForSlot(ctx, "client.footer")
	if len(afterRemove) != 1 {
		t.Errorf("Expected 1 widget after removal, got %d", len(afterRemove))
	}
}
