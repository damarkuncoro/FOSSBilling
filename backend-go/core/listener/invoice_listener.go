package listener

import (
	"context"
	"log"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/notification"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
)

type InvoiceListener struct {
	emailService *notification.EmailService
	invoiceRepo  domain.InvoiceRepository
	clientRepo   domain.ClientRepository
}

func NewInvoiceListener(
	emailService *notification.EmailService,
	invoiceRepo domain.InvoiceRepository,
	clientRepo domain.ClientRepository,
) *InvoiceListener {
	return &InvoiceListener{
		emailService: emailService,
		invoiceRepo:  invoiceRepo,
		clientRepo:   clientRepo,
	}
}

func (l *InvoiceListener) HandleInvoicePaid(ctx context.Context, e events.Event) error {
	var invID int64

	switch p := e.Payload.(type) {
	case domain.InvoicePaidPayload:
		invID = p.InvoiceID
	case *domain.InvoicePaidPayload:
		if p != nil {
			invID = p.InvoiceID
		}
	case map[string]interface{}:
		if id, ok := p["invoice_id"].(int64); ok {
			invID = id
		}
	}

	if invID == 0 {
		return nil
	}

	log.Printf("📢 [Event Listener] Invoice #%d marked as paid. Dispatching receipt...", invID)

	inv, err := l.invoiceRepo.GetByID(ctx, invID)
	if err != nil || inv == nil {
		return err
	}

	client, err := l.clientRepo.GetByID(ctx, inv.ClientID)
	if err != nil || client == nil {
		return err
	}

	return l.emailService.SendPaymentReceipt(ctx, client, inv, nil)
}
