package memory

import (
	"context"
	"sync"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
)

type MockNotificationRepository struct {
	mu                 sync.RWMutex
	notifications      map[int64]*domain.Notification
	adminNotifications map[int64]*domain.AdminNotification
	nextID             int64
	nextAdminID        int64
}

func NewMockNotificationRepository() *MockNotificationRepository {
	return &MockNotificationRepository{
		notifications:      make(map[int64]*domain.Notification),
		adminNotifications: make(map[int64]*domain.AdminNotification),
		nextID:             1,
		nextAdminID:        1,
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

// --- Admin Notifications ---

func (r *MockNotificationRepository) CreateAdmin(ctx context.Context, n *domain.AdminNotification) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	n.ID = r.nextAdminID
	r.nextAdminID++
	if n.CreatedAt.IsZero() {
		n.CreatedAt = time.Now()
	}
	r.adminNotifications[n.ID] = n
	return nil
}

func (r *MockNotificationRepository) ListAdmin(ctx context.Context, limit, offset int, unreadOnly bool) ([]*domain.AdminNotification, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var matched []*domain.AdminNotification
	for _, n := range r.adminNotifications {
		if unreadOnly && n.IsRead {
			continue
		}
		matched = append(matched, n)
	}

	total := len(matched)
	if offset >= total {
		return []*domain.AdminNotification{}, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return matched[offset:end], total, nil
}

func (r *MockNotificationRepository) MarkAdminAsRead(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	n, ok := r.adminNotifications[id]
	if !ok {
		return appErrors.ErrNotFound
	}
	n.IsRead = true
	return nil
}

func (r *MockNotificationRepository) MarkAdminAllAsRead(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, n := range r.adminNotifications {
		n.IsRead = true
	}
	return nil
}

func (r *MockNotificationRepository) DeleteAdminOld(ctx context.Context, days int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cutoff := time.Now().AddDate(0, 0, -days)
	for id, n := range r.adminNotifications {
		if n.CreatedAt.Before(cutoff) {
			delete(r.adminNotifications, id)
		}
	}
	return nil
}
