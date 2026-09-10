package tasks

import (
	"context"
	"time"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/logger"
)

func RunSystemMaintenanceTask(ctx context.Context) {
	start := time.Now()
	logger.Info("Running maintenance")
	logger.Info("Maintenance completed", "took", time.Since(start))
}
