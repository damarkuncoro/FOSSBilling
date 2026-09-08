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

type HestiaConfig struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	AccessKey string `json:"access_key"` // Hestia Access Key ID
	SecretKey string `json:"secret_key"` // Hestia Secret Key
	Insecure  bool   `json:"insecure"`
	Package   string `json:"package"`
}

type HestiaProvisioner struct {
	config HestiaConfig
	client *http.Client
}

func NewHestiaProvisioner(config HestiaConfig) *HestiaProvisioner {
	if config.Port <= 0 {
		config.Port = 8083
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: config.Insecure},
	}
	return &HestiaProvisioner{
		config: config,
		client: &http.Client{
			Transport: tr,
			Timeout:   30 * time.Second,
		},
	}
}

func (p *HestiaProvisioner) Type() domain.ProductType {
	return domain.ProductTypeHosting
}

func (p *HestiaProvisioner) call(ctx context.Context, cmd string, args ...string) (string, error) {
	apiURL := fmt.Sprintf("https://%s:%d/api/", p.config.Host, p.config.Port)

	formData := url.Values{}
	formData.Set("user", p.config.AccessKey)
	formData.Set("password", p.config.SecretKey)
	formData.Set("returncode", "yes")
	formData.Set("cmd", cmd)

	for i, arg := range args {
		formData.Set(fmt.Sprintf("arg%d", i+1), arg)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("hestia connection error: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	res := strings.TrimSpace(string(bodyBytes))
	if res != "0" && res != "OK" && !strings.HasPrefix(res, "{") && !strings.HasPrefix(res, "[") {
		if resp.StatusCode >= 400 || (len(res) > 0 && res != "0") {
			return res, fmt.Errorf("hestia returned code: %s", res)
		}
	}

	return res, nil
}

func (p *HestiaProvisioner) GenerateUsername(domainName string) string {
	clean := strings.ToLower(strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			return r
		}
		return -1
	}, domainName))

	if len(clean) > 8 {
		clean = clean[:8]
	}
	if len(clean) == 0 || (clean[0] >= '0' && clean[0] <= '9') {
		clean = "a" + clean
	}
	return clean
}

func (p *HestiaProvisioner) getOrderDomain(order *domain.Order) string {
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

func (p *HestiaProvisioner) Create(ctx context.Context, order *domain.Order) (*domain.ProvisionResult, error) {
	domainName := p.getOrderDomain(order)
	username := p.GenerateUsername(domainName)
	pkg := p.config.Package
	if pkg == "" {
		pkg = "default"
	}

	// 1. Create Hestia User: v-add-user <user> <password> <email> [package] [first_name] [last_name]
	password := fmt.Sprintf("HestiaPass_%d!", time.Now().Unix())
	email := fmt.Sprintf("client_%d@fossbilling.org", order.ClientID)

	_, err := p.call(ctx, "v-add-user", username, password, email, pkg)
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return &domain.ProvisionResult{
			Success:      false,
			ErrorMessage: err.Error(),
		}, err
	}

	// 2. Add Web Domain: v-add-web-domain <user> <domain>
	if domainName != "" {
		_, _ = p.call(ctx, "v-add-web-domain", username, domainName)
	}

	details, _ := json.Marshal(map[string]interface{}{
		"control_panel": "HestiaCP",
		"server":        p.config.Host,
		"port":          p.config.Port,
		"login_url":     fmt.Sprintf("https://%s:%d", p.config.Host, p.config.Port),
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

func (p *HestiaProvisioner) Suspend(ctx context.Context, order *domain.Order, reason string) error {
	username := p.GenerateUsername(p.getOrderDomain(order))
	_, err := p.call(ctx, "v-suspend-user", username)
	return err
}

func (p *HestiaProvisioner) Unsuspend(ctx context.Context, order *domain.Order) error {
	username := p.GenerateUsername(p.getOrderDomain(order))
	_, err := p.call(ctx, "v-unsuspend-user", username)
	return err
}

func (p *HestiaProvisioner) Renew(ctx context.Context, order *domain.Order) error {
	return nil // Hosting renewals are billing-handled, server account stays alive
}

func (p *HestiaProvisioner) Terminate(ctx context.Context, order *domain.Order) error {
	username := p.GenerateUsername(p.getOrderDomain(order))
	_, err := p.call(ctx, "v-delete-user", username)
	return err
}

func (p *HestiaProvisioner) Sync(ctx context.Context, order *domain.Order) (*domain.ServiceStatus, error) {
	username := p.GenerateUsername(p.getOrderDomain(order))
	out, err := p.call(ctx, "v-list-user", username, "json")
	if err != nil {
		return &domain.ServiceStatus{IsActive: false, RemoteState: "error"}, nil
	}

	isActive := !strings.Contains(out, `"SUSPENDED":"yes"`)
	return &domain.ServiceStatus{
		IsActive:    isActive,
		RemoteState: "active",
	}, nil
}

func (p *HestiaProvisioner) ChangePassword(ctx context.Context, order *domain.Order, newPassword string) error {
	if strings.TrimSpace(newPassword) == "" {
		return errors.New("new password cannot be empty")
	}
	username := p.GenerateUsername(p.getOrderDomain(order))
	_, err := p.call(ctx, "v-change-user-password", username, newPassword)
	return err
}

func (p *HestiaProvisioner) TestConnection(ctx context.Context) error {
	_, err := p.call(ctx, "v-list-sys-services", "json")
	return err
}
