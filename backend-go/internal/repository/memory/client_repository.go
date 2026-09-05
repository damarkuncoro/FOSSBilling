package memory

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/fossbilling/backend-go/internal/domain"
	appErrors "github.com/fossbilling/backend-go/pkg/errors"
	"github.com/fossbilling/backend-go/pkg/decimal"
)

type MockClientRepository struct {
	mu       sync.RWMutex
	clients  map[int64]*domain.Client
	balances map[int64][]*domain.ClientBalance
	nextID   int64
}

func NewMockClientRepository() *MockClientRepository {
	return &MockClientRepository{
		clients:  make(map[int64]*domain.Client),
		balances: make(map[int64][]*domain.ClientBalance),
		nextID:   1,
	}
}

func (r *MockClientRepository) GetByID(ctx context.Context, id int64) (*domain.Client, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	client, ok := r.clients[id]
	if !ok {
		return nil, appErrors.ErrNotFound
	}
	cp := *client
	return &cp, nil
}

func (r *MockClientRepository) GetByEmail(ctx context.Context, email string) (*domain.Client, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, c := range r.clients {
		if strings.EqualFold(c.Email, email) {
			cp := *c
			return &cp, nil
		}
	}
	return nil, appErrors.ErrNotFound
}

func (r *MockClientRepository) Create(ctx context.Context, c *domain.Client) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, existing := range r.clients {
		if strings.EqualFold(existing.Email, c.Email) {
			return appErrors.ErrDuplicateEntry
		}
	}

	c.ID = r.nextID
	r.nextID++
	now := time.Now().UTC()
	c.CreatedAt = now
	c.UpdatedAt = now
	if c.Status == "" {
		c.Status = domain.ClientStatusActive
	}

	cp := *c
	r.clients[c.ID] = &cp
	return nil
}

func (r *MockClientRepository) Update(ctx context.Context, c *domain.Client) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.clients[c.ID]; !ok {
		return appErrors.ErrNotFound
	}
	c.UpdatedAt = time.Now().UTC()
	cp := *c
	r.clients[c.ID] = &cp
	return nil
}

func (r *MockClientRepository) GetBalance(ctx context.Context, clientID int64) (decimal.Money, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var total int64
	for _, b := range r.balances[clientID] {
		if b.Type == domain.BalanceTypeCredit {
			total += int64(b.Amount)
		} else {
			total -= int64(b.Amount)
		}
	}
	return decimal.Money(total), nil
}

func (r *MockClientRepository) AddBalanceTransaction(ctx context.Context, b *domain.ClientBalance) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	b.ID = int64(len(r.balances[b.ClientID]) + 1)
	b.CreatedAt = time.Now().UTC()
	cp := *b
	r.balances[b.ClientID] = append(r.balances[b.ClientID], &cp)
	return nil
}
