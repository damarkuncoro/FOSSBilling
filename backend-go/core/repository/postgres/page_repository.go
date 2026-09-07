package postgres

import (
	"context"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PageRepository struct {
	pool *pgxpool.Pool
}

func NewPageRepository(pool *pgxpool.Pool) *PageRepository {
	return &PageRepository{pool: pool}
}

func (r *PageRepository) GetBySlug(ctx context.Context, slug string) (*domain.Page, error) {
	var p domain.Page
	err := r.pool.QueryRow(ctx, `SELECT id, title, slug, content, published, created_at, updated_at FROM pages WHERE slug = $1`, slug).
		Scan(&p.ID, &p.Title, &p.Slug, &p.Content, &p.Published, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PageRepository) List(ctx context.Context, limit, offset int) ([]*domain.Page, int, error) {
	var total int
	_ = r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM pages`).Scan(&total)

	rows, err := r.pool.Query(ctx, `SELECT id, title, slug, content, published, created_at, updated_at FROM pages ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*domain.Page
	for rows.Next() {
		p := &domain.Page{}
		if err := rows.Scan(&p.ID, &p.Title, &p.Slug, &p.Content, &p.Published, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		list = append(list, p)
	}
	return list, total, nil
}

func (r *PageRepository) Create(ctx context.Context, p *domain.Page) error {
	return r.pool.QueryRow(ctx, `INSERT INTO pages (title, slug, content, published, created_at, updated_at) VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id, created_at, updated_at`,
		p.Title, p.Slug, p.Content, p.Published).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *PageRepository) Update(ctx context.Context, p *domain.Page) error {
	_, err := r.pool.Exec(ctx, `UPDATE pages SET title = $1, slug = $2, content = $3, published = $4, updated_at = CURRENT_TIMESTAMP WHERE id = $5`,
		p.Title, p.Slug, p.Content, p.Published, p.ID)
	return err
}

func (r *PageRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM pages WHERE id = $1`, id)
	return err
}
