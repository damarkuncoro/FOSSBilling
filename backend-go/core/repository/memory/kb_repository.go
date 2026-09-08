package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type MockKBRepository struct {
	mu         sync.RWMutex
	categories map[int64]*domain.KBCategory
	articles   map[int64]*domain.KBArticle
}

func NewMockKBRepository() *MockKBRepository {
	return &MockKBRepository{
		categories: make(map[int64]*domain.KBCategory),
		articles:   make(map[int64]*domain.KBArticle),
	}
}

func (m *MockKBRepository) ListCategories(ctx context.Context) ([]*domain.KBCategory, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var list []*domain.KBCategory
	for _, c := range m.categories {
		list = append(list, c)
	}
	return list, nil
}

func (m *MockKBRepository) GetCategoryBySlug(ctx context.Context, slug string) (*domain.KBCategory, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, c := range m.categories {
		if c.Slug == slug {
			return c, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *MockKBRepository) CreateCategory(ctx context.Context, cat *domain.KBCategory) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cat.ID = int64(len(m.categories) + 1)
	m.categories[cat.ID] = cat
	return nil
}

func (m *MockKBRepository) ListArticles(ctx context.Context, categoryID int64, limit, offset int) ([]*domain.KBArticle, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var list []*domain.KBArticle
	for _, a := range m.articles {
		if categoryID == 0 || a.CategoryID == categoryID {
			list = append(list, a)
		}
	}
	return list, len(list), nil
}

func (m *MockKBRepository) GetArticleBySlug(ctx context.Context, slug string) (*domain.KBArticle, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, a := range m.articles {
		if a.Slug == slug {
			return a, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *MockKBRepository) CreateArticle(ctx context.Context, art *domain.KBArticle) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	art.ID = int64(len(m.articles) + 1)
	m.articles[art.ID] = art
	return nil
}

func (m *MockKBRepository) UpdateArticle(ctx context.Context, art *domain.KBArticle) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.articles[art.ID] = art
	return nil
}

func (m *MockKBRepository) DeleteArticle(ctx context.Context, id int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.articles, id)
	return nil
}

func (m *MockKBRepository) IncrementViews(ctx context.Context, id int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if a, ok := m.articles[id]; ok {
		a.Views++
	}
	return nil
}
