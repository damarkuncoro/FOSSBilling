package domain

import (
	"context"
	"encoding/json"
)

type ProvisionResult struct {
	Success       bool            `json:"success"`
	RemoteID      string          `json:"remote_id"`
	AccountDetails json.RawMessage `json:"account_details,omitempty"`
	ErrorMessage  string          `json:"error_message,omitempty"`
}

type ServiceStatus struct {
	IsActive    bool   `json:"is_active"`
	RemoteState string `json:"remote_state"`
	DiskUsageMB int64  `json:"disk_usage_mb,omitempty"`
	BandwidthMB int64  `json:"bandwidth_mb,omitempty"`
}

// ServiceProvisioner defines the SPI for external service providers (cPanel, Plesk, Registrars)
type ServiceProvisioner interface {
	Type() ProductType
	Create(ctx context.Context, order *Order) (*ProvisionResult, error)
	Suspend(ctx context.Context, order *Order, reason string) error
	Unsuspend(ctx context.Context, order *Order) error
	Renew(ctx context.Context, order *Order) error
	Terminate(ctx context.Context, order *Order) error
	Sync(ctx context.Context, order *Order) (*ServiceStatus, error)
	ChangePassword(ctx context.Context, order *Order, newPassword string) error
}
