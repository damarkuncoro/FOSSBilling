package memory

import (
	"context"
	"sync"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type MockActivityRepository struct {
	mu   sync.RWMutex
	logs []*domain.Activity
}

func NewMockActivityRepository() *MockActivityRepository {
	return &MockActivityRepository{
		logs: []*domain.Activity{},
	}
}

func (m *MockActivityRepository) Log(ctx context.Context, a *domain.Activity) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	a.ID = int64(len(m.logs) + 1)
	a.CreatedAt = time.Now()
	m.logs = append(m.logs, a)
	return nil
}

func (m *MockActivityRepository) List(ctx context.Context, limit, offset int) ([]*domain.Activity, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.logs, len(m.logs), nil
}

func (m *MockActivityRepository) ListByClientID(ctx context.Context, clientID int64, limit, offset int) ([]*domain.Activity, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []*domain.Activity
	for _, l := range m.logs {
		if l.ClientID != nil && *l.ClientID == clientID {
			result = append(result, l)
		}
	}
	return result, len(result), nil
}

func (m *MockActivityRepository) DeleteOld(ctx context.Context, days int) error {
	return nil
}

func (m *MockActivityRepository) GetTrend(ctx context.Context, days int) (map[string]int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := make(map[string]int)
	now := time.Now()
	for i := 0; i < days; i++ {
		res[now.AddDate(0, 0, -i).Format("2006-01-02")] = 0
	}
	for _, l := range m.logs {
		d := l.CreatedAt.Format("2006-01-02")
		if _, ok := res[d]; ok {
			res[d]++
		}
	}
	return res, nil
}
