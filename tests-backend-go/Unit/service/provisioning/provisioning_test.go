package provisioning_test

import (
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/provisioning"
)

func TestDirectAdminProvisioner(t *testing.T) {
	_ = provisioning.NewDirectAdminProvisioner("da.example.com", 2222, "admin", "secret")
	// Test basic instantiation
}

func TestCpanelProvisioner(t *testing.T) {
	cpanel := provisioning.NewCpanelProvisioner(provisioning.CpanelConfig{
		Host:     "cpanel.example.com",
		Username: "root",
		APIToken: "token123",
		Insecure: false,
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

	registry.Register("cpanel", cpanel)
	registry.Register("license", license)

	list := registry.List()
	if len(list) != 2 {
		t.Fatalf("expected 2 provisioners, got %d", len(list))
	}

	p, err := registry.Get("cpanel")
	if err != nil || p == nil {
		t.Fatalf("failed to get cpanel provisioner: %v", err)
	}

	_, err = registry.Get("nonexistent")
	if err == nil {
		t.Fatalf("expected error for nonexistent driver")
	}
}
