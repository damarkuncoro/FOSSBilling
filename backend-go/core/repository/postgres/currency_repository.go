package postgres

import (
	"context"
	"errors"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const curCols = `id, code, title, conversion_rate, format, price_format, is_default, created_at, updated_at`

type CurrencyRepository struct{ pool *pgxpool.Pool }

func NewCurrencyRepository(p *pgxpool.Pool) *CurrencyRepository { return &CurrencyRepository{p} }

func scanCur(r pgx.Row) (*domain.Currency, error) {
	var c domain.Currency; err := r.Scan(&c.ID, &c.Code, &c.Title, &c.ConversionRate, &c.Format, &c.PriceFormat, &c.IsDefault, &c.CreatedAt, &c.UpdatedAt); return &c, err
}

func (r *CurrencyRepository) GetByCode(ctx context.Context, c string) (*domain.Currency, error) {
	curr, err := scanCur(r.pool.QueryRow(ctx, "SELECT "+curCols+" FROM currencies WHERE code = $1", c))
	if err != nil && errors.Is(err, pgx.ErrNoRows) { return nil, appErrors.ErrNotFound }; return curr, err
}

func (r *CurrencyRepository) GetDefault(ctx context.Context) (*domain.Currency, error) {
	return scanCur(r.pool.QueryRow(ctx, "SELECT "+curCols+" FROM currencies WHERE is_default = true LIMIT 1"))
}

func (r *CurrencyRepository) List(ctx context.Context) ([]*domain.Currency, error) {
	return list(ctx, r.pool, "SELECT "+curCols+" FROM currencies ORDER BY is_default DESC, code ASC", scanCur)
}

func (r *CurrencyRepository) Create(ctx context.Context, c *domain.Currency) error {
	return r.pool.QueryRow(ctx, `INSERT INTO currencies (code, title, conversion_rate, format, price_format, is_default, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id, created_at, updated_at`, c.Code, c.Title, c.ConversionRate, c.Format, c.PriceFormat, c.IsDefault).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *CurrencyRepository) Update(ctx context.Context, c *domain.Currency) error {
	t, err := r.pool.Exec(ctx, `UPDATE currencies SET title = $2, conversion_rate = $3, format = $4, price_format = $5, is_default = $6, updated_at = CURRENT_TIMESTAMP WHERE code = $1`, c.Code, c.Title, c.ConversionRate, c.Format, c.PriceFormat, c.IsDefault)
	if err == nil && t.RowsAffected() == 0 { return appErrors.ErrNotFound }; return err
}

func (r *CurrencyRepository) SetDefault(ctx context.Context, c string) error {
	tx, err := r.pool.Begin(ctx); if err != nil { return err }; defer tx.Rollback(ctx)
	_, _ = tx.Exec(ctx, `UPDATE currencies SET is_default = false WHERE is_default = true`)
	t, err := tx.Exec(ctx, `UPDATE currencies SET is_default = true, conversion_rate = 1.0, updated_at = CURRENT_TIMESTAMP WHERE code = $1`, c)
	if err == nil && t.RowsAffected() == 0 { return appErrors.ErrNotFound }
	if err != nil { return err }; return tx.Commit(ctx)
}

func (r *CurrencyRepository) Delete(ctx context.Context, c string) error {
	t, err := r.pool.Exec(ctx, `DELETE FROM currencies WHERE code = $1`, c)
	if err == nil && t.RowsAffected() == 0 { return appErrors.ErrNotFound }; return err
}
