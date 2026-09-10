package listener

import (
	"context"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/notification"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
)

type InvoiceListener struct {
	es *notification.EmailService; ir domain.InvoiceRepository; cr domain.ClientRepository
}

func NewInvoiceListener(es *notification.EmailService, ir domain.InvoiceRepository, cr domain.ClientRepository) *InvoiceListener {
	return &InvoiceListener{es, ir, cr}
}

func (l *InvoiceListener) HandleInvoicePaid(ctx context.Context, e events.Event) error {
	p, ok := e.Payload.(domain.InvoicePaidPayload); if !ok { return nil }
	inv, err := l.ir.GetByID(ctx, p.InvoiceID); if err != nil { return err }

	dep := false; for _, it := range inv.Items { if it.OrderID == nil && (it.Title == "Account Balance Deposit / Top-up" || it.Title == "Topup Saldo") { dep = true; break } }
	if dep { _ = l.cr.AddBalanceTransaction(ctx, &domain.ClientBalance{ClientID: inv.ClientID, Type: domain.BalanceTypeCredit, Amount: inv.Total, Description: "Deposit #" + inv.Serie + inv.Nr, RelID: &inv.ID}) }

	c, _ := l.cr.GetByID(ctx, inv.ClientID)
	if c != nil { _ = l.es.SendPaymentReceipt(ctx, c, inv, nil) }
	return nil
}
