package postgres

import (
	"context"
	"fmt"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const actCols = `id, client_id, admin_id, type, event, message, ip_address, created_at`

type ActivityRepository struct{ pool *pgxpool.Pool }

func NewActivityRepository(p *pgxpool.Pool) *ActivityRepository { return &ActivityRepository{p} }

func scanAct(r pgx.Row) (*domain.Activity, error) {
	var a domain.Activity; err := r.Scan(&a.ID, &a.ClientID, &a.AdminID, &a.Type, &a.Event, &a.Message, &a.IPAddress, &a.CreatedAt); return &a, err
}

func (r *ActivityRepository) Log(ctx context.Context, a *domain.Activity) error {
	return r.pool.QueryRow(ctx, `INSERT INTO activity_logs (client_id, admin_id, type, event, message, ip_address, created_at) VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP) RETURNING id, created_at`, a.ClientID, a.AdminID, a.Type, a.Event, a.Message, a.IPAddress).Scan(&a.ID, &a.CreatedAt)
}

func (r *ActivityRepository) List(ctx context.Context, l, o int) ([]*domain.Activity, int, error) {
	res, err := list(ctx, r.pool, "SELECT "+actCols+" FROM activity_logs ORDER BY id DESC LIMIT $1 OFFSET $2", scanAct, l, o)
	return res, total(ctx, r.pool, "SELECT COUNT(*) FROM activity_logs"), err
}

func (r *ActivityRepository) ListByClientID(ctx context.Context, cid int64, l, o int) ([]*domain.Activity, int, error) {
	res, err := list(ctx, r.pool, "SELECT "+actCols+" FROM activity_logs WHERE client_id = $1 ORDER BY id DESC LIMIT $2 OFFSET $3", scanAct, cid, l, o)
	return res, total(ctx, r.pool, "SELECT COUNT(*) FROM activity_logs WHERE client_id = $1", cid), err
}

func (r *ActivityRepository) DeleteOld(ctx context.Context, days int) error {
	_, err := r.pool.Exec(ctx, fmt.Sprintf(`DELETE FROM activity_logs WHERE created_at < CURRENT_TIMESTAMP - INTERVAL '%d days'`, days)); return err
}

func (r *ActivityRepository) GetTrend(ctx context.Context, days int) (map[string]int, error) {
	q := `SELECT TO_CHAR(created_at, 'YYYY-MM-DD') as day, COUNT(*) as count
	      FROM activity_logs
	      WHERE created_at > CURRENT_DATE - INTERVAL '1 day' * $1
	      GROUP BY day ORDER BY day ASC`
	rows, err := r.pool.Query(ctx, q, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make(map[string]int)
	for rows.Next() {
		var day string
		var count int
		if err := rows.Scan(&day, &count); err != nil {
			return nil, err
		}
		res[day] = count
	}
	return res, nil
}
