package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RedirectRepository struct {
	pool *pgxpool.Pool
}

func NewRedirectRepository(pool *pgxpool.Pool) *RedirectRepository {
	return &RedirectRepository{pool: pool}
}

func (r *RedirectRepository) List(ctx context.Context, limit, offset int) ([]*domain.Redirect, int, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	var total int
	countQuery := `SELECT COUNT(*) FROM redirects`
	if err := r.pool.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count redirects: %w", err)
	}

	query := `
		SELECT id, path, target, status_code, is_enabled, hit_count, created_at, updated_at
		FROM redirects
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list redirects: %w", err)
	}
	defer rows.Close()

	var list []*domain.Redirect
	for rows.Next() {
		item := &domain.Redirect{}
		if err := rows.Scan(
			&item.ID,
			&item.Path,
			&item.Target,
			&item.StatusCode,
			&item.IsEnabled,
			&item.HitCount,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan redirect row: %w", err)
		}
		list = append(list, item)
	}

	return list, total, nil
}

func (r *RedirectRepository) GetByID(ctx context.Context, id int64) (*domain.Redirect, error) {
	query := `
		SELECT id, path, target, status_code, is_enabled, hit_count, created_at, updated_at
		FROM redirects
		WHERE id = $1
	`
	item := &domain.Redirect{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&item.ID,
		&item.Path,
		&item.Target,
		&item.StatusCode,
		&item.IsEnabled,
		&item.HitCount,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get redirect by id: %w", err)
	}
	return item, nil
}

func (r *RedirectRepository) GetByPath(ctx context.Context, path string) (*domain.Redirect, error) {
	query := `
		SELECT id, path, target, status_code, is_enabled, hit_count, created_at, updated_at
		FROM redirects
		WHERE path = $1
	`
	item := &domain.Redirect{}
	err := r.pool.QueryRow(ctx, query, path).Scan(
		&item.ID,
		&item.Path,
		&item.Target,
		&item.StatusCode,
		&item.IsEnabled,
		&item.HitCount,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get redirect by path: %w", err)
	}
	return item, nil
}

func (r *RedirectRepository) Create(ctx context.Context, redirect *domain.Redirect) error {
	query := `
		INSERT INTO redirects (path, target, status_code, is_enabled, hit_count, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	now := time.Now().UTC()
	redirect.CreatedAt = now
	redirect.UpdatedAt = now

	if redirect.StatusCode == 0 {
		redirect.StatusCode = 301
	}

	err := r.pool.QueryRow(
		ctx,
		query,
		redirect.Path,
		redirect.Target,
		redirect.StatusCode,
		redirect.IsEnabled,
		redirect.HitCount,
		redirect.CreatedAt,
		redirect.UpdatedAt,
	).Scan(&redirect.ID)
	if err != nil {
		return fmt.Errorf("failed to create redirect: %w", err)
	}
	return nil
}

func (r *RedirectRepository) Update(ctx context.Context, redirect *domain.Redirect) error {
	query := `
		UPDATE redirects
		SET path = $1, target = $2, status_code = $3, is_enabled = $4, updated_at = $5
		WHERE id = $6
	`
	redirect.UpdatedAt = time.Now().UTC()
	cmd, err := r.pool.Exec(
		ctx,
		query,
		redirect.Path,
		redirect.Target,
		redirect.StatusCode,
		redirect.IsEnabled,
		redirect.UpdatedAt,
		redirect.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update redirect: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return appErrors.ErrNotFound
	}
	return nil
}

func (r *RedirectRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM redirects WHERE id = $1`
	cmd, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete redirect: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return appErrors.ErrNotFound
	}
	return nil
}

func (r *RedirectRepository) IncrementHitCount(ctx context.Context, id int64) error {
	query := `UPDATE redirects SET hit_count = hit_count + 1 WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}
