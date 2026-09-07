package listener

import (
	"context"
	"fmt"
	"log"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/activity"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
)

type ActivityListener struct {
	activityService *activity.ActivityService
}

func NewActivityListener(activityService *activity.ActivityService) *ActivityListener {
	return &ActivityListener{
		activityService: activityService,
	}
}

func (l *ActivityListener) HandleOrderActivated(ctx context.Context, e events.Event) error {
	payload := e.Payload.(domain.OrderActivatedPayload)
	msg := fmt.Sprintf("Order #%d (%s) has been successfully activated.", payload.OrderID, payload.Title)

	log.Printf("📝 [Activity] Logging Order Activation: %s", msg)
	return l.activityService.LogClientEvent(ctx, payload.ClientID, "order_activated", msg, "")
}

func (l *ActivityListener) HandleInvoicePaid(ctx context.Context, e events.Event) error {
	payload := e.Payload.(domain.InvoicePaidPayload)
	msg := fmt.Sprintf("Invoice #%d paid in full (%s %s).", payload.InvoiceID, payload.Currency, payload.Amount.String())

	log.Printf("📝 [Activity] Logging Payment: %s", msg)
	return l.activityService.LogClientEvent(ctx, payload.ClientID, "invoice_paid", msg, "")
}

func (l *ActivityListener) HandleTicketOpened(ctx context.Context, e events.Event) error {
	payload := e.Payload.(domain.TicketOpenedPayload)
	msg := fmt.Sprintf("New support ticket opened: %s", payload.Subject)

	return l.activityService.LogClientEvent(ctx, payload.ClientID, "ticket_opened", msg, "")
}
