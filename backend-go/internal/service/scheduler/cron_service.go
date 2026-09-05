package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/usecase/billing"
	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/usecase/order"
)

type CronService struct {
	orderRepo      domain.OrderRepository
	orderService   *order.OrderService
	invoiceService *billing.InvoiceService
}

func NewCronService(
	orderRepo domain.OrderRepository,
	orderService *order.OrderService,
	invoiceService *billing.InvoiceService,
) *CronService {
	return &CronService{
		orderRepo:      orderRepo,
		orderService:   orderService,
		invoiceService: invoiceService,
	}
}

// GenerateRenewalInvoicesBatch finds active orders due within issueDaysBefore and generates invoices
func (s *CronService) GenerateRenewalInvoicesBatch(ctx context.Context, issueDaysBefore int) (*domain.CronTaskResult, error) {
	start := time.Now()
	now := time.Now().UTC()
	cutoffDate := now.AddDate(0, 0, issueDaysBefore)

	dueOrders, err := s.orderRepo.ListDueOrders(ctx, cutoffDate)
	if err != nil {
		return nil, err
	}

	result := &domain.CronTaskResult{
		TaskName:       "BatchInvoiceGenerator",
		ProcessedCount: len(dueOrders),
	}

	for _, ord := range dueOrders {
		item := billing.CreateInvoiceItemDTO{
			OrderID:  &ord.ID,
			Title:    fmt.Sprintf("Renewal: %s (%s)", ord.Title, ord.Period),
			Period:   &ord.Period,
			Price:    ord.Price,
			Quantity: 1,
			Taxable:  true,
		}

		inv, err := s.invoiceService.CreateInvoice(ctx, billing.CreateInvoiceDTO{
			ClientID: ord.ClientID,
			Currency: ord.Currency,
			DueDays:  7,
			Items:    []billing.CreateInvoiceItemDTO{item},
		})
		if err != nil {
			result.ErrorCount++
			result.Errors = append(result.Errors, fmt.Sprintf("Order #%d invoice error: %v", ord.ID, err))
			continue
		}

		ord.InvoiceID = &inv.ID
		_ = s.orderRepo.Update(ctx, ord)
		result.SuccessCount++
	}

	result.Duration = time.Since(start)
	return result, nil
}

// AutoSuspendOverdueOrdersBatch finds overdue orders exceeding grace period and suspends them
func (s *CronService) AutoSuspendOverdueOrdersBatch(ctx context.Context, gracePeriodDays int) (*domain.CronTaskResult, error) {
	start := time.Now()
	overdueOrders, err := s.orderRepo.ListOverdueSuspensions(ctx, gracePeriodDays)
	if err != nil {
		return nil, err
	}

	result := &domain.CronTaskResult{
		TaskName:       "BatchAutoSuspension",
		ProcessedCount: len(overdueOrders),
	}

	reason := fmt.Sprintf("Auto-suspended by system: payment overdue past %d days grace period", gracePeriodDays)

	for _, ord := range overdueOrders {
		_, err := s.orderService.Suspend(ctx, ord.ID, reason)
		if err != nil {
			result.ErrorCount++
			result.Errors = append(result.Errors, fmt.Sprintf("Order #%d suspend error: %v", ord.ID, err))
			continue
		}
		result.SuccessCount++
	}

	result.Duration = time.Since(start)
	return result, nil
}
