package memory

import (
	"context"
	"sync"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
)

type MockNotificationRepository struct {
	mu            sync.RWMutex
	notifications map[int64]*domain.Notification
	nextID        int64
}

func NewMockNotificationRepository() *MockNotificationRepository {
	return &MockNotificationRepository{
		notifications: make(map[int64]*domain.Notification),
		nextID:        1,
	}
}

func (r *MockNotificationRepository) Create(ctx context.Context, n *domain.Notification) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	n.ID = r.nextID
	r.nextID++
	if n.CreatedAt.IsZero() {
		n.CreatedAt = time.Now()
	}
	r.notifications[n.ID] = n
	return nil
}

func (r *MockNotificationRepository) ListByClientID(ctx context.Context, clientID int64, limit, offset int) ([]*domain.Notification, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var matched []*domain.Notification
	for _, n := range r.notifications {
		if n.ClientID == clientID {
			matched = append(matched, n)
		}
	}
	total := len(matched)
	if offset >= total {
		return []*domain.Notification{}, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return matched[offset:end], total, nil
}

func (r *MockNotificationRepository) MarkAsRead(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	n, ok := r.notifications[id]
	if !ok {
		return appErrors.ErrNotFound
	}
	n.IsRead = true
	return nil
}

func (r *MockNotificationRepository) MarkAllAsRead(ctx context.Context, clientID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, n := range r.notifications {
		if n.ClientID == clientID {
			n.IsRead = true
		}
	}
	return nil
}

func (r *MockNotificationRepository) DeleteOld(ctx context.Context, days int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cutoff := time.Now().AddDate(0, 0, -days)
	for id, n := range r.notifications {
		if n.CreatedAt.Before(cutoff) {
			delete(r.notifications, id)
		}
	}
	return nil
}
