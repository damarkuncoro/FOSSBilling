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
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/scheduler"
	billingUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/billing"
	orderUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/order"
)

func main() {
	cfg := config.Load()
	log.Printf("⚙️ Starting FOSSBilling Background Worker & Scheduler (%s)...", cfg.AppEnv)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Database Connection Pool
	pgPool, err := postgres.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Printf("⚠️ Database connection failed (%v). Worker running with pool handle.", err)
	} else {
		defer pgPool.Close()
		log.Println("✅ Connected to PostgreSQL database pool.")
	}

	// 2. Repositories & Domain Services
	orderRepo := postgres.NewOrderRepository(pgPool)
	clientRepo := postgres.NewClientRepository(pgPool)
	invoiceRepo := postgres.NewInvoiceRepository(pgPool)

	taxCalculator := billingUsecase.NewTaxCalculator(nil)
	orderService := orderUsecase.NewOrderService(orderRepo)
	invoiceService := billingUsecase.NewInvoiceService(invoiceRepo, clientRepo, taxCalculator)

	cronService := scheduler.NewCronService(orderRepo, orderService, invoiceService)

	// 3. Periodic Scheduler Loop
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	// Execute initial job batch on startup
	go ExecuteCronBatch(cronService)

	go func() {
		for {
			select {
			case <-ticker.C:
				ExecuteCronBatch(cronService)
			case <-ctx.Done():
				return
			}
		}
	}()

	// 4. Graceful Shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	<-stopChan
	log.Println("🛑 Shutting down background worker gracefully...")
}
