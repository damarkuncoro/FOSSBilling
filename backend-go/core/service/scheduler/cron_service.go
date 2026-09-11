package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/notification"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/billing"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/order"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/system"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/lock"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/mailer"
)

type CronService struct {
	orderRepo      domain.OrderRepository
	orderService   *order.OrderService
	invoiceService *billing.InvoiceService
	invoiceRepo    domain.InvoiceRepository
	emailService   *notification.EmailService
	systemService  *system.SystemService
	supportRepo    domain.SupportRepository
	massMailRepo   domain.MassMailRepository
	clientRepo     domain.ClientRepository
	locker         lock.Locker
	mailer         mailer.Mailer
	dbURL          string
	concurrency    int
}

func NewCronService(or domain.OrderRepository, os *order.OrderService, is *billing.InvoiceService, ir domain.InvoiceRepository, es *notification.EmailService, sys *system.SystemService, sr domain.SupportRepository, mr domain.MassMailRepository, cr domain.ClientRepository, l lock.Locker, m mailer.Mailer, dbURL string) *CronService {
	return &CronService{
		orderRepo:      or,
		orderService:   os,
		invoiceService: is,
		invoiceRepo:    ir,
		emailService:   es,
		systemService:  sys,
		supportRepo:    sr,
		massMailRepo:   mr,
		clientRepo:     cr,
		locker:         l,
		mailer:         m,
		dbURL:          dbURL,
		concurrency:    20,
	}
}

// runTaskInParallel: Generic Worker Pool (SRP & DRY)
func runTaskInParallel[T any](ctx context.Context, items []T, concurrency int, taskName string, fn func(item T) error) *domain.CronTaskResult {
	start := time.Now()
	res := &domain.CronTaskResult{TaskName: taskName, ProcessedCount: len(items)}
	if len(items) == 0 {
		res.Duration = time.Since(start)
		return res
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	jobs := make(chan T, len(items))
	workerCount := min(concurrency, len(items))

	for w := 0; w < workerCount; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range jobs {
				select {
				case <-ctx.Done():
					return
				default:
					err := fn(item)
					mu.Lock()
					if err != nil {
						res.ErrorCount++
						res.Errors = append(res.Errors, err.Error())
					} else {
						res.SuccessCount++
					}
					mu.Unlock()
				}
			}
		}()
	}
	for _, item := range items {
		jobs <- item
	}
	close(jobs)
	wg.Wait()
	res.Duration = time.Since(start)
	return res
}

func (s *CronService) GenerateRenewalInvoicesBatch(ctx context.Context, days int) (*domain.CronTaskResult, error) {
	orders, err := s.orderRepo.ListDueOrders(ctx, time.Now().UTC().AddDate(0, 0, days))
	if err != nil {
		return nil, err
	}

	return runTaskInParallel(ctx, orders, s.concurrency, "Renewals", func(ord *domain.Order) error {
		if ord.InvoiceID != nil {
			inv, err := s.invoiceService.GetInvoice(ctx, *ord.InvoiceID)
			if err == nil && inv.Status == domain.InvoiceStatusUnpaid {
				return nil
			}
		}
		newInv, err := s.invoiceService.CreateInvoice(ctx, billing.CreateInvoiceDTO{
			ClientID: ord.ClientID, Currency: ord.Currency, DueDays: 7,
			Items: []billing.CreateInvoiceItemDTO{{OrderID: &ord.ID, Title: "Renewal: " + ord.Title, Period: &ord.Period, Price: ord.Price, Quantity: 1, Taxable: true}},
		})
		if err != nil {
			return err
		}
		ord.InvoiceID = &newInv.ID
		return s.orderRepo.Update(ctx, ord)
	}), nil
}

func (s *CronService) AutoSuspendOverdueOrdersBatch(ctx context.Context, grace int) (*domain.CronTaskResult, error) {
	orders, err := s.orderRepo.ListOverdueSuspensions(ctx, grace)
	if err != nil {
		return nil, err
	}
	reason := fmt.Sprintf("Auto-suspended: overdue > %d days", grace)
	return runTaskInParallel(ctx, orders, s.concurrency, "Suspensions", func(o *domain.Order) error {
		_, err := s.orderService.Suspend(ctx, o.ID, reason)
		return err
	}), nil
}

