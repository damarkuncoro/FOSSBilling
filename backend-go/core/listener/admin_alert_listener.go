package listener

import (
	"context"
	"fmt"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/notification"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/centralalerts"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/notifications"
)

type AdminAlertListener struct {
	ans   *notification.AdminNotificationService
	ts    *centralalerts.TelegramService
	wsHub *notifications.WSHub
}

func NewAdminAlertListener(s *notification.AdminNotificationService, t *centralalerts.TelegramService, ws *notifications.WSHub) *AdminAlertListener {
	return &AdminAlertListener{s, t, ws}
}

func (l *AdminAlertListener) HandleInvoicePaid(ctx context.Context, e events.Event) error {
	p, ok := e.Payload.(domain.InvoicePaidPayload); if !ok { return nil }
	m := fmt.Sprintf("Client #%d paid Invoice #%d: %s %s", p.ClientID, p.InvoiceID, p.Currency, p.Amount.String())

	// 1. WebSocket Live Push
	if l.wsHub != nil {
		_ = l.wsHub.Broadcast(ctx, map[string]any{"type": "invoice_paid", "title": "💰 Payment Received", "message": m})
	}

	// 2. Telegram Alert
	if l.ts != nil {
		priority := p.Amount >= 1000000 // High priority for large payments
		_ = l.ts.SendAlert("Payment Received", m, priority)
	}

	// 3. Database Alert (if large amount)
	if p.Amount >= 1000000 {
		return l.ans.CreateAlert(ctx, "💰 Payment Received", m, "success", "billing")
	}
	return nil
}

func (l *AdminAlertListener) HandleTicketOpened(ctx context.Context, e events.Event) error {
	p, ok := e.Payload.(domain.TicketOpenedPayload); if !ok { return nil }
	m := fmt.Sprintf("Ticket #%d by Client #%d: %s (Priority: %s)", p.TicketID, p.ClientID, p.Subject, p.Priority)

	// 1. WebSocket Live Push
	if l.wsHub != nil {
		_ = l.wsHub.Broadcast(ctx, map[string]any{"type": "ticket_opened", "title": "🎫 New Support Ticket", "message": m, "priority": p.Priority})
	}

	// 2. Telegram Alert
	if l.ts != nil {
		priority := p.Priority == string(domain.PriorityHigh) || p.Priority == string(domain.PriorityUrgent)
		_ = l.ts.SendAlert("New Support Ticket", m, priority)
	}

	// 3. Database Alert (if high priority)
	if p.Priority == string(domain.PriorityHigh) || p.Priority == string(domain.PriorityUrgent) {
		return l.ans.CreateAlert(ctx, "🔥 Urgent Ticket", m, "danger", "support")
	}
	return nil
}

func (l *AdminAlertListener) HandleOrderProvisioningFailed(ctx context.Context, e events.Event) error {
	p, ok := e.Payload.(domain.OrderProvisioningFailedPayload); if !ok { return nil }
	m := fmt.Sprintf("Order #%d failed to provision: %s", p.OrderID, p.Error)

	if l.wsHub != nil {
		_ = l.wsHub.Broadcast(ctx, map[string]any{"type": "provisioning_failed", "title": "❌ Provisioning Failed", "message": m, "priority": "high"})
	}

	if l.ts != nil {
		_ = l.ts.SendAlert("Provisioning Failed", m, true)
	}

	return l.ans.CreateAlert(ctx, "❌ Provisioning Error", m, "danger", "orders")
}
