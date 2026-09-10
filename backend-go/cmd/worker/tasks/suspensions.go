package tasks

import (
	"context"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/scheduler"
)

func RunOverdueSuspensionsTask(ctx context.Context, s *scheduler.CronService, grace int) {
	run("Suspensions", func() (*domain.CronTaskResult, error) { return s.AutoSuspendOverdueOrdersBatch(ctx, grace) })
}
