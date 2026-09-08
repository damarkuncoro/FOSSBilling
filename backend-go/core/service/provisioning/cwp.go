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

type CWPConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	APIKey   string `json:"api_key"`
	Insecure bool   `json:"insecure"`
	Package  string `json:"package"`
}

type CWPProvisioner struct {
	config CWPConfig
	client *http.Client
}

func NewCWPProvisioner(config CWPConfig) *CWPProvisioner {
	if config.Port <= 0 {
		config.Port = 2304 // CWP API port
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: config.Insecure},
	}
	return &CWPProvisioner{
		config: config,
		client: &http.Client{
			Transport: tr,
			Timeout:   30 * time.Second,
		},
	}
}

func (p *CWPProvisioner) Type() domain.ProductType {
	return domain.ProductTypeHosting
}

func (p *CWPProvisioner) call(ctx context.Context, action string, params url.Values) (map[string]interface{}, error) {
	apiURL := fmt.Sprintf("https://%s:%d/v1/%s", p.config.Host, p.config.Port, action)
	if params == nil {
		params = url.Values{}
	}
	params.Set("key", p.config.APIKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cwp connection error: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &res); err != nil {
		return map[string]interface{}{
			"raw":    string(bodyBytes),
			"status": "OK",
		}, nil
	}

	return res, nil
}

func (p *CWPProvisioner) GenerateUsername(domainName string) string {
	clean := strings.ToLower(strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			return r
		}
		return -1
	}, domainName))

	if len(clean) > 8 {
		clean = clean[:8]
	}
	if len(clean) == 0 {
		clean = "cwpuser"
	}
	return clean
}

func (p *CWPProvisioner) getOrderDomain(order *domain.Order) string {
	var cfg struct {
		Domain string `json:"domain"`
	}
	if len(order.Config) > 0 {
		_ = json.Unmarshal(order.Config, &cfg)
	}
	if cfg.Domain != "" {
		return cfg.Domain
	}
	return fmt.Sprintf("client%d-srv.com", order.ClientID)
}

func (p *CWPProvisioner) Create(ctx context.Context, order *domain.Order) (*domain.ProvisionResult, error) {
	domainName := p.getOrderDomain(order)
	username := p.GenerateUsername(domainName)
	pkg := p.config.Package
	if pkg == "" {
		pkg = "default"
	}

	password := fmt.Sprintf("CWPPass_%d!", time.Now().Unix())
	email := fmt.Sprintf("client_%d@fossbilling.org", order.ClientID)

	params := url.Values{}
	params.Set("domain", domainName)
	params.Set("user", username)
	params.Set("pass", password)
	params.Set("email", email)
	params.Set("package", pkg)

	_, err := p.call(ctx, "account", params)
	if err != nil {
		return &domain.ProvisionResult{
			Success:      false,
			ErrorMessage: err.Error(),
		}, err
	}

	details, _ := json.Marshal(map[string]interface{}{
		"control_panel": "CentOS Web Panel (CWP)",
		"server":        p.config.Host,
		"login_url":     fmt.Sprintf("https://%s:2083", p.config.Host),
		"username":      username,
		"password":      password,
		"package":       pkg,
		"domain":        domainName,
	})

	return &domain.ProvisionResult{
		Success:        true,
		RemoteID:       username,
		AccountDetails: details,
	}, nil
}

func (p *CWPProvisioner) Suspend(ctx context.Context, order *domain.Order, reason string) error {
	username := p.GenerateUsername(p.getOrderDomain(order))
	params := url.Values{}
	params.Set("user", username)
	_, err := p.call(ctx, "account/suspend", params)
	return err
}

func (p *CWPProvisioner) Unsuspend(ctx context.Context, order *domain.Order) error {
	username := p.GenerateUsername(p.getOrderDomain(order))
	params := url.Values{}
	params.Set("user", username)
	_, err := p.call(ctx, "account/unsuspend", params)
	return err
}

func (p *CWPProvisioner) Renew(ctx context.Context, order *domain.Order) error {
	return nil
}

func (p *CWPProvisioner) Terminate(ctx context.Context, order *domain.Order) error {
	username := p.GenerateUsername(p.getOrderDomain(order))
	params := url.Values{}
	params.Set("user", username)
	_, err := p.call(ctx, "account/delete", params)
	return err
}

func (p *CWPProvisioner) Sync(ctx context.Context, order *domain.Order) (*domain.ServiceStatus, error) {
	username := p.GenerateUsername(p.getOrderDomain(order))
	params := url.Values{}
	params.Set("user", username)
	res, err := p.call(ctx, "account/status", params)
	if err != nil {
		return &domain.ServiceStatus{IsActive: false, RemoteState: "unknown"}, nil
	}

	state := "active"
	if statusVal, ok := res["status"].(string); ok && strings.ToLower(statusVal) == "suspended" {
		state = "suspended"
	}

	return &domain.ServiceStatus{
		IsActive:    state == "active",
		RemoteState: state,
	}, nil
}

func (p *CWPProvisioner) ChangePassword(ctx context.Context, order *domain.Order, newPassword string) error {
	if strings.TrimSpace(newPassword) == "" {
		return errors.New("new password cannot be empty")
	}
	username := p.GenerateUsername(p.getOrderDomain(order))
	params := url.Values{}
	params.Set("user", username)
	params.Set("pass", newPassword)
	_, err := p.call(ctx, "account/update_pass", params)
	return err
}

func (p *CWPProvisioner) TestConnection(ctx context.Context) error {
	_, err := p.call(ctx, "status", nil)
	return err
}
