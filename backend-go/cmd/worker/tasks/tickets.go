package tasks

import (
	"context"
	"log"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/scheduler"
)

// RunTicketAutoCloseTask executes batch closing of inactive support tickets
func RunTicketAutoCloseTask(ctx context.Context, cronService *scheduler.CronService, inactiveDays int) {
	log.Printf("⏱️ [Task: Tickets] Running batch auto-close for inactive tickets (>%d days)...", inactiveDays)
	res, err := cronService.AutoCloseInactiveTicketsBatch(ctx, inactiveDays)
	if err != nil {
		log.Printf("❌ [Task: Tickets] Ticket auto-close failed: %v", err)
	} else if res.ProcessedCount > 0 {
		log.Printf("✅ [Task: Tickets] Inactive tickets closed: %d processed, %d succeeded, %d failed (took %v)",
			res.ProcessedCount, res.SuccessCount, res.ErrorCount, res.Duration)
	} else {
		log.Printf("✅ [Task: Tickets] No inactive tickets to close.")
	}
}
