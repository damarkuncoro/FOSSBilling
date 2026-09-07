package memory

import (
	"context"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type MockCatalogRepository struct{}

func NewMockCatalogRepository() *MockCatalogRepository {
	return &MockCatalogRepository{}
}

func (m *MockCatalogRepository) ListCategories(ctx context.Context) ([]*domain.ProductCategory, error) {
	return []*domain.ProductCategory{
		{ID: 1, Title: "Hosting", Slug: "hosting"},
	}, nil
}

func (m *MockCatalogRepository) GetCategoryByID(ctx context.Context, id int64) (*domain.ProductCategory, error) {
	return &domain.ProductCategory{ID: 1, Title: "Hosting", Slug: "hosting"}, nil
}

func (m *MockCatalogRepository) ListServers(ctx context.Context) ([]*domain.Server, error) {
	return []*domain.Server{
		{ID: 1, Name: "Mock Server", Hostname: "localhost"},
	}, nil
}

func (m *MockCatalogRepository) GetServerByID(ctx context.Context, id int64) (*domain.Server, error) {
	return &domain.Server{ID: 1, Name: "Mock Server", Hostname: "localhost"}, nil
}

func (m *MockCatalogRepository) ListTlds(ctx context.Context) ([]*domain.TLD, error) {
	return []*domain.TLD{
		{ID: 1, Tld: ".com", RegistrarID: "mock"},
	}, nil
}
