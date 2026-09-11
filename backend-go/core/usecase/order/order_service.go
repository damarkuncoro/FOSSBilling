package order

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/provisioning"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/metrics"
)

type OrderService struct {
	orderRepo   domain.OrderRepository
	productRepo domain.ProductRepository
	provReg     *provisioning.ProvisionerRegistry
	regReg      *provisioning.RegistrarRegistry
	eventBus    *events.EventBus
}

func NewOrderService(or domain.OrderRepository, pr domain.ProductRepository, prv *provisioning.ProvisionerRegistry, reg *provisioning.RegistrarRegistry, eb ...*events.EventBus) *OrderService {
	var bus *events.EventBus
	if len(eb) > 0 { bus = eb[0] }
	return &OrderService{or, pr, prv, reg, bus}
}

func (s *OrderService) callProv(ctx context.Context, o *domain.Order, fn func(domain.ServiceProvisioner) error) error {
	p, _ := s.productRepo.GetByID(ctx, o.ProductID)
	if p == nil || p.Type != domain.ProductTypeHosting || s.provReg == nil { return nil }
	var cfg map[string]interface{}; _ = json.Unmarshal(o.Config, &cfg)
	drv, _ := cfg["server_type"].(string)
	if drv == "" { drv = "cpanel" }
	prov, err := s.provReg.Get(drv)
	if err != nil { return nil }
	return fn(prov)
}

func (s *OrderService) Activate(ctx context.Context, id int64, from time.Time) (*domain.Order, error) {
	o, err := s.orderRepo.GetByID(ctx, id)
	if err != nil || o.Status == domain.OrderStatusActive { return o, err }
	if o.Status != domain.OrderStatusPendingSetup && o.Status != domain.OrderStatusSuspended { return nil, errors.New("invalid transition") }

	per, err := decimal.ParsePeriod(o.Period)
	if err != nil { return nil, err }
	if from.IsZero() { from = time.Now().UTC() }
	exp := per.CalculateNextDueDate(from)

	o.Status, o.ActivatedAt, o.ExpiresAt, o.NextDueDate = domain.OrderStatusActive, pointer(time.Now().UTC()), &exp, &exp
	_ = s.callProv(ctx, o, func(p domain.ServiceProvisioner) error {
		res, err := p.Create(ctx, o)
		if err == nil && res != nil && res.Success {
			var cfg, det map[string]interface{}
			_ = json.Unmarshal(o.Config, &cfg); _ = json.Unmarshal(res.AccountDetails, &det)
			for k, v := range det { cfg[k] = v }
			o.Config, _ = json.Marshal(cfg)
		}
		return err
	})

	if err := s.orderRepo.Update(ctx, o); err != nil { return nil, err }
	metrics.ActiveOrders.Inc()
	if s.eventBus != nil { s.eventBus.PublishAsync(ctx, events.Event{Type: events.EventOrderActivated, Payload: domain.OrderActivatedPayload{OrderID: o.ID, ClientID: o.ClientID, ProductID: o.ProductID, Title: o.Title, ActivatedAt: *o.ActivatedAt}}) }
	return o, nil
}

func (s *OrderService) Suspend(ctx context.Context, id int64, reason string) (*domain.Order, error) {
	o, err := s.orderRepo.GetByID(ctx, id)
	if err != nil || o.Status != domain.OrderStatusActive { return nil, errors.New("cannot suspend") }
	_ = s.callProv(ctx, o, func(p domain.ServiceProvisioner) error { return p.Suspend(ctx, o, reason) })
	if err := s.orderRepo.UpdateStatus(ctx, id, domain.OrderStatusSuspended, &reason); err != nil { return nil, err }
	metrics.ActiveOrders.Dec()
	if s.eventBus != nil { s.eventBus.PublishAsync(ctx, events.Event{Type: events.EventOrderSuspended, Payload: domain.OrderSuspendedPayload{OrderID: o.ID, ClientID: o.ClientID, Reason: reason, SuspendedAt: time.Now().UTC()}}) }
	return s.orderRepo.GetByID(ctx, id)
}

func (s *OrderService) Unsuspend(ctx context.Context, id int64) (*domain.Order, error) {
	o, err := s.orderRepo.GetByID(ctx, id)
	if err != nil || o.Status != domain.OrderStatusSuspended { return nil, errors.New("cannot unsuspend") }
	_ = s.callProv(ctx, o, func(p domain.ServiceProvisioner) error { return p.Unsuspend(ctx, o) })
	if err := s.orderRepo.UpdateStatus(ctx, id, domain.OrderStatusActive, nil); err != nil { return nil, err }
	return s.orderRepo.GetByID(ctx, id)
}

