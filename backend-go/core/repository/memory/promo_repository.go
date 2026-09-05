package memory

import (
	"context"
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
