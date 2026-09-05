package postgres

import (
	"context"
	"errors"

	"github.com/fossbilling/backend-go/internal/domain"
	appErrors "github.com/fossbilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PromoRepository struct {
	pool *pgxpool.Pool
}

func NewPromoRepository(pool *pgxpool.Pool) *PromoRepository {
	return &PromoRepository{pool: pool}
}

func (r *PromoRepository) GetByCode(ctx context.Context, code string) (*domain.Promo, error) {
	query := `
		SELECT id, code, description, type, value, max_uses, used_count, once_per_client, start_date, end_date, active, created_at, updated_at
		FROM promos
		WHERE LOWER(code) = LOWER($1)
	`
	var p domain.Promo
	err := r.pool.QueryRow(ctx, query, code).Scan(
		&p.ID, &p.Code, &p.Description, &p.Type, &p.Value, &p.MaxUses, &p.UsedCount,
		&p.OncePerClient, &p.StartDate, &p.EndDate, &p.Active, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *PromoRepository) GetRedemptionCount(ctx context.Context, promoID int64, clientID int64) (int, error) {
	query := `SELECT COUNT(*) FROM promo_redemptions WHERE promo_id = $1 AND client_id = $2`
	var count int
	if err := r.pool.QueryRow(ctx, query, promoID, clientID).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *PromoRepository) IncrementUsed(ctx context.Context, promoID int64, clientID int64, orderID *int64) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	updateQuery := `UPDATE promos SET used_count = used_count + 1, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	if _, err := tx.Exec(ctx, updateQuery, promoID); err != nil {
		return err
	}

	insertRedemption := `
		INSERT INTO promo_redemptions (promo_id, client_id, order_id, created_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
	`
	if _, err := tx.Exec(ctx, insertRedemption, promoID, clientID, orderID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *PromoRepository) Create(ctx context.Context, p *domain.Promo) error {
	query := `
		INSERT INTO promos (
			code, description, type, value, max_uses, used_count, once_per_client, start_date, end_date, active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at
	`
	return r.pool.QueryRow(ctx, query,
		p.Code, p.Description, p.Type, p.Value, p.MaxUses, p.UsedCount,
		p.OncePerClient, p.StartDate, p.EndDate, p.Active,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}
