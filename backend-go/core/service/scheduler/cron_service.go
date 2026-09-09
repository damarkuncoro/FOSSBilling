package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/billing"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/order"
)

type CronService struct {
	orderRepo      domain.OrderRepository
	orderService   *order.OrderService
	invoiceService *billing.InvoiceService
	supportRepo    domain.SupportRepository
	concurrency    int
}

func NewCronService(
	orderRepo domain.OrderRepository,
	orderService *order.OrderService,
	invoiceService *billing.InvoiceService,
	supportRepo ...domain.SupportRepository,
) *CronService {
	var sRepo domain.SupportRepository
	if len(supportRepo) > 0 {
		sRepo = supportRepo[0]
	}
	return &CronService{
		orderRepo:      orderRepo,
		orderService:   orderService,
		invoiceService: invoiceService,
		supportRepo:    sRepo,
		concurrency:    20, // Default worker pool size
	}
}

func (s *CronService) SetConcurrency(n int) {
	if n > 0 {
		s.concurrency = n
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

	if len(dueOrders) == 0 {
		result.Duration = time.Since(start)
		return result, nil
	}

	// Use Worker Pool for high-performance concurrent processing
	var wg sync.WaitGroup
	var mu sync.Mutex
	jobs := make(chan *domain.Order, len(dueOrders))

	// Start workers
	workerCount := s.concurrency
	if workerCount > len(dueOrders) {
		workerCount = len(dueOrders)
	}

	for w := 1; w <= workerCount; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ord := range jobs {
				// BUG-33 FIX: Skip if order already has an active renewal invoice
				if ord.InvoiceID != nil {
					existingInv, err := s.invoiceService.GetInvoice(ctx, *ord.InvoiceID)
					if err == nil && (existingInv.Status == domain.InvoiceStatusUnpaid) {
						// Already has an unpaid invoice, skip to prevent duplicates
						mu.Lock()
						result.ProcessedCount-- // Adjust count since we skipped
						mu.Unlock()
						continue
					}
				}

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

				if err == nil {
					ord.InvoiceID = &inv.ID
					err = s.orderRepo.Update(ctx, ord)
				}

				mu.Lock()
				if err != nil {
					result.ErrorCount++
					result.Errors = append(result.Errors, fmt.Sprintf("Order #%d processing error: %v", ord.ID, err))
				} else {
					result.SuccessCount++
				}
				mu.Unlock()
			}
		}()
	}

	// Send jobs
	for _, ord := range dueOrders {
		jobs <- ord
	}
	close(jobs)
	wg.Wait()

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

	if len(overdueOrders) == 0 {
		result.Duration = time.Since(start)
		return result, nil
	}

	reason := fmt.Sprintf("Auto-suspended by system: payment overdue past %d days grace period", gracePeriodDays)

	// Concurrent Processing
	var wg sync.WaitGroup
	var mu sync.Mutex
	jobs := make(chan *domain.Order, len(overdueOrders))

	workerCount := s.concurrency
	if workerCount > len(overdueOrders) {
		workerCount = len(overdueOrders)
	}

	for w := 1; w <= workerCount; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ord := range jobs {
				_, err := s.orderService.Suspend(ctx, ord.ID, reason)
				mu.Lock()
				if err != nil {
					result.ErrorCount++
					result.Errors = append(result.Errors, fmt.Sprintf("Order #%d suspend error: %v", ord.ID, err))
				} else {
					result.SuccessCount++
				}
				mu.Unlock()
			}
		}()
	}

	for _, ord := range overdueOrders {
		jobs <- ord
	}
	close(jobs)
	wg.Wait()

	result.Duration = time.Since(start)
	return result, nil
}

// AutoCloseInactiveTicketsBatch closes resolved or abandoned tickets exceeding inactiveDays
func (s *CronService) AutoCloseInactiveTicketsBatch(ctx context.Context, inactiveDays int) (*domain.CronTaskResult, error) {
	start := time.Now()
	result := &domain.CronTaskResult{
		TaskName: "BatchAutoCloseTickets",
	}

	if s.supportRepo == nil {
		result.Duration = time.Since(start)
		return result, nil
	}

	cutoff := time.Now().UTC().AddDate(0, 0, -inactiveDays)
	closedCount, err := s.supportRepo.CloseInactiveTickets(ctx, cutoff)
	if err != nil {
		result.ErrorCount++
		result.Errors = append(result.Errors, fmt.Sprintf("Close inactive tickets error: %v", err))
	} else {
		result.ProcessedCount = closedCount
		result.SuccessCount = closedCount
	}

	result.Duration = time.Since(start)
	return result, nil
}
