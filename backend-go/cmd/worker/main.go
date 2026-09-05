package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/config"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/postgres"
	billingUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/billing"
	orderUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/order"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/scheduler"
)

func main() {
	cfg := config.Load()
	log.Printf("⚙️ Starting FOSSBilling Background Worker & Scheduler (%s)...", cfg.AppEnv)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Initialize PostgreSQL Connection Pool
	pgPool, err := postgres.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Printf("⚠️ Database connection failed (%v). Worker running with pool handle.", err)
	} else {
		defer pgPool.Close()
		log.Println("✅ Connected to PostgreSQL database pool.")
	}

	// 2. Repositories & Services
	orderRepo := postgres.NewOrderRepository(pgPool)
	clientRepo := postgres.NewClientRepository(pgPool)
	invoiceRepo := postgres.NewInvoiceRepository(pgPool)


	taxCalculator := billingUsecase.NewTaxCalculator(nil)
	orderService := orderUsecase.NewOrderService(orderRepo)
	invoiceService := billingUsecase.NewInvoiceService(invoiceRepo, clientRepo, taxCalculator)

	cronService := scheduler.NewCronService(orderRepo, orderService, invoiceService)


	// 3. Setup Periodic Batch Tickers
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	runScheduledTasks := func() {
		jobCtx, jobCancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer jobCancel()

		log.Println("⏱️ [Cron] Running batch renewal invoices generation (14 days advance)...")
		invRes, err := cronService.GenerateRenewalInvoicesBatch(jobCtx, 14)
		if err != nil {
			log.Printf("❌ [Cron] Invoice generation failed: %v", err)
		} else {
			log.Printf("✅ [Cron] Invoices generated: %d processed, %d succeeded, %d failed (took %v)",
				invRes.ProcessedCount, invRes.SuccessCount, invRes.ErrorCount, invRes.Duration)
		}

		log.Println("⏱️ [Cron] Running batch overdue order suspension (7 days grace period)...")
		suspRes, err := cronService.AutoSuspendOverdueOrdersBatch(jobCtx, 7)
		if err != nil {
			log.Printf("❌ [Cron] Auto-suspension failed: %v", err)
		} else {
			log.Printf("✅ [Cron] Orders suspended: %d processed, %d succeeded, %d failed (took %v)",
				suspRes.ProcessedCount, suspRes.SuccessCount, suspRes.ErrorCount, suspRes.Duration)
		}
	}

	// Run initial pass on startup
	go runScheduledTasks()

	go func() {
		for {
			select {
			case <-ticker.C:
				runScheduledTasks()
			case <-ctx.Done():
				return
			}
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	<-stopChan
	log.Println("🛑 Shutting down background worker gracefully...")
}
