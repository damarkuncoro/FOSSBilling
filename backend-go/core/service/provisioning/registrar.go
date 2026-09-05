package provisioning

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// DomainAvailability represents the result of a domain lookup check
type DomainAvailability struct {
	DomainName   string    `json:"domain_name"`
	TLD          string    `json:"tld"`
	IsAvailable  bool      `json:"is_available"`
	Price        int64     `json:"price"` // in hundredths of cent (domain.Money)
	Currency     string    `json:"currency"`
	CheckedAt    time.Time `json:"checked_at"`
	RegistrarRef string    `json:"registrar_ref,omitempty"`
}

// DomainRegistrationRequest contains payload for registering a domain
type DomainRegistrationRequest struct {
	DomainName string            `json:"domain_name"`
	Years      int               `json:"years"`
	Nameservers []string         `json:"nameservers"`
	ContactInfo map[string]string `json:"contact_info"`
}

// DomainRegistrationResult represents outcome of registrar API action
type DomainRegistrationResult struct {
	DomainName   string    `json:"domain_name"`
	Status       string    `json:"status"` // "active", "pending", "failed"
	RegisteredAt time.Time `json:"registered_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	Nameservers  []string  `json:"nameservers"`
	AuthCode     string    `json:"auth_code,omitempty"`
	TransactionID string   `json:"transaction_id"`
}

// RegistrarDriver defines the interface for domain registrars (e.g. Namecheap, ResellerClub, OpenSRS)
type RegistrarDriver interface {
	CheckAvailability(ctx context.Context, domainName string) (*DomainAvailability, error)
	RegisterDomain(ctx context.Context, req DomainRegistrationRequest) (*DomainRegistrationResult, error)
	RenewDomain(ctx context.Context, domainName string, years int) (*DomainRegistrationResult, error)
}

// MockRegistrarDriver implements a deterministic registrar driver for testing and default installs
type MockRegistrarDriver struct {
	tldPricing map[string]int64 // TLD -> Price (e.g. "com" -> 129900)
	takenDomains map[string]bool
}

// NewMockRegistrarDriver creates a new registrar driver
func NewMockRegistrarDriver() *MockRegistrarDriver {
	return &MockRegistrarDriver{
		tldPricing: map[string]int64{
			"com": 129900, // $12.99
			"net": 149900, // $14.99
			"org": 139900, // $13.99
			"id":  189900, // $18.99
			"io":  399900, // $39.99
		},
		takenDomains: map[string]bool{
			"google.com":    true,
			"github.com":    true,
			"fossbilling.org": true,
			"example.com":   true,
		},
	}
}

func (d *MockRegistrarDriver) CheckAvailability(ctx context.Context, domainName string) (*DomainAvailability, error) {
	domainName = strings.ToLower(strings.TrimSpace(domainName))
	parts := strings.Split(domainName, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid domain format: %s", domainName)
	}

	tld := parts[len(parts)-1]
	price, exists := d.tldPricing[tld]
	if !exists {
		price = 199900 // Default $19.99
	}

	isAvailable := !d.takenDomains[domainName]

	return &DomainAvailability{
		DomainName:   domainName,
		TLD:          tld,
		IsAvailable:  isAvailable,
		Price:        price,
		Currency:     "USD",
		CheckedAt:    time.Now().UTC(),
		RegistrarRef: "mock-registrar",
	}, nil
}

func (d *MockRegistrarDriver) RegisterDomain(ctx context.Context, req DomainRegistrationRequest) (*DomainRegistrationResult, error) {
	domainName := strings.ToLower(strings.TrimSpace(req.DomainName))
	if d.takenDomains[domainName] {
		return nil, fmt.Errorf("domain %s is already taken", domainName)
	}

	if req.Years <= 0 {
		req.Years = 1
	}

	now := time.Now().UTC()
	expiresAt := now.AddDate(req.Years, 0, 0)
	d.takenDomains[domainName] = true

	ns := req.Nameservers
	if len(ns) == 0 {
		ns = []string{"ns1.fossbilling.org", "ns2.fossbilling.org"}
	}

	return &DomainRegistrationResult{
		DomainName:    domainName,
		Status:        "active",
		RegisteredAt:  now,
		ExpiresAt:     expiresAt,
		Nameservers:   ns,
		AuthCode:      fmt.Sprintf("EPP-%d-FOSS", time.Now().UnixNano()%1000000),
		TransactionID: fmt.Sprintf("REG-%d", time.Now().UnixNano()),
	}, nil
}

func (d *MockRegistrarDriver) RenewDomain(ctx context.Context, domainName string, years int) (*DomainRegistrationResult, error) {
	domainName = strings.ToLower(strings.TrimSpace(domainName))
	if years <= 0 {
		years = 1
	}

	now := time.Now().UTC()
	expiresAt := now.AddDate(years, 0, 0)

	return &DomainRegistrationResult{
		DomainName:    domainName,
		Status:        "active",
		RegisteredAt:  now,
		ExpiresAt:     expiresAt,
		Nameservers:   []string{"ns1.fossbilling.org", "ns2.fossbilling.org"},
		TransactionID: fmt.Sprintf("REN-%d", time.Now().UnixNano()),
	}, nil
}
