package memory

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/fossbilling/backend-go/internal/domain"
	appErrors "github.com/fossbilling/backend-go/pkg/errors"
)

type MockNewsRepository struct {
	mu     sync.RWMutex
	posts  map[int64]*domain.NewsPost
	nextID int64
}

func NewMockNewsRepository() *MockNewsRepository {
	return &MockNewsRepository{
		posts:  make(map[int64]*domain.NewsPost),
		nextID: 1,
	}
}

func (r *MockNewsRepository) GetByID(ctx context.Context, id int64) (*domain.NewsPost, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.posts[id]
	if !ok {
		return nil, appErrors.ErrNotFound
	}
	cp := *p
	return &cp, nil
}

func (r *MockNewsRepository) GetBySlug(ctx context.Context, slug string) (*domain.NewsPost, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, p := range r.posts {
		if strings.EqualFold(p.Slug, slug) {
			cp := *p
			return &cp, nil
		}
	}
	return nil, appErrors.ErrNotFound
}

func (r *MockNewsRepository) ListPublished(ctx context.Context, limit, offset int) ([]*domain.NewsPost, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var published []*domain.NewsPost
	for _, p := range r.posts {
		if p.Status == domain.NewsStatusPublished {
			cp := *p
			published = append(published, &cp)
		}
	}

	total := len(published)
	if offset >= total {
		return []*domain.NewsPost{}, total, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}
	return published[offset:end], total, nil
}

func (r *MockNewsRepository) ListAll(ctx context.Context, limit, offset int) ([]*domain.NewsPost, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var all []*domain.NewsPost
	for _, p := range r.posts {
		cp := *p
		all = append(all, &cp)
	}

	total := len(all)
	if offset >= total {
		return []*domain.NewsPost{}, total, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}
	return all[offset:end], total, nil
}

func (r *MockNewsRepository) Create(ctx context.Context, p *domain.NewsPost) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	p.ID = r.nextID
	r.nextID++
	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now
	if p.Status == domain.NewsStatusPublished && p.PublishedAt == nil {
		p.PublishedAt = &now
	}

	cp := *p
	r.posts[p.ID] = &cp
	return nil
}

func (r *MockNewsRepository) Update(ctx context.Context, p *domain.NewsPost) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.posts[p.ID]; !ok {
		return appErrors.ErrNotFound
	}

	p.UpdatedAt = time.Now().UTC()
	cp := *p
	r.posts[p.ID] = &cp
	return nil
}

func (r *MockNewsRepository) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.posts[id]; !ok {
		return appErrors.ErrNotFound
	}
	delete(r.posts, id)
	return nil
}
