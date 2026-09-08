package integration_test

import (
	"context"
	"fmt"
	"sync"
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

func TestWorkerStress_BatchInvoiceGeneration(t *testing.T) {
	// We use Mock Repositories to test the Service Logic overhead without DB network latency bottlenecking
	// This helps us see how the Go code handles scaling.
	orderRepo := memory.NewMockOrderRepository()
	invRepo := memory.NewMockInvoiceRepository()
	clientRepo := memory.NewMockClientRepository()
	prodRepo := memory.NewMockProductRepository()

	taxCalc := billing.NewTaxCalculator(nil)
	invService := billing.NewInvoiceService(invRepo, clientRepo, taxCalc, plugins.NewHookManager())
	orderService := order.NewOrderService(orderRepo, prodRepo, nil, nil)
	cronService := scheduler.NewCronService(orderRepo, orderService, invService)

	ctx := context.Background()

	// 1. Seeding 5,000 Orders due for renewal
	count := 5000
	fmt.Printf("🏗️  Seeding %d orders into memory...\n", count)

	_ = clientRepo.Create(ctx, &domain.Client{ID: 1, Email: "admin@example.com", Currency: "USD"})
	dueAt := time.Now().UTC().AddDate(0, 0, 5)

	for i := 1; i <= count; i++ {
		_ = orderRepo.Create(ctx, &domain.Order{
			ClientID:    1,
			Title:       fmt.Sprintf("Enterprise VPS #%d", i),
			Period:      "1M",
			Price:       decimal.FromFloat(49.99),
			Currency:    "USD",
			Status:      domain.OrderStatusActive,
			NextDueDate: &dueAt,
		})
	}

	// 2. Run Stress Test (Sequential Baseline)
	fmt.Println("🚀 Starting Sequential Batch Processing...")
	start := time.Now()
	res, err := cronService.GenerateRenewalInvoicesBatch(ctx, 14)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Batch processing failed: %v", err)
	}

	fmt.Printf("📊 RESULTS (SEQUENTIAL):\n")
	fmt.Printf("   - Processed: %d orders\n", res.ProcessedCount)
	fmt.Printf("   - Success:   %d\n", res.SuccessCount)
	fmt.Printf("   - Errors:    %d\n", res.ErrorCount)
	fmt.Printf("   - Duration:  %v\n", duration)
	fmt.Printf("   - Speed:     %.2f orders/sec\n", float64(res.ProcessedCount)/duration.Seconds())

	if res.SuccessCount != count {
		t.Errorf("Expected %d successes, got %d", count, res.SuccessCount)
	}
}

func TestWorkerStress_ConcurrentSimulated(t *testing.T) {
	// This test simulates how a concurrent version would behave
	count := 5000
	fmt.Printf("\n🚀 Starting Concurrent Simulation (Goroutines)...\n")

	start := time.Now()
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 50) // Limit concurrency to 50 workers

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			semaphore <- struct{}{}
			// Simulate a standard DB + Logic task (5ms)
			time.Sleep(5 * time.Millisecond)
			<-semaphore
		}()
	}
	wg.Wait()
	duration := time.Since(start)

	fmt.Printf("📊 RESULTS (CONCURRENT SIM):\n")
	fmt.Printf("   - Processed: %d tasks\n", count)
	fmt.Printf("   - Duration:  %v\n", duration)
	fmt.Printf("   - Speed:     %.2f tasks/sec\n", float64(count)/duration.Seconds())
}
