package payment

import (
	"context"
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/usecase/order"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

func setupWebhookService() (*WebhookService, *memory.MockTransactionRepository, *memory.MockInvoiceRepository, *memory.MockOrderRepository) {
	txnRepo := memory.NewMockTransactionRepository()
	invRepo := memory.NewMockInvoiceRepository()
	orderRepo := memory.NewMockOrderRepository()

	orderService := order.NewOrderService(orderRepo)
	webhookService := NewWebhookService(txnRepo, invRepo, orderService, orderRepo)

	return webhookService, txnRepo, invRepo, orderRepo
}

func TestWebhookService_HandlePaymentWebhookAndAutoActivate(t *testing.T) {
	ctx := context.Background()
	service, _, invRepo, orderRepo := setupWebhookService()

	// 1. Create Pending Order
	testOrder := &domain.Order{
		ClientID:  1,
		ProductID: 10,
		Title:     "cPanel Web Hosting",
		Period:    "1M",
		Price:     decimal.FromFloat(25.00),
		Currency:  "USD",
		Status:    domain.OrderStatusPendingSetup,
	}
	_ = orderRepo.Create(ctx, testOrder)

	// 2. Create Unpaid Invoice with linked Order
	invoice := &domain.Invoice{
		ClientID: 1,
		Status:   domain.InvoiceStatusUnpaid,
		Currency: "USD",
		Total:    decimal.FromFloat(25.00),
	}
	items := []domain.InvoiceItem{
		{OrderID: &testOrder.ID, Title: "cPanel Web Hosting", Price: decimal.FromFloat(25.00), Quantity: 1},
	}
	_ = invRepo.Create(ctx, invoice, items)

	// 3. Receive Inbound Webhook (e.g. from Stripe / Midtrans)
	webhook := WebhookPayload{
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

	// 4. Verify Invoice is marked as Paid
	updatedInv, _ := invRepo.GetByID(ctx, invoice.ID)
	if updatedInv.Status != domain.InvoiceStatusPaid {
		t.Errorf("Invoice status = %s; want paid", updatedInv.Status)
	}
	if updatedInv.PaidAt == nil {
		t.Error("Expected Invoice PaidAt timestamp to be set")
	}

	// 5. Verify Order was AUTOMATICALLY Activated
	activatedOrder, _ := orderRepo.GetByID(ctx, testOrder.ID)
	if activatedOrder.Status != domain.OrderStatusActive {
		t.Errorf("Order status = %s; want active", activatedOrder.Status)
	}
	if activatedOrder.ActivatedAt == nil || activatedOrder.ExpiresAt == nil {
		t.Error("Expected Order ActivatedAt and ExpiresAt to be set")
	}

	// 6. Test Idempotency: Duplicate webhook with same TxnID should be rejected
	_, duplicateErr := service.HandlePaymentWebhook(ctx, webhook)
	if duplicateErr != ErrDuplicateTransaction {
		t.Errorf("Expected ErrDuplicateTransaction on replayed webhook, got: %v", duplicateErr)
	}
}
