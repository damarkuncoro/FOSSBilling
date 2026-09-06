package listener_test

import (
	"context"
	"testing"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/listener"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/notification"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/mailer"
)

func TestListeners_Flows(t *testing.T) {
	ctx := context.Background()
	mockMailer := mailer.NewMockMailer()
	emailService := notification.NewEmailService(mockMailer, "no-reply@fossbilling.org", "FOSSBilling")

	clientRepo := memory.NewMockClientRepository()
	orderRepo := memory.NewMockOrderRepository()
	invoiceRepo := memory.NewMockInvoiceRepository()

	// Seed client
	_ = clientRepo.Create(ctx, &domain.Client{
		ID:        1,
		Email:     "listener.user@example.com",
		FirstName: "Listener",
		LastName:  "Tester",
		Currency:  "USD",
		Status:    domain.ClientStatusActive,
	})

	// Seed invoice
	now := time.Now().UTC()
	_ = invoiceRepo.Create(ctx, &domain.Invoice{
		ID:        1,
		Serie:     "INV",
		Nr:        "2026-0001",
		ClientID:  1,
		Status:    domain.InvoiceStatusPaid,
		Currency:  "USD",
		Total:     decimal.FromFloat(100.0),
		DueAt:     now,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil)

	// 1. Client Listener
	clientListener := listener.NewClientListener(emailService, clientRepo)
	err := clientListener.HandleClientRegistered(ctx, events.Event{
		Type:    events.EventClientRegistered,
		Payload: map[string]interface{}{"client_id": int64(1)},
	})
	if err != nil {
		t.Fatalf("ClientListener failed: %v", err)
	}

	// 2. Invoice Listener
	invoiceListener := listener.NewInvoiceListener(emailService, invoiceRepo, clientRepo)
	err = invoiceListener.HandleInvoicePaid(ctx, events.Event{
		Type:    events.EventInvoicePaid,
		Payload: map[string]interface{}{"invoice_id": int64(1)},
	})
	if err != nil {
		t.Fatalf("InvoiceListener failed: %v", err)
	}

	// Seed order
	ord := &domain.Order{
		ClientID:  1,
		ProductID: 101,
		Title:     "Cloud VPS Premium",
		Status:    domain.OrderStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}
	_ = orderRepo.Create(ctx, ord)

	// 3. Order Listener
	orderListener := listener.NewOrderListener(emailService, orderRepo, clientRepo)
	err = orderListener.HandleOrderActivated(ctx, events.Event{
		Type:    events.EventOrderActivated,
		Payload: map[string]interface{}{"order_id": ord.ID},
	})
	if err != nil {
		t.Fatalf("OrderListener failed: %v", err)
	}

	if mockMailer.GetSentCount() < 3 {
		t.Fatalf("Expected at least 3 emails sent via listeners, got %d", mockMailer.GetSentCount())
	}
}
