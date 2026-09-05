package provisioning

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type LicenseProvisioner struct {
	secretKey string
}

func NewLicenseProvisioner(secretKey string) *LicenseProvisioner {
	return &LicenseProvisioner{secretKey: secretKey}
}

func (p *LicenseProvisioner) Type() domain.ProductType {
	return domain.ProductTypeLicense
}

// GenerateLicenseKey creates a formatted serial key: FOSS-XXXX-XXXX-XXXX-XXXX
func (p *LicenseProvisioner) GenerateLicenseKey(clientID int64, orderID int64, issuedAt time.Time) string {
	raw := fmt.Sprintf("%d-%d-%d-%s", clientID, orderID, issuedAt.Unix(), p.secretKey)
	hash := sha256.Sum256([]byte(raw))
	hexStr := strings.ToUpper(hex.EncodeToString(hash[:8])) // 16 chars

	return fmt.Sprintf("FOSS-%s-%s-%s-%s",
		hexStr[0:4],
		hexStr[4:8],
		hexStr[8:12],
		hexStr[12:16],
	)
}

func (p *LicenseProvisioner) Create(ctx context.Context, order *domain.Order) (*domain.ProvisionResult, error) {
	now := time.Now().UTC()
	key := p.GenerateLicenseKey(order.ClientID, order.ID, now)

	details := map[string]string{
		"license_key": key,
		"status":      "active",
		"issued_at":   now.Format(time.RFC3339),
	}
	detailsJSON, _ := json.Marshal(details)

	return &domain.ProvisionResult{
		Success:        true,
		RemoteID:       key,
		AccountDetails: detailsJSON,
	}, nil
}

func (p *LicenseProvisioner) Suspend(ctx context.Context, order *domain.Order, reason string) error {
	return nil
}

func (p *LicenseProvisioner) Unsuspend(ctx context.Context, order *domain.Order) error {
	return nil
}

func (p *LicenseProvisioner) Renew(ctx context.Context, order *domain.Order) error {
	return nil
}

func (p *LicenseProvisioner) Terminate(ctx context.Context, order *domain.Order) error {
	return nil
}

func (p *LicenseProvisioner) Sync(ctx context.Context, order *domain.Order) (*domain.ServiceStatus, error) {
	return &domain.ServiceStatus{
		IsActive:    order.Status == domain.OrderStatusActive,
		RemoteState: string(order.Status),
	}, nil
}

func (p *LicenseProvisioner) ChangePassword(ctx context.Context, order *domain.Order, newPassword string) error {
	return nil
}
