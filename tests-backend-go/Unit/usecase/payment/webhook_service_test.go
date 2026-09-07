package payment_test

import (
	"context"
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/payment"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
)

func setupWebhookService() (*payment.WebhookService, *memory.MockTransactionRepository, *memory.MockInvoiceRepository, *events.EventBus) {
	txnRepo := memory.NewMockTransactionRepository()
	invRepo := memory.NewMockInvoiceRepository()
	eventBus := events.NewEventBus()

	webhookService := payment.NewWebhookService(txnRepo, invRepo, eventBus)

	return webhookService, txnRepo, invRepo, eventBus
}

func TestWebhookService_HandlePaymentWebhookAndEventPublish(t *testing.T) {
	ctx := context.Background()
	service, _, invRepo, eventBus := setupWebhookService()

	// Capture event
	eventReceived := false
	eventBus.Subscribe(events.EventInvoicePaid, func(ctx context.Context, e events.Event) error {
		payload := e.Payload.(domain.InvoicePaidPayload)
		if payload.InvoiceID != 0 {
			eventReceived = true
		}
		return nil
	})

	// 1. Create Unpaid Invoice
	invoice := &domain.Invoice{
		ClientID: 1,
		Status:   domain.InvoiceStatusUnpaid,
		Currency: "USD",
		Total:    decimal.FromFloat(25.00),
	}
	_ = invRepo.Create(ctx, invoice, []domain.InvoiceItem{})

	// 2. Receive Inbound Webhook
	webhook := payment.WebhookPayload{
		GatewayID: "stripe",
		TxnID:     "ch_3M456xyz789",
		InvoiceID: invoice.ID,
		Amount:    decimal.FromFloat(25.00),
		Currency:  "USD",
		Raw:       []byte(`{"status":"succeeded"}`),
	}

	txn, err := service.HandlePaymentWebhook(ctx, webhook)
	if err != nil {
		t.Fatalf("HandlePaymentWebhook failed: %v", err)
	}

	if txn.Status != domain.TransactionStatusComplete {
		t.Errorf("Txn status = %s; want complete", txn.Status)
	}

	// 3. Verify Invoice is marked as Paid
	updatedInv, _ := invRepo.GetByID(ctx, invoice.ID)
	if updatedInv.Status != domain.InvoiceStatusPaid {
		t.Errorf("Invoice status = %s; want paid", updatedInv.Status)
	}

	// 4. Verify Event was Published
	if !eventReceived {
		t.Error("Expected EventInvoicePaid to be published")
	}

	// 5. Test Idempotency
	_, duplicateErr := service.HandlePaymentWebhook(ctx, webhook)
	if duplicateErr != payment.ErrDuplicateTransaction {
		t.Errorf("Expected ErrDuplicateTransaction on replayed webhook, got: %v", duplicateErr)
	}
}
