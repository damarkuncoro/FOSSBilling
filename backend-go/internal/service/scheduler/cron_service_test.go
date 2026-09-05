package scheduler

import (
	"context"
	"testing"
	"time"

	"github.com/fossbilling/backend-go/internal/domain"
	"github.com/fossbilling/backend-go/internal/repository/memory"
	"github.com/fossbilling/backend-go/internal/usecase/billing"
	"github.com/fossbilling/backend-go/internal/usecase/order"
	"github.com/fossbilling/backend-go/pkg/decimal"
)

func setupCronService() (*CronService, *memory.MockOrderRepository, *memory.MockInvoiceRepository, *memory.MockClientRepository) {
	orderRepo := memory.NewMockOrderRepository()
	invRepo := memory.NewMockInvoiceRepository()
	clientRepo := memory.NewMockClientRepository()

	taxCalc := billing.NewTaxCalculator(nil)
	invService := billing.NewInvoiceService(invRepo, clientRepo, taxCalc)
	orderService := order.NewOrderService(orderRepo)

	cronService := NewCronService(orderRepo, orderService, invService)
	return cronService, orderRepo, invRepo, clientRepo
}

func TestCronService_GenerateRenewalInvoicesBatch(t *testing.T) {
	ctx := context.Background()
	service, orderRepo, invRepo, clientRepo := setupCronService()

	client := &domain.Client{Email: "cron.client@example.com", Country: "US", Currency: "USD"}
	_ = clientRepo.Create(ctx, client)

	// Order 1: Due in 5 days (Should be picked up if looking 14 days ahead)
	dueIn5Days := time.Now().UTC().AddDate(0, 0, 5)
	ord1 := &domain.Order{
		ClientID:    client.ID,
		ProductID:   1,
		Title:       "Server VPS A",
		Period:      "1M",
		Price:       decimal.FromFloat(30.00),
		Currency:    "USD",
		Status:      domain.OrderStatusActive,
		NextDueDate: &dueIn5Days,
	}
	_ = orderRepo.Create(ctx, ord1)

	// Order 2: Due in 30 days (Should NOT be picked up)
	dueIn30Days := time.Now().UTC().AddDate(0, 0, 30)
	ord2 := &domain.Order{
		ClientID:    client.ID,
		ProductID:   2,
		Title:       "Domain B",
		Period:      "1Y",
		Price:       decimal.FromFloat(15.00),
		Currency:    "USD",
		Status:      domain.OrderStatusActive,
		NextDueDate: &dueIn30Days,
	}
	_ = orderRepo.Create(ctx, ord2)

	// Run Batch: Issue invoices 14 days ahead
	res, err := service.GenerateRenewalInvoicesBatch(ctx, 14)
	if err != nil {
		t.Fatalf("GenerateRenewalInvoicesBatch failed: %v", err)
	}

	if res.ProcessedCount != 1 {
		t.Errorf("ProcessedCount = %d; want 1", res.ProcessedCount)
	}
	if res.SuccessCount != 1 {
		t.Errorf("SuccessCount = %d; want 1", res.SuccessCount)
	}

	// Verify an invoice was created for client
	invoices, total, _ := invRepo.ListByClientID(ctx, client.ID, 10, 0)
	if total != 1 {
		t.Fatalf("Expected 1 invoice created, got: %d", total)
	}
	if invoices[0].Subtotal.String() != "30.00" {
		t.Errorf("Invoice subtotal = %s; want 30.00", invoices[0].Subtotal.String())
	}
}

func TestCronService_AutoSuspendOverdueOrdersBatch(t *testing.T) {
	ctx := context.Background()
	service, orderRepo, _, _ := setupCronService()

	// Order 1: Expired 10 days ago (Exceeds 7 days grace period -> Should be suspended)
	expired10DaysAgo := time.Now().UTC().AddDate(0, 0, -10)
	ord1 := &domain.Order{
		ClientID:  1,
		ProductID: 1,
		Title:     "Overdue Service",
		Period:    "1M",
		Status:    domain.OrderStatusActive,
		ExpiresAt: &expired10DaysAgo,
	}
	_ = orderRepo.Create(ctx, ord1)

	// Order 2: Expired 2 days ago (Within 7 days grace period -> Should NOT be suspended)
	expired2DaysAgo := time.Now().UTC().AddDate(0, 0, -2)
	ord2 := &domain.Order{
		ClientID:  1,
		ProductID: 2,
		Title:     "Recently Expired Service",
		Period:    "1M",
		Status:    domain.OrderStatusActive,
		ExpiresAt: &expired2DaysAgo,
	}
	_ = orderRepo.Create(ctx, ord2)

	// Run Batch Auto Suspension with 7 days grace period
	res, err := service.AutoSuspendOverdueOrdersBatch(ctx, 7)
	if err != nil {
		t.Fatalf("AutoSuspendOverdueOrdersBatch failed: %v", err)
	}

	if res.ProcessedCount != 1 {
		t.Errorf("ProcessedCount = %d; want 1", res.ProcessedCount)
	}

	// Verify Ord1 is suspended
	updatedOrd1, _ := orderRepo.GetByID(ctx, ord1.ID)
	if updatedOrd1.Status != domain.OrderStatusSuspended {
		t.Errorf("Ord1 status = %s; want suspended", updatedOrd1.Status)
	}

	// Verify Ord2 is still active
	updatedOrd2, _ := orderRepo.GetByID(ctx, ord2.ID)
	if updatedOrd2.Status != domain.OrderStatusActive {
		t.Errorf("Ord2 status = %s; want active", updatedOrd2.Status)
	}
}
