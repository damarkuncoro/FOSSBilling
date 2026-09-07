package extension_test

import (
	"context"
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/extension"
)

func setupExtensionService() (*extension.ExtensionService, *memory.MockExtensionRepository) {
	repo := memory.NewMockExtensionRepository()
	svc := extension.NewExtensionService(repo)
	return svc, repo
}

func TestExtensionService_ListAndFilter(t *testing.T) {
	ctx := context.Background()
	svc, _ := setupExtensionService()

	// 1. List all
	all, err := svc.ListExtensions(ctx, domain.ExtensionFilter{})
	if err != nil || len(all) == 0 {
		t.Fatalf("Expected extensions list, got err=%v, count=%d", err, len(all))
	}

	// 2. Filter by Type
	gateways, err := svc.ListExtensions(ctx, domain.ExtensionFilter{Type: "gateway"})
	if err != nil {
		t.Fatalf("Filter by type failed: %v", err)
	}
	if len(gateways) != 2 { // midtrans, stripe
		t.Errorf("Expected 2 gateways, got %d", len(gateways))
	}

	// 3. Search filter
	searched, err := svc.ListExtensions(ctx, domain.ExtensionFilter{Search: "midtrans"})
	if err != nil || len(searched) != 1 {
		t.Fatalf("Search failed, got count=%d, err=%v", len(searched), err)
	}
	if searched[0].ID != "midtrans" {
		t.Errorf("Searched ID = %s; want midtrans", searched[0].ID)
	}
}

func TestExtensionService_ActivationAndDeactivation(t *testing.T) {
	ctx := context.Background()
	svc, _ := setupExtensionService()

	// 1. Deactivate extension
	deactivated, err := svc.DeactivateExtension(ctx, "midtrans")
	if err != nil {
		t.Fatalf("Deactivate failed: %v", err)
	}
	if deactivated.Status != domain.ExtensionStatusInactive {
		t.Errorf("Status = %s; want inactive", deactivated.Status)
	}

	// 2. Activate extension
	activated, err := svc.ActivateExtension(ctx, "midtrans")
	if err != nil {
		t.Fatalf("Activate failed: %v", err)
	}
	if activated.Status != domain.ExtensionStatusActive {
		t.Errorf("Status = %s; want active", activated.Status)
	}
}

func TestExtensionService_MarketplaceAndInstall(t *testing.T) {
	ctx := context.Background()
	svc, _ := setupExtensionService()

	// 1. Fetch marketplace catalog
	marketItems, err := svc.FetchMarketplace(ctx, "all")
	if err != nil || len(marketItems) == 0 {
		t.Fatalf("FetchMarketplace failed, got count=%d, err=%v", len(marketItems), err)
	}

	// 2. Filter marketplace by type
	gateways, _ := svc.FetchMarketplace(ctx, "gateway")
	if len(gateways) < 2 {
		t.Errorf("Expected at least 2 marketplace gateways, got %d", len(gateways))
	}

	// 3. Readme
	readme, err := svc.GetMarketplaceReadme(ctx, "tripay")
	if err != nil || readme == "" {
		t.Fatalf("GetMarketplaceReadme failed: %v", err)
	}

	// 4. Install from marketplace
	installed, err := svc.InstallExtension(ctx, "tripay")
	if err != nil {
		t.Fatalf("InstallExtension failed: %v", err)
	}
	if installed.ID != "tripay" || installed.Status != domain.ExtensionStatusActive {
		t.Errorf("Unexpected installed extension: %+v", installed)
	}

	// Verify installed exists in list
	ext, err := svc.GetExtension(ctx, "tripay")
	if err != nil || ext == nil {
		t.Fatalf("Expected tripay to be registered in system")
	}

	// 5. Uninstall extension
	err = svc.UninstallExtension(ctx, "tripay")
	if err != nil {
		t.Fatalf("UninstallExtension failed: %v", err)
	}
	_, err = svc.GetExtension(ctx, "tripay")
	if err == nil {
		t.Error("Expected error fetching uninstalled extension, got nil")
	}
}

func TestExtensionService_Configuration(t *testing.T) {
	ctx := context.Background()
	svc, _ := setupExtensionService()

	// Update configuration
	newCfg := map[string]interface{}{
		"server_key": "midtrans-server-key-12345",
		"is_sandbox": true,
	}
	err := svc.UpdateExtensionConfig(ctx, "midtrans", newCfg)
	if err != nil {
		t.Fatalf("UpdateExtensionConfig failed: %v", err)
	}

	// Retrieve configuration
	cfg, err := svc.GetExtensionConfig(ctx, "midtrans")
	if err != nil {
		t.Fatalf("GetExtensionConfig failed: %v", err)
	}
	if cfg["server_key"] != "midtrans-server-key-12345" || cfg["is_sandbox"] != true {
		t.Errorf("Unexpected config retrieved: %+v", cfg)
	}
}
