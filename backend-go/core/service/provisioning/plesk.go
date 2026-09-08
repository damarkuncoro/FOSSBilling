package provisioning

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type PleskConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	APIKey   string `json:"api_key"`
	Insecure bool   `json:"insecure"`
}

type PleskProvisioner struct {
	config PleskConfig
	client *http.Client
}

func NewPleskProvisioner(config PleskConfig) *PleskProvisioner {
	if config.Port == 0 {
		config.Port = 8443
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: config.Insecure},
	}
	return &PleskProvisioner{
		config: config,
		client: &http.Client{
			Transport: tr,
			Timeout:   30 * time.Second,
		},
	}
}

func (p *PleskProvisioner) Type() domain.ProductType {
	return domain.ProductTypeHosting
}

func (p *PleskProvisioner) call(ctx context.Context, method, path string, payload interface{}) (map[string]interface{}, error) {
	apiURL := fmt.Sprintf("https://%s:%d/api/v2%s", p.config.Host, p.config.Port, path)

	var body io.Reader
	if payload != nil {
		b, _ := json.Marshal(payload)
		body = bytes.NewBuffer(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, apiURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-Key", p.config.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("plesk api error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result map[string]interface{}
	if len(respBody) > 0 {
		if err := json.Unmarshal(respBody, &result); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (p *PleskProvisioner) Create(ctx context.Context, order *domain.Order) (*domain.ProvisionResult, error) {
	var accountConfig struct {
		Domain string `json:"domain"`
		Plan   string `json:"plan"`
	}
	if len(order.Config) > 0 {
		_ = json.Unmarshal(order.Config, &accountConfig)
	}
	if accountConfig.Domain == "" {
		return nil, errors.New("domain is required for Plesk provisioning")
	}

	// 1. Create Client
	username := "pl_" + strings.ToLower(fmt.Sprintf("%d", order.ID))
	password := fmt.Sprintf("Plsk!%s#%d", username, time.Now().Year())

	clientPayload := map[string]interface{}{
		"name":     "Client " + fmt.Sprint(order.ClientID),
		"login":    username,
		"password": password,
		"email":    "client@example.com",
	}

	cliRes, err := p.call(ctx, http.MethodPost, "/clients", clientPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to create plesk client: %w", err)
	}

	clientGUID := cliRes["guid"].(string)

	// 2. Create Domain/Subscription
	domainPayload := map[string]interface{}{
		"name":         accountConfig.Domain,
		"owner_client": map[string]string{"guid": clientGUID},
		"hosting_type": "virtual",
		"hosting_settings": map[string]interface{}{
			"ftp_login":    username,
			"ftp_password": password,
		},
	}

	if _, err := p.call(ctx, http.MethodPost, "/domains", domainPayload); err != nil {
		return nil, fmt.Errorf("failed to create plesk domain: %w", err)
	}

	details := map[string]string{
		"username":  username,
		"password":  password,
		"domain":    accountConfig.Domain,
		"server":    p.config.Host,
		"plesk_url": fmt.Sprintf("https://%s:%d", p.config.Host, p.config.Port),
	}
	detailsJSON, _ := json.Marshal(details)

	return &domain.ProvisionResult{
		Success:        true,
		RemoteID:       username,
		AccountDetails: detailsJSON,
	}, nil
}

func (p *PleskProvisioner) Suspend(ctx context.Context, order *domain.Order, reason string) error {
	var details map[string]string
	_ = json.Unmarshal(order.Config, &details)
	domainName := details["domain"]

	payload := map[string]interface{}{"status": "suspended"}
	_, err := p.call(ctx, http.MethodPut, "/domains/"+domainName, payload)
	return err
}

func (p *PleskProvisioner) Unsuspend(ctx context.Context, order *domain.Order) error {
	var details map[string]string
	_ = json.Unmarshal(order.Config, &details)
	domainName := details["domain"]

	payload := map[string]interface{}{"status": "active"}
	_, err := p.call(ctx, http.MethodPut, "/domains/"+domainName, payload)
	return err
}

func (p *PleskProvisioner) Renew(ctx context.Context, order *domain.Order) error {
	return nil
}

func (p *PleskProvisioner) Terminate(ctx context.Context, order *domain.Order) error {
	var details map[string]string
	_ = json.Unmarshal(order.Config, &details)
	domainName := details["domain"]

	_, err := p.call(ctx, http.MethodDelete, "/domains/"+domainName, nil)
	return err
}

func (p *PleskProvisioner) Sync(ctx context.Context, order *domain.Order) (*domain.ServiceStatus, error) {
	var details map[string]string
	_ = json.Unmarshal(order.Config, &details)
	domainName := details["domain"]

	res, err := p.call(ctx, http.MethodGet, "/domains/"+domainName, nil)
	if err != nil {
		return nil, err
	}

	status, _ := res["status"].(string)
	return &domain.ServiceStatus{
		IsActive:    status == "active",
		RemoteState: status,
	}, nil
}

func (p *PleskProvisioner) ChangePassword(ctx context.Context, order *domain.Order, newPassword string) error {
	var details map[string]string
	_ = json.Unmarshal(order.Config, &details)
	domainName := details["domain"]

	payload := map[string]interface{}{
		"hosting_settings": map[string]interface{}{
			"ftp_password": newPassword,
		},
	}
	_, err := p.call(ctx, http.MethodPut, "/domains/"+domainName, payload)
	return err
}

func (p *PleskProvisioner) TestConnection(ctx context.Context) error {
	_, err := p.call(ctx, http.MethodGet, "/server", nil)
	return err
}
