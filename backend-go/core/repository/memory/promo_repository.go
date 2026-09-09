package memory

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
)

type MockPromoRepository struct {
	mu          sync.RWMutex
	promos      map[int64]*domain.Promo
	redemptions []*domain.PromoRedemption
	nextID      int64
}

func NewMockPromoRepository() *MockPromoRepository {
	return &MockPromoRepository{
		promos: make(map[int64]*domain.Promo),
		nextID: 1,
	}
}

func (r *MockPromoRepository) GetByID(ctx context.Context, id int64) (*domain.Promo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.promos[id]
	if !ok {
		return nil, appErrors.ErrNotFound
	}
	cp := *p
	return &cp, nil
}

func (r *MockPromoRepository) GetByCode(ctx context.Context, code string) (*domain.Promo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, p := range r.promos {
		if strings.EqualFold(p.Code, code) {
			cp := *p
			return &cp, nil
		}
	}
	return nil, appErrors.ErrNotFound
}

func (r *MockPromoRepository) GetRedemptionCount(ctx context.Context, promoID int64, clientID int64) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, red := range r.redemptions {
		if red.PromoID == promoID && red.ClientID == clientID {
			count++
		}
	}
	return count, nil
}

func (r *MockPromoRepository) IncrementUsed(ctx context.Context, promoID int64, clientID int64, orderID *int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	p, ok := r.promos[promoID]
	if !ok {
		return appErrors.ErrNotFound
	}

	// Atomic integrity check
	if p.MaxUses > 0 && p.UsedCount >= p.MaxUses {
		return errors.New("maximum promo usage limit reached")
	}

	// BUG-18 Fix: Check OncePerClient inside the lock
	if p.OncePerClient {
		count := 0
		for _, red := range r.redemptions {
			if red.PromoID == promoID && red.ClientID == clientID {
				count++
			}
		}
		if count > 0 {
			return errors.New("promo code can only be used once per client")
		}
	}

	p.UsedCount++
	p.UpdatedAt = time.Now().UTC()

	r.redemptions = append(r.redemptions, &domain.PromoRedemption{
		ID:        int64(len(r.redemptions) + 1),
		PromoID:   promoID,
		ClientID:  clientID,
		OrderID:   orderID,
		CreatedAt: time.Now().UTC(),
	})

	return nil
}

func (r *MockPromoRepository) List(ctx context.Context, limit, offset int) ([]*domain.Promo, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var all []*domain.Promo
	for _, p := range r.promos {
		cp := *p
		all = append(all, &cp)
	}

	total := len(all)
	if offset >= total {
		return []*domain.Promo{}, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return all[offset:end], total, nil
}

func (r *MockPromoRepository) Create(ctx context.Context, promo *domain.Promo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	promo.ID = r.nextID
	r.nextID++
	now := time.Now().UTC()
	promo.CreatedAt = now
	promo.UpdatedAt = now

	cp := *promo
	r.promos[promo.ID] = &cp
	return nil
}

func (r *MockPromoRepository) Update(ctx context.Context, promo *domain.Promo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.promos[promo.ID]; !ok {
		return appErrors.ErrNotFound
	}
	promo.UpdatedAt = time.Now().UTC()
	cp := *promo
	r.promos[promo.ID] = &cp
	return nil
}

func (r *MockPromoRepository) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.promos, id)
	return nil
}
