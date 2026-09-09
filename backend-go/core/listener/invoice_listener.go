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

	log.Printf("📢 [Event Listener] Invoice #%d marked as paid. Processing side effects...", invID)

	inv, err := l.invoiceRepo.GetByID(ctx, invID)
	if err != nil || inv == nil {
		return err
	}

	// BUG-19 Fix: Check if this is a deposit invoice and add balance
	isDeposit := false
	for _, item := range inv.Items {
		// Pattern matching for deposit items
		if item.OrderID == nil && (item.Title == "Account Balance Deposit / Top-up" || item.Title == "Topup Saldo") {
			isDeposit = true
			break
		}
	}

	if isDeposit {
		log.Printf("   💰 Invoice #%d recognized as Deposit. Adding %s to client #%d balance...", invID, inv.Total.String(), inv.ClientID)
		err := l.clientRepo.AddBalanceTransaction(ctx, &domain.ClientBalance{
			ClientID:    inv.ClientID,
			Type:        domain.BalanceTypeCredit,
			Amount:      inv.Total,
			Description: "Deposit via Invoice #" + inv.Serie + inv.Nr,
			RelID:       &inv.ID,
		})
		if err != nil {
			log.Printf("   ❌ FAILED to add balance for invoice #%d: %v", invID, err)
		}
	}

	client, err := l.clientRepo.GetByID(ctx, inv.ClientID)
	if err != nil || client == nil {
		return err
	}

	return l.emailService.SendPaymentReceipt(ctx, client, inv, nil)
}
