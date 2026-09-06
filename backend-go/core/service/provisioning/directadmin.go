package provisioning

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type DirectAdminProvisioner struct {
	Host       string
	Port       int
	Username   string
	Password   string
	HTTPClient *http.Client
}

func NewDirectAdminProvisioner(host string, port int, username, password string) *DirectAdminProvisioner {
	if port == 0 {
		port = 2222
	}
	return &DirectAdminProvisioner{
		Host:       host,
		Port:       port,
		Username:   username,
		Password:   password,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (p *DirectAdminProvisioner) Type() domain.ProductType {
	return domain.ProductTypeHosting
}

type DirectAdminAccount struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Domain   string `json:"domain"`
	Package  string `json:"package"`
	Email    string `json:"email"`
	IP       string `json:"ip"`
	Status   string `json:"status"`
}

func (p *DirectAdminProvisioner) CreateAccount(ctx context.Context, acc DirectAdminAccount) (*DirectAdminAccount, error) {
	if acc.Username == "" {
		acc.Username = "da" + strings.ToLower(fmt.Sprintf("%d", time.Now().UnixNano()%1000000))
	}
	if acc.Password == "" {
		acc.Password = fmt.Sprintf("DA!%s#%d", acc.Username, time.Now().Year())
	}
	if acc.IP == "" {
		acc.IP = p.Host
	}

	acc.Status = "active"
	return &acc, nil
}

func (p *DirectAdminProvisioner) Create(ctx context.Context, order *domain.Order) (*domain.ProvisionResult, error) {
	if order == nil {
		return nil, errors.New("invalid order for directadmin provisioning")
	}

	acc, err := p.CreateAccount(ctx, DirectAdminAccount{
		Domain:  fmt.Sprintf("client%d-service.com", order.ClientID),
		Package: "default",
	})
	if err != nil {
		return nil, err
	}

	details := map[string]string{
		"username": acc.Username,
		"password": acc.Password,
		"domain":   acc.Domain,
		"server":   p.Host,
		"da_url":   fmt.Sprintf("https://%s:%d", p.Host, p.Port),
	}
	detailsJSON, _ := json.Marshal(details)

	return &domain.ProvisionResult{
		Success:        true,
		RemoteID:       acc.Username,
		AccountDetails: detailsJSON,
	}, nil
}

func (p *DirectAdminProvisioner) Suspend(ctx context.Context, order *domain.Order, reason string) error {
	return nil
}

func (p *DirectAdminProvisioner) Unsuspend(ctx context.Context, order *domain.Order) error {
	return nil
}

func (p *DirectAdminProvisioner) Renew(ctx context.Context, order *domain.Order) error {
	return nil
}

func (p *DirectAdminProvisioner) Terminate(ctx context.Context, order *domain.Order) error {
	return nil
}

func (p *DirectAdminProvisioner) Sync(ctx context.Context, order *domain.Order) (*domain.ServiceStatus, error) {
	return &domain.ServiceStatus{
		IsActive:    true,
		RemoteState: "active",
	}, nil
}

func (p *DirectAdminProvisioner) ChangePassword(ctx context.Context, order *domain.Order, newPassword string) error {
	return nil
}
