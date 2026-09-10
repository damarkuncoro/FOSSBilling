package tasks

import (
	"context"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/scheduler"
)

func RunTicketAutoCloseTask(ctx context.Context, s *scheduler.CronService, days int) {
	run("Tickets", func() (*domain.CronTaskResult, error) { return s.AutoCloseInactiveTicketsBatch(ctx, days) })
}
