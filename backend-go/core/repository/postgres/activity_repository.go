package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ActivityRepository struct {
	pool *pgxpool.Pool
}

func NewActivityRepository(pool *pgxpool.Pool) *ActivityRepository {
	return &ActivityRepository{pool: pool}
}

func (r *ActivityRepository) Log(ctx context.Context, a *domain.Activity) error {
	query := `
		INSERT INTO activity_logs (client_id, admin_id, type, event, message, ip_address, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP)
		RETURNING id, created_at
	`
	return r.pool.QueryRow(ctx, query, a.ClientID, a.AdminID, a.Type, a.Event, a.Message, a.IPAddress).Scan(&a.ID, &a.CreatedAt)
}

func (r *ActivityRepository) List(ctx context.Context, limit, offset int) ([]*domain.Activity, int, error) {
	var total int
	_ = r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM activity_logs`).Scan(&total)

	query := `
		SELECT id, client_id, admin_id, type, event, message, ip_address, created_at
		FROM activity_logs
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []*domain.Activity
	for rows.Next() {
		a := &domain.Activity{}
		if err := rows.Scan(&a.ID, &a.ClientID, &a.AdminID, &a.Type, &a.Event, &a.Message, &a.IPAddress, &a.CreatedAt); err != nil {
			return nil, 0, err
		}
		logs = append(logs, a)
	}

	return logs, total, nil
}

func (r *ActivityRepository) ListByClientID(ctx context.Context, clientID int64, limit, offset int) ([]*domain.Activity, int, error) {
	var total int
	_ = r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM activity_logs WHERE client_id = $1`, clientID).Scan(&total)

	query := `
		SELECT id, client_id, admin_id, type, event, message, ip_address, created_at
		FROM activity_logs
		WHERE client_id = $1
		ORDER BY id DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.pool.Query(ctx, query, clientID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []*domain.Activity
	for rows.Next() {
		a := &domain.Activity{}
		if err := rows.Scan(&a.ID, &a.ClientID, &a.AdminID, &a.Type, &a.Event, &a.Message, &a.IPAddress, &a.CreatedAt); err != nil {
			return nil, 0, err
		}
		logs = append(logs, a)
	}

	return logs, total, nil
}

func (r *ActivityRepository) DeleteOld(ctx context.Context, days int) error {
	query := fmt.Sprintf(`DELETE FROM activity_logs WHERE created_at < CURRENT_TIMESTAMP - INTERVAL '%d days'`, days)
	_, err := r.pool.Exec(ctx, query)
	return err
}
