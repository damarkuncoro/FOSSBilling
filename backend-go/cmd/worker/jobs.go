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
	_ = cronService.RunLocked(jobCtx, "task:renewals", time.Minute, func() error {
		tasks.RunInvoiceRenewalsTask(jobCtx, cronService, renewalDays)
		return nil
	})

	// 2. Automated Provisioning Job
	_ = cronService.RunLocked(jobCtx, "task:provisioning", 30*time.Second, func() error {
		_, err := cronService.ProcessPendingProvisioningBatch(jobCtx)
		return err
	})

	// 3. Overdue Order Auto-Suspension Job
	_ = cronService.RunLocked(jobCtx, "task:suspensions", time.Minute, func() error {
		tasks.RunOverdueSuspensionsTask(jobCtx, cronService, suspensionGrace)
		return nil
	})

	// 3. Inactive Support Tickets Auto-Close
	_ = cronService.RunLocked(jobCtx, "task:tickets_close", time.Hour, func() error {
		tasks.RunTicketAutoCloseTask(jobCtx, cronService, ticketAutoClose)
		return nil
	})

	// 4. Housekeeping & Maintenance
	tasks.RunSystemMaintenanceTask(jobCtx)

	// 5. Automated Daily Backup (Runs at 02:00 UTC)
	if time.Now().Hour() == 2 {
		_ = cronService.RunLocked(jobCtx, "task:daily_backup", 2*time.Hour, func() error {
			_, err := cronService.PerformAutomatedBackup(jobCtx)
			return err
		})
	}
}
