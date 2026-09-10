package memory

import (
	"context"
	"sync"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type MockEmailTemplateRepository struct {
	mu   sync.RWMutex
	tpls map[string]*domain.EmailTemplate
}

func NewMockEmailTemplateRepository() *MockEmailTemplateRepository {
	return &MockEmailTemplateRepository{
		tpls: make(map[string]*domain.EmailTemplate),
	}
}

func (m *MockEmailTemplateRepository) GetByCode(ctx context.Context, code string) (*domain.EmailTemplate, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	t, ok := m.tpls[code]
	if !ok {
		return nil, nil
	}
	return t, nil
}

func (m *MockEmailTemplateRepository) List(ctx context.Context) ([]*domain.EmailTemplate, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []*domain.EmailTemplate
	for _, t := range m.tpls {
		res = append(res, t)
	}
	return res, nil
}

func (m *MockEmailTemplateRepository) Update(ctx context.Context, t *domain.EmailTemplate) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	t.UpdatedAt = time.Now()
	m.tpls[t.Code] = t
	return nil
}
