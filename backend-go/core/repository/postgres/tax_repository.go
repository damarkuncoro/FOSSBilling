package postgres

import (
	"context"
	"errors"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const txCols = `id, name, country, state, rate, is_active, tax_exempt, created_at, updated_at`

type TaxRepository struct{ pool *pgxpool.Pool }

func NewTaxRepository(p *pgxpool.Pool) *TaxRepository { return &TaxRepository{p} }

func scanTax(r pgx.Row) (*domain.TaxRule, error) {
	var t domain.TaxRule; err := r.Scan(&t.ID, &t.Name, &t.Country, &t.State, &t.Rate, &t.IsActive, &t.TaxExempt, &t.CreatedAt, &t.UpdatedAt); return &t, err
}

func (r *TaxRepository) GetByID(ctx context.Context, id int64) (*domain.TaxRule, error) {
	tr, err := scanTax(r.pool.QueryRow(ctx, "SELECT "+txCols+" FROM tax_rules WHERE id = $1", id))
	if err != nil && errors.Is(err, pgx.ErrNoRows) { return nil, appErrors.ErrNotFound }; return tr, err
}

func (r *TaxRepository) List(ctx context.Context) ([]*domain.TaxRule, error) {
	return list(ctx, r.pool, "SELECT "+txCols+" FROM tax_rules ORDER BY id ASC", scanTax)
}

func (r *TaxRepository) GetByLocation(ctx context.Context, c, s string) (*domain.TaxRule, error) {
	q := "SELECT "+txCols+" FROM tax_rules WHERE is_active = true AND ((country = $1 AND state = $2) OR (country = $1 AND state IS NULL) OR (country IS NULL)) ORDER BY country ASC NULLS LAST, state ASC NULLS LAST LIMIT 1"
	tr, err := scanTax(r.pool.QueryRow(ctx, q, c, s))
	if err != nil && errors.Is(err, pgx.ErrNoRows) { return nil, nil }; return tr, err
}

func (r *TaxRepository) Create(ctx context.Context, t *domain.TaxRule) error {
	return r.pool.QueryRow(ctx, `INSERT INTO tax_rules (name, country, state, rate, is_active, tax_exempt, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id, created_at, updated_at`, t.Name, t.Country, t.State, t.Rate, t.IsActive, t.TaxExempt).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *TaxRepository) Update(ctx context.Context, t *domain.TaxRule) error {
	_, err := r.pool.Exec(ctx, `UPDATE tax_rules SET name = $1, country = $2, state = $3, rate = $4, is_active = $5, tax_exempt = $6, updated_at = CURRENT_TIMESTAMP WHERE id = $7`, t.Name, t.Country, t.State, t.Rate, t.IsActive, t.TaxExempt, t.ID); return err
}

func (r *TaxRepository) Delete(ctx context.Context, id int64) error { _, err := r.pool.Exec(ctx, `DELETE FROM tax_rules WHERE id = $1`, id); return err }
