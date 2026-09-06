package main

import (
	"context"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/cmd/worker/tasks"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/scheduler"
)

// ExecuteCronBatch executes all periodic billing and maintenance jobs via modular task handlers
func ExecuteCronBatch(cronService *scheduler.CronService) {
	jobCtx, jobCancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer jobCancel()

	// 1. Invoice Renewal Job (14 days advance notice)
	tasks.RunInvoiceRenewalsTask(jobCtx, cronService, 14)

	// 2. Overdue Order Auto-Suspension Job (7 days grace period)
	tasks.RunOverdueSuspensionsTask(jobCtx, cronService, 7)
}
