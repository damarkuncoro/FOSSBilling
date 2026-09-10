package postgres

import (
	"context"
	"errors"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const akCols = `id, client_id, name, key, secret, expires_at, created_at, updated_at`

type APIKeyRepository struct{ pool *pgxpool.Pool }

func NewAPIKeyRepository(p *pgxpool.Pool) *APIKeyRepository { return &APIKeyRepository{p} }

func scanAK(r pgx.Row) (*domain.APIKey, error) {
	var k domain.APIKey; err := r.Scan(&k.ID, &k.ClientID, &k.Name, &k.Key, &k.Secret, &k.ExpiresAt, &k.CreatedAt, &k.UpdatedAt); return &k, err
}

func (r *APIKeyRepository) GetByID(ctx context.Context, id int64) (*domain.APIKey, error) {
	ak, err := scanAK(r.pool.QueryRow(ctx, "SELECT "+akCols+" FROM api_keys WHERE id = $1", id))
	if err != nil && errors.Is(err, pgx.ErrNoRows) { return nil, appErrors.ErrNotFound }; return ak, err
}

func (r *APIKeyRepository) GetByKey(ctx context.Context, k string) (*domain.APIKey, error) {
	ak, err := scanAK(r.pool.QueryRow(ctx, "SELECT "+akCols+" FROM api_keys WHERE key = $1", k))
	if err != nil && errors.Is(err, pgx.ErrNoRows) { return nil, appErrors.ErrNotFound }; return ak, err
}

func (r *APIKeyRepository) ListByClientID(ctx context.Context, cid int64) ([]*domain.APIKey, error) {
	return list(ctx, r.pool, "SELECT "+akCols+" FROM api_keys WHERE client_id = $1 ORDER BY created_at DESC", scanAK, cid)
}

func (r *APIKeyRepository) Create(ctx context.Context, k *domain.APIKey) error {
	return r.pool.QueryRow(ctx, `INSERT INTO api_keys (client_id, name, key, secret, expires_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id, created_at, updated_at`, k.ClientID, k.Name, k.Key, k.Secret, k.ExpiresAt).Scan(&k.ID, &k.CreatedAt, &k.UpdatedAt)
}

func (r *APIKeyRepository) Delete(ctx context.Context, id, cid int64) error {
	t, err := r.pool.Exec(ctx, `DELETE FROM api_keys WHERE id = $1 AND client_id = $2`, id, cid)
	if err == nil && t.RowsAffected() == 0 { return appErrors.ErrNotFound }; return err
}
