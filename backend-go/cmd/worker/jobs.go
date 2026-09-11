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

	// Fetch dynamic settings for billing cycles
	renewalDays := cronService.GetSystemService().GetIntSetting(jobCtx, "billing", "invoice_renewal_days", 14)
	suspensionGrace := cronService.GetSystemService().GetIntSetting(jobCtx, "billing", "suspension_grace_days", 7)
	ticketAutoClose := cronService.GetSystemService().GetIntSetting(jobCtx, "support", "auto_close_days", 7)

	// 1. Invoice Renewal Job
	tasks.RunInvoiceRenewalsTask(jobCtx, cronService, renewalDays)

	// 2. Automated Provisioning Job
	cronService.ProcessPendingProvisioningBatch(jobCtx)

	// 3. Overdue Order Auto-Suspension Job
	tasks.RunOverdueSuspensionsTask(jobCtx, cronService, suspensionGrace)

	// 3. Inactive Support Tickets Auto-Close
	tasks.RunTicketAutoCloseTask(jobCtx, cronService, ticketAutoClose)

	// 4. Housekeeping & Maintenance
	tasks.RunSystemMaintenanceTask(jobCtx)

	// 5. Automated Daily Backup (Runs at 02:00 UTC)
	if time.Now().Hour() == 2 {
		_, _ = cronService.PerformAutomatedBackup(jobCtx)
	}
}
