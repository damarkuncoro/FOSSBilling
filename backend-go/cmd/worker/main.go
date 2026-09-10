package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/config"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/postgres"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/scheduler"
	billing "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/billing"
	order "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/order"
	system "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/system"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/logger"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/plugins"
)

func main() {
	cfg := config.Load(); logger.Init(cfg.AppEnv); logger.Info("Worker Start")
	ctx, cancel := context.WithCancel(context.Background()); defer cancel()

	pool, _ := postgres.NewPostgresPool(ctx, cfg.DatabaseURL)
	if pool != nil {
		defer pool.Close()
		_ = postgres.RunMigrations(ctx, pool, "migrations")
	}

	eb, hm := events.NewEventBus(), plugins.NewHookManager()
	or, pr, cr, ir, sr, mr, tr, sysR := postgres.NewOrderRepository(pool), postgres.NewProductRepository(pool), postgres.NewClientRepository(pool), postgres.NewInvoiceRepository(pool), postgres.NewSupportRepository(pool), postgres.NewMassMailRepository(pool), postgres.NewTaxRepository(pool), postgres.NewSystemRepository(pool)
	cor := postgres.NewCompanyRepository(pool)
	ordSvc := order.NewOrderService(or, pr, nil, nil, eb)
	is := billing.NewInvoiceService(ir, cr, cor, billing.NewTaxCalculator(tr), hm, eb)
	sysSvc := system.NewSystemService(sysR)
	cs := scheduler.NewCronService(or, ordSvc, is, sysSvc, sr, mr, cr, nil, cfg.DatabaseURL)

	tick := time.NewTicker(time.Minute); defer tick.Stop()
	go ExecuteCronBatch(cs)
	go func() {
		for {
			select {
			case <-tick.C: ExecuteCronBatch(cs)
			case <-ctx.Done(): return
			}
		}
	}()

	sig := make(chan os.Signal, 1); signal.Notify(sig, os.Interrupt, syscall.SIGTERM); <-sig
	logger.Warn("Worker Stop")
}
