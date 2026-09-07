package memory

import (
	"context"
	"sync"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
)

type MockPageRepository struct {
	mu     sync.RWMutex
	pages  map[int64]*domain.Page
	nextID int64
}

func NewMockPageRepository() *MockPageRepository {
	return &MockPageRepository{
		pages:  make(map[int64]*domain.Page),
		nextID: 1,
	}
}

func (r *MockPageRepository) GetBySlug(ctx context.Context, slug string) (*domain.Page, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, p := range r.pages {
		if p.Slug == slug {
			return p, nil
		}
	}
	return nil, appErrors.ErrNotFound
}

func (r *MockPageRepository) List(ctx context.Context, limit, offset int) ([]*domain.Page, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var list []*domain.Page
	for _, p := range r.pages {
		list = append(list, p)
	}
	total := len(list)
	if offset >= total {
		return []*domain.Page{}, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return list[offset:end], total, nil
}

func (r *MockPageRepository) Create(ctx context.Context, page *domain.Page) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	page.ID = r.nextID
	r.nextID++
	if page.CreatedAt.IsZero() {
		page.CreatedAt = time.Now()
	}
	page.UpdatedAt = time.Now()
	r.pages[page.ID] = page
	return nil
}

func (r *MockPageRepository) Update(ctx context.Context, page *domain.Page) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.pages[page.ID]; !exists {
		return appErrors.ErrNotFound
	}
	page.UpdatedAt = time.Now()
	r.pages[page.ID] = page
	return nil
}

func (r *MockPageRepository) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.pages[id]; !exists {
		return appErrors.ErrNotFound
	}
	delete(r.pages, id)
	return nil
}
