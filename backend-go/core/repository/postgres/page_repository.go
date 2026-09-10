package postgres

import (
	"context"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const pgCols = `id, title, slug, content, published, created_at, updated_at`

type PageRepository struct{ pool *pgxpool.Pool }

func NewPageRepository(p *pgxpool.Pool) *PageRepository { return &PageRepository{p} }

func scanPage(r pgx.Row) (*domain.Page, error) {
	var p domain.Page; err := r.Scan(&p.ID, &p.Title, &p.Slug, &p.Content, &p.Published, &p.CreatedAt, &p.UpdatedAt); return &p, err
}

func (r *PageRepository) GetBySlug(ctx context.Context, sl string) (*domain.Page, error) {
	return scanPage(r.pool.QueryRow(ctx, "SELECT "+pgCols+" FROM pages WHERE slug = $1", sl))
}

func (r *PageRepository) List(ctx context.Context, l, o int) ([]*domain.Page, int, error) {
	res, err := list(ctx, r.pool, "SELECT "+pgCols+" FROM pages ORDER BY id DESC LIMIT $1 OFFSET $2", scanPage, l, o)
	return res, total(ctx, r.pool, "SELECT COUNT(*) FROM pages"), err
}

func (r *PageRepository) Create(ctx context.Context, p *domain.Page) error {
	return r.pool.QueryRow(ctx, `INSERT INTO pages (title, slug, content, published, created_at, updated_at) VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id, created_at, updated_at`, p.Title, p.Slug, p.Content, p.Published).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *PageRepository) Update(ctx context.Context, p *domain.Page) error {
	_, err := r.pool.Exec(ctx, `UPDATE pages SET title = $1, slug = $2, content = $3, published = $4, updated_at = CURRENT_TIMESTAMP WHERE id = $5`, p.Title, p.Slug, p.Content, p.Published, p.ID); return err
}

func (r *PageRepository) Delete(ctx context.Context, id int64) error { _, err := r.pool.Exec(ctx, `DELETE FROM pages WHERE id = $1`, id); return err }
