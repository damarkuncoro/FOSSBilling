package memory

import (
	"context"
	"sync"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
)

type MockAPIKeyRepository struct {
	mu     sync.RWMutex
	keys   map[int64]*domain.APIKey
	nextID int64
}

func NewMockAPIKeyRepository() *MockAPIKeyRepository {
	return &MockAPIKeyRepository{
		keys:   make(map[int64]*domain.APIKey),
		nextID: 1,
	}
}

func (r *MockAPIKeyRepository) GetByID(ctx context.Context, id int64) (*domain.APIKey, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	k, ok := r.keys[id]
	if !ok {
		return nil, appErrors.ErrNotFound
	}
	cp := *k
	return &cp, nil
}

func (r *MockAPIKeyRepository) GetByKey(ctx context.Context, key string) (*domain.APIKey, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, k := range r.keys {
		if k.Key == key {
			cp := *k
			return &cp, nil
		}
	}
	return nil, appErrors.ErrNotFound
}

func (r *MockAPIKeyRepository) ListByClientID(ctx context.Context, clientID int64) ([]*domain.APIKey, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var list []*domain.APIKey
	for _, k := range r.keys {
		if k.ClientID == clientID {
			cp := *k
			list = append(list, &cp)
		}
	}
	return list, nil
}

func (r *MockAPIKeyRepository) Create(ctx context.Context, k *domain.APIKey) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	k.ID = r.nextID
	r.nextID++
	now := time.Now().UTC()
	k.CreatedAt = now
	k.UpdatedAt = now

	cp := *k
	r.keys[k.ID] = &cp
	return nil
}

func (r *MockAPIKeyRepository) Delete(ctx context.Context, id, clientID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	k, ok := r.keys[id]
	if !ok {
		return appErrors.ErrNotFound
	}
	if k.ClientID != clientID {
		return appErrors.ErrForbidden
	}

	delete(r.keys, id)
	return nil
}
