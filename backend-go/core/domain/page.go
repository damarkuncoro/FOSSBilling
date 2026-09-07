package domain

import (
	"context"
	"time"
)

type Page struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Content   string    `json:"content"`
	Published bool      `json:"published"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PageRepository interface {
	GetBySlug(ctx context.Context, slug string) (*Page, error)
	List(ctx context.Context, limit, offset int) ([]*Page, int, error)
	Create(ctx context.Context, page *Page) error
	Update(ctx context.Context, page *Page) error
	Delete(ctx context.Context, id int64) error
}
