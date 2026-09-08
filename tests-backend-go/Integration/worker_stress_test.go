package integration_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/scheduler"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/billing"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/order"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/plugins"
)

// SlowOrderRepository simulates a DB with 5ms latency
type SlowOrderRepository struct {
	*memory.MockOrderRepository
	latency time.Duration
}

func (r *SlowOrderRepository) Update(ctx context.Context, o *domain.Order) error {
	time.Sleep(r.latency)
	return r.MockOrderRepository.Update(ctx, o)
}

func TestWorkerStress_ComparePerformance(t *testing.T) {
	ctx := context.Background()
	count := 200 // Reduced for faster test execution in dev, but enough to see difference
	latency := 5 * time.Millisecond

	// Setup Services
	baseOrderRepo := memory.NewMockOrderRepository()
	slowRepo := &SlowOrderRepository{MockOrderRepository: baseOrderRepo, latency: latency}
	invRepo := memory.NewMockInvoiceRepository()
	clientRepo := memory.NewMockClientRepository()
	prodRepo := memory.NewMockProductRepository()

	taxCalc := billing.NewTaxCalculator(nil)
	invService := billing.NewInvoiceService(invRepo, clientRepo, taxCalc, plugins.NewHookManager())
	orderService := order.NewOrderService(slowRepo, prodRepo, nil, nil)
	cronService := scheduler.NewCronService(slowRepo, orderService, invService)

	// Seed
	_ = clientRepo.Create(ctx, &domain.Client{ID: 1, Email: "admin@example.com", Currency: "USD"})
	dueAt := time.Now().UTC().AddDate(0, 0, 5)
	for i := 1; i <= count; i++ {
		_ = baseOrderRepo.Create(ctx, &domain.Order{
			ClientID: 1, Title: "VPS", Period: "1M", Price: decimal.FromFloat(10),
			Currency: "USD", Status: domain.OrderStatusActive, NextDueDate: &dueAt,
		})
	}

	// 1. Sequential Mode (Concurrency = 1)
	fmt.Printf("⏱️  Testing %d orders with %v DB latency...\n", count, latency)
	cronService.SetConcurrency(1)
	start := time.Now()
	_, _ = cronService.GenerateRenewalInvoicesBatch(ctx, 14)
	durSeq := time.Since(start)
	fmt.Printf("   [SEQUENTIAL] Duration: %v (%.2f/sec)\n", durSeq, float64(count)/durSeq.Seconds())

	// Reset for next run (clear invoice_id from orders)
	orders, _ := baseOrderRepo.ListDueOrders(ctx, time.Now().AddDate(0, 0, 20))
	for _, o := range orders {
		o.InvoiceID = nil
		_ = baseOrderRepo.Update(ctx, o)
	}

	// 2. Concurrent Mode (Concurrency = 20)
	cronService.SetConcurrency(20)
	start = time.Now()
	resCon, _ := cronService.GenerateRenewalInvoicesBatch(ctx, 14)
	durCon := time.Since(start)
	fmt.Printf("   [CONCURRENT] Duration: %v (%.2f/sec) [20 workers]\n", durCon, float64(count)/durCon.Seconds())

	if durCon >= durSeq {
		t.Errorf("Concurrent mode should be faster than sequential. Seq: %v, Con: %v", durSeq, durCon)
	}

	improvement := (durSeq.Seconds() / durCon.Seconds())
	fmt.Printf("🚀 PERFORMANCE BOOST: %.1fx faster\n", improvement)

	if resCon.SuccessCount != count {
		t.Errorf("Concurrent run lost some data: %d/%d", resCon.SuccessCount, count)
	}
}

func TestWorkerStress_HighVolume(t *testing.T) {
	ctx := context.Background()
	count := 10000 // 10,000 orders

	orderRepo := memory.NewMockOrderRepository()
	invRepo := memory.NewMockInvoiceRepository()
	clientRepo := memory.NewMockClientRepository()
	prodRepo := memory.NewMockProductRepository()

	taxCalc := billing.NewTaxCalculator(nil)
	invService := billing.NewInvoiceService(invRepo, clientRepo, taxCalc, plugins.NewHookManager())
	orderService := order.NewOrderService(orderRepo, prodRepo, nil, nil)
	cronService := scheduler.NewCronService(orderRepo, orderService, invService)
	cronService.SetConcurrency(50) // Use 50 workers for 10k items

	// Seed
	_ = clientRepo.Create(ctx, &domain.Client{ID: 1, Email: "admin@example.com", Currency: "USD"})
	dueAt := time.Now().UTC().AddDate(0, 0, 5)
	for i := 1; i <= count; i++ {
		_ = orderRepo.Create(ctx, &domain.Order{
			ClientID: 1, Title: "VPS", Period: "1M", Price: decimal.FromFloat(10),
			Currency: "USD", Status: domain.OrderStatusActive, NextDueDate: &dueAt,
		})
	}

	fmt.Printf("\n🔥 STRESS TEST: Processing %d renewal invoices with 50 concurrent workers...\n", count)
	start := time.Now()
	res, err := cronService.GenerateRenewalInvoicesBatch(ctx, 14)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Stress test failed: %v", err)
	}

	fmt.Printf("📊 RESULTS:\n")
	fmt.Printf("   - Duration:  %v\n", duration)
	fmt.Printf("   - Speed:     %.2f invoices/sec\n", float64(count)/duration.Seconds())

	if res.SuccessCount != count {
		t.Errorf("Incomplete run: %d/%d", res.SuccessCount, count)
	}
}
