package listener

import (
	"context"
	"encoding/json"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/notification"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/provisioning"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/order"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
)

type OrderListener struct {
	es *notification.EmailService; or domain.OrderRepository; pr domain.ProductRepository; cr domain.ClientRepository; os *order.OrderService; rr *provisioning.RegistrarRegistry; prr *provisioning.ProvisionerRegistry
}

func NewOrderListener(es *notification.EmailService, or domain.OrderRepository, pr domain.ProductRepository, cr domain.ClientRepository, os *order.OrderService, rr *provisioning.RegistrarRegistry, prr *provisioning.ProvisionerRegistry) *OrderListener {
	return &OrderListener{es, or, pr, cr, os, rr, prr}
}

func (l *OrderListener) HandleOrderActivated(ctx context.Context, e events.Event) error {
	var id int64
	if p, ok := e.Payload.(domain.OrderActivatedPayload); ok { id = p.OrderID } else if p, ok := e.Payload.(*domain.OrderActivatedPayload); ok && p != nil { id = p.OrderID }
	if id == 0 { return nil }

	o, err := l.or.GetByID(ctx, id); if err != nil { return err }
	p, err := l.pr.GetByID(ctx, o.ProductID); if err != nil { return err }

	var cfg map[string]any; _ = json.Unmarshal(o.Config, &cfg); if cfg == nil { cfg = make(map[string]any) }
	if len(p.Config) > 0 { var pc map[string]any; if err := json.Unmarshal(p.Config, &pc); err == nil { for k, v := range pc { if _, ex := cfg[k]; !ex { cfg[k] = v } } } }

	ok := true
	if p.Type == domain.ProductTypeDomain && l.rr != nil {
		rid, _ := cfg["registrar"].(string); if rid == "" { rid = "rdap" }
		if reg, _ := l.rr.Get(rid); reg != nil {
			c, _ := l.cr.GetByID(ctx, o.ClientID)
			con := map[string]string{"first_name": c.FirstName, "last_name": c.LastName, "email": c.Email}
			dnm, _ := cfg["domain_name"].(string)
			if res, err := reg.RegisterDomain(ctx, provisioning.DomainRegistrationRequest{DomainName: dnm, Years: 1, ContactInfo: con}); err == nil {
				cfg["remote_id"], cfg["status"] = res.AuthCode, "active"
			} else { ok = false }
		}
	} else if p.Type == domain.ProductTypeHosting && l.prr != nil {
		did, _ := cfg["server_type"].(string); if did == "" { did = "cpanel" }
		if prov, _ := l.prr.Get(did); prov != nil {
			if res, err := prov.Create(ctx, o); err == nil {
				cfg["remote_id"], cfg["account_details"] = res.RemoteID, res.AccountDetails
			} else { ok = false }
		}
	}

	if !ok { return l.or.UpdateStatus(ctx, o.ID, domain.OrderStatusPendingSetup, nil) }
	o.Config, _ = json.Marshal(cfg); _ = l.or.Update(ctx, o)
	c, _ := l.cr.GetByID(ctx, o.ClientID); if c != nil { _ = l.es.SendServiceActivatedEmail(ctx, c, o) }
	return nil
}

func (l *OrderListener) HandleInvoicePaid(ctx context.Context, e events.Event) error {
	var id int64
	if p, ok := e.Payload.(domain.InvoicePaidPayload); ok { id = p.InvoiceID } else if p, ok := e.Payload.(*domain.InvoicePaidPayload); ok && p != nil { id = p.InvoiceID }
	if id == 0 { return nil }
	return l.os.ActivateOrdersByInvoiceID(ctx, id)
}
