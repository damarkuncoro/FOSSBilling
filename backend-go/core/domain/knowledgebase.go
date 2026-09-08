package domain

import (
	"context"
	"time"
)

type KBCategory struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Description string    `json:"description,omitempty"`
	Icon        string    `json:"icon,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type KBArticle struct {
	ID         int64     `json:"id"`
	CategoryID int64     `json:"category_id"`
	Title      string    `json:"title"`
	Slug       string    `json:"slug"`
	Content    string    `json:"content"`
	Status     string    `json:"status"` // "published", "draft"
	Views      int       `json:"views"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type KBRepository interface {
	ListCategories(ctx context.Context) ([]*KBCategory, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*KBCategory, error)
	CreateCategory(ctx context.Context, cat *KBCategory) error

	ListArticles(ctx context.Context, categoryID int64, limit, offset int) ([]*KBArticle, int, error)
	GetArticleBySlug(ctx context.Context, slug string) (*KBArticle, error)
	CreateArticle(ctx context.Context, art *KBArticle) error
	UpdateArticle(ctx context.Context, art *KBArticle) error
	DeleteArticle(ctx context.Context, id int64) error
	IncrementViews(ctx context.Context, id int64) error
}
