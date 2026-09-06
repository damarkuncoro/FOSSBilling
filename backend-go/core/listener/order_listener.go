package listener

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/notification"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/provisioning"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
)

type OrderListener struct {
	emailService *notification.EmailService
	orderRepo    domain.OrderRepository
	clientRepo   domain.ClientRepository
	registrar    provisioning.RegistrarDriver
}

func NewOrderListener(
	emailService *notification.EmailService,
	orderRepo domain.OrderRepository,
	clientRepo domain.ClientRepository,
	registrar ...provisioning.RegistrarDriver,
) *OrderListener {
	var reg provisioning.RegistrarDriver
	if len(registrar) > 0 {
		reg = registrar[0]
	}
	return &OrderListener{
		emailService: emailService,
		orderRepo:    orderRepo,
		clientRepo:   clientRepo,
		registrar:    reg,
	}
}

func (l *OrderListener) HandleOrderActivated(ctx context.Context, e events.Event) error {
	var orderID int64

	switch p := e.Payload.(type) {
	case domain.OrderActivatedPayload:
		orderID = p.OrderID
	case *domain.OrderActivatedPayload:
		if p != nil {
			orderID = p.OrderID
		}
	case map[string]interface{}:
		if id, ok := p["order_id"].(int64); ok {
			orderID = id
		}
	}

	if orderID == 0 {
		return nil
	}

	log.Printf("📢 [Event Listener] Order #%d successfully activated", orderID)

	order, err := l.orderRepo.GetByID(ctx, orderID)
	if err != nil || order == nil {
		return err
	}

	// 1. Provision Domain if this is a domain registration order
	var cfg map[string]interface{}
	if len(order.Config) > 0 {
		_ = json.Unmarshal(order.Config, &cfg)
	}
	if cfg == nil {
		cfg = make(map[string]interface{})
	}

	domainName, _ := cfg["domain_name"].(string)
	if domainName != "" && l.registrar != nil {
		var ns []string
		if rawNs, ok := cfg["nameservers"].([]interface{}); ok {
			for _, n := range rawNs {
				if s, ok := n.(string); ok && s != "" {
					ns = append(ns, s)
				}
			}
		}

		regRes, err := l.registrar.RegisterDomain(ctx, provisioning.DomainRegistrationRequest{
			DomainName:  domainName,
			Years:       1,
			Nameservers: ns,
		})
		if err == nil && regRes != nil {
			cfg["status"] = "active"
			cfg["registered_at"] = regRes.RegisteredAt.Format(time.RFC3339)
			cfg["expires_at"] = regRes.ExpiresAt.Format(time.RFC3339)
			cfg["epp_code"] = regRes.AuthCode
			if len(regRes.Nameservers) > 0 {
				cfg["nameservers"] = regRes.Nameservers
			}
			newCfg, _ := json.Marshal(cfg)
			order.Config = newCfg
			_ = l.orderRepo.Update(ctx, order)
			log.Printf("🌐 [Registrar] Successfully registered domain %s (EPP: %s)", domainName, regRes.AuthCode)
		}
	}

	// 2. Send activation confirmation email
	if l.clientRepo != nil && l.emailService != nil {
		client, err := l.clientRepo.GetByID(ctx, order.ClientID)
		if err == nil && client != nil {
			_ = l.emailService.SendServiceActivatedEmail(ctx, client, order)
		}
	}

	return nil
}
