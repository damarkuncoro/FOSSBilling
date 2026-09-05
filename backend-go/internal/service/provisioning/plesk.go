package provisioning

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
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

func (p *PleskProvisioner) SuspendSubscription(ctx context.Context, domainName string) error {
	return nil
}

func (p *PleskProvisioner) UnsuspendSubscription(ctx context.Context, domainName string) error {
	return nil
}
