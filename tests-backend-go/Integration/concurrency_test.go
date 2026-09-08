package integration_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
)

// ThreadSafeMockClientRepo is a concurrent-safe mock repository for testing wallet race conditions
type ThreadSafeMockClientRepo struct {
	mu      sync.Mutex
	balance decimal.Money
	history []*domain.ClientBalance
}

func NewThreadSafeMockClientRepo(initial decimal.Money) *ThreadSafeMockClientRepo {
	return &ThreadSafeMockClientRepo{
		balance: initial,
		history: make([]*domain.ClientBalance, 0),
	}
}

func (m *ThreadSafeMockClientRepo) GetByID(ctx context.Context, id int64) (*domain.Client, error) {
	return &domain.Client{
		ID:        id,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Currency:  "USD",
		Status:    domain.ClientStatusActive,
	}, nil
}

func (m *ThreadSafeMockClientRepo) GetBalance(ctx context.Context, clientID int64) (decimal.Money, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.balance, nil
}

func (m *ThreadSafeMockClientRepo) AddBalanceTransaction(ctx context.Context, tx *domain.ClientBalance) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if tx.Type == domain.BalanceTypeDebit {
		if m.balance < tx.Amount {
			return appErrors.ErrInsufficientFunds
		}
		m.balance -= tx.Amount
	} else {
		m.balance += tx.Amount
	}

	m.history = append(m.history, tx)
	return nil
}

// 1. TDD / Integration: High Concurrency Wallet Debit & Double-Spending Prevention Test
func TestConcurrency_WalletDebitDoubleSpendPrevention(t *testing.T) {
	initialBal := decimal.FromFloat(100.00)

	repo := NewThreadSafeMockClientRepo(initialBal)
	ctx := context.Background()

	debitAmt := decimal.FromFloat(10.00)
	numRoutines := 50

	var successCount int32
	var failedCount int32

	var wg sync.WaitGroup
	wg.Add(numRoutines)

	startBarrier := make(chan struct{})

	for i := 0; i < numRoutines; i++ {
		go func(idx int) {
			defer wg.Done()
			<-startBarrier // Synchronize simultaneous launch of all goroutines

			tx := &domain.ClientBalance{
				ClientID:    1,
				Amount:      debitAmt,
				Type:        domain.BalanceTypeDebit,
				Description: "Concurrent debit test",
				CreatedAt:   time.Now(),
			}

			err := repo.AddBalanceTransaction(ctx, tx)
			if err == nil {
				atomic.AddInt32(&successCount, 1)
			} else {
				atomic.AddInt32(&failedCount, 1)
			}
		}(i)
	}

	// Trigger all 50 goroutines simultaneously
	close(startBarrier)
	wg.Wait()

	finalBal, err := repo.GetBalance(ctx, 1)
	if err != nil {
		t.Fatalf("failed to get final balance: %v", err)
	}

	zeroMoney := decimal.FromFloat(0.00)

	// Assertions:
	// Exactly 10 debits should succeed ($100.00 / $10.00 = 10)
	if successCount != 10 {
		t.Errorf("expected exactly 10 successful transactions, got %d", successCount)
	}

	// 40 debits must fail with ErrInsufficientFunds
	if failedCount != 40 {
		t.Errorf("expected 40 failed transactions, got %d", failedCount)
	}

	// Final balance must be exactly 0.00 (no negative balance or double spending)
	if finalBal != zeroMoney {
		t.Errorf("expected final balance $0.00, got %s", finalBal.String())
	}
}

// 2. TDD / Integration: Idempotent Concurrent Invoice Settlement
func TestConcurrency_InvoiceSettlementIdempotency(t *testing.T) {
	initialBal := decimal.FromFloat(50.00)
	clientRepo := NewThreadSafeMockClientRepo(initialBal)

	invTotal := decimal.FromFloat(50.00)
	inv := &domain.Invoice{
		ID:       42,
		ClientID: 1,
		Status:   domain.InvoiceStatusUnpaid,
		Total:    invTotal,
	}

	var mu sync.Mutex
	var payCount int32

	errAlreadyPaid := errors.New("invoice already paid")

	// Custom invoice settlement handler for idempotency test
	payInvoiceOnce := func(ctx context.Context, invoiceID int64) error {
		mu.Lock()
		defer mu.Unlock()

		if inv.Status == domain.InvoiceStatusPaid {
			return errAlreadyPaid
		}

		bal, err := clientRepo.GetBalance(ctx, inv.ClientID)
		if err != nil {
			return err
		}

		if bal < inv.Total {
			return appErrors.ErrInsufficientFunds
		}

		deduct := &domain.ClientBalance{
			ClientID: inv.ClientID,
			Amount:   inv.Total,
			Type:     domain.BalanceTypeDebit,
		}

		if err := clientRepo.AddBalanceTransaction(ctx, deduct); err != nil {
			return err
		}

		inv.Status = domain.InvoiceStatusPaid
		atomic.AddInt32(&payCount, 1)
		return nil
	}

	numWorkers := 20
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	startBarrier := make(chan struct{})

	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			<-startBarrier
			_ = payInvoiceOnce(context.Background(), 42)
		}()
	}

	close(startBarrier)
	wg.Wait()

	// Assertions:
	// The invoice should be paid exactly once
	if atomic.LoadInt32(&payCount) != 1 {
		t.Errorf("expected invoice to be settled exactly 1 time, got %d", payCount)
	}

	if inv.Status != domain.InvoiceStatusPaid {
		t.Errorf("expected invoice status to be PAID, got %s", inv.Status)
	}

	finalBal, _ := clientRepo.GetBalance(context.Background(), 1)
	zeroMoney := decimal.FromFloat(0.00)
	if finalBal != zeroMoney {
		t.Errorf("expected client balance to be $0.00 after 1 deduction, got %s", finalBal.String())
	}
}
