package listener_test

import (
	"context"
	"testing"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/listener"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/notification"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/provisioning"
	orderUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/order"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/mailer"
)

func TestListeners_Flows(t *testing.T) {
	ctx := context.Background()
	mockMailer := mailer.NewMockMailer()
	tplRepo := memory.NewMockEmailTemplateRepository()
	emailService := notification.NewEmailService(mockMailer, tplRepo, "no-reply@fossbilling.org", "FOSSBilling")

	clientRepo := memory.NewMockClientRepository()
	orderRepo := memory.NewMockOrderRepository()
	invoiceRepo := memory.NewMockInvoiceRepository()
	productRepo := memory.NewMockProductRepository()

	eventBus := events.NewEventBus()

	regRegistry := provisioning.NewRegistrarRegistry()
	registrar := provisioning.NewMockRegistrarDriver()
	regRegistry.Register("rdap", registrar)
	provRegistry := provisioning.NewProvisionerRegistry()

	orderService := orderUsecase.NewOrderService(orderRepo, productRepo, provRegistry, regRegistry, eventBus)

	// Seed client
	_ = clientRepo.Create(ctx, &domain.Client{
		ID:        1,
		Email:     "listener.user@example.com",
		FirstName: "Listener",
		LastName:  "Tester",
		Currency:  "USD",
		Status:    domain.ClientStatusActive,
	})

	// Seed product
	_ = productRepo.Create(ctx, &domain.Product{
		ID:   10,
		Type: domain.ProductTypeDomain,
		Name: "Domain",
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

	// Seed domain order
	domainOrd := &domain.Order{
		ClientID:  1,
		ProductID: 10,
		Title:     "Domain Registration: myawesomeapp2026.com",
		Status:    domain.OrderStatusActive,
		Config:    []byte(`{"domain_name":"myawesomeapp2026.com","nameservers":["ns1.fossbilling.org","ns2.fossbilling.org"]}`),
		CreatedAt: now,
		UpdatedAt: now,
	}
	_ = orderRepo.Create(ctx, domainOrd)

	// 1. Client Listener
	clientListener := listener.NewClientListener(emailService, clientRepo)
	err := clientListener.HandleClientRegistered(ctx, events.Event{
		Type: events.EventClientRegistered,
		Payload: domain.ClientRegisteredPayload{
			ClientID: 1,
			Email:    "listener.user@example.com",
		},
	})
	if err != nil {
		t.Fatalf("ClientListener failed: %v", err)
	}

	// 2. Invoice Listener
	invoiceListener := listener.NewInvoiceListener(emailService, invoiceRepo, clientRepo)
	err = invoiceListener.HandleInvoicePaid(ctx, events.Event{
		Type:    events.EventInvoicePaid,
		Payload: domain.InvoicePaidPayload{InvoiceID: 1},
	})
	if err != nil {
		t.Fatalf("InvoiceListener failed: %v", err)
	}

	// 3. Order Listener
	orderListenerWithReg := listener.NewOrderListener(emailService, orderRepo, productRepo, clientRepo, orderService, regRegistry, provRegistry)
	err = orderListenerWithReg.HandleOrderActivated(ctx, events.Event{
		Type:    events.EventOrderActivated,
		Payload: domain.OrderActivatedPayload{OrderID: domainOrd.ID},
	})
	if err != nil {
		t.Fatalf("OrderListener with registrar failed: %v", err)
	}

	if mockMailer.GetSentCount() < 3 {
		t.Fatalf("Expected at least 3 emails sent via listeners, got %d", mockMailer.GetSentCount())
	}
}
