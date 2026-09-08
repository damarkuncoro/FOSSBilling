package memory

import (
	"context"
	"sync"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
)

type MockTaxRepository struct {
	mu    sync.RWMutex
	rules map[int64]*domain.TaxRule
	nextID int64
}

func NewMockTaxRepository() *MockTaxRepository {
	return &MockTaxRepository{
		rules: make(map[int64]*domain.TaxRule),
		nextID: 1,
	}
}

func (r *MockTaxRepository) GetByID(ctx context.Context, id int64) (*domain.TaxRule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rule, ok := r.rules[id]
	if !ok {
		return nil, appErrors.ErrNotFound
	}
	return rule, nil
}

func (r *MockTaxRepository) List(ctx context.Context) ([]*domain.TaxRule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*domain.TaxRule
	for _, rule := range r.rules {
		list = append(list, rule)
	}
	return list, nil
}

func (r *MockTaxRepository) GetByLocation(ctx context.Context, country, state string) (*domain.TaxRule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, rule := range r.rules {
		if !rule.IsActive {
			continue
		}
		c := ""
		if rule.Country != nil {
			c = *rule.Country
		}
		s := ""
		if rule.State != nil {
			s = *rule.State
		}
		if c == country && s == state {
			return rule, nil
		}
		if c == country && s == "" {
			return rule, nil
		}
		if c == "" && s == "" {
			return rule, nil
		}
	}
	return nil, nil
}

func (r *MockTaxRepository) Create(ctx context.Context, rule *domain.TaxRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	rule.ID = r.nextID
	r.nextID++
	r.rules[rule.ID] = rule
	return nil
}

func (r *MockTaxRepository) Update(ctx context.Context, rule *domain.TaxRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.rules[rule.ID]; !ok {
		return appErrors.ErrNotFound
	}
	r.rules[rule.ID] = rule
	return nil
}

func (r *MockTaxRepository) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.rules, id)
	return nil
}
