package catalog

import (
	"context"
	"fmt"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/cache"
)

type ProductService struct {
	repo  domain.ProductRepository; cache cache.Cache
}

func NewProductService(r domain.ProductRepository, c ...cache.Cache) *ProductService {
	var cc cache.Cache; if len(c) > 0 { cc = c[0] }
	return &ProductService{r, cc}
}

func (s *ProductService) ListProducts(ctx context.Context, l, o int) ([]*domain.Product, int, error) {
	if l <= 0 { l = 20 }
	key := fmt.Sprintf("products:list:%d:%d", l, o)
	if s.cache != nil {
		var res struct { Items []*domain.Product; Total int }
		if err := s.cache.Get(ctx, key, &res); err == nil { return res.Items, res.Total, nil }
	}
	is, t, err := s.repo.List(ctx, l, o); if err != nil { return nil, 0, err }
	if s.cache != nil { _ = s.cache.Set(ctx, key, struct { Items []*domain.Product; Total int }{is, t}, 10*time.Minute) }
	return is, t, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id int64) (*domain.Product, error) { return s.repo.GetByID(ctx, id) }
func (s *ProductService) CreateProduct(ctx context.Context, p *domain.Product) error { return s.repo.Create(ctx, p) }
func (s *ProductService) UpdateProduct(ctx context.Context, p *domain.Product) error { return s.repo.Update(ctx, p) }
func (s *ProductService) DeleteProduct(ctx context.Context, id int64) error { return s.repo.Delete(ctx, id) }
