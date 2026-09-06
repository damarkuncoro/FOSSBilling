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

type PleskProvisioner struct {
	Host       string
	Port       int
	APIKey     string
	HTTPClient *http.Client
}

func NewPleskProvisioner(host string, port int, apiKey string) *PleskProvisioner {
	if port == 0 {
		port = 8443
	}
	return &PleskProvisioner{
		Host:       host,
		Port:       port,
		APIKey:     apiKey,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (p *PleskProvisioner) Type() domain.ProductType {
	return domain.ProductTypeHosting
}

type PleskSubscription struct {
	DomainName string `json:"domain_name"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	PlanName   string `json:"plan_name"`
	Status     string `json:"status"`
}

func (p *PleskProvisioner) CreateSubscription(ctx context.Context, sub PleskSubscription) (*PleskSubscription, error) {
	if sub.Username == "" {
		sub.Username = "pl_" + strings.ToLower(fmt.Sprintf("%d", time.Now().UnixNano()%1000000))
	}
	if sub.Password == "" {
		sub.Password = fmt.Sprintf("Plsk!%s#%d", sub.Username, time.Now().Year())
	}
	sub.Status = "active"
	return &sub, nil
}

func (p *PleskProvisioner) Create(ctx context.Context, order *domain.Order) (*domain.ProvisionResult, error) {
	if order == nil {
		return nil, errors.New("invalid order for plesk provisioning")
	}

	sub, err := p.CreateSubscription(ctx, PleskSubscription{
		DomainName: fmt.Sprintf("client%d-service.com", order.ClientID),
		PlanName:   "default",
	})
	if err != nil {
		return nil, err
	}

	details := map[string]string{
		"username":  sub.Username,
		"password":  sub.Password,
		"domain":    sub.DomainName,
		"server":    p.Host,
		"plesk_url": fmt.Sprintf("https://%s:%d", p.Host, p.Port),
	}
	detailsJSON, _ := json.Marshal(details)

	return &domain.ProvisionResult{
		Success:        true,
		RemoteID:       sub.Username,
		AccountDetails: detailsJSON,
	}, nil
}

func (p *PleskProvisioner) Suspend(ctx context.Context, order *domain.Order, reason string) error {
	return nil
}

func (p *PleskProvisioner) Unsuspend(ctx context.Context, order *domain.Order) error {
	return nil
}

func (p *PleskProvisioner) Renew(ctx context.Context, order *domain.Order) error {
	return nil
}

func (p *PleskProvisioner) Terminate(ctx context.Context, order *domain.Order) error {
	return nil
}

func (p *PleskProvisioner) Sync(ctx context.Context, order *domain.Order) (*domain.ServiceStatus, error) {
	return &domain.ServiceStatus{
		IsActive:    true,
		RemoteState: "active",
	}, nil
}

func (p *PleskProvisioner) ChangePassword(ctx context.Context, order *domain.Order, newPassword string) error {
	return nil
}

func (p *PleskProvisioner) SuspendSubscription(ctx context.Context, domainName string) error {
	return nil
}

func (p *PleskProvisioner) UnsuspendSubscription(ctx context.Context, domainName string) error {
	return nil
}
