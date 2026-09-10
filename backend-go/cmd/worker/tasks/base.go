package tasks

import (
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/logger"
)

func run(n string, fn func() (*domain.CronTaskResult, error)) {
	logger.Info("Running task", "name", n)
	res, err := fn()
	if err != nil {
		logger.Error("Task failed", err, "name", n)
		return
	}
	logger.Info("Task completed", "name", n, "processed", res.ProcessedCount, "success", res.SuccessCount, "errors", res.ErrorCount, "took", res.Duration)
}
