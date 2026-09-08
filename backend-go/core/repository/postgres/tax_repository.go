package postgres

import (
	"context"
	"errors"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaxRepository struct {
	pool *pgxpool.Pool
}

func NewTaxRepository(pool *pgxpool.Pool) *TaxRepository {
	return &TaxRepository{pool: pool}
}

func (r *TaxRepository) GetByID(ctx context.Context, id int64) (*domain.TaxRule, error) {
	query := `SELECT id, name, country, state, rate, is_active, tax_exempt, created_at, updated_at FROM tax_rules WHERE id = $1`
	var rule domain.TaxRule
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&rule.ID, &rule.Name, &rule.Country, &rule.State, &rule.Rate, &rule.IsActive, &rule.TaxExempt, &rule.CreatedAt, &rule.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, err
	}
	return &rule, nil
}

func (r *TaxRepository) List(ctx context.Context) ([]*domain.TaxRule, error) {
	query := `SELECT id, name, country, state, rate, is_active, tax_exempt, created_at, updated_at FROM tax_rules ORDER BY id ASC`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []*domain.TaxRule
	for rows.Next() {
		var rule domain.TaxRule
		if err := rows.Scan(&rule.ID, &rule.Name, &rule.Country, &rule.State, &rule.Rate, &rule.IsActive, &rule.TaxExempt, &rule.CreatedAt, &rule.UpdatedAt); err != nil {
			return nil, err
		}
		rules = append(rules, &rule)
	}
	return rules, nil
}

func (r *TaxRepository) GetByLocation(ctx context.Context, country, state string) (*domain.TaxRule, error) {
	// Try state specific first, then country, then global
	query := `
		SELECT id, name, country, state, rate, is_active, tax_exempt, created_at, updated_at
		FROM tax_rules
		WHERE is_active = true AND (
			(country = $1 AND state = $2) OR
			(country = $1 AND state IS NULL) OR
			(country IS NULL)
		)
		ORDER BY country ASC NULLS LAST, state ASC NULLS LAST
		LIMIT 1
	`
	var rule domain.TaxRule
	err := r.pool.QueryRow(ctx, query, country, state).Scan(
		&rule.ID, &rule.Name, &rule.Country, &rule.State, &rule.Rate, &rule.IsActive, &rule.TaxExempt, &rule.CreatedAt, &rule.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // No tax rule applies
		}
		return nil, err
	}
	return &rule, nil
}

func (r *TaxRepository) Create(ctx context.Context, rule *domain.TaxRule) error {
	query := `INSERT INTO tax_rules (name, country, state, rate, is_active, tax_exempt, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id, created_at, updated_at`
	return r.pool.QueryRow(ctx, query, rule.Name, rule.Country, rule.State, rule.Rate, rule.IsActive, rule.TaxExempt).Scan(&rule.ID, &rule.CreatedAt, &rule.UpdatedAt)
}

func (r *TaxRepository) Update(ctx context.Context, rule *domain.TaxRule) error {
	query := `UPDATE tax_rules SET name = $1, country = $2, state = $3, rate = $4, is_active = $5, tax_exempt = $6, updated_at = CURRENT_TIMESTAMP WHERE id = $7`
	_, err := r.pool.Exec(ctx, query, rule.Name, rule.Country, rule.State, rule.Rate, rule.IsActive, rule.TaxExempt, rule.ID)
	return err
}

func (r *TaxRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM tax_rules WHERE id = $1`, id)
	return err
}
