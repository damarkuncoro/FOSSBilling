package provisioning

import (
	"fmt"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type ServerConfig struct {
	Type     string `json:"type"` // "cpanel", "directadmin", "plesk", "license"
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	APIToken string `json:"api_token"`
	UseSSL   bool   `json:"use_ssl"`
}

// ProvisionerFactory instantiates the appropriate domain.ServiceProvisioner based on server configurations
type ProvisionerFactory struct{}

func NewProvisionerFactory() *ProvisionerFactory {
	return &ProvisionerFactory{}
}

func (f *ProvisionerFactory) CreateProvisioner(cfg ServerConfig) (domain.ServiceProvisioner, error) {
	switch cfg.Type {
	case "cpanel":
		return NewCpanelProvisioner(CpanelConfig{
			Host:     cfg.Host,
			Username: cfg.Username,
			APIToken: cfg.APIToken,
			UseSSL:   cfg.UseSSL,
		}), nil
	case "directadmin":
		return NewDirectAdminProvisioner(cfg.Host, cfg.Port, cfg.Username, cfg.Password), nil
	case "plesk":
		return NewPleskProvisioner(cfg.Host, cfg.Port, cfg.APIToken), nil
	case "license":
		salt := cfg.Password
		if salt == "" {
			salt = "FOSSBILLING_DEFAULT_SALT"
		}
		return NewLicenseProvisioner(salt), nil
	default:
		return nil, fmt.Errorf("unsupported server provisioner type: %s", cfg.Type)
	}
}
