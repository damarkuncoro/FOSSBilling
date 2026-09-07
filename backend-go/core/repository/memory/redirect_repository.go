package memory

import (
	"context"
	"sync"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
)

type RedirectRepository struct {
	mu        sync.RWMutex
	redirects map[int64]*domain.Redirect
	nextID    int64
}

func NewRedirectRepository() *RedirectRepository {
	return &RedirectRepository{
		redirects: make(map[int64]*domain.Redirect),
		nextID:    1,
	}
}

func (r *RedirectRepository) List(ctx context.Context, limit, offset int) ([]*domain.Redirect, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var all []*domain.Redirect
	for _, item := range r.redirects {
		cpy := *item
		all = append(all, &cpy)
	}

	total := len(all)
	if offset >= total {
		return []*domain.Redirect{}, total, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}

	return all[offset:end], total, nil
}

func (r *RedirectRepository) GetByID(ctx context.Context, id int64) (*domain.Redirect, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.redirects[id]
	if !ok {
		return nil, appErrors.ErrNotFound
	}
	cpy := *item
	return &cpy, nil
}

func (r *RedirectRepository) GetByPath(ctx context.Context, path string) (*domain.Redirect, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, item := range r.redirects {
		if item.Path == path {
			cpy := *item
			return &cpy, nil
		}
	}
	return nil, appErrors.ErrNotFound
}

func (r *RedirectRepository) Create(ctx context.Context, redirect *domain.Redirect) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	redirect.ID = r.nextID
	r.nextID++
	now := time.Now().UTC()
	redirect.CreatedAt = now
	redirect.UpdatedAt = now
	if redirect.StatusCode == 0 {
		redirect.StatusCode = 301
	}

	cpy := *redirect
	r.redirects[redirect.ID] = &cpy
	return nil
}

func (r *RedirectRepository) Update(ctx context.Context, redirect *domain.Redirect) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.redirects[redirect.ID]; !ok {
		return appErrors.ErrNotFound
	}
	redirect.UpdatedAt = time.Now().UTC()
	cpy := *redirect
	r.redirects[redirect.ID] = &cpy
	return nil
}

func (r *RedirectRepository) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.redirects[id]; !ok {
		return appErrors.ErrNotFound
	}
	delete(r.redirects, id)
	return nil
}

func (r *RedirectRepository) IncrementHitCount(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, ok := r.redirects[id]
	if !ok {
		return appErrors.ErrNotFound
	}
	item.HitCount++
	return nil
}
