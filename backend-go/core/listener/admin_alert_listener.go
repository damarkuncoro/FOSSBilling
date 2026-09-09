package listener

import (
	"context"
	"fmt"
	"log"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/notification"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/centralalerts"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
)

type AdminAlertListener struct {
	adminNotifService *notification.AdminNotificationService
	telegramService   *centralalerts.TelegramService
}

func NewAdminAlertListener(svc *notification.AdminNotificationService, tg *centralalerts.TelegramService) *AdminAlertListener {
	return &AdminAlertListener{
		adminNotifService: svc,
		telegramService:   tg,
	}
}

func (l *AdminAlertListener) HandleInvoicePaid(ctx context.Context, e events.Event) error {
	payload := e.Payload.(domain.InvoicePaidPayload)

	title := "💰 Payment Received"
	msg := fmt.Sprintf("Client #%d has paid Invoice #%d with total %s %s.",
		payload.ClientID, payload.InvoiceID, payload.Currency, payload.Amount.String())

	// Push to Telegram if configured
	if l.telegramService != nil {
		tgMsg := fmt.Sprintf("<b>%s</b>\n%s", title, msg)
		_ = l.telegramService.SendMessage(tgMsg)
	}

	// Large payment alert (e.g. > $100)
	if payload.Amount >= 1000000 { // $100.00
		return l.adminNotifService.CreateAlert(ctx, title, msg, "success", "billing")
	}
	return nil
}

func (l *AdminAlertListener) HandleTicketOpened(ctx context.Context, e events.Event) error {
	payload := e.Payload.(domain.TicketOpenedPayload)

	title := "🎫 New Support Ticket"
	msg := fmt.Sprintf("Ticket #%d was opened by Client #%d: %s", payload.TicketID, payload.ClientID, payload.Subject)

	// Push to Telegram if configured
	if l.telegramService != nil {
		tgMsg := fmt.Sprintf("<b>%s</b>\n%s\nPriority: %s", title, msg, payload.Priority)
		_ = l.telegramService.SendMessage(tgMsg)
	}

	if payload.Priority == string(domain.PriorityHigh) || payload.Priority == string(domain.PriorityUrgent) {
		title = "🔥 Urgent Support Ticket"
		return l.adminNotifService.CreateAlert(ctx, title, msg, "danger", "support")
	}
	return nil
}
