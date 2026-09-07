package postgres

import (
	"context"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationRepository struct {
	pool *pgxpool.Pool
}

func NewNotificationRepository(pool *pgxpool.Pool) *NotificationRepository {
	return &NotificationRepository{pool: pool}
}

func (r *NotificationRepository) Create(ctx context.Context, n *domain.Notification) error {
	query := `INSERT INTO notifications (client_id, title, message, type, is_read, created_at) VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP) RETURNING id, created_at`
	return r.pool.QueryRow(ctx, query, n.ClientID, n.Title, n.Message, n.Type, n.IsRead).Scan(&n.ID, &n.CreatedAt)
}

func (r *NotificationRepository) ListByClientID(ctx context.Context, clientID int64, limit, offset int) ([]*domain.Notification, int, error) {
	var total int
	_ = r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM notifications WHERE client_id = $1`, clientID).Scan(&total)

	query := `SELECT id, client_id, title, message, type, is_read, created_at FROM notifications WHERE client_id = $1 ORDER BY id DESC LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, query, clientID, limit, offset)
	if err != nil { return nil, 0, err }
	defer rows.Close()

	var list []*domain.Notification
	for rows.Next() {
		n := &domain.Notification{}
		if err := rows.Scan(&n.ID, &n.ClientID, &n.Title, &n.Message, &n.Type, &n.IsRead, &n.CreatedAt); err != nil { return nil, 0, err }
		list = append(list, n)
	}
	return list, total, nil
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, id int64) error {
	_, err := r.pool.Exec(ctx, `UPDATE notifications SET is_read = true WHERE id = $1`, id)
	return err
}

func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, clientID int64) error {
	_, err := r.pool.Exec(ctx, `UPDATE notifications SET is_read = true WHERE client_id = $1`, clientID)
	return err
}

func (r *NotificationRepository) DeleteOld(ctx context.Context, days int) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM notifications WHERE created_at < CURRENT_TIMESTAMP - (INTERVAL '1 day' * $1)`, days)
	return err
}
