package provisioning_test

import (
	"context"
	"strings"
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/provisioning"
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

func TestCpanelProvisioner(t *testing.T) {
	cpanel := provisioning.NewCpanelProvisioner(provisioning.CpanelConfig{
		Host:     "cpanel.example.com",
		Username: "root",
		APIToken: "token123",
		UseSSL:   true,
	})

	user, pass := cpanel.GenerateAccountCredentials("my-awesome-domain.com")
	if user == "" || pass == "" {
		t.Fatalf("failed to generate cPanel credentials")
	}
}

func TestProvisionerRegistry(t *testing.T) {
	registry := provisioning.NewProvisionerRegistry()

	cpanel := provisioning.NewCpanelProvisioner(provisioning.CpanelConfig{Host: "cpanel.example.com"})
	license := provisioning.NewLicenseProvisioner("SECRET_SIGNING_SALT")

	registry.Register(cpanel)
	registry.Register(license)

	list := registry.List()
	if len(list) != 2 {
		t.Fatalf("expected 2 provisioners, got %d", len(list))
	}

	p, err := registry.Get(cpanel.Type())
	if err != nil || p == nil {
		t.Fatalf("failed to get cpanel provisioner: %v", err)
	}

	_, err = registry.Get("nonexistent_type")
	if err == nil {
		t.Fatalf("expected error for nonexistent product type")
	}
}

