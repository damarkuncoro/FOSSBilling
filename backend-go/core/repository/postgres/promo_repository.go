package postgres

import (
	"context"
	"errors"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const promoCols = `id, code, description, type, value, max_uses, used_count, once_per_client, start_date, end_date, active, created_at, updated_at`

type PromoRepository struct{ pool *pgxpool.Pool }

func NewPromoRepository(p *pgxpool.Pool) *PromoRepository { return &PromoRepository{p} }

func scanPromo(r pgx.Row) (*domain.Promo, error) {
	var p domain.Promo
	err := r.Scan(&p.ID, &p.Code, &p.Description, &p.Type, &p.Value, &p.MaxUses, &p.UsedCount, &p.OncePerClient, &p.StartDate, &p.EndDate, &p.Active, &p.CreatedAt, &p.UpdatedAt)
	return &p, err
}

func (r *PromoRepository) GetByID(ctx context.Context, id int64) (*domain.Promo, error) {
	p, err := scanPromo(r.pool.QueryRow(ctx, "SELECT "+promoCols+" FROM promos WHERE id = $1", id))
	if err != nil && errors.Is(err, pgx.ErrNoRows) { return nil, appErrors.ErrNotFound }; return p, err
}

func (r *PromoRepository) GetByCode(ctx context.Context, c string) (*domain.Promo, error) {
	p, err := scanPromo(r.pool.QueryRow(ctx, "SELECT "+promoCols+" FROM promos WHERE LOWER(code) = LOWER($1)", c))
	if err != nil && errors.Is(err, pgx.ErrNoRows) { return nil, appErrors.ErrNotFound }; return p, err
}

func (r *PromoRepository) GetRedemptionCount(ctx context.Context, pID, cID int64) (int, error) {
	var cnt int; err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM promo_redemptions WHERE promo_id = $1 AND client_id = $2`, pID, cID).Scan(&cnt); return cnt, err
}

func (r *PromoRepository) IncrementUsed(ctx context.Context, pID, cID int64, oID *int64) error {
	tx, err := r.pool.Begin(ctx); if err != nil { return err }; defer tx.Rollback(ctx)
	var p domain.Promo
	if err := tx.QueryRow(ctx, `SELECT id, max_uses, used_count, once_per_client, active FROM promos WHERE id = $1 FOR UPDATE`, pID).Scan(&p.ID, &p.MaxUses, &p.UsedCount, &p.OncePerClient, &p.Active); err != nil { return err }
	if !p.Active || (p.MaxUses > 0 && p.UsedCount >= p.MaxUses) { return errors.New("promo limit reached") }
	if p.OncePerClient {
		var cnt int; _ = tx.QueryRow(ctx, `SELECT COUNT(*) FROM promo_redemptions WHERE promo_id = $1 AND client_id = $2`, pID, cID).Scan(&cnt)
		if cnt > 0 { return errors.New("promo already used") }
	}
	_, _ = tx.Exec(ctx, `UPDATE promos SET used_count = used_count + 1, updated_at = CURRENT_TIMESTAMP WHERE id = $1`, pID)
	_, err = tx.Exec(ctx, `INSERT INTO promo_redemptions (promo_id, client_id, order_id, created_at) VALUES ($1, $2, $3, CURRENT_TIMESTAMP)`, pID, cID, oID)
	if err != nil { return err }; return tx.Commit(ctx)
}

func (r *PromoRepository) List(ctx context.Context, l, o int) ([]*domain.Promo, int, error) {
	res, err := list(ctx, r.pool, "SELECT "+promoCols+" FROM promos ORDER BY id DESC LIMIT $1 OFFSET $2", scanPromo, l, o)
	return res, total(ctx, r.pool, "SELECT COUNT(*) FROM promos"), err
}

func (r *PromoRepository) Create(ctx context.Context, p *domain.Promo) error {
	return r.pool.QueryRow(ctx, `INSERT INTO promos (code, description, type, value, max_uses, used_count, once_per_client, start_date, end_date, active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id, created_at, updated_at`, p.Code, p.Description, p.Type, p.Value, p.MaxUses, p.UsedCount, p.OncePerClient, p.StartDate, p.EndDate, p.Active).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *PromoRepository) Update(ctx context.Context, p *domain.Promo) error {
	_, err := r.pool.Exec(ctx, `UPDATE promos SET code = $1, description = $2, type = $3, value = $4, max_uses = $5, used_count = $6, once_per_client = $7, start_date = $8, end_date = $9, active = $10, updated_at = CURRENT_TIMESTAMP WHERE id = $11`, p.Code, p.Description, p.Type, p.Value, p.MaxUses, p.UsedCount, p.OncePerClient, p.StartDate, p.EndDate, p.Active, p.ID)
	return err
}

func (r *PromoRepository) Delete(ctx context.Context, id int64) error { _, err := r.pool.Exec(ctx, `DELETE FROM promos WHERE id = $1`, id); return err }
