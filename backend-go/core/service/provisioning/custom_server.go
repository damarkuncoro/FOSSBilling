package provisioning

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type CustomServerConfig struct {
	EndpointURL string            `json:"endpoint_url"`
	AuthHeader  string            `json:"auth_header"`
	AuthToken   string            `json:"auth_token"`
	Headers     map[string]string `json:"headers"`
	TimeoutSec  int               `json:"timeout_sec"`
}

type CustomServerProvisioner struct {
	config CustomServerConfig
	client *http.Client
}

func NewCustomServerProvisioner(config CustomServerConfig) *CustomServerProvisioner {
	timeout := 30 * time.Second
	if config.TimeoutSec > 0 {
		timeout = time.Duration(config.TimeoutSec) * time.Second
	}
	return &CustomServerProvisioner{
		config: config,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (p *CustomServerProvisioner) Type() domain.ProductType {
	return domain.ProductTypeCustom
}

func (p *CustomServerProvisioner) call(ctx context.Context, action string, payload interface{}) ([]byte, error) {
	if p.config.EndpointURL == "" {
		// Mock offline pass if no remote endpoint configured
		return []byte(`{"status":"success","mock":true}`), nil
	}

	url := fmt.Sprintf("%s/%s", strings.TrimRight(p.config.EndpointURL, "/"), action)
	bodyData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if p.config.AuthToken != "" {
		authHeader := p.config.AuthHeader
		if authHeader == "" {
			authHeader = "Authorization"
		}
		req.Header.Set(authHeader, p.config.AuthToken)
	}

	for k, v := range p.config.Headers {
		req.Header.Set(k, v)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("custom server webhook error: %w", err)
	}
	defer resp.Body.Close()

	resBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return resBytes, fmt.Errorf("custom server returned HTTP %d: %s", resp.StatusCode, string(resBytes))
	}

	return resBytes, nil
}

func (p *CustomServerProvisioner) Create(ctx context.Context, order *domain.Order) (*domain.ProvisionResult, error) {
	reqPayload := map[string]interface{}{
		"action":     "create",
		"order_id":   order.ID,
		"client_id":  order.ClientID,
		"product_id": order.ProductID,
		"title":      order.Title,
		"config":     order.Config,
		"created_at": time.Now(),
	}

	res, err := p.call(ctx, "create", reqPayload)
	if err != nil {
		return &domain.ProvisionResult{
			Success:      false,
			ErrorMessage: err.Error(),
		}, err
	}

	return &domain.ProvisionResult{
		Success:        true,
		RemoteID:       fmt.Sprintf("custom_%d", order.ID),
		AccountDetails: res,
	}, nil
}

func (p *CustomServerProvisioner) Suspend(ctx context.Context, order *domain.Order, reason string) error {
	_, err := p.call(ctx, "suspend", map[string]interface{}{
		"order_id": order.ID,
		"reason":   reason,
	})
	return err
}

func (p *CustomServerProvisioner) Unsuspend(ctx context.Context, order *domain.Order) error {
	_, err := p.call(ctx, "unsuspend", map[string]interface{}{
		"order_id": order.ID,
	})
	return err
}

func (p *CustomServerProvisioner) Renew(ctx context.Context, order *domain.Order) error {
	_, err := p.call(ctx, "renew", map[string]interface{}{
		"order_id": order.ID,
	})
	return err
}

func (p *CustomServerProvisioner) Terminate(ctx context.Context, order *domain.Order) error {
	_, err := p.call(ctx, "terminate", map[string]interface{}{
		"order_id": order.ID,
	})
	return err
}

func (p *CustomServerProvisioner) Sync(ctx context.Context, order *domain.Order) (*domain.ServiceStatus, error) {
	_, err := p.call(ctx, "status", map[string]interface{}{
		"order_id": order.ID,
	})
	if err != nil {
		return &domain.ServiceStatus{IsActive: false, RemoteState: "error"}, nil
	}

	return &domain.ServiceStatus{
		IsActive:    true,
		RemoteState: "active",
	}, nil
}

func (p *CustomServerProvisioner) ChangePassword(ctx context.Context, order *domain.Order, newPassword string) error {
	if strings.TrimSpace(newPassword) == "" {
		return errors.New("new password cannot be empty")
	}
	_, err := p.call(ctx, "change_password", map[string]interface{}{
		"order_id":     order.ID,
		"new_password": newPassword,
	})
	return err
}
