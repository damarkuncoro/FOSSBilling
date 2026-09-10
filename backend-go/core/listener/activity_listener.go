package listener

import (
	"context"
	"fmt"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/activity"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
)

type ActivityListener struct { svc *activity.ActivityService }

func NewActivityListener(s *activity.ActivityService) *ActivityListener { return &ActivityListener{s} }

func (l *ActivityListener) HandleOrderActivated(ctx context.Context, e events.Event) error {
	p, ok := e.Payload.(domain.OrderActivatedPayload); if !ok { return nil }
	return l.svc.LogClientEvent(ctx, p.ClientID, "order_activated", fmt.Sprintf("Order #%d (%s) activated.", p.OrderID, p.Title), "")
}

func (l *ActivityListener) HandleInvoicePaid(ctx context.Context, e events.Event) error {
	p, ok := e.Payload.(domain.InvoicePaidPayload); if !ok { return nil }
	return l.svc.LogClientEvent(ctx, p.ClientID, "invoice_paid", fmt.Sprintf("Invoice #%d paid (%s %s).", p.InvoiceID, p.Currency, p.Amount.String()), "")
}

func (l *ActivityListener) HandleTicketOpened(ctx context.Context, e events.Event) error {
	p, ok := e.Payload.(domain.TicketOpenedPayload); if !ok { return nil }
	return l.svc.LogClientEvent(ctx, p.ClientID, "ticket_opened", fmt.Sprintf("Ticket opened: %s", p.Subject), "")
}
