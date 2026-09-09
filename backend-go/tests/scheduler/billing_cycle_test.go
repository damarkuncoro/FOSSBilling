package scheduler_test

import (
	"context"
	"testing"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/scheduler"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/billing"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/order"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

func TestBillingCycle_RenewalAndSuspension(t *testing.T) {
	ctx := context.Background()

	// 1. Setup
	orderRepo := memory.NewMockOrderRepository()
	invoiceRepo := memory.NewMockInvoiceRepository()
	clientRepo := memory.NewMockClientRepository()
	taxRepo := memory.NewMockTaxRepository()
	productRepo := memory.NewMockProductRepository()

	taxCalc := billing.NewTaxCalculator(taxRepo)
	orderService := order.NewOrderService(orderRepo, productRepo, nil, nil, nil)
	invoiceService := billing.NewInvoiceService(invoiceRepo, clientRepo, taxCalc, nil, nil)

	cron := scheduler.NewCronService(orderRepo, orderService, invoiceService)

	// 2. Prepare Data
	_ = clientRepo.Create(ctx, &domain.Client{ID: 1, Email: "test@client.com", Currency: "USD", Country: "US"})

	now := time.Now().UTC()
	dueDate := now.AddDate(0, 0, 5) // Due in 5 days
	expiryDate := dueDate

	ord := &domain.Order{
		ID: 1, ClientID: 1, Status: domain.OrderStatusActive,
		Title: "Test VPS", Price: decimal.FromFloat(10), Currency: "USD",
		Period: "1M", NextDueDate: &dueDate, ExpiresAt: &expiryDate,
	}
	_ = orderRepo.Create(ctx, ord)

	t.Run("Generate Renewal Invoice", func(t *testing.T) {
		res, err := cron.GenerateRenewalInvoicesBatch(ctx, 7)
		if err != nil {
			t.Fatalf("Cron failed: %v", err)
		}

		if res.SuccessCount != 1 {
			t.Errorf("Expected 1 success, got %d. Errors: %v", res.SuccessCount, res.Errors)
		}

		invoices, _, _ := invoiceRepo.ListByClientID(ctx, 1, 10, 0)
		if len(invoices) != 1 {
			t.Errorf("Expected 1 invoice to be generated, got %d", len(invoices))
		}
	})

	t.Run("Idempotency - Do not duplicate invoice", func(t *testing.T) {
		// Run again immediately
		res, _ := cron.GenerateRenewalInvoicesBatch(ctx, 7)

		if res.SuccessCount != 0 {
			t.Errorf("Idempotency FAIL: Expected 0 new invoices, got %d. BUG-33 confirm.", res.SuccessCount)
		}

		invoices, _, _ := invoiceRepo.ListByClientID(ctx, 1, 10, 0)
		if len(invoices) > 1 {
			t.Errorf("Invoices duplicated! Found %d", len(invoices))
		}
	})

	t.Run("Auto Suspension after Grace Period", func(t *testing.T) {
		// Move order to overdue status
		pastDue := now.AddDate(0, 0, -10) // 10 days ago
		ord.ExpiresAt = &pastDue
		ord.NextDueDate = &pastDue
		_ = orderRepo.Update(ctx, ord)

		res, err := cron.AutoSuspendOverdueOrdersBatch(ctx, 5)
		if err != nil {
			t.Fatalf("Suspension cron failed: %v", err)
		}

		if res.SuccessCount != 1 {
			t.Errorf("Expected 1 suspension, got %d", res.SuccessCount)
		}

		updatedOrd, _ := orderRepo.GetByID(ctx, 1)
		if updatedOrd.Status != domain.OrderStatusSuspended {
			t.Errorf("Expected status SUSPENDED, got %s", updatedOrd.Status)
		}
	})
}
