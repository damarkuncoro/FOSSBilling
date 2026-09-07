package catalog

import (
	"context"
	"fmt"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/cache"
	"time"
)

type ProductService struct {
	productRepo domain.ProductRepository
	cache       cache.Cache
}

func NewProductService(productRepo domain.ProductRepository, c ...cache.Cache) *ProductService {
	var cc cache.Cache
	if len(c) > 0 {
		cc = c[0]
	}
	return &ProductService{
		productRepo: productRepo,
		cache:       cc,
	}
}

func (s *ProductService) ListProducts(ctx context.Context, limit, offset int) ([]*domain.Product, int, error) {
	if limit <= 0 {
		limit = 20
	}

	cacheKey := fmt.Sprintf("products:list:%d:%d", limit, offset)
	if s.cache != nil {
		var cached struct {
			Items []*domain.Product
			Total int
		}
		if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
			return cached.Items, cached.Total, nil
		}
	}

	items, total, err := s.productRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	if s.cache != nil {
		_ = s.cache.Set(ctx, cacheKey, struct {
			Items []*domain.Product
			Total int
		}{Items: items, Total: total}, 10*time.Minute)
	}

	return items, total, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id int64) (*domain.Product, error) {
	return s.productRepo.GetByID(ctx, id)
}

func (s *ProductService) CreateProduct(ctx context.Context, product *domain.Product) error {
	return s.productRepo.Create(ctx, product)
}

func (s *ProductService) UpdateProduct(ctx context.Context, product *domain.Product) error {
	return s.productRepo.Update(ctx, product)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int64) error {
	return s.productRepo.Delete(ctx, id)
}
