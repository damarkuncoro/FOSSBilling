package listener

import (
	"context"
	"fmt"
	"log"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/notification"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
)

type AdminAlertListener struct {
	adminNotifService *notification.AdminNotificationService
}

func NewAdminAlertListener(svc *notification.AdminNotificationService) *AdminAlertListener {
	return &AdminAlertListener{
		adminNotifService: svc,
	}
}

func (l *AdminAlertListener) HandleInvoicePaid(ctx context.Context, e events.Event) error {
	payload := e.Payload.(domain.InvoicePaidPayload)

	// Large payment alert (e.g. > $100)
	if payload.Amount >= 1000000 { // $100.00
		title := "💰 High Value Payment Received"
		msg := fmt.Sprintf("Client #%d has paid Invoice #%d with total %s %s.",
			payload.ClientID, payload.InvoiceID, payload.Currency, payload.Amount.String())

		log.Printf("📢 [Admin Alert] %s", title)
		return l.adminNotifService.CreateAlert(ctx, title, msg, "success", "billing")
	}
	return nil
}

func (l *AdminAlertListener) HandleTicketOpened(ctx context.Context, e events.Event) error {
	payload := e.Payload.(domain.TicketOpenedPayload)

	if payload.Priority == domain.PriorityHigh || payload.Priority == domain.PriorityUrgent {
		title := "🔥 Urgent Support Ticket"
		msg := fmt.Sprintf("A new high-priority ticket was opened by Client #%d: %s", payload.ClientID, payload.Subject)

		return l.adminNotifService.CreateAlert(ctx, title, msg, "danger", "support")
	}
	return nil
}
