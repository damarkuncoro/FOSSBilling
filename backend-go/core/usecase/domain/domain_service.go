package domain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/provisioning"
)

var (
	ErrDomainNotFound      = errors.New("domain order not found")
	ErrUnauthorizedDomain  = errors.New("unauthorized domain access")
	ErrInvalidDomainName   = errors.New("invalid domain name")
	ErrInvalidNameservers  = errors.New("at least one valid nameserver is required")
)

type DomainAvailabilityDTO struct {
	DomainName  string  `json:"domain"`
	TLD         string  `json:"tld"`
	IsAvailable bool    `json:"available"`
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`
}

type DomainRecordDTO struct {
	ID          int64     `json:"id"`
	DomainName  string    `json:"domain_name"`
	TLD         string    `json:"tld"`
	Status      string    `json:"status"`
	Nameservers []string  `json:"nameservers"`
	EPPCode     string    `json:"epp_code"`
	AutoRenew   bool      `json:"auto_renew"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type DomainConfig struct {
	DomainName  string   `json:"domain_name"`
	Nameservers []string `json:"nameservers"`
	AutoRenew   bool     `json:"auto_renew"`
	EPPCode     string   `json:"epp_code"`
}

type DomainService struct {
	orderRepo domain.OrderRepository
	registrar provisioning.RegistrarDriver
}

func NewDomainService(orderRepo domain.OrderRepository, registrar provisioning.RegistrarDriver) *DomainService {
	return &DomainService{
		orderRepo: orderRepo,
		registrar: registrar,
	}
}

// CheckAvailability queries WHOIS / RDAP registry and calculates pricing
func (s *DomainService) CheckAvailability(ctx context.Context, domainName string) (*DomainAvailabilityDTO, error) {
	clean := strings.ToLower(strings.TrimSpace(domainName))
	if clean == "" || !strings.Contains(clean, ".") {
		return nil, ErrInvalidDomainName
	}

	result, err := s.registrar.CheckAvailability(ctx, clean)
	if err != nil {
		return nil, fmt.Errorf("registry check failed: %w", err)
	}

	return &DomainAvailabilityDTO{
		DomainName:  result.DomainName,
		TLD:         result.TLD,
		IsAvailable: result.IsAvailable,
		Price:       float64(result.Price) / 10000.0,
		Currency:    result.Currency,
	}, nil
}

// ListClientDomains returns all active and pending domain orders for a client
func (s *DomainService) ListClientDomains(ctx context.Context, clientID int64) ([]DomainRecordDTO, error) {
	orders, _, err := s.orderRepo.ListByClientID(ctx, clientID, 100, 0)
	if err != nil {
		return nil, err
	}

	list := make([]DomainRecordDTO, 0)
	for _, ord := range orders {
		var cfg DomainConfig
		if len(ord.Config) > 0 {
			_ = json.Unmarshal(ord.Config, &cfg)
		}

		if cfg.DomainName != "" || strings.Contains(strings.ToLower(ord.Title), "domain") {
			domainName := cfg.DomainName
			if domainName == "" {
				domainName = ord.Title
			}
			parts := strings.Split(domainName, ".")
			tld := ""
			if len(parts) > 1 {
				tld = "." + parts[len(parts)-1]
			}

			ns := cfg.Nameservers
			if len(ns) == 0 {
				ns = []string{"ns1.fossbilling.org", "ns2.fossbilling.org"}
			}

			exp := time.Now().AddDate(1, 0, 0)
			if ord.ExpiresAt != nil {
				exp = *ord.ExpiresAt
			}

			epp := cfg.EPPCode
			if epp == "" {
				epp = "EPP-" + strconv.FormatInt(ord.ID*1337%9999+1000, 10)
			}

			list = append(list, DomainRecordDTO{
				ID:          ord.ID,
				DomainName:  domainName,
				TLD:         tld,
				Status:      string(ord.Status),
				Nameservers: ns,
				EPPCode:     epp,
				AutoRenew:   cfg.AutoRenew,
				ExpiresAt:   exp,
			})
		}
	}

	return list, nil
}

// UpdateNameservers updates DNS nameservers for a client domain order
func (s *DomainService) UpdateNameservers(ctx context.Context, clientID int64, orderID int64, nameservers []string) error {
	cleanNS := make([]string, 0, len(nameservers))
	for _, ns := range nameservers {
		trimmed := strings.TrimSpace(ns)
		if trimmed != "" {
			cleanNS = append(cleanNS, trimmed)
		}
	}

	if len(cleanNS) == 0 {
		return ErrInvalidNameservers
	}

	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return ErrDomainNotFound
	}

	if order.ClientID != clientID {
		return ErrUnauthorizedDomain
	}

	var cfg DomainConfig
	if len(order.Config) > 0 {
		_ = json.Unmarshal(order.Config, &cfg)
	}
	cfg.Nameservers = cleanNS

	cfgBytes, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	order.Config = cfgBytes
	return s.orderRepo.Update(ctx, order)
}

// ToggleAutoRenew switches the auto-renew flag for a domain order
func (s *DomainService) ToggleAutoRenew(ctx context.Context, clientID int64, orderID int64) (bool, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return false, ErrDomainNotFound
	}

	if order.ClientID != clientID {
		return false, ErrUnauthorizedDomain
	}

	var cfg DomainConfig
	if len(order.Config) > 0 {
		_ = json.Unmarshal(order.Config, &cfg)
	}
	cfg.AutoRenew = !cfg.AutoRenew

	cfgBytes, err := json.Marshal(cfg)
	if err != nil {
		return false, err
	}

	order.Config = cfgBytes
	if err := s.orderRepo.Update(ctx, order); err != nil {
		return false, err
	}

	return cfg.AutoRenew, nil
}
