package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const rdCols = `id, path, target, status_code, is_enabled, hit_count, created_at, updated_at`

type RedirectRepository struct{ pool *pgxpool.Pool }

func NewRedirectRepository(p *pgxpool.Pool) *RedirectRepository { return &RedirectRepository{p} }

func scanRd(r pgx.Row) (*domain.Redirect, error) {
	var item domain.Redirect; err := r.Scan(&item.ID, &item.Path, &item.Target, &item.StatusCode, &item.IsEnabled, &item.HitCount, &item.CreatedAt, &item.UpdatedAt); return &item, err
}

func (r *RedirectRepository) List(ctx context.Context, l, o int) ([]*domain.Redirect, int, error) {
	res, err := list(ctx, r.pool, "SELECT "+rdCols+" FROM redirects ORDER BY id DESC LIMIT $1 OFFSET $2", scanRd, l, o)
	return res, total(ctx, r.pool, "SELECT COUNT(*) FROM redirects"), err
}

func (r *RedirectRepository) GetByID(ctx context.Context, id int64) (*domain.Redirect, error) {
	rd, err := scanRd(r.pool.QueryRow(ctx, "SELECT "+rdCols+" FROM redirects WHERE id = $1", id))
	if err != nil && errors.Is(err, pgx.ErrNoRows) { return nil, appErrors.ErrNotFound }; return rd, err
}

func (r *RedirectRepository) GetByPath(ctx context.Context, p string) (*domain.Redirect, error) {
	rd, err := scanRd(r.pool.QueryRow(ctx, "SELECT "+rdCols+" FROM redirects WHERE path = $1", p))
	if err != nil && errors.Is(err, pgx.ErrNoRows) { return nil, appErrors.ErrNotFound }; return rd, err
}

func (r *RedirectRepository) Create(ctx context.Context, rd *domain.Redirect) error {
	rd.CreatedAt, rd.UpdatedAt = time.Now().UTC(), time.Now().UTC(); if rd.StatusCode == 0 { rd.StatusCode = 301 }
	return r.pool.QueryRow(ctx, `INSERT INTO redirects (path, target, status_code, is_enabled, hit_count, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`, rd.Path, rd.Target, rd.StatusCode, rd.IsEnabled, rd.HitCount, rd.CreatedAt, rd.UpdatedAt).Scan(&rd.ID)
}

func (r *RedirectRepository) Update(ctx context.Context, rd *domain.Redirect) error {
	rd.UpdatedAt = time.Now().UTC()
	t, err := r.pool.Exec(ctx, `UPDATE redirects SET path = $1, target = $2, status_code = $3, is_enabled = $4, updated_at = $5 WHERE id = $6`, rd.Path, rd.Target, rd.StatusCode, rd.IsEnabled, rd.UpdatedAt, rd.ID)
	if err == nil && t.RowsAffected() == 0 { return appErrors.ErrNotFound }; return err
}

func (r *RedirectRepository) Delete(ctx context.Context, id int64) error {
	t, err := r.pool.Exec(ctx, `DELETE FROM redirects WHERE id = $1`, id)
	if err == nil && t.RowsAffected() == 0 { return appErrors.ErrNotFound }; return err
}

func (r *RedirectRepository) IncrementHitCount(ctx context.Context, id int64) error { _, err := r.pool.Exec(ctx, `UPDATE redirects SET hit_count = hit_count + 1 WHERE id = $1`, id); return err }
