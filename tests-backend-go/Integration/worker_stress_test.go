package integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/scheduler"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/billing"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/order"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/plugins"
)

type SlowOrderRepo struct { *memory.MockOrderRepository; lat time.Duration }
func (r *SlowOrderRepo) Update(ctx context.Context, o *domain.Order) error { time.Sleep(r.lat); return r.MockOrderRepository.Update(ctx, o) }

func TestWorkerStress_ComparePerformance(t *testing.T) {
	ctx := context.Background(); count, lat := 50, 2*time.Millisecond
	bor := memory.NewMockOrderRepository(); slow := &SlowOrderRepo{MockOrderRepository: bor, lat: lat}
	ir, cr, pr, cor := memory.NewMockInvoiceRepository(), memory.NewMockClientRepository(), memory.NewMockProductRepository(), memory.NewMockCompanyRepository()
	is := billing.NewInvoiceService(ir, cr, cor, billing.NewTaxCalculator(nil), plugins.NewHookManager())
	os := order.NewOrderService(slow, pr, nil, nil); cs := scheduler.NewCronService(slow, os, is, nil, memory.NewMockSupportRepository(), memory.NewMockMassMailRepository(), cr, nil, "")

	_ = cr.Create(ctx, &domain.Client{ID: 1, Email: "a@e.com", Currency: "USD"})
	due := time.Now().UTC().AddDate(0, 0, 5)
	for i := 1; i <= count; i++ { _ = bor.Create(ctx, &domain.Order{ClientID: 1, Title: "V", Period: "1M", Price: 1000, Currency: "USD", Status: domain.OrderStatusActive, NextDueDate: &due}) }

	s1 := time.Now(); _, _ = cs.GenerateRenewalInvoicesBatch(ctx, 14); d1 := time.Since(s1)
	if d1 == 0 { t.Error("Should take time") }
}
