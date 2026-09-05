package provisioning

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/fossbilling/backend-go/internal/domain"
)

func TestCpanelProvisioner_GenerateAccountCredentials(t *testing.T) {
	prov := NewCpanelProvisioner(CpanelConfig{Host: "cpanel1.hoster.com"})

	username, password := prov.GenerateAccountCredentials("my-awesome-domain.com")
	if len(username) > 8 {
		t.Errorf("Username length = %d; want <= 8", len(username))
	}
	if !strings.HasPrefix(password, "Sec!") {
		t.Errorf("Password should start with Sec!, got: %s", password)
	}
}

func TestCpanelProvisioner_Create(t *testing.T) {
	ctx := context.Background()
	prov := NewCpanelProvisioner(CpanelConfig{Host: "cpanel1.hoster.com"})

	cfg, _ := json.Marshal(map[string]string{
		"domain": "acme-corp.com",
		"plan":   "starter_unlimited",
	})
	order := &domain.Order{
		ID:        99,
		ClientID:  10,
		ProductID: 5,
		Config:    cfg,
	}

	result, err := prov.Create(ctx, order)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected Success = true")
	}
	if result.RemoteID == "" {
		t.Error("Expected non-empty RemoteID")
	}

	var details map[string]string
	_ = json.Unmarshal(result.AccountDetails, &details)
	if details["domain"] != "acme-corp.com" {
		t.Errorf("domain = %s; want acme-corp.com", details["domain"])
	}
	if details["cpanel_url"] != "https://cpanel1.hoster.com:2083" {
		t.Errorf("cpanel_url = %s; want https://cpanel1.hoster.com:2083", details["cpanel_url"])
	}
}

func TestLicenseProvisioner_GenerateLicenseKey(t *testing.T) {
	prov := NewLicenseProvisioner("super-secret-licensing-master-key")

	now := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	key1 := prov.GenerateLicenseKey(10, 100, now)
	key2 := prov.GenerateLicenseKey(10, 100, now)
	key3 := prov.GenerateLicenseKey(11, 101, now)

	// Deterministic hash with same inputs
	if key1 != key2 {
		t.Errorf("key1 (%s) != key2 (%s)", key1, key2)
	}
	// Different inputs -> different keys
	if key1 == key3 {
		t.Errorf("key1 (%s) == key3 (%s)", key1, key3)
	}

	if !strings.HasPrefix(key1, "FOSS-") {
		t.Errorf("License key format should start with FOSS-, got: %s", key1)
	}
	parts := strings.Split(key1, "-")
	if len(parts) != 5 {
		t.Errorf("License key parts count = %d; want 5", len(parts))
	}
}
