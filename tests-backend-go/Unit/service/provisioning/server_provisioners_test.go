package provisioning_test

import (
	"context"
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/provisioning"
)

func TestProvisionerFactory_AllServerDrivers(t *testing.T) {
	factory := provisioning.NewProvisionerFactory()

	tests := []struct {
		name        string
		serverType  string
		expectedTyp domain.ProductType
	}{
		{"cPanel WHM", "cpanel", domain.ProductTypeHosting},
		{"DirectAdmin", "directadmin", domain.ProductTypeHosting},
		{"Plesk", "plesk", domain.ProductTypeHosting},
		{"HestiaCP", "hestia", domain.ProductTypeHosting},
		{"CentOS Web Panel", "cwp", domain.ProductTypeHosting},
		{"CyberPanel", "cyberpanel", domain.ProductTypeHosting},
		{"Custom Server Webhook", "custom", domain.ProductTypeCustom},
		{"License Key Generator", "license", domain.ProductTypeLicense},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := factory.CreateProvisioner(provisioning.ServerConfig{
				Type:     tt.serverType,
				Host:     "127.0.0.1",
				Port:     8083,
				Username: "admin",
				Password: "SecretPassword123!",
				APIToken: "token_abc123",
				UseSSL:   false,
			})
			if err != nil {
				t.Fatalf("Failed to create provisioner for %s: %v", tt.serverType, err)
			}
			if p.Type() != tt.expectedTyp {
				t.Errorf("Type = %s; want %s", p.Type(), tt.expectedTyp)
			}
		})
	}
}

func TestHestiaProvisioner_UsernameGeneration(t *testing.T) {
	p := provisioning.NewHestiaProvisioner(provisioning.HestiaConfig{
		Host: "localhost",
	})

	username := p.GenerateUsername("example-blog.com")
	if len(username) == 0 {
		t.Fatal("Expected non-empty username")
	}
	// Username must not start with a digit in HestiaCP
	if username[0] >= '0' && username[0] <= '9' {
		t.Errorf("Username %s starts with digit", username)
	}

	// Starting with digit domain
	usernameNum := p.GenerateUsername("123domain.net")
	if usernameNum[0] >= '0' && usernameNum[0] <= '9' {
		t.Errorf("Expected leading letter for numeric domain, got %s", usernameNum)
	}
}

func TestCWPProvisioner_UsernameGeneration(t *testing.T) {
	p := provisioning.NewCWPProvisioner(provisioning.CWPConfig{
		Host: "localhost",
	})

	username := p.GenerateUsername("mycompany.id")
	if len(username) == 0 || len(username) > 8 {
		t.Errorf("Unexpected CWP username length: %s", username)
	}
}

func TestCyberPanelProvisioner_UsernameGeneration(t *testing.T) {
	p := provisioning.NewCyberPanelProvisioner(provisioning.CyberPanelConfig{
		Host: "localhost",
	})

	username := p.GenerateUsername("test-website.com")
	if len(username) == 0 || len(username) > 8 {
		t.Errorf("Unexpected CyberPanel username length: %s", username)
	}
}

func TestCustomServerProvisioner_OfflineFallback(t *testing.T) {
	ctx := context.Background()
	p := provisioning.NewCustomServerProvisioner(provisioning.CustomServerConfig{
		EndpointURL: "", // offline mode
	})

	order := &domain.Order{
		ID:        1001,
		ClientID:  42,
		ProductID: 5,
		Title:     "Dedicated Docker Container",
	}

	res, err := p.Create(ctx, order)
	if err != nil {
		t.Fatalf("Expected offline custom provisioner to succeed, got: %v", err)
	}
	if !res.Success {
		t.Errorf("Expected success=true, got %+v", res)
	}
}