func (s *CronService) ProcessPendingProvisioningBatch(ctx context.Context) (*domain.CronTaskResult, error) {
	orders, err := s.orderRepo.ListPendingProvisioning(ctx)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	return runTaskInParallel(ctx, orders, s.concurrency, "Provisioning", func(o *domain.Order) error {
		_, err := s.orderService.Activate(ctx, o.ID, now)
		return err
	}), nil
}

func (s *CronService) ProcessPendingMassMailBatch(ctx context.Context) (*domain.CronTaskResult, error) {
	campaigns, _, err := s.massMailRepo.List(ctx, 5, 0)
	if err != nil || len(campaigns) == 0 {
		return nil, err
	}
	clients, _, _ := s.clientRepo.List(ctx, 5000, 0)

	return runTaskInParallel(ctx, campaigns, 1, "MassMail", func(cp *domain.MassMailCampaign) error {
		if cp.Status != domain.CampaignStatusSending {
			return nil
		}
		sent := 0
		for _, c := range clients {
			if c.Email != "" && s.mailer.Send(ctx, mailer.Message{To: []string{c.Email}, Subject: cp.Subject, HTMLBody: cp.Content}) == nil {
				sent++
			}
		}
		cp.SentCount, cp.Status, cp.SentAt = sent, domain.CampaignStatusCompleted, pointer(time.Now().UTC())
		return s.massMailRepo.Update(ctx, cp)
	}), nil
}

func (s *CronService) AutoCloseInactiveTicketsBatch(ctx context.Context, days int) (*domain.CronTaskResult, error) {
	start := time.Now()
	if s.supportRepo == nil {
		return &domain.CronTaskResult{TaskName: "TicketClose", Duration: time.Since(start)}, nil
	}
	count, err := s.supportRepo.CloseInactiveTickets(ctx, time.Now().UTC().AddDate(0, 0, -days))
	if err != nil {
		return nil, err
	}
	return &domain.CronTaskResult{TaskName: "TicketClose", ProcessedCount: count, SuccessCount: count, Duration: time.Since(start)}, nil
}

func (s *CronService) SendInvoiceRemindersBatch(ctx context.Context) (*domain.CronTaskResult, error) {
	invoices, _, err := s.invoiceRepo.List(ctx, 10000, 0)
	if err != nil { return nil, err }

	return runTaskInParallel(ctx, invoices, s.concurrency, "Reminders", func(inv *domain.Invoice) error {
		if inv.Status != domain.InvoiceStatusUnpaid { return nil }
		c, _ := s.clientRepo.GetByID(ctx, inv.ClientID)
		if c != nil && s.emailService != nil {
			return s.emailService.SendInvoiceReminderEmail(ctx, c, inv)
		}
		return nil
	}), nil
}

func (s *CronService) PerformAutomatedBackup(ctx context.Context) (*domain.CronTaskResult, error) {
	start := time.Now()
	if s.systemService == nil || s.dbURL == "" {
		return &domain.CronTaskResult{TaskName: "DBBackup", Duration: time.Since(start)}, nil
	}
	_, err := s.systemService.CreateFullBackup(ctx, s.dbURL)
	if err != nil {
		return nil, err
	}
	return &domain.CronTaskResult{TaskName: "DBBackup", SuccessCount: 1, Duration: time.Since(start)}, nil
}

func (s *CronService) GetSystemService() *system.SystemService { return s.systemService }

func (s *CronService) RunLocked(ctx context.Context, key string, ttl time.Duration, fn func() error) error {
	if s.locker == nil {
		return fn()
	}

	ok, err := s.locker.Acquire(ctx, key, ttl)
	if err != nil {
		return err
	}
	if !ok {
		return nil // Silently skip if locked
	}

	defer s.locker.Release(ctx, key)
	return fn()
}

func pointer[T any](v T) *T { return &v }
func min(a, b int) int      { if a < b { return a }; return b }
