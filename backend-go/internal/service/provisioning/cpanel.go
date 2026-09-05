package provisioning

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/domain"
)

type CpanelConfig struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	APIToken string `json:"api_token"`
	UseSSL   bool   `json:"use_ssl"`
}

type CpanelAccountPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Domain   string `json:"domain"`
	Plan     string `json:"plan"`
}

type CpanelProvisioner struct {
	config CpanelConfig
}

func NewCpanelProvisioner(config CpanelConfig) *CpanelProvisioner {
	return &CpanelProvisioner{config: config}
}

func (p *CpanelProvisioner) Type() domain.ProductType {
	return domain.ProductTypeHosting
}

// GenerateAccountCredentials generates cPanel compliant username and password
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
	if order == nil || order.ProductID == 0 {
		return nil, errors.New("invalid order for hosting provisioning")
	}

	var accountConfig struct {
		Domain string `json:"domain"`
		Plan   string `json:"plan"`
	}
	if len(order.Config) > 0 {
		_ = json.Unmarshal(order.Config, &accountConfig)
	}
	if accountConfig.Domain == "" {
		accountConfig.Domain = fmt.Sprintf("client%d-service.com", order.ClientID)
	}
	if accountConfig.Plan == "" {
		accountConfig.Plan = "default_package"
	}

	username, password := p.GenerateAccountCredentials(accountConfig.Domain)

	details := map[string]string{
		"username": username,
		"password": password,
		"domain":   accountConfig.Domain,
		"plan":     accountConfig.Plan,
		"server":   p.config.Host,
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
	// In production, sends WHM API: /json-api/suspendacct?api.version=1&user={username}&reason={reason}
	return nil
}

func (p *CpanelProvisioner) Unsuspend(ctx context.Context, order *domain.Order) error {
	// In production, sends WHM API: /json-api/unsuspendacct?api.version=1&user={username}
	return nil
}

func (p *CpanelProvisioner) Renew(ctx context.Context, order *domain.Order) error {
	return nil
}

func (p *CpanelProvisioner) Terminate(ctx context.Context, order *domain.Order) error {
	// In production, sends WHM API: /json-api/removeacct?api.version=1&user={username}
	return nil
}

func (p *CpanelProvisioner) Sync(ctx context.Context, order *domain.Order) (*domain.ServiceStatus, error) {
	return &domain.ServiceStatus{
		IsActive:    true,
		RemoteState: "active",
		DiskUsageMB: 128,
		BandwidthMB: 1024,
	}, nil
}

func (p *CpanelProvisioner) ChangePassword(ctx context.Context, order *domain.Order, newPassword string) error {
	// In production, sends WHM API: /json-api/passwd?api.version=1&user={username}&password={newPassword}
	return nil
}
