package catalog

import (
	"context"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type ProductService struct {
	productRepo domain.ProductRepository
}

func NewProductService(productRepo domain.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

func (s *ProductService) ListProducts(ctx context.Context, limit, offset int) ([]*domain.Product, int, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.productRepo.List(ctx, limit, offset)
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
