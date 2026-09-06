package provisioning

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// RDAPRegistrarDriver implements live RFC 7480/7484 RDAP availability checks with DNS fallback
type RDAPRegistrarDriver struct {
	httpClient   *http.Client
	tldPricing   map[string]int64
	bootstrapURL string
	bootstrapMap map[string][]string
	mu           sync.RWMutex
	takenCache   map[string]bool
}

// NewRDAPRegistrarDriver constructs an RDAP registry driver
func NewRDAPRegistrarDriver() *RDAPRegistrarDriver {
	return &RDAPRegistrarDriver{
		httpClient: &http.Client{
			Timeout: 6 * time.Second,
		},
		bootstrapURL: "https://data.iana.org/rdap/dns.json",
		tldPricing: map[string]int64{
			"com": 129900, // $12.99
			"net": 149900, // $14.99
			"org": 139900, // $13.99
			"id":  189900, // $18.99
			"io":  399900, // $39.99
			"co":  249900, // $24.99
			"xyz": 99900,  // $9.99
			"app": 169900, // $16.99
			"dev": 159900, // $15.99
		},
		bootstrapMap: map[string][]string{
			"com": {"https://rdap.verisign.com/com/v1/"},
			"net": {"https://rdap.verisign.com/net/v1/"},
			"org": {"https://rdap.publicinterestregistry.org/rdap/"},
			"id":  {"https://rdap.pandi.id/rdap/"},
			"io":  {"https://rdap.identitydigital.services/rdap/"},
			"co":  {"https://rdap.nic.co/"},
			"xyz": {"https://rdap.centralnic.com/xyz/"},
			"app": {"https://rdap.registry.google/rdap/"},
			"dev": {"https://rdap.registry.google/rdap/"},
		},
		takenCache: map[string]bool{
			"google.com":      true,
			"github.com":      true,
			"fossbilling.org": true,
			"example.com":     true,
		},
	}
}

// CheckAvailability queries authoritative RDAP registry endpoints and falls back to DNS lookup
func (d *RDAPRegistrarDriver) CheckAvailability(ctx context.Context, domainName string) (*DomainAvailability, error) {
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

	// 1. Fast cache check
	d.mu.RLock()
	if d.takenCache[domainName] {
		d.mu.RUnlock()
		return &DomainAvailability{
			DomainName:   domainName,
			TLD:          tld,
			IsAvailable:  false,
			Price:        price,
			Currency:     "USD",
			CheckedAt:    time.Now().UTC(),
			RegistrarRef: "rdap-cache",
		}, nil
	}
	d.mu.RUnlock()

	// 2. Perform live RDAP lookup
	isAvailable, err := d.checkRDAP(ctx, domainName, tld)
	if err != nil {
		// Fallback to live DNS query (NS / Host resolution)
		isAvailable = d.checkDNSFallback(ctx, domainName)
	}

	// Cache taken domains
	if !isAvailable {
		d.mu.Lock()
		d.takenCache[domainName] = true
		d.mu.Unlock()
	}

	return &DomainAvailability{
		DomainName:   domainName,
		TLD:          tld,
		IsAvailable:  isAvailable,
		Price:        price,
		Currency:     "USD",
		CheckedAt:    time.Now().UTC(),
		RegistrarRef: "rdap-iana",
	}, nil
}

func (d *RDAPRegistrarDriver) checkRDAP(ctx context.Context, domainName, tld string) (bool, error) {
	d.mu.RLock()
	servers, ok := d.bootstrapMap[tld]
	d.mu.RUnlock()

	if !ok || len(servers) == 0 {
		return false, fmt.Errorf("no RDAP server configured for .%s", tld)
	}

	for _, server := range servers {
		reqURL := fmt.Sprintf("%sdomain/%s", strings.TrimRight(server, "/")+"/", domainName)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
		if err != nil {
			continue
		}
		req.Header.Set("Accept", "application/rdap+json, application/json")

		resp, err := d.httpClient.Do(req)
		if err != nil {
			continue
		}
		_ = resp.Body.Close()

		// RDAP RFC 7480: 404 indicates domain is not registered (available)
		if resp.StatusCode == http.StatusNotFound {
			return true, nil
		}
		// 200 / 2xx indicates domain exists (taken)
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return false, nil
		}
	}

	return false, fmt.Errorf("all RDAP queries failed for %s", domainName)
}

func (d *RDAPRegistrarDriver) checkDNSFallback(ctx context.Context, domainName string) bool {
	resolver := net.DefaultResolver

	// Check if NS records exist
	ns, err := resolver.LookupNS(ctx, domainName)
	if err == nil && len(ns) > 0 {
		return false // Domain is taken and has nameservers
	}

	// Check if A/AAAA host records exist
	ips, err := resolver.LookupIP(ctx, "ip", domainName)
	if err == nil && len(ips) > 0 {
		return false // Domain is active
	}

	// If lookup returns standard NXDOMAIN error, domain is likely available
	return true
}

func (d *RDAPRegistrarDriver) RegisterDomain(ctx context.Context, req DomainRegistrationRequest) (*DomainRegistrationResult, error) {
	domainName := strings.ToLower(strings.TrimSpace(req.DomainName))

	d.mu.Lock()
	d.takenCache[domainName] = true
	d.mu.Unlock()

	if req.Years <= 0 {
		req.Years = 1
	}

	now := time.Now().UTC()
	expiresAt := now.AddDate(req.Years, 0, 0)
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

func (d *RDAPRegistrarDriver) RenewDomain(ctx context.Context, domainName string, years int) (*DomainRegistrationResult, error) {
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

// FetchBootstrap syncs with IANA DNS bootstrap registry
func (d *RDAPRegistrarDriver) FetchBootstrap(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, d.bootstrapURL, nil)
	if err != nil {
		return err
	}
	resp, err := d.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var payload struct {
		Services [][][]string `json:"services"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return err
	}

	d.mu.Lock()
	defer d.mu.Unlock()
	for _, entry := range payload.Services {
		if len(entry) >= 2 {
			tlds := entry[0]
			servers := entry[1]
			for _, t := range tlds {
				d.bootstrapMap[strings.ToLower(t)] = servers
			}
		}
	}
	return nil
}
