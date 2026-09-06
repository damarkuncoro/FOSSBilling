package tasks

import (
	"context"
	"log"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/scheduler"
)

// RunOverdueSuspensionsTask executes batch overdue order suspensions
func RunOverdueSuspensionsTask(ctx context.Context, cronService *scheduler.CronService, gracePeriodDays int) {
	log.Printf("⏱️ [Task: Suspensions] Running batch overdue order suspension (%d days grace period)...", gracePeriodDays)
	suspRes, err := cronService.AutoSuspendOverdueOrdersBatch(ctx, gracePeriodDays)
	if err != nil {
		log.Printf("❌ [Task: Suspensions] Auto-suspension failed: %v", err)
	} else {
		log.Printf("✅ [Task: Suspensions] Orders suspended: %d processed, %d succeeded, %d failed (took %v)",
			suspRes.ProcessedCount, suspRes.SuccessCount, suspRes.ErrorCount, suspRes.Duration)
	}
}
