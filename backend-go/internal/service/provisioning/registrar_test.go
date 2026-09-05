package provisioning_test

import (
	"context"
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/service/provisioning"
)

func TestMockRegistrarDriver_CheckAvailability(t *testing.T) {
	driver := provisioning.NewMockRegistrarDriver()
	ctx := context.Background()

	// 1. Check taken domain
	res, err := driver.CheckAvailability(ctx, "google.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.IsAvailable {
		t.Errorf("expected google.com to be unavailable")
	}
	if res.Price != 129900 {
		t.Errorf("expected price 129900 for .com, got %d", res.Price)
	}

	// 2. Check available domain
	res2, err := driver.CheckAvailability(ctx, "my-fresh-startup-2026.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res2.IsAvailable {
		t.Errorf("expected new domain to be available")
	}

	// 3. Register domain
	regRes, err := driver.RegisterDomain(ctx, provisioning.DomainRegistrationRequest{
		DomainName: "my-fresh-startup-2026.com",
		Years:      2,
	})
	if err != nil {
		t.Fatalf("failed to register domain: %v", err)
	}
	if regRes.Status != "active" {
		t.Errorf("expected status active, got %s", regRes.Status)
	}

	// 4. Verify domain is now taken
	res3, err := driver.CheckAvailability(ctx, "my-fresh-startup-2026.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res3.IsAvailable {
		t.Errorf("expected registered domain to now be unavailable")
	}
}
