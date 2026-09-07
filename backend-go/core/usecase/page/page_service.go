package page

import (
	"context"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type PageService struct {
	pageRepo domain.PageRepository
}

func NewPageService(pageRepo domain.PageRepository) *PageService {
	return &PageService{pageRepo: pageRepo}
}

func (s *PageService) GetPage(ctx context.Context, slug string) (*domain.Page, error) {
	return s.pageRepo.GetBySlug(ctx, slug)
}

func (s *PageService) ListPages(ctx context.Context, limit, offset int) ([]*domain.Page, int, error) {
	return s.pageRepo.List(ctx, limit, offset)
}

func (s *PageService) CreatePage(ctx context.Context, page *domain.Page) error {
	return s.pageRepo.Create(ctx, page)
}

func (s *PageService) UpdatePage(ctx context.Context, page *domain.Page) error {
	return s.pageRepo.Update(ctx, page)
}

func (s *PageService) DeletePage(ctx context.Context, id int64) error {
	return s.pageRepo.Delete(ctx, id)
}
