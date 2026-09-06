package tasks

import (
	"context"
	"log"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/scheduler"
)

// RunInvoiceRenewalsTask executes batch renewal invoices generation
func RunInvoiceRenewalsTask(ctx context.Context, cronService *scheduler.CronService, advanceDays int) {
	log.Printf("⏱️ [Task: Renewals] Running batch renewal invoices generation (%d days advance)...", advanceDays)
	invRes, err := cronService.GenerateRenewalInvoicesBatch(ctx, advanceDays)
	if err != nil {
		log.Printf("❌ [Task: Renewals] Invoice generation failed: %v", err)
	} else {
		log.Printf("✅ [Task: Renewals] Invoices generated: %d processed, %d succeeded, %d failed (took %v)",
			invRes.ProcessedCount, invRes.SuccessCount, invRes.ErrorCount, invRes.Duration)
	}
}
