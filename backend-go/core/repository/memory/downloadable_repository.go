package memory

import (
	"context"
	"sync"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
)

type MockDownloadableRepository struct {
	mu     sync.RWMutex
	files  map[int64]*domain.DownloadableFile
	nextID int64
}

func NewMockDownloadableRepository() *MockDownloadableRepository {
	return &MockDownloadableRepository{
		files:  make(map[int64]*domain.DownloadableFile),
		nextID: 1,
	}
}

func (r *MockDownloadableRepository) GetByID(ctx context.Context, id int64) (*domain.DownloadableFile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	f, ok := r.files[id]
	if !ok {
		return nil, appErrors.ErrNotFound
	}
	cp := *f
	return &cp, nil
}

func (r *MockDownloadableRepository) GetByProductID(ctx context.Context, productID int64) (*domain.DownloadableFile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, f := range r.files {
		if f.ProductID == productID {
			cp := *f
			return &cp, nil
		}
	}
	return nil, appErrors.ErrNotFound
}

func (r *MockDownloadableRepository) ListByClientID(ctx context.Context, clientID int64) ([]*domain.DownloadableFile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var list []*domain.DownloadableFile
	for _, f := range r.files {
		cp := *f
		list = append(list, &cp)
	}
	return list, nil
}

func (r *MockDownloadableRepository) Create(ctx context.Context, f *domain.DownloadableFile) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	f.ID = r.nextID
	r.nextID++
	now := time.Now().UTC()
	f.CreatedAt = now
	f.UpdatedAt = now

	cp := *f
	r.files[f.ID] = &cp
	return nil
}

func (r *MockDownloadableRepository) IncrementDownloads(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	f, ok := r.files[id]
	if !ok {
		return appErrors.ErrNotFound
	}
	f.Downloads++
	f.UpdatedAt = time.Now().UTC()
	return nil
}

func (r *MockDownloadableRepository) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.files[id]; !ok {
		return appErrors.ErrNotFound
	}
	delete(r.files, id)
	return nil
}
