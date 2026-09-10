package listener

import (
	"context"
	"fmt"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/notification"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
)

type NotificationListener struct { ns *notification.NotificationService }

func NewNotificationListener(s *notification.NotificationService) *NotificationListener { return &NotificationListener{s} }

func (l *NotificationListener) HandleInvoicePaid(ctx context.Context, e events.Event) error {
	p, ok := e.Payload.(domain.InvoicePaidPayload); if !ok { return nil }
	return l.ns.CreateNotification(ctx, p.ClientID, "Payment Received", fmt.Sprintf("Payment of %s %s for Invoice #%d confirmed.", p.Currency, p.Amount.String(), p.InvoiceID), "success")
}

func (l *NotificationListener) HandleOrderActivated(ctx context.Context, e events.Event) error {
	p, ok := e.Payload.(domain.OrderActivatedPayload); if !ok { return nil }
	return l.ns.CreateNotification(ctx, p.ClientID, "Service Activated", fmt.Sprintf("Order #%d (%s) is now active.", p.OrderID, p.Title), "info")
}
