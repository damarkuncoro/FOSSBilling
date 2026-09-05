package provisioning_test

import (
	"context"
	"strings"
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/service/provisioning"
)

func TestDirectAdminProvisioner(t *testing.T) {
	da := provisioning.NewDirectAdminProvisioner("da.example.com", 2222, "admin", "secret")
	ctx := context.Background()

	acc, err := da.CreateAccount(ctx, provisioning.DirectAdminAccount{
		Domain:  "myclientwebsite.com",
		Package: "Standard",
		Email:   "user@myclientwebsite.com",
	})
	if err != nil {
		t.Fatalf("CreateAccount failed: %v", err)
	}

	if acc.Status != "active" {
		t.Errorf("expected status active, got %s", acc.Status)
	}
	if !strings.HasPrefix(acc.Username, "da") {
		t.Errorf("expected username prefix 'da', got %s", acc.Username)
	}
	if acc.IP != "da.example.com" {
		t.Errorf("expected IP da.example.com, got %s", acc.IP)
	}
}

func TestPleskProvisioner(t *testing.T) {
	plesk := provisioning.NewPleskProvisioner("plesk.example.com", 8443, "api-key-123")
	ctx := context.Background()

	sub, err := plesk.CreateSubscription(ctx, provisioning.PleskSubscription{
		DomainName: "plesk-site.org",
		PlanName:   "Unlimited Plan",
	})
	if err != nil {
		t.Fatalf("CreateSubscription failed: %v", err)
	}

	if sub.Status != "active" {
		t.Errorf("expected status active, got %s", sub.Status)
	}
	if !strings.HasPrefix(sub.Username, "pl_") {
		t.Errorf("expected username prefix 'pl_', got %s", sub.Username)
	}
}
