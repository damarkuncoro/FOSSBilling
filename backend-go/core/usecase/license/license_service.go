package license

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

var (
	ErrLicenseNotFound     = errors.New("license order not found")
	ErrUnauthorizedLicense = errors.New("unauthorized license access")
)

type LicenseRecordDTO struct {
	ID             int64     `json:"id"`
	ProductName    string    `json:"product_title"`
	LicenseKey     string    `json:"license_key"`
	Status         string    `json:"status"`
	MaxInstances   int       `json:"max_instances"`
	LicensedDomain string    `json:"licensed_domain"`
	LicensedIP     string    `json:"licensed_ip"`
	ExpiresAt      time.Time `json:"expires_at"`
}

type LicenseConfig struct {
	LicenseKey     string `json:"license_key"`
	MaxInstances   int    `json:"max_instances"`
	LicensedDomain string `json:"licensed_domain"`
	LicensedIP     string `json:"licensed_ip"`
}

type LicenseService struct {
	orderRepo domain.OrderRepository
}

func NewLicenseService(orderRepo domain.OrderRepository) *LicenseService {
	return &LicenseService{
		orderRepo: orderRepo,
	}
}

// ListClientLicenses returns all software licenses owned by a client
func (s *LicenseService) ListClientLicenses(ctx context.Context, clientID int64) ([]LicenseRecordDTO, error) {
	orders, _, err := s.orderRepo.ListByClientID(ctx, clientID, 100, 0)
	if err != nil {
		return nil, err
	}

	list := make([]LicenseRecordDTO, 0)
	for _, ord := range orders {
		var cfg LicenseConfig
		if len(ord.Config) > 0 {
			_ = json.Unmarshal(ord.Config, &cfg)
		}

		if cfg.LicenseKey != "" || strings.Contains(strings.ToLower(ord.Title), "license") || strings.Contains(strings.ToLower(ord.Title), "edition") {
			key := cfg.LicenseKey
			if key == "" {
				key = fmt.Sprintf("FB-ENT-%X-%X", ord.ID*48271%65535, ord.ID*12345%65535)
			}

			max := cfg.MaxInstances
			if max <= 0 {
				max = 5
			}

			exp := time.Now().AddDate(1, 0, 0)
			if ord.ExpiresAt != nil {
				exp = *ord.ExpiresAt
			}

			list = append(list, LicenseRecordDTO{
				ID:             ord.ID,
				ProductName:    ord.Title,
				LicenseKey:     key,
				Status:         string(ord.Status),
				MaxInstances:   max,
				LicensedDomain: cfg.LicensedDomain,
				LicensedIP:     cfg.LicensedIP,
				ExpiresAt:      exp,
			})
		}
	}

	return list, nil
}

// ResetLicenseKey generates and assigns a new cryptographically secure license key
func (s *LicenseService) ResetLicenseKey(ctx context.Context, clientID int64, orderID int64) (string, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return "", ErrLicenseNotFound
	}

	if order.ClientID != clientID {
		return "", ErrUnauthorizedLicense
	}

	var cfg LicenseConfig
	if len(order.Config) > 0 {
		_ = json.Unmarshal(order.Config, &cfg)
	}

	b := make([]byte, 8)
	_, _ = rand.Read(b)
	newKey := fmt.Sprintf("FB-ENT-%s-%s", strings.ToUpper(hex.EncodeToString(b[:4])), strings.ToUpper(hex.EncodeToString(b[4:])))
	cfg.LicenseKey = newKey

	cfgBytes, err := json.Marshal(cfg)
	if err != nil {
		return "", err
	}

	order.Config = cfgBytes
	if err := s.orderRepo.Update(ctx, order); err != nil {
		return "", err
	}

	return newKey, nil
}

// ResetLicenseLock resets the domain and IP lock for a license order
func (s *LicenseService) ResetLicenseLock(ctx context.Context, clientID int64, orderID int64) error {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return ErrLicenseNotFound
	}

	if order.ClientID != clientID {
		return ErrUnauthorizedLicense
	}

	var cfg LicenseConfig
	if len(order.Config) > 0 {
		_ = json.Unmarshal(order.Config, &cfg)
	}

	cfg.LicensedDomain = ""
	cfg.LicensedIP = ""

	cfgBytes, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	order.Config = cfgBytes
	return s.orderRepo.Update(ctx, order)
}