func (s *OrderService) Renew(ctx context.Context, id int64) (*domain.Order, error) {
	o, err := s.orderRepo.GetByID(ctx, id)
	if err != nil || o.Status == domain.OrderStatusTerminated { return nil, appErrors.ErrExpired }
	per, err := decimal.ParsePeriod(o.Period)
	if err != nil { return nil, err }
	base := time.Now().UTC(); if o.ExpiresAt != nil { base = *o.ExpiresAt }
	exp := per.CalculateNextDueDate(base)
	o.ExpiresAt, o.NextDueDate, o.Status, o.SuspendedAt = &exp, &exp, domain.OrderStatusActive, nil
	if err := s.orderRepo.Update(ctx, o); err != nil { return nil, err }
	return o, nil
}

func (s *OrderService) Cancel(ctx context.Context, id int64, reason string) (*domain.Order, error) {
	o, err := s.orderRepo.GetByID(ctx, id)
	if err != nil || o.Status == domain.OrderStatusTerminated { return o, err }
	_ = s.callProv(ctx, o, func(p domain.ServiceProvisioner) error { return p.Terminate(ctx, o) })
	if err := s.orderRepo.UpdateStatus(ctx, id, domain.OrderStatusCanceled, &reason); err != nil { return nil, err }
	return s.orderRepo.GetByID(ctx, id)
}

func (s *OrderService) GetByIDForClient(ctx context.Context, cID, oID int64) (*domain.Order, error) {
	o, err := s.orderRepo.GetByID(ctx, oID)
	if err != nil || o.ClientID != cID { return nil, appErrors.ErrNotFound }
	return o, nil
}

func (s *OrderService) SyncRemote(ctx context.Context, o *domain.Order) (*domain.ServiceStatus, error) {
	var res *domain.ServiceStatus
	err := s.callProv(ctx, o, func(p domain.ServiceProvisioner) error {
		r, e := p.Sync(ctx, o); res = r; return e
	})
	if res != nil { return res, err }
	return &domain.ServiceStatus{IsActive: o.Status == domain.OrderStatusActive, RemoteState: string(o.Status)}, nil
}

func (s *OrderService) ChangePasswordRemote(ctx context.Context, o *domain.Order, pw string) error {
	return s.callProv(ctx, o, func(p domain.ServiceProvisioner) error {
		if err := p.ChangePassword(ctx, o, pw); err != nil { return err }
		var cfg map[string]interface{}; _ = json.Unmarshal(o.Config, &cfg); cfg["password"] = pw; o.Config, _ = json.Marshal(cfg)
		return s.orderRepo.Update(ctx, o)
	})
}

func (s *OrderService) ActivateOrdersByInvoiceID(ctx context.Context, invID int64) error {
	ords, err := s.orderRepo.ListByInvoiceID(ctx, invID)
	if err != nil { return err }
	for _, o := range ords {
		if o.Status == domain.OrderStatusPendingSetup || o.Status == domain.OrderStatusSuspended { s.Activate(ctx, o.ID, time.Now().UTC()) } else { s.Renew(ctx, o.ID) }
	}
	return nil
}

func (s *OrderService) CheckGracePeriodOverdue(o *domain.Order, grace int, now time.Time) bool {
	if o.Status != domain.OrderStatusActive || o.ExpiresAt == nil { return false }
	return now.After(o.ExpiresAt.AddDate(0, 0, grace))
}

func (s *OrderService) ListByClientID(ctx context.Context, cID int64, l, o int) ([]*domain.Order, int, error) { return s.orderRepo.ListByClientID(ctx, cID, l, o) }

func (s *OrderService) PublishProvisioningFailure(ctx context.Context, orderID int64, errStr string) {
	if s.eventBus == nil { return }
	o, _ := s.orderRepo.GetByID(ctx, orderID)
	if o == nil { return }
	metrics.ProvisioningFailures.Inc()
	s.eventBus.PublishAsync(ctx, events.Event{
		Type: events.EventOrderProvisioningFailed,
		Payload: domain.OrderProvisioningFailedPayload{
			OrderID:   o.ID,
			ProductID: o.ProductID,
			ClientID:  o.ClientID,
			Error:     errStr,
		},
	})
}

func pointer[T any](v T) *T { return &v }
