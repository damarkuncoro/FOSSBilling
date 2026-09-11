package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type PostgresAffiliateRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresAffiliateRepository(pool *pgxpool.Pool) *PostgresAffiliateRepository {
	return &PostgresAffiliateRepository{pool: pool}
}

func (r *PostgresAffiliateRepository) GetByClientID(ctx context.Context, clientID int64) (*domain.Affiliate, error) {
	var a domain.Affiliate
	err := r.pool.QueryRow(ctx,
		"SELECT id, client_id, commission_rate, status, balance, total_earned, created_at, updated_at FROM affiliates WHERE client_id = $1",
		clientID).Scan(&a.ID, &a.ClientID, &a.CommissionRate, &a.Status, &a.Balance, &a.TotalEarned, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *PostgresAffiliateRepository) Create(ctx context.Context, aff *domain.Affiliate) error {
	aff.CreatedAt = time.Now().UTC()
	aff.UpdatedAt = aff.CreatedAt
	return r.pool.QueryRow(ctx,
		"INSERT INTO affiliates (client_id, commission_rate, status, balance, total_earned, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		aff.ClientID, aff.CommissionRate, aff.Status, aff.Balance, aff.TotalEarned, aff.CreatedAt, aff.UpdatedAt).Scan(&aff.ID)
}

func (r *PostgresAffiliateRepository) Update(ctx context.Context, aff *domain.Affiliate) error {
	aff.UpdatedAt = time.Now().UTC()
	_, err := r.pool.Exec(ctx,
		"UPDATE affiliates SET commission_rate = $1, status = $2, balance = $3, total_earned = $4, updated_at = $5 WHERE id = $6",
		aff.CommissionRate, aff.Status, aff.Balance, aff.TotalEarned, aff.UpdatedAt, aff.ID)
	return err
}

func (r *PostgresAffiliateRepository) AddReferral(ctx context.Context, ref *domain.AffiliateReferral) error {
	ref.CreatedAt = time.Now().UTC()
	return r.pool.QueryRow(ctx,
		"INSERT INTO affiliate_referrals (affiliate_id, client_id, order_id, amount, status, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		ref.AffiliateID, ref.ClientID, ref.OrderID, ref.Amount, ref.Status, ref.CreatedAt).Scan(&ref.ID)
}

func (r *PostgresAffiliateRepository) ListReferrals(ctx context.Context, affiliateID int64) ([]*domain.AffiliateReferral, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, affiliate_id, client_id, order_id, amount, status, created_at FROM affiliate_referrals WHERE affiliate_id = $1 ORDER BY created_at DESC", affiliateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []*domain.AffiliateReferral
	for rows.Next() {
		var ref domain.AffiliateReferral
		if err := rows.Scan(&ref.ID, &ref.AffiliateID, &ref.ClientID, &ref.OrderID, &ref.Amount, &ref.Status, &ref.CreatedAt); err != nil {
			return nil, err
		}
		refs = append(refs, &ref)
	}
	return refs, nil
}
