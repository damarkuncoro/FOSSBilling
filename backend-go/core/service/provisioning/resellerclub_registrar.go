package provisioning

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ResellerClubConfig struct {
	AuthUserID string `json:"auth_user_id"`
	APIKey     string `json:"api_key"`
	IsTest     bool   `json:"is_test"`
}

type ResellerClubRegistrarDriver struct {
	config     ResellerClubConfig
	httpClient *http.Client
}

func NewResellerClubRegistrarDriver(config ResellerClubConfig) *ResellerClubRegistrarDriver {
	return &ResellerClubRegistrarDriver{
		config: config,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (d *ResellerClubRegistrarDriver) getBaseURL() string {
	if d.config.IsTest {
		return "https://test.httpapi.com/api/"
	}
	return "https://httpapi.com/api/"
}

func (d *ResellerClubRegistrarDriver) makeRequest(ctx context.Context, method, path string, params url.Values) (json.RawMessage, error) {
	params.Set("auth-userid", d.config.AuthUserID)
	params.Set("api-key", d.config.APIKey)

	fullURL := d.getBaseURL() + path + ".json"
	if method == "GET" {
		fullURL += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, nil)
	if err != nil {
		return nil, err
	}

	if method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Body = http.NoBody // Standard POST in ResellerClub uses query params often, but let's be careful
		// If real body needed: req.Body = io.NopCloser(strings.NewReader(params.Encode()))
		// For ResellerClub, many POSTs actually take params in the body
		req, _ = http.NewRequestWithContext(ctx, method, d.getBaseURL()+path+".json", strings.NewReader(params.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Basic error handling based on ResellerClub response format
	var errorResp struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Error   string `json:"error"`
	}
	_ = json.Unmarshal(result, &errorResp)
	if strings.ToUpper(errorResp.Status) == "ERROR" || errorResp.Error != "" {
		msg := errorResp.Message
		if msg == "" {
			msg = errorResp.Error
		}
		return nil, fmt.Errorf("resellerclub error: %s", msg)
	}

	return result, nil
}

func (d *ResellerClubRegistrarDriver) CheckAvailability(ctx context.Context, domainName string) (*DomainAvailability, error) {
	parts := strings.Split(domainName, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid domain: %s", domainName)
	}
	sld := parts[0]
	tld := parts[len(parts)-1]

	params := url.Values{}
	params.Set("domain-name", sld)
	params.Set("tlds", tld)
	params.Set("suggest-alternative", "false")

	data, err := d.makeRequest(ctx, "GET", "domains/available", params)
	if err != nil {
		return nil, err
	}

	var availabilityMap map[string]struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(data, &availabilityMap); err != nil {
		return nil, err
	}

	info, ok := availabilityMap[domainName]
	if !ok {
		// Try case-insensitive lookup
		for k, v := range availabilityMap {
			if strings.EqualFold(k, domainName) {
				info = v
				ok = true
				break
			}
		}
	}

	if !ok {
		return nil, fmt.Errorf("domain %s not found in response", domainName)
	}

	return &DomainAvailability{
		DomainName:   domainName,
		TLD:          tld,
		IsAvailable:  strings.ToLower(info.Status) == "available",
		Currency:     "USD",
		CheckedAt:    time.Now().UTC(),
		RegistrarRef: "resellerclub",
	}, nil
}

func (d *ResellerClubRegistrarDriver) RegisterDomain(ctx context.Context, req DomainRegistrationRequest) (*DomainRegistrationResult, error) {
	// Implementation would involve creating customer, adding contacts, and then calling domains/register
	// This is a simplified version returning a pending status for now
	return &DomainRegistrationResult{
		DomainName:   req.DomainName,
		Status:       "pending",
		RegisteredAt: time.Now().UTC(),
		ExpiresAt:    time.Now().AddDate(req.Years, 0, 0),
		Nameservers:  req.Nameservers,
	}, nil
}

func (d *ResellerClubRegistrarDriver) RenewDomain(ctx context.Context, domainName string, years int) (*DomainRegistrationResult, error) {
	return nil, fmt.Errorf("renew not implemented for ResellerClub in Go yet")
}
