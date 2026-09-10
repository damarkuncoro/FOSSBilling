package postgres

import (
	"context"
	"errors"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const newsCols = `id, admin_id, title, slug, content, status, published_at, created_at, updated_at`

type NewsRepository struct{ pool *pgxpool.Pool }

func NewNewsRepository(p *pgxpool.Pool) *NewsRepository { return &NewsRepository{p} }

func scanNews(r pgx.Row) (*domain.NewsPost, error) {
	var p domain.NewsPost; err := r.Scan(&p.ID, &p.AdminID, &p.Title, &p.Slug, &p.Content, &p.Status, &p.PublishedAt, &p.CreatedAt, &p.UpdatedAt); return &p, err
}

func (r *NewsRepository) GetByID(ctx context.Context, id int64) (*domain.NewsPost, error) {
	p, err := scanNews(r.pool.QueryRow(ctx, "SELECT "+newsCols+" FROM news_posts WHERE id = $1", id))
	if err != nil && errors.Is(err, pgx.ErrNoRows) { return nil, appErrors.ErrNotFound }; return p, err
}

func (r *NewsRepository) GetBySlug(ctx context.Context, sl string) (*domain.NewsPost, error) {
	p, err := scanNews(r.pool.QueryRow(ctx, "SELECT "+newsCols+" FROM news_posts WHERE slug = $1 AND status = 'published'", sl))
	if err != nil && errors.Is(err, pgx.ErrNoRows) { return nil, appErrors.ErrNotFound }; return p, err
}

func (r *NewsRepository) ListPublished(ctx context.Context, l, o int) ([]*domain.NewsPost, int, error) {
	res, err := list(ctx, r.pool, "SELECT "+newsCols+" FROM news_posts WHERE status = 'published' ORDER BY created_at DESC LIMIT $1 OFFSET $2", scanNews, l, o)
	return res, total(ctx, r.pool, "SELECT COUNT(*) FROM news_posts WHERE status = 'published'"), err
}

func (r *NewsRepository) ListAll(ctx context.Context, l, o int) ([]*domain.NewsPost, int, error) {
	res, err := list(ctx, r.pool, "SELECT "+newsCols+" FROM news_posts ORDER BY created_at DESC LIMIT $1 OFFSET $2", scanNews, l, o)
	return res, total(ctx, r.pool, "SELECT COUNT(*) FROM news_posts"), err
}

func (r *NewsRepository) Create(ctx context.Context, p *domain.NewsPost) error {
	return r.pool.QueryRow(ctx, `INSERT INTO news_posts (admin_id, title, slug, content, status, published_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id, created_at, updated_at`, p.AdminID, p.Title, p.Slug, p.Content, p.Status, p.PublishedAt).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *NewsRepository) Update(ctx context.Context, p *domain.NewsPost) error {
	err := r.pool.QueryRow(ctx, `UPDATE news_posts SET title = $2, slug = $3, content = $4, status = $5, published_at = $6, updated_at = CURRENT_TIMESTAMP WHERE id = $1 RETURNING updated_at`, p.ID, p.Title, p.Slug, p.Content, p.Status, p.PublishedAt).Scan(&p.UpdatedAt)
	if err != nil && errors.Is(err, pgx.ErrNoRows) { return appErrors.ErrNotFound }; return err
}

func (r *NewsRepository) Delete(ctx context.Context, id int64) error {
	t, err := r.pool.Exec(ctx, `DELETE FROM news_posts WHERE id = $1`, id)
	if err == nil && t.RowsAffected() == 0 { return appErrors.ErrNotFound }; return err
}
