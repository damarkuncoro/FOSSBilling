package postgres

import (
	"context"
	"errors"

	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NewsRepository struct {
	pool *pgxpool.Pool
}

func NewNewsRepository(pool *pgxpool.Pool) *NewsRepository {
	return &NewsRepository{pool: pool}
}

func (r *NewsRepository) GetByID(ctx context.Context, id int64) (*domain.NewsPost, error) {
	query := `
		SELECT id, admin_id, title, slug, content, status, published_at, created_at, updated_at
		FROM news_posts
		WHERE id = $1
	`
	var p domain.NewsPost
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.AdminID, &p.Title, &p.Slug, &p.Content, &p.Status, &p.PublishedAt, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *NewsRepository) GetBySlug(ctx context.Context, slug string) (*domain.NewsPost, error) {
	query := `
		SELECT id, admin_id, title, slug, content, status, published_at, created_at, updated_at
		FROM news_posts
		WHERE slug = $1 AND status = 'published'
	`
	var p domain.NewsPost
	err := r.pool.QueryRow(ctx, query, slug).Scan(
		&p.ID, &p.AdminID, &p.Title, &p.Slug, &p.Content, &p.Status, &p.PublishedAt, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *NewsRepository) ListPublished(ctx context.Context, limit, offset int) ([]*domain.NewsPost, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM news_posts WHERE status = 'published'`).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, admin_id, title, slug, content, status, published_at, created_at, updated_at
		FROM news_posts
		WHERE status = 'published'
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*domain.NewsPost
	for rows.Next() {
		var p domain.NewsPost
		if err := rows.Scan(&p.ID, &p.AdminID, &p.Title, &p.Slug, &p.Content, &p.Status, &p.PublishedAt, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		list = append(list, &p)
	}
	return list, total, rows.Err()
}

func (r *NewsRepository) ListAll(ctx context.Context, limit, offset int) ([]*domain.NewsPost, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM news_posts`).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, admin_id, title, slug, content, status, published_at, created_at, updated_at
		FROM news_posts
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*domain.NewsPost
	for rows.Next() {
		var p domain.NewsPost
		if err := rows.Scan(&p.ID, &p.AdminID, &p.Title, &p.Slug, &p.Content, &p.Status, &p.PublishedAt, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		list = append(list, &p)
	}
	return list, total, rows.Err()
}

func (r *NewsRepository) Create(ctx context.Context, p *domain.NewsPost) error {
	query := `
		INSERT INTO news_posts (admin_id, title, slug, content, status, published_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at
	`
	return r.pool.QueryRow(ctx, query, p.AdminID, p.Title, p.Slug, p.Content, p.Status, p.PublishedAt).
		Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *NewsRepository) Update(ctx context.Context, p *domain.NewsPost) error {
	query := `
		UPDATE news_posts
		SET title = $2, slug = $3, content = $4, status = $5, published_at = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at
	`
	err := r.pool.QueryRow(ctx, query, p.ID, p.Title, p.Slug, p.Content, p.Status, p.PublishedAt).Scan(&p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return appErrors.ErrNotFound
		}
		return err
	}
	return nil
}

func (r *NewsRepository) Delete(ctx context.Context, id int64) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM news_posts WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return appErrors.ErrNotFound
	}
	return nil
}
