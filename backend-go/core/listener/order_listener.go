package listener

import (
	"context"
	"log"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/notification"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
)

type OrderListener struct {
	emailService *notification.EmailService
	orderRepo    domain.OrderRepository
	clientRepo   domain.ClientRepository
}

func NewOrderListener(
	emailService *notification.EmailService,
	orderRepo domain.OrderRepository,
	clientRepo domain.ClientRepository,
) *OrderListener {
	return &OrderListener{
		emailService: emailService,
		orderRepo:    orderRepo,
		clientRepo:   clientRepo,
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

	if l.clientRepo != nil {
		client, err := l.clientRepo.GetByID(ctx, order.ClientID)
		if err == nil && client != nil {
			_ = l.emailService.SendServiceActivatedEmail(ctx, client, order)
		}
	}

	return nil
}
