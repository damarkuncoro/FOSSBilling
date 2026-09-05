package memory

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
)

type MockCurrencyRepository struct {
	mu         sync.RWMutex
	currencies map[string]*domain.Currency
	nextID     int64
}

func NewMockCurrencyRepository() *MockCurrencyRepository {
	repo := &MockCurrencyRepository{
		currencies: make(map[string]*domain.Currency),
		nextID:     1,
	}
	now := time.Now().UTC()
	repo.currencies["USD"] = &domain.Currency{
		ID: 1, Code: "USD", Title: "US Dollar", ConversionRate: 1.0,
		Format: "$ {{price}}", PriceFormat: "2", IsDefault: true, CreatedAt: now, UpdatedAt: now,
	}
	repo.currencies["IDR"] = &domain.Currency{
		ID: 2, Code: "IDR", Title: "Indonesian Rupiah", ConversionRate: 15500.0,
		Format: "Rp {{price}}", PriceFormat: "0", IsDefault: false, CreatedAt: now, UpdatedAt: now,
	}
	repo.currencies["EUR"] = &domain.Currency{
		ID: 3, Code: "EUR", Title: "Euro", ConversionRate: 0.92,
		Format: "€ {{price}}", PriceFormat: "2", IsDefault: false, CreatedAt: now, UpdatedAt: now,
	}
	repo.nextID = 4
	return repo
}

func (r *MockCurrencyRepository) GetByCode(ctx context.Context, code string) (*domain.Currency, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	c, ok := r.currencies[strings.ToUpper(code)]
	if !ok {
		return nil, appErrors.ErrNotFound
	}
	cp := *c
	return &cp, nil
}

func (r *MockCurrencyRepository) GetDefault(ctx context.Context) (*domain.Currency, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, c := range r.currencies {
		if c.IsDefault {
			cp := *c
			return &cp, nil
		}
	}
	return nil, appErrors.ErrNotFound
}

func (r *MockCurrencyRepository) List(ctx context.Context) ([]*domain.Currency, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var list []*domain.Currency
	for _, c := range r.currencies {
		cp := *c
		list = append(list, &cp)
	}
	return list, nil
}

func (r *MockCurrencyRepository) Create(ctx context.Context, c *domain.Currency) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	code := strings.ToUpper(c.Code)
	if _, ok := r.currencies[code]; ok {
		return appErrors.ErrDuplicateEntry
	}

	c.ID = r.nextID
	r.nextID++
	c.Code = code
	now := time.Now().UTC()
	c.CreatedAt = now
	c.UpdatedAt = now

	cp := *c
	r.currencies[code] = &cp
	return nil
}

func (r *MockCurrencyRepository) Update(ctx context.Context, c *domain.Currency) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	code := strings.ToUpper(c.Code)
	if _, ok := r.currencies[code]; !ok {
		return appErrors.ErrNotFound
	}

	c.UpdatedAt = time.Now().UTC()
	cp := *c
	r.currencies[code] = &cp
	return nil
}

func (r *MockCurrencyRepository) SetDefault(ctx context.Context, code string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	targetCode := strings.ToUpper(code)
	if _, ok := r.currencies[targetCode]; !ok {
		return appErrors.ErrNotFound
	}

	for _, c := range r.currencies {
		c.IsDefault = (c.Code == targetCode)
	}
	return nil
}

func (r *MockCurrencyRepository) Delete(ctx context.Context, code string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	targetCode := strings.ToUpper(code)
	if _, ok := r.currencies[targetCode]; !ok {
		return appErrors.ErrNotFound
	}
	delete(r.currencies, targetCode)
	return nil
}
