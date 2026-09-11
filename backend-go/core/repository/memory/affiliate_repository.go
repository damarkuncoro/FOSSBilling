package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type MockAffiliateRepository struct {
	mu   sync.RWMutex
	affs map[int64]*domain.Affiliate
	refs []*domain.AffiliateReferral
	next int64
}

func NewMockAffiliateRepository() *MockAffiliateRepository {
	return &MockAffiliateRepository{affs: make(map[int64]*domain.Affiliate), next: 1}
}

func (r *MockAffiliateRepository) GetByClientID(ctx context.Context, clientID int64) (*domain.Affiliate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, a := range r.affs {
		if a.ClientID == clientID {
			return a, nil
		}
	}
	return nil, errors.New("not found")
}

func (r *MockAffiliateRepository) Create(ctx context.Context, aff *domain.Affiliate) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	aff.ID = r.next
	r.next++
	r.affs[aff.ID] = aff
	return nil
}

func (r *MockAffiliateRepository) Update(ctx context.Context, aff *domain.Affiliate) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.affs[aff.ID] = aff
	return nil
}

func (r *MockAffiliateRepository) AddReferral(ctx context.Context, ref *domain.AffiliateReferral) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	ref.ID = int64(len(r.refs) + 1)
	r.refs = append(r.refs, ref)
	return nil
}

func (r *MockAffiliateRepository) ListReferrals(ctx context.Context, affiliateID int64) ([]*domain.AffiliateReferral, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []*domain.AffiliateReferral
	for _, ref := range r.refs {
		if ref.AffiliateID == affiliateID {
			res = append(res, ref)
		}
	}
	return res, nil
}
