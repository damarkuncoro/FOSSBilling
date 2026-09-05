package memory

import (
	"context"
	"sync"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
)

type MockOrderRepository struct {
	mu     sync.RWMutex
	orders map[int64]*domain.Order
	nextID int64
}

func NewMockOrderRepository() *MockOrderRepository {
	return &MockOrderRepository{
		orders: make(map[int64]*domain.Order),
		nextID: 1,
	}
}

func (r *MockOrderRepository) GetByID(ctx context.Context, id int64) (*domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[id]
	if !ok {
		return nil, appErrors.ErrNotFound
	}
	cp := *order
	return &cp, nil
}

func (r *MockOrderRepository) ListByClientID(ctx context.Context, clientID int64, limit, offset int) ([]*domain.Order, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var matched []*domain.Order
	for _, o := range r.orders {
		if o.ClientID == clientID {
			cp := *o
			matched = append(matched, &cp)
		}
	}

	total := len(matched)
	if offset >= total {
		return []*domain.Order{}, total, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}
	return matched[offset:end], total, nil
}

func (r *MockOrderRepository) List(ctx context.Context, limit, offset int) ([]*domain.Order, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var all []*domain.Order
	for _, o := range r.orders {
		cp := *o
		all = append(all, &cp)
	}

	total := len(all)
	if offset >= total {
		return []*domain.Order{}, total, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}
	return all[offset:end], total, nil
}


func (r *MockOrderRepository) ListDueOrders(ctx context.Context, dueBefore time.Time) ([]*domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Order
	for _, o := range r.orders {
		if o.Status == domain.OrderStatusActive && o.NextDueDate != nil && (o.NextDueDate.Before(dueBefore) || o.NextDueDate.Equal(dueBefore)) {
			cp := *o
			result = append(result, &cp)
		}
	}
	return result, nil
}

func (r *MockOrderRepository) ListOverdueSuspensions(ctx context.Context, overdueDays int) ([]*domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	threshold := time.Now().UTC().AddDate(0, 0, -overdueDays)
	var result []*domain.Order
	for _, o := range r.orders {
		if o.Status == domain.OrderStatusActive && o.ExpiresAt != nil && o.ExpiresAt.Before(threshold) {
			cp := *o
			result = append(result, &cp)
		}
	}
	return result, nil
}

func (r *MockOrderRepository) Create(ctx context.Context, o *domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	o.ID = r.nextID
	r.nextID++
	now := time.Now().UTC()
	o.CreatedAt = now
	o.UpdatedAt = now
	if o.Status == "" {
		o.Status = domain.OrderStatusPendingSetup
	}

	cp := *o
	r.orders[o.ID] = &cp
	return nil
}

func (r *MockOrderRepository) Update(ctx context.Context, o *domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.orders[o.ID]; !ok {
		return appErrors.ErrNotFound
	}
	o.UpdatedAt = time.Now().UTC()
	cp := *o
	r.orders[o.ID] = &cp
	return nil
}

func (r *MockOrderRepository) UpdateStatus(ctx context.Context, id int64, status domain.OrderStatus, reason *string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	o, ok := r.orders[id]
	if !ok {
		return appErrors.ErrNotFound
	}

	o.Status = status
	o.UpdatedAt = time.Now().UTC()
	if status == domain.OrderStatusSuspended {
		now := time.Now().UTC()
		o.SuspendedAt = &now
		o.SuspensionReason = reason
	} else if status == domain.OrderStatusActive {
		o.SuspendedAt = nil
		o.SuspensionReason = nil
	}

	return nil
}
