package provisioning

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type CpanelConfig struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	APIToken string `json:"api_token"`
	Insecure bool   `json:"insecure"`
}

type CpanelProvisioner struct {
	config CpanelConfig
	client *http.Client
}

func NewCpanelProvisioner(config CpanelConfig) *CpanelProvisioner {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: config.Insecure},
	}
	return &CpanelProvisioner{
		config: config,
		client: &http.Client{
			Transport: tr,
			Timeout:   30 * time.Second,
		},
	}
}

func (p *CpanelProvisioner) Type() domain.ProductType {
	return domain.ProductTypeHosting
}

func (p *CpanelProvisioner) call(ctx context.Context, function string, params url.Values) (map[string]interface{}, error) {
	baseURL := fmt.Sprintf("https://%s:2087/json-api/%s", p.config.Host, function)
	params.Set("api.version", "1")

	reqURL := baseURL + "?" + params.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("whm %s:%s", p.config.Username, p.config.APIToken))

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode WHM response: %w (body: %s)", err, string(body))
	}

	return result, nil
}

func (p *CpanelProvisioner) GenerateAccountCredentials(domainName string) (username, password string) {
	clean := strings.ToLower(domainName)
	clean = strings.ReplaceAll(clean, ".", "")
	clean = strings.ReplaceAll(clean, "-", "")
	if len(clean) > 8 {
		clean = clean[:8]
	}
	if len(clean) < 4 {
		clean = clean + "host"
	}
	username = clean
	password = fmt.Sprintf("Sec!%s#%d9", clean, len(domainName)*7)
	return username, password
}

func (p *CpanelProvisioner) Create(ctx context.Context, order *domain.Order) (*domain.ProvisionResult, error) {
	var accountConfig struct {
		Domain string `json:"domain"`
		Plan   string `json:"plan"`
	}
	if len(order.Config) > 0 {
		_ = json.Unmarshal(order.Config, &accountConfig)
	}
	if accountConfig.Domain == "" {
		return nil, errors.New("domain is required for cPanel provisioning")
	}

	username, password := p.GenerateAccountCredentials(accountConfig.Domain)
	params := url.Values{}
	params.Set("username", username)
	params.Set("password", password)
	params.Set("domain", accountConfig.Domain)
	if accountConfig.Plan != "" {
		params.Set("plan", accountConfig.Plan)
	}

	res, err := p.call(ctx, "createacct", params)
	if err != nil {
		return nil, err
	}

	metadata, ok := res["metadata"].(map[string]interface{})
	if !ok || metadata["result"].(float64) != 1 {
		reason := "unknown error"
		if metadata != nil && metadata["reason"] != nil {
			reason = metadata["reason"].(string)
		}
		return nil, fmt.Errorf("cPanel account creation failed: %s", reason)
	}

	details := map[string]string{
		"username":   username,
		"password":   password,
		"domain":     accountConfig.Domain,
		"server":     p.config.Host,
		"cpanel_url": fmt.Sprintf("https://%s:2083", p.config.Host),
	}
	detailsJSON, _ := json.Marshal(details)

	return &domain.ProvisionResult{
		Success:        true,
		RemoteID:       username,
		AccountDetails: detailsJSON,
	}, nil
}

func (p *CpanelProvisioner) Suspend(ctx context.Context, order *domain.Order, reason string) error {
	var details map[string]string
	_ = json.Unmarshal(order.Config, &details)
	username := details["username"]
	if username == "" {
		return errors.New("cPanel username not found in order config")
	}

	params := url.Values{}
	params.Set("user", username)
	params.Set("reason", reason)

	res, err := p.call(ctx, "suspendacct", params)
	if err != nil {
		return err
	}

	metadata := res["metadata"].(map[string]interface{})
	if metadata["result"].(float64) != 1 {
		return fmt.Errorf("cPanel suspension failed: %s", metadata["reason"])
	}
	return nil
}

func (p *CpanelProvisioner) Unsuspend(ctx context.Context, order *domain.Order) error {
	var details map[string]string
	_ = json.Unmarshal(order.Config, &details)
	username := details["username"]

	params := url.Values{}
	params.Set("user", username)

	res, err := p.call(ctx, "unsuspendacct", params)
	if err != nil {
		return err
	}

	metadata := res["metadata"].(map[string]interface{})
	if metadata["result"].(float64) != 1 {
		return fmt.Errorf("cPanel unsuspension failed: %s", metadata["reason"])
	}
	return nil
}

func (p *CpanelProvisioner) Renew(ctx context.Context, order *domain.Order) error {
	return nil // cPanel renewal is usually handled by billing
}

func (p *CpanelProvisioner) Terminate(ctx context.Context, order *domain.Order) error {
	var details map[string]string
	_ = json.Unmarshal(order.Config, &details)
	username := details["username"]

	params := url.Values{}
	params.Set("user", username)

	res, err := p.call(ctx, "removeacct", params)
	if err != nil {
		return err
	}

	metadata := res["metadata"].(map[string]interface{})
	if metadata["result"].(float64) != 1 {
		return fmt.Errorf("cPanel termination failed: %s", metadata["reason"])
	}
	return nil
}

func (p *CpanelProvisioner) Sync(ctx context.Context, order *domain.Order) (*domain.ServiceStatus, error) {
	var details map[string]string
	_ = json.Unmarshal(order.Config, &details)
	username := details["username"]

	params := url.Values{}
	params.Set("user", username)

	res, err := p.call(ctx, "accountsummary", params)
	if err != nil {
		return nil, err
	}

	metadata := res["metadata"].(map[string]interface{})
	if metadata["result"].(float64) != 1 {
		return nil, fmt.Errorf("cPanel sync failed: %s", metadata["reason"])
	}

	data := res["data"].(map[string]interface{})
	acct := data["acct"].([]interface{})[0].(map[string]interface{})

	return &domain.ServiceStatus{
		IsActive:    acct["suspended"].(float64) == 0,
		RemoteState: acct["status"].(string),
		DiskUsageMB: int64(acct["diskused"].(string)[0]), // Simplified parsing
	}, nil
}

func (p *CpanelProvisioner) ChangePassword(ctx context.Context, order *domain.Order, newPassword string) error {
	var details map[string]string
	_ = json.Unmarshal(order.Config, &details)
	username := details["username"]

	params := url.Values{}
	params.Set("user", username)
	params.Set("password", newPassword)

	res, err := p.call(ctx, "passwd", params)
	if err != nil {
		return err
	}

	metadata := res["metadata"].(map[string]interface{})
	if metadata["result"].(float64) != 1 {
		return fmt.Errorf("cPanel password change failed: %s", metadata["reason"])
	}
	return nil
}
