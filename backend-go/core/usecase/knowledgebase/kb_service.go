package knowledgebase

import (
	"context"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type Service struct {
	repo domain.KBRepository
}

func NewKBService(repo domain.KBRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListCategories(ctx context.Context) ([]*domain.KBCategory, error) {
	return s.repo.ListCategories(ctx)
}

func (s *Service) ListArticles(ctx context.Context, categoryID int64, limit, offset int) ([]*domain.KBArticle, int, error) {
	return s.repo.ListArticles(ctx, categoryID, limit, offset)
}

func (s *Service) GetArticle(ctx context.Context, slug string) (*domain.KBArticle, error) {
	art, err := s.repo.GetArticleBySlug(ctx, slug)
	if err == nil && art != nil {
		_ = s.repo.IncrementViews(ctx, art.ID)
	}
	return art, err
}

func (s *Service) CreateArticle(ctx context.Context, art *domain.KBArticle) error {
	return s.repo.CreateArticle(ctx, art)
}

func (s *Service) UpdateArticle(ctx context.Context, art *domain.KBArticle) error {
	return s.repo.UpdateArticle(ctx, art)
}

func (s *Service) DeleteArticle(ctx context.Context, id int64) error {
	return s.repo.DeleteArticle(ctx, id)
}

func (s *Service) CreateCategory(ctx context.Context, cat *domain.KBCategory) error {
	return s.repo.CreateCategory(ctx, cat)
}
