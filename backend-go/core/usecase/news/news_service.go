package news

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

var (
	ErrEmptyTitle   = errors.New("news title cannot be empty")
	ErrEmptyContent = errors.New("news content cannot be empty")
)

type CreateNewsDTO struct {
	AdminID int64             `json:"admin_id"`
	Title   string            `json:"title"`
	Slug    string            `json:"slug"`
	Content string            `json:"content"`
	Status  domain.NewsStatus `json:"status"`
}

type UpdateNewsDTO struct {
	Title   string            `json:"title"`
	Slug    string            `json:"slug"`
	Content string            `json:"content"`
	Status  domain.NewsStatus `json:"status"`
}

type NewsService struct {
	newsRepo domain.NewsRepository
}

func NewNewsService(newsRepo domain.NewsRepository) *NewsService {
	return &NewsService{newsRepo: newsRepo}
}

func slugify(text string) string {
	text = strings.ToLower(strings.TrimSpace(text))
	reg := regexp.MustCompile("[^a-z0-9]+")
	slug := reg.ReplaceAllString(text, "-")
	return strings.Trim(slug, "-")
}

func (s *NewsService) ListPublished(ctx context.Context, limit, offset int) ([]*domain.NewsPost, int, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.newsRepo.ListPublished(ctx, limit, offset)
}

func (s *NewsService) GetBySlug(ctx context.Context, slug string) (*domain.NewsPost, error) {
	return s.newsRepo.GetBySlug(ctx, slug)
}

func (s *NewsService) ListAll(ctx context.Context, limit, offset int) ([]*domain.NewsPost, int, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.newsRepo.ListAll(ctx, limit, offset)
}

func (s *NewsService) GetByID(ctx context.Context, id int64) (*domain.NewsPost, error) {
	return s.newsRepo.GetByID(ctx, id)
}

func (s *NewsService) Create(ctx context.Context, dto CreateNewsDTO) (*domain.NewsPost, error) {
	if strings.TrimSpace(dto.Title) == "" {
		return nil, ErrEmptyTitle
	}
	if strings.TrimSpace(dto.Content) == "" {
		return nil, ErrEmptyContent
	}

	slug := dto.Slug
	if slug == "" {
		slug = slugify(dto.Title)
	}

	if dto.Status == "" {
		dto.Status = domain.NewsStatusDraft
	}

	post := &domain.NewsPost{
		AdminID: dto.AdminID,
		Title:   dto.Title,
		Slug:    slug,
		Content: dto.Content,
		Status:  dto.Status,
	}

	if err := s.newsRepo.Create(ctx, post); err != nil {
		return nil, err
	}

	return s.newsRepo.GetByID(ctx, post.ID)
}

func (s *NewsService) Update(ctx context.Context, id int64, dto UpdateNewsDTO) (*domain.NewsPost, error) {
	post, err := s.newsRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if dto.Title != "" {
		post.Title = dto.Title
	}
	if dto.Slug != "" {
		post.Slug = dto.Slug
	}
	if dto.Content != "" {
		post.Content = dto.Content
	}
	if dto.Status != "" {
		post.Status = dto.Status
		if dto.Status == domain.NewsStatusPublished && post.PublishedAt == nil {
			now := time.Now().UTC()
			post.PublishedAt = &now
		}
	}

	if err := s.newsRepo.Update(ctx, post); err != nil {
		return nil, err
	}

	return s.newsRepo.GetByID(ctx, id)
}

func (s *NewsService) Delete(ctx context.Context, id int64) error {
	return s.newsRepo.Delete(ctx, id)
}
