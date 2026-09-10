package knowledgebase

import (
	"context"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type Service struct{ repo domain.KBRepository }

func NewKBService(r domain.KBRepository) *Service { return &Service{r} }

func (s *Service) ListCategories(ctx context.Context) ([]*domain.KBCategory, error) { return s.repo.ListCategories(ctx) }
func (s *Service) ListArticles(ctx context.Context, cid int64, l, o int) ([]*domain.KBArticle, int, error) { return s.repo.ListArticles(ctx, cid, l, o) }

func (s *Service) GetArticle(ctx context.Context, sl string) (*domain.KBArticle, error) {
	a, err := s.repo.GetArticleBySlug(ctx, sl); if err == nil && a != nil { _ = s.repo.IncrementViews(ctx, a.ID) }; return a, err
}

func (s *Service) CreateArticle(ctx context.Context, a *domain.KBArticle) error { return s.repo.CreateArticle(ctx, a) }
func (s *Service) UpdateArticle(ctx context.Context, a *domain.KBArticle) error { return s.repo.UpdateArticle(ctx, a) }
func (s *Service) DeleteArticle(ctx context.Context, id int64) error { return s.repo.DeleteArticle(ctx, id) }
func (s *Service) CreateCategory(ctx context.Context, c *domain.KBCategory) error { return s.repo.CreateCategory(ctx, c) }
