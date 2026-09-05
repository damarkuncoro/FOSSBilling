package postgres

import (
	"context"
	"errors"

	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type APIKeyRepository struct {
	pool *pgxpool.Pool
}

func NewAPIKeyRepository(pool *pgxpool.Pool) *APIKeyRepository {
	return &APIKeyRepository{pool: pool}
}

func (r *APIKeyRepository) GetByID(ctx context.Context, id int64) (*domain.APIKey, error) {
	query := `
		SELECT id, client_id, name, key, secret, expires_at, created_at, updated_at
		FROM api_keys
		WHERE id = $1
	`
	var k domain.APIKey
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&k.ID, &k.ClientID, &k.Name, &k.Key, &k.Secret, &k.ExpiresAt, &k.CreatedAt, &k.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, err
	}
	return &k, nil
}

func (r *APIKeyRepository) GetByKey(ctx context.Context, key string) (*domain.APIKey, error) {
	query := `
		SELECT id, client_id, name, key, secret, expires_at, created_at, updated_at
		FROM api_keys
		WHERE key = $1
	`
	var k domain.APIKey
	err := r.pool.QueryRow(ctx, query, key).Scan(
		&k.ID, &k.ClientID, &k.Name, &k.Key, &k.Secret, &k.ExpiresAt, &k.CreatedAt, &k.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, err
	}
	return &k, nil
}

func (r *APIKeyRepository) ListByClientID(ctx context.Context, clientID int64) ([]*domain.APIKey, error) {
	query := `
		SELECT id, client_id, name, key, secret, expires_at, created_at, updated_at
		FROM api_keys
		WHERE client_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.APIKey
	for rows.Next() {
		var k domain.APIKey
		if err := rows.Scan(&k.ID, &k.ClientID, &k.Name, &k.Key, &k.Secret, &k.ExpiresAt, &k.CreatedAt, &k.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, &k)
	}
	return list, rows.Err()
}

func (r *APIKeyRepository) Create(ctx context.Context, k *domain.APIKey) error {
	query := `
		INSERT INTO api_keys (client_id, name, key, secret, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at
	`
	return r.pool.QueryRow(ctx, query,
		k.ClientID, k.Name, k.Key, k.Secret, k.ExpiresAt,
	).Scan(&k.ID, &k.CreatedAt, &k.UpdatedAt)
}

func (r *APIKeyRepository) Delete(ctx context.Context, id, clientID int64) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM api_keys WHERE id = $1 AND client_id = $2`, id, clientID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return appErrors.ErrNotFound
	}
	return nil
}
