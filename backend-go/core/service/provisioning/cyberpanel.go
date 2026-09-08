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

type CyberPanelConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	AdminPass string `json:"admin_password"`
	Insecure bool   `json:"insecure"`
}

type CyberPanelProvisioner struct {
	config CyberPanelConfig
	client *http.Client
}

func NewCyberPanelProvisioner(config CyberPanelConfig) *CyberPanelProvisioner {
	if config.Port <= 0 {
		config.Port = 8090
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: config.Insecure},
	}
	return &CyberPanelProvisioner{
		config: config,
		client: &http.Client{
			Transport: tr,
			Timeout:   30 * time.Second,
		},
	}
}

func (p *CyberPanelProvisioner) Type() domain.ProductType {
	return domain.ProductTypeHosting
}

func (p *CyberPanelProvisioner) call(ctx context.Context, action string, payload map[string]interface{}) (map[string]interface{}, error) {
	apiURL := fmt.Sprintf("https://%s:%d/api/%s", p.config.Host, p.config.Port, action)

	if payload == nil {
		payload = make(map[string]interface{})
	}

	// CyberPanel API usually requires admin password in the payload
	payload["adminPass"] = p.config.AdminPass

	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cyberpanel connection error: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var res map[string]interface{}
	if err := json.Unmarshal(respBody, &res); err != nil {
		return nil, fmt.Errorf("failed to parse cyberpanel response: %s", string(respBody))
	}

	// CyberPanel usually returns status codes like "1" for success or "status" field
	// Based on docs, it often uses "status" : 1 or similar
	return res, nil
}

func (p *CyberPanelProvisioner) Create(ctx context.Context, order *domain.Order) (*domain.ProvisionResult, error) {
	var cfg struct {
		Domain string `json:"domain"`
		Plan   string `json:"plan"`
	}
	if len(order.Config) > 0 {
		_ = json.Unmarshal(order.Config, &cfg)
	}
	if cfg.Domain == "" {
		return nil, errors.New("domain name is required")
	}

	username := "cp" + fmt.Sprintf("%d", order.ID)
	password := "Cp!" + username + "2026#"

	payload := map[string]interface{}{
		"websiteName": cfg.Domain,
		"ownerEmail":  "client@fossbilling.org",
		"packageName": cfg.Plan,
		"template":    "Default",
		"phpSelection": "7.4",
		"acl":         "user",
	}
	if payload["packageName"] == "" {
		payload["packageName"] = "Default"
	}

	res, err := p.call(ctx, "createWebsite", payload)
	if err != nil {
		return nil, err
	}

	// Simple check (CyberPanel returns "status" 1 for success)
	if fmt.Sprintf("%v", res["status"]) != "1" {
		return nil, fmt.Errorf("cyberpanel creation failed: %v", res["error_message"])
	}

	details, _ := json.Marshal(map[string]interface{}{
		"control_panel": "CyberPanel",
		"host":          p.config.Host,
		"username":      username,
		"password":      password,
		"domain":        cfg.Domain,
		"login_url":     fmt.Sprintf("https://%s:%d", p.config.Host, p.config.Port),
	})

	return &domain.ProvisionResult{
		Success:        true,
		RemoteID:       cfg.Domain,
		AccountDetails: details,
	}, nil
}

func (p *CyberPanelProvisioner) Suspend(ctx context.Context, order *domain.Order, reason string) error {
	var cfg struct {
		Domain string `json:"domain"`
	}
	_ = json.Unmarshal(order.Config, &cfg)

	_, err := p.call(ctx, "suspendWebsite", map[string]interface{}{
		"websiteName": cfg.Domain,
		"state":       "Suspend",
	})
	return err
}

func (p *CyberPanelProvisioner) Unsuspend(ctx context.Context, order *domain.Order) error {
	var cfg struct {
		Domain string `json:"domain"`
	}
	_ = json.Unmarshal(order.Config, &cfg)

	_, err := p.call(ctx, "suspendWebsite", map[string]interface{}{
		"websiteName": cfg.Domain,
		"state":       "Unsuspend",
	})
	return err
}

func (p *CyberPanelProvisioner) Renew(ctx context.Context, order *domain.Order) error {
	return nil
}

func (p *CyberPanelProvisioner) Terminate(ctx context.Context, order *domain.Order) error {
	var cfg struct {
		Domain string `json:"domain"`
	}
	_ = json.Unmarshal(order.Config, &cfg)

	_, err := p.call(ctx, "deleteWebsite", map[string]interface{}{
		"websiteName": cfg.Domain,
	})
	return err
}

func (p *CyberPanelProvisioner) Sync(ctx context.Context, order *domain.Order) (*domain.ServiceStatus, error) {
	return &domain.ServiceStatus{
		IsActive:    order.Status == domain.OrderStatusActive,
		RemoteState: string(order.Status),
	}, nil
}

func (p *CyberPanelProvisioner) ChangePassword(ctx context.Context, order *domain.Order, newPassword string) error {
	// CyberPanel change password API for users
	return nil
}

func (p *CyberPanelProvisioner) TestConnection(ctx context.Context) error {
	_, err := p.call(ctx, "verifyLogin", nil)
	return err
}
