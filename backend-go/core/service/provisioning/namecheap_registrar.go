package provisioning

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type NamecheapConfig struct {
	ApiUser   string `json:"api_user"`
	ApiKey    string `json:"api_key"`
	UserName  string `json:"username"`
	ClientIp  string `json:"client_ip"`
	IsSandbox bool   `json:"is_sandbox"`
}

type NamecheapRegistrarDriver struct {
	config     NamecheapConfig
	httpClient *http.Client
}

func NewNamecheapRegistrarDriver(config NamecheapConfig) *NamecheapRegistrarDriver {
	return &NamecheapRegistrarDriver{
		config: config,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (d *NamecheapRegistrarDriver) getBaseURL() string {
	if d.config.IsSandbox {
		return "https://api.sandbox.namecheap.com/xml.response"
	}
	return "https://api.namecheap.com/xml.response"
}

type namecheapResponse struct {
	XMLName xml.Name `xml:"ApiResponse"`
	Status  string   `xml:"Status,attr"`
	Errors  []struct {
		Message string `xml:",chardata"`
		Number  string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse struct {
		DomainCheckResult struct {
			Domain    string `xml:"Domain,attr"`
			Available bool   `xml:"Available,attr"`
		} `xml:"DomainCheckResult"`
	} `xml:"CommandResponse"`
}

func (d *NamecheapRegistrarDriver) makeRequest(ctx context.Context, command string, params url.Values) (*namecheapResponse, error) {
	params.Set("ApiUser", d.config.ApiUser)
	params.Set("ApiKey", d.config.ApiKey)
	params.Set("UserName", d.config.UserName)
	params.Set("ClientIp", d.config.ClientIp)
	params.Set("Command", command)

	fullURL := d.getBaseURL() + "?" + params.Encode()

	resp, err := d.httpClient.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result namecheapResponse
	if err := xml.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if strings.ToLower(result.Status) == "error" && len(result.Errors) > 0 {
		return nil, fmt.Errorf("namecheap error: %s", result.Errors[0].Message)
	}

	return &result, nil
}

func (d *NamecheapRegistrarDriver) CheckAvailability(ctx context.Context, domainName string) (*DomainAvailability, error) {
	params := url.Values{}
	params.Set("DomainList", domainName)

	res, err := d.makeRequest(ctx, "namecheap.domains.check", params)
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
		IsAvailable:  res.CommandResponse.DomainCheckResult.Available,
		Currency:     "USD",
		CheckedAt:    time.Now().UTC(),
		RegistrarRef: "namecheap",
	}, nil
}

func (d *NamecheapRegistrarDriver) RegisterDomain(ctx context.Context, req DomainRegistrationRequest) (*DomainRegistrationResult, error) {
	params := url.Values{}
	params.Set("DomainName", req.DomainName)
	params.Set("Years", fmt.Sprintf("%d", req.Years))

	if len(req.Nameservers) > 0 {
		params.Set("Nameservers", strings.Join(req.Nameservers, ","))
	}

	// Mapping Contact Info to Namecheap fields
	contactTypes := []string{"Registrant", "Admin", "Tech", "AuxBilling"}
	for _, ct := range contactTypes {
		params.Set(ct+"FirstName", req.ContactInfo["first_name"])
		params.Set(ct+"LastName", req.ContactInfo["last_name"])
		params.Set(ct+"EmailAddress", req.ContactInfo["email"])
		params.Set(ct+"Phone", "+"+req.ContactInfo["phone_cc"]+"."+req.ContactInfo["phone"])
		params.Set(ct+"Address1", req.ContactInfo["address1"])
		params.Set(ct+"City", req.ContactInfo["city"])
		params.Set(ct+"StateProvince", req.ContactInfo["state"])
		params.Set(ct+"Country", req.ContactInfo["country"])
		params.Set(ct+"PostalCode", req.ContactInfo["postcode"])
	}

	res, err := d.makeRequest(ctx, "namecheap.domains.create", params)
	if err != nil {
		if d.config.IsSandbox {
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

	// Check if registered successfully in XML response (simplified check here)
	if strings.ToLower(res.Status) != "ok" {
		return nil, fmt.Errorf("namecheap registration failed with status: %s", res.Status)
	}

	now := time.Now().UTC()
	return &DomainRegistrationResult{
		DomainName:   req.DomainName,
		Status:       "active",
		RegisteredAt: now,
		ExpiresAt:    now.AddDate(req.Years, 0, 0),
		Nameservers:  req.Nameservers,
	}, nil
}

func (d *NamecheapRegistrarDriver) RenewDomain(ctx context.Context, domainName string, years int) (*DomainRegistrationResult, error) {
	return nil, fmt.Errorf("renew not implemented for Namecheap in Go yet")
}
