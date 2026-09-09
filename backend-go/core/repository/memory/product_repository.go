package memory

import (
	"context"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"sync"
)

type MockProductRepository struct {
	mu       sync.RWMutex
	products map[int64]*domain.Product
}

func NewMockProductRepository() *MockProductRepository {
	return &MockProductRepository{
		products: make(map[int64]*domain.Product),
	}
}

func (m *MockProductRepository) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if p, ok := m.products[id]; ok {
		return p, nil
	}
	return nil, nil
}

func (m *MockProductRepository) GetBySlug(ctx context.Context, slug string) (*domain.Product, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, p := range m.products {
		if p.Slug == slug {
			return p, nil
		}
	}
	return nil, nil
}

func (m *MockProductRepository) List(ctx context.Context, limit, offset int) ([]*domain.Product, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var list []*domain.Product
	for _, p := range m.products {
		list = append(list, p)
	}
	return list, len(list), nil
}

func (m *MockProductRepository) Create(ctx context.Context, p *domain.Product) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if p.ID == 0 {
		p.ID = int64(len(m.products) + 1)
	}
	m.products[p.ID] = p
	return nil
}

func (m *MockProductRepository) Update(ctx context.Context, p *domain.Product) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.products[p.ID] = p
	return nil
}

func (m *MockProductRepository) DecrementStock(ctx context.Context, id int64, quantity int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	p, ok := m.products[id]
	if !ok {
		return nil
	}
	if p.Stock < quantity {
		return nil // Simplified for mock
	}
	p.Stock -= quantity
	return nil
}

func (m *MockProductRepository) Delete(ctx context.Context, id int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.products, id)
	return nil
}
