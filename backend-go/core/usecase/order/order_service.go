package order

import (
	"context"
	"errors"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
)

var (
	ErrInvalidStatusTransition = errors.New("invalid order status transition")
)

type OrderService struct {
	orderRepo domain.OrderRepository
	eventBus  *events.EventBus
}

func NewOrderService(orderRepo domain.OrderRepository, eventBus ...*events.EventBus) *OrderService {
	var bus *events.EventBus
	if len(eventBus) > 0 {
		bus = eventBus[0]
	}
	return &OrderService{
		orderRepo: orderRepo,
		eventBus:  bus,
	}
}

// Activate transitions order from PendingSetup/Suspended to Active and calculates expiry & due date
func (s *OrderService) Activate(ctx context.Context, orderID int64, fromDate time.Time) (*domain.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order.Status == domain.OrderStatusActive {
		return order, nil
	}

	if order.Status != domain.OrderStatusPendingSetup && order.Status != domain.OrderStatusSuspended {
		return nil, ErrInvalidStatusTransition
	}

	period, err := decimal.ParsePeriod(order.Period)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	if fromDate.IsZero() {
		fromDate = now
	}

	expiresAt := period.CalculateNextDueDate(fromDate)
	nextDueDate := expiresAt

	order.Status = domain.OrderStatusActive
	order.ActivatedAt = &now
	order.ExpiresAt = &expiresAt
	order.NextDueDate = &nextDueDate
	order.SuspendedAt = nil
	order.SuspensionReason = nil

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	// Publish Event
	if s.eventBus != nil {
		s.eventBus.PublishAsync(ctx, events.Event{
			Type: events.EventOrderActivated,
			Payload: domain.OrderActivatedPayload{
				OrderID:     order.ID,
				ClientID:    order.ClientID,
				ProductID:   order.ProductID,
				Title:       order.Title,
				ActivatedAt: now,
			},
		})
	}

	return order, nil
}

// Suspend transitions order from Active to Suspended with a reason
func (s *OrderService) Suspend(ctx context.Context, orderID int64, reason string) (*domain.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order.Status != domain.OrderStatusActive {
		return nil, ErrInvalidStatusTransition
	}

	if err := s.orderRepo.UpdateStatus(ctx, orderID, domain.OrderStatusSuspended, &reason); err != nil {
		return nil, err
	}

	// Publish Event
	if s.eventBus != nil {
		s.eventBus.PublishAsync(ctx, events.Event{
			Type: events.EventOrderSuspended,
			Payload: domain.OrderSuspendedPayload{
				OrderID:     order.ID,
				ClientID:    order.ClientID,
				Reason:      reason,
				SuspendedAt: time.Now().UTC(),
			},
		})
	}

	return s.orderRepo.GetByID(ctx, orderID)
}

// Unsuspend transitions order from Suspended back to Active
func (s *OrderService) Unsuspend(ctx context.Context, orderID int64) (*domain.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order.Status != domain.OrderStatusSuspended {
		return nil, ErrInvalidStatusTransition
	}

	if err := s.orderRepo.UpdateStatus(ctx, orderID, domain.OrderStatusActive, nil); err != nil {
		return nil, err
	}

	return s.orderRepo.GetByID(ctx, orderID)
}

// Renew extends order expiry date and next due date by its billing period
func (s *OrderService) Renew(ctx context.Context, orderID int64) (*domain.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order.Status == domain.OrderStatusTerminated || order.Status == domain.OrderStatusCanceled {
		return nil, appErrors.ErrOrderExpired
	}

	period, err := decimal.ParsePeriod(order.Period)
	if err != nil {
		return nil, err
	}

	baseDate := time.Now().UTC()
	if order.ExpiresAt != nil {
		baseDate = *order.ExpiresAt
	}

	newExpiry := period.CalculateNextDueDate(baseDate)
	order.ExpiresAt = &newExpiry
	order.NextDueDate = &newExpiry
	order.Status = domain.OrderStatusActive
	order.SuspendedAt = nil
	order.SuspensionReason = nil

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

// CheckGracePeriodOverdue returns true if order expiration + gracePeriodDays is before now
func (s *OrderService) CheckGracePeriodOverdue(order *domain.Order, gracePeriodDays int, now time.Time) bool {
	if order.Status != domain.OrderStatusActive || order.ExpiresAt == nil {
		return false
	}
	graceDeadline := order.ExpiresAt.AddDate(0, 0, gracePeriodDays)
	return now.After(graceDeadline)
}

// ListByClientID retrieves orders belonging to a specific client with pagination
func (s *OrderService) ListByClientID(ctx context.Context, clientID int64, limit, offset int) ([]*domain.Order, int, error) {
	return s.orderRepo.ListByClientID(ctx, clientID, limit, offset)
}

// GetByIDForClient retrieves an order ensuring ownership matches clientID
func (s *OrderService) GetByIDForClient(ctx context.Context, clientID, orderID int64) (*domain.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order.ClientID != clientID {
		return nil, appErrors.ErrNotFound
	}
	return order, nil
}

// Cancel transitions order to Canceled status with a reason
func (s *OrderService) Cancel(ctx context.Context, orderID int64, reason string) (*domain.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order.Status == domain.OrderStatusTerminated || order.Status == domain.OrderStatusCanceled {
		return order, nil
	}

	if err := s.orderRepo.UpdateStatus(ctx, orderID, domain.OrderStatusCanceled, &reason); err != nil {
		return nil, err
	}

	return s.orderRepo.GetByID(ctx, orderID)
}

// CancelForClient cancels an order ensuring ownership matches clientID
func (s *OrderService) CancelForClient(ctx context.Context, clientID, orderID int64, reason string) (*domain.Order, error) {
	order, err := s.GetByIDForClient(ctx, clientID, orderID)
	if err != nil {
		return nil, err
	}
	return s.Cancel(ctx, order.ID, reason)
}

// ActivateOrdersByInvoiceID activates all orders linked to the given invoice
func (s *OrderService) ActivateOrdersByInvoiceID(ctx context.Context, invoiceID int64) error {
	orders, err := s.orderRepo.ListByInvoiceID(ctx, invoiceID)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	for _, ord := range orders {
		if ord.Status == domain.OrderStatusPendingSetup || ord.Status == domain.OrderStatusSuspended {
			_, _ = s.Activate(ctx, ord.ID, now)
		} else if ord.Status == domain.OrderStatusActive {
			_, _ = s.Renew(ctx, ord.ID)
		}
	}
	return nil
}
