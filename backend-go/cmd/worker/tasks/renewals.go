package tasks

import (
	"context"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/scheduler"
)

func RunInvoiceRenewalsTask(ctx context.Context, s *scheduler.CronService, days int) {
	run("Renewals", func() (*domain.CronTaskResult, error) { return s.GenerateRenewalInvoicesBatch(ctx, days) })
}
