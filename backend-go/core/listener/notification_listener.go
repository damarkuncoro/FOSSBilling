package listener

import (
	"context"
	"fmt"
	"log"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/notification"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
)

type NotificationListener struct {
	notifService *notification.NotificationService
}

func NewNotificationListener(notifService *notification.NotificationService) *NotificationListener {
	return &NotificationListener{
		notifService: notifService,
	}
}

func (l *NotificationListener) HandleInvoicePaid(ctx context.Context, e events.Event) error {
	payload := e.Payload.(domain.InvoicePaidPayload)
	title := "Payment Received"
	msg := fmt.Sprintf("Your payment of %s %s for Invoice #%d has been confirmed.", payload.Currency, payload.Amount.String(), payload.InvoiceID)

	log.Printf("🔔 [Notification] Sending to Client #%d: %s", payload.ClientID, title)
	return l.notifService.CreateNotification(ctx, payload.ClientID, title, msg, "success")
}

func (l *NotificationListener) HandleOrderActivated(ctx context.Context, e events.Event) error {
	payload := e.Payload.(domain.OrderActivatedPayload)
	title := "Service Activated"
	msg := fmt.Sprintf("Your order #%d (%s) is now active and ready to use.", payload.OrderID, payload.Title)

	return l.notifService.CreateNotification(ctx, payload.ClientID, title, msg, "info")
}
