package factory_test

import (
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/payment/gateways"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/provisioning"
)

func TestProvisionerFactory(t *testing.T) {
	factory := provisioning.NewProvisionerFactory()

	// 1. cPanel
	cpanel, err := factory.CreateProvisioner(provisioning.ServerConfig{
		Type:     "cpanel",
		Host:     "cpanel.example.com",
		Username: "root",
		APIToken: "token",
	})
	if err != nil || cpanel == nil {
		t.Fatalf("failed to create cpanel provisioner: %v", err)
	}

	// 2. DirectAdmin
	da, err := factory.CreateProvisioner(provisioning.ServerConfig{
		Type:     "directadmin",
		Host:     "da.example.com",
		Username: "admin",
	})
	if err != nil || da == nil {
		t.Fatalf("failed to create directadmin provisioner: %v", err)
	}

	// 3. Unknown
	_, err = factory.CreateProvisioner(provisioning.ServerConfig{
		Type: "unknown_engine",
	})
	if err == nil {
		t.Fatalf("expected error for unknown provisioner engine")
	}
}

func TestPaymentGatewayFactory(t *testing.T) {
	factory := gateways.NewPaymentGatewayFactory()

	// 1. Stripe
	stripe, err := factory.CreateGateway(gateways.GatewayConfig{
		ID:      "stripe",
		Enabled: true,
		Settings: map[string]interface{}{
			"secret_key":      "sk_test_123",
			"publishable_key": "pk_test_123",
			"webhook_secret":  "whsec_123",
		},
	})
	if err != nil || stripe == nil || stripe.ID() != "stripe" {
		t.Fatalf("failed to create stripe gateway: %v", err)
	}

	// 2. Midtrans
	midtrans, err := factory.CreateGateway(gateways.GatewayConfig{
		ID:      "midtrans",
		Enabled: true,
		Settings: map[string]interface{}{
			"server_key": "server_123",
			"client_key": "client_123",
		},
	})
	if err != nil || midtrans == nil || midtrans.ID() != "midtrans" {
		t.Fatalf("failed to create midtrans gateway: %v", err)
	}

	// 3. Unknown
	_, err = factory.CreateGateway(gateways.GatewayConfig{
		ID: "crypto_unknown",
	})
	if err == nil {
		t.Fatalf("expected error for unknown payment gateway")
	}
}
