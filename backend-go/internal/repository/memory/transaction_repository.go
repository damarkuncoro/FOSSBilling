package memory

import (
	"context"
	"sync"
	"time"

	"github.com/fossbilling/backend-go/internal/domain"
	appErrors "github.com/fossbilling/backend-go/pkg/errors"
)

type MockTransactionRepository struct {
	mu           sync.RWMutex
	transactions map[int64]*domain.Transaction
	nextID       int64
}

func NewMockTransactionRepository() *MockTransactionRepository {
	return &MockTransactionRepository{
		transactions: make(map[int64]*domain.Transaction),
		nextID:       1,
	}
}

func (r *MockTransactionRepository) GetByID(ctx context.Context, id int64) (*domain.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	txn, ok := r.transactions[id]
	if !ok {
		return nil, appErrors.ErrNotFound
	}
	cp := *txn
	return &cp, nil
}

func (r *MockTransactionRepository) GetByTxnID(ctx context.Context, gatewayID, txnID string) (*domain.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, txn := range r.transactions {
		if txn.GatewayID == gatewayID && txn.TxnID == txnID {
			cp := *txn
			return &cp, nil
		}
	}
	return nil, appErrors.ErrNotFound
}

func (r *MockTransactionRepository) Create(ctx context.Context, txn *domain.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	txn.ID = r.nextID
	r.nextID++
	txn.CreatedAt = time.Now().UTC()

	cp := *txn
	r.transactions[txn.ID] = &cp
	return nil
}

func (r *MockTransactionRepository) UpdateStatus(ctx context.Context, id int64, status domain.TransactionStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	txn, ok := r.transactions[id]
	if !ok {
		return appErrors.ErrNotFound
	}
	txn.Status = status
	return nil
}
