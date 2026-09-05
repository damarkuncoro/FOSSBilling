package postgres

import (
	"context"
	"errors"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CurrencyRepository struct {
	pool *pgxpool.Pool
}

func NewCurrencyRepository(pool *pgxpool.Pool) *CurrencyRepository {
	return &CurrencyRepository{pool: pool}
}

func (r *CurrencyRepository) GetByCode(ctx context.Context, code string) (*domain.Currency, error) {
	query := `
		SELECT id, code, title, conversion_rate, format, price_format, is_default, created_at, updated_at
		FROM currencies
		WHERE code = $1
	`
	var c domain.Currency
	err := r.pool.QueryRow(ctx, query, code).Scan(
		&c.ID, &c.Code, &c.Title, &c.ConversionRate, &c.Format, &c.PriceFormat, &c.IsDefault, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *CurrencyRepository) GetDefault(ctx context.Context) (*domain.Currency, error) {
	query := `
		SELECT id, code, title, conversion_rate, format, price_format, is_default, created_at, updated_at
		FROM currencies
		WHERE is_default = true
		LIMIT 1
	`
	var c domain.Currency
	err := r.pool.QueryRow(ctx, query).Scan(
		&c.ID, &c.Code, &c.Title, &c.ConversionRate, &c.Format, &c.PriceFormat, &c.IsDefault, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *CurrencyRepository) List(ctx context.Context) ([]*domain.Currency, error) {
	query := `
		SELECT id, code, title, conversion_rate, format, price_format, is_default, created_at, updated_at
		FROM currencies
		ORDER BY is_default DESC, code ASC
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.Currency
	for rows.Next() {
		var c domain.Currency
		if err := rows.Scan(&c.ID, &c.Code, &c.Title, &c.ConversionRate, &c.Format, &c.PriceFormat, &c.IsDefault, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, &c)
	}
	return list, rows.Err()
}

func (r *CurrencyRepository) Create(ctx context.Context, c *domain.Currency) error {
	query := `
		INSERT INTO currencies (code, title, conversion_rate, format, price_format, is_default, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at
	`
	return r.pool.QueryRow(ctx, query, c.Code, c.Title, c.ConversionRate, c.Format, c.PriceFormat, c.IsDefault).
		Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *CurrencyRepository) Update(ctx context.Context, c *domain.Currency) error {
	query := `
		UPDATE currencies
		SET title = $2, conversion_rate = $3, format = $4, price_format = $5, is_default = $6, updated_at = CURRENT_TIMESTAMP
		WHERE code = $1
	`
	cmd, err := r.pool.Exec(ctx, query, c.Code, c.Title, c.ConversionRate, c.Format, c.PriceFormat, c.IsDefault)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return appErrors.ErrNotFound
	}
	return nil
}

func (r *CurrencyRepository) SetDefault(ctx context.Context, code string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `UPDATE currencies SET is_default = false WHERE is_default = true`); err != nil {
		return err
	}
	cmd, err := tx.Exec(ctx, `UPDATE currencies SET is_default = true, conversion_rate = 1.0, updated_at = CURRENT_TIMESTAMP WHERE code = $1`, code)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return appErrors.ErrNotFound
	}
	return tx.Commit(ctx)
}

func (r *CurrencyRepository) Delete(ctx context.Context, code string) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM currencies WHERE code = $1`, code)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return appErrors.ErrNotFound
	}
	return nil
}
