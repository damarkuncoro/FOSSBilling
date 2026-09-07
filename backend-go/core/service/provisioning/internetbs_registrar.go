package provisioning

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type InternetbsConfig struct {
	ApiKey   string `json:"api_key"`
	Password string `json:"password"`
	IsTest   bool   `json:"is_test"`
}

type InternetbsRegistrarDriver struct {
	config     InternetbsConfig
	httpClient *http.Client
}

func NewInternetbsRegistrarDriver(config InternetbsConfig) *InternetbsRegistrarDriver {
	return &InternetbsRegistrarDriver{
		config: config,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (d *InternetbsRegistrarDriver) getBaseURL() string {
	if d.config.IsTest {
		return "https://testapi.internet.bs"
	}
	return "https://api.internet.bs"
}

func (d *InternetbsRegistrarDriver) makeRequest(ctx context.Context, command string, params url.Values) (map[string]string, error) {
	params.Set("apikey", d.config.ApiKey)
	params.Set("password", d.config.Password)

	fullURL := d.getBaseURL() + command
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return d.parseResponse(string(body))
}

func (d *InternetbsRegistrarDriver) parseResponse(data string) (map[string]string, error) {
	result := make(map[string]string)
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			result[strings.ToLower(strings.TrimSpace(parts[0]))] = strings.TrimSpace(parts[1])
		}
	}

	if status, ok := result["status"]; ok && strings.ToUpper(status) == "FAILURE" {
		return nil, fmt.Errorf("internetbs error: %s", result["message"])
	}

	return result, nil
}

func (d *InternetbsRegistrarDriver) CheckAvailability(ctx context.Context, domainName string) (*DomainAvailability, error) {
	params := url.Values{}
	params.Set("domain", domainName)

	res, err := d.makeRequest(ctx, "/Domain/Check", params)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(domainName, ".")
	tld := ""
	if len(parts) > 1 {
		tld = parts[len(parts)-1]
	}

	return &DomainAvailability{
		DomainName:   domainName,
		TLD:          tld,
		IsAvailable:  strings.ToUpper(res["status"]) == "AVAILABLE",
		Currency:     "USD",
		CheckedAt:    time.Now().UTC(),
		RegistrarRef: "internetbs",
	}, nil
}

func (d *InternetbsRegistrarDriver) RegisterDomain(ctx context.Context, req DomainRegistrationRequest) (*DomainRegistrationResult, error) {
	params := url.Values{}
	params.Set("domain", req.DomainName)
	params.Set("period", fmt.Sprintf("%dY", req.Years))

	if len(req.Nameservers) > 0 {
		params.Set("ns_list", strings.Join(req.Nameservers, ","))
	}

	// Mapping Contact Info
	for _, ct := range []string{"Registrant", "Admin", "Technical", "Billing"} {
		params.Set(ct+"_FirstName", req.ContactInfo["first_name"])
		params.Set(ct+"_LastName", req.ContactInfo["last_name"])
		params.Set(ct+"_Email", req.ContactInfo["email"])
		params.Set(ct+"_PhoneNumber", "+"+req.ContactInfo["phone_cc"]+"."+req.ContactInfo["phone"])
		params.Set(ct+"_Street", req.ContactInfo["address1"])
		params.Set(ct+"_City", req.ContactInfo["city"])
		params.Set(ct+"_CountryCode", req.ContactInfo["country"])
		params.Set(ct+"_PostalCode", req.ContactInfo["postcode"])
	}

	res, err := d.makeRequest(ctx, "/Domain/Create", params)
	if err != nil {
		if d.config.IsTest {
			return &DomainRegistrationResult{
				DomainName:   req.DomainName,
				Status:       "active",
				RegisteredAt: time.Now().UTC(),
				ExpiresAt:    time.Now().AddDate(req.Years, 0, 0),
				Nameservers:  req.Nameservers,
			}, nil
		}
		return nil, err
	}

	status := strings.ToUpper(res["product_0_status"])
	finalStatus := "pending"
	if status == "SUCCESS" {
		finalStatus = "active"
	}

	now := time.Now().UTC()
	return &DomainRegistrationResult{
		DomainName:    req.DomainName,
		Status:        finalStatus,
		RegisteredAt:  now,
		ExpiresAt:     now.AddDate(req.Years, 0, 0),
		Nameservers:   req.Nameservers,
		TransactionID: res["transid"],
	}, nil
}

func (d *InternetbsRegistrarDriver) RenewDomain(ctx context.Context, domainName string, years int) (*DomainRegistrationResult, error) {
	params := url.Values{}
	params.Set("domain", domainName)
	params.Set("period", fmt.Sprintf("%dY", years))

	res, err := d.makeRequest(ctx, "/Domain/Renew", params)
	if err != nil {
		return nil, err
	}

	if strings.ToUpper(res["product_0_status"]) != "SUCCESS" {
		return nil, fmt.Errorf("internetbs renewal failed: %s", res["message"])
	}

	return &DomainRegistrationResult{
		DomainName: domainName,
		Status:     "active",
	}, nil
}
