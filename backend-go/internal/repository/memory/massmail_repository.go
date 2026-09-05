package memory

import (
	"context"
	"sync"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
)

type MockMassMailRepository struct {
	mu        sync.RWMutex
	campaigns map[int64]*domain.MassMailCampaign
	nextID    int64
}

func NewMockMassMailRepository() *MockMassMailRepository {
	return &MockMassMailRepository{
		campaigns: make(map[int64]*domain.MassMailCampaign),
		nextID:    1,
	}
}

func (r *MockMassMailRepository) GetByID(ctx context.Context, id int64) (*domain.MassMailCampaign, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	c, ok := r.campaigns[id]
	if !ok {
		return nil, appErrors.ErrNotFound
	}
	cp := *c
	return &cp, nil
}

func (r *MockMassMailRepository) List(ctx context.Context, limit, offset int) ([]*domain.MassMailCampaign, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var all []*domain.MassMailCampaign
	for _, c := range r.campaigns {
		cp := *c
		all = append(all, &cp)
	}

	total := len(all)
	if offset >= total {
		return []*domain.MassMailCampaign{}, total, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}
	return all[offset:end], total, nil
}

func (r *MockMassMailRepository) Create(ctx context.Context, c *domain.MassMailCampaign) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	c.ID = r.nextID
	r.nextID++
	now := time.Now().UTC()
	c.CreatedAt = now
	if c.Status == "" {
		c.Status = domain.CampaignStatusDraft
	}

	cp := *c
	r.campaigns[c.ID] = &cp
	return nil
}

func (r *MockMassMailRepository) Update(ctx context.Context, c *domain.MassMailCampaign) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.campaigns[c.ID]; !ok {
		return appErrors.ErrNotFound
	}

	cp := *c
	r.campaigns[c.ID] = &cp
	return nil
}
