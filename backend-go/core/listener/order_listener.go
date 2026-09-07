package listener

import (
	"context"
	"encoding/json"
	"log"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/notification"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/provisioning"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/order"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
)

type OrderListener struct {
	emailService        *notification.EmailService
	orderRepo           domain.OrderRepository
	productRepo         domain.ProductRepository
	clientRepo          domain.ClientRepository
	orderService        *order.OrderService
	registrarRegistry   *provisioning.RegistrarRegistry
	provisionerRegistry *provisioning.ProvisionerRegistry
}

func NewOrderListener(
	emailService *notification.EmailService,
	orderRepo domain.OrderRepository,
	productRepo domain.ProductRepository,
	clientRepo domain.ClientRepository,
	orderService *order.OrderService,
	registrarRegistry *provisioning.RegistrarRegistry,
	provisionerRegistry *provisioning.ProvisionerRegistry,
) *OrderListener {
	return &OrderListener{
		emailService:        emailService,
		orderRepo:           orderRepo,
		productRepo:         productRepo,
		clientRepo:          clientRepo,
		orderService:        orderService,
		registrarRegistry:   registrarRegistry,
		provisionerRegistry: provisionerRegistry,
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
	}

	if orderID == 0 {
		return nil
	}

	order, err := l.orderRepo.GetByID(ctx, orderID)
	if err != nil || order == nil {
		return err
	}

	product, err := l.productRepo.GetByID(ctx, order.ProductID)
	if err != nil || product == nil {
		return err
	}

	log.Printf("📢 [OrderListener] Processing activation for Order #%d (Product Type: %s)", order.ID, product.Type)

	var cfg map[string]interface{}
	if len(order.Config) > 0 {
		_ = json.Unmarshal(order.Config, &cfg)
	}
	if cfg == nil {
		cfg = make(map[string]interface{})
	}

	// Merge product config into temporary cfg if keys don't exist
	if len(product.Config) > 0 {
		var prodCfg map[string]interface{}
		if err := json.Unmarshal(product.Config, &prodCfg); err == nil {
			for k, v := range prodCfg {
				if _, exists := cfg[k]; !exists {
					cfg[k] = v
				}
			}
		}
	}

	// 1. Provisioning based on Product Type
	switch product.Type {
	case domain.ProductTypeDomain:
		if l.registrarRegistry != nil {
			domainName, _ := cfg["domain_name"].(string)
			registrarID, _ := cfg["registrar"].(string)
			if registrarID == "" {
				registrarID = "rdap" // Default
			}

			if domainName != "" {
				reg, err := l.registrarRegistry.Get(registrarID)
				if err != nil {
					// Fallback to email if specified driver fails
					reg, _ = l.registrarRegistry.Get("email")
				}

				if reg != nil {
					// Fetch client for contact info
					var contactInfo map[string]string
					if client, err := l.clientRepo.GetByID(ctx, order.ClientID); err == nil && client != nil {
						contactInfo = map[string]string{
							"first_name": client.FirstName,
							"last_name":  client.LastName,
							"email":      client.Email,
							"company":    client.Company,
							"address1":   client.Address1,
							"address2":   client.Address2,
							"city":       client.City,
							"state":      client.State,
							"postcode":   client.Postcode,
							"country":    client.Country,
							"phone_cc":   client.PhoneCC,
							"phone":      client.Phone,
						}
					}

					regRes, err := reg.RegisterDomain(ctx, provisioning.DomainRegistrationRequest{
						DomainName:  domainName,
						Years:       1,
						ContactInfo: contactInfo,
					})
					if err == nil && regRes != nil {
						cfg["remote_id"] = regRes.AuthCode
						cfg["status"] = "active"
						log.Printf("🌐 [Registrar] Domain %s registered via %s", domainName, registrarID)
					}
				}
			}
		}

	case domain.ProductTypeHosting:
		if l.provisionerRegistry != nil {
			driverID, _ := cfg["server_type"].(string)
			if driverID == "" {
				driverID = "cpanel" // Default
			}

			prov, err := l.provisionerRegistry.Get(driverID)
			if err == nil {
				res, err := prov.Create(ctx, order)
				if err == nil && res.Success {
					cfg["remote_id"] = res.RemoteID
					cfg["account_details"] = res.AccountDetails
					log.Printf("🖥️ [Hosting] Account provisioned: %s via %s", res.RemoteID, driverID)
				} else if err != nil {
					log.Printf("❌ [Hosting] Provisioning failed: %v", err)
				}
			}
		}
	}

	// Update order config with provisioning results
	newCfg, _ := json.Marshal(cfg)
	order.Config = newCfg
	_ = l.orderRepo.Update(ctx, order)

	// 2. Send activation confirmation email
	if l.clientRepo != nil && l.emailService != nil {
		client, err := l.clientRepo.GetByID(ctx, order.ClientID)
		if err == nil && client != nil {
			_ = l.emailService.SendServiceActivatedEmail(ctx, client, order)
		}
	}

	return nil
}

func (l *OrderListener) HandleInvoicePaid(ctx context.Context, e events.Event) error {
	var invID int64
	switch p := e.Payload.(type) {
	case domain.InvoicePaidPayload:
		invID = p.InvoiceID
	case *domain.InvoicePaidPayload:
		if p != nil {
			invID = p.InvoiceID
		}
	}

	if invID == 0 {
		return nil
	}

	log.Printf("📢 [OrderListener] Invoice #%d paid. Activating associated orders...", invID)
	return l.orderService.ActivateOrdersByInvoiceID(ctx, invID)
}
