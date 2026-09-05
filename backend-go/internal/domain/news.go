package domain

import (
	"context"
	"time"
)

type NewsStatus string

const (
	NewsStatusDraft     NewsStatus = "draft"
	NewsStatusPublished NewsStatus = "published"
)

type NewsPost struct {
	ID          int64      `json:"id"`
	AdminID     int64      `json:"admin_id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Content     string     `json:"content"`
	Status      NewsStatus `json:"status"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type NewsRepository interface {
	GetByID(ctx context.Context, id int64) (*NewsPost, error)
	GetBySlug(ctx context.Context, slug string) (*NewsPost, error)
	ListPublished(ctx context.Context, limit, offset int) ([]*NewsPost, int, error)
	ListAll(ctx context.Context, limit, offset int) ([]*NewsPost, int, error)
	Create(ctx context.Context, post *NewsPost) error
	Update(ctx context.Context, post *NewsPost) error
	Delete(ctx context.Context, id int64) error
}
