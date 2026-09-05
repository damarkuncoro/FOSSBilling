package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/fossbilling/backend-go/internal/domain"
	appErrors "github.com/fossbilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{pool: pool}
}

func (r *OrderRepository) GetByID(ctx context.Context, id int64) (*domain.Order, error) {
	query := `
		SELECT id, client_id, product_id, invoice_id, status, title, period, price, currency,
		       config, activated_at, expires_at, next_due_date, suspended_at, suspension_reason, created_at, updated_at
		FROM client_orders
		WHERE id = $1
	`
	var o domain.Order
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&o.ID, &o.ClientID, &o.ProductID, &o.InvoiceID, &o.Status, &o.Title, &o.Period, &o.Price, &o.Currency,
		&o.Config, &o.ActivatedAt, &o.ExpiresAt, &o.NextDueDate, &o.SuspendedAt, &o.SuspensionReason, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, err
	}
	return &o, nil
}

func (r *OrderRepository) ListByClientID(ctx context.Context, clientID int64, limit, offset int) ([]*domain.Order, int, error) {
	countQuery := `SELECT COUNT(*) FROM client_orders WHERE client_id = $1`
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, clientID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, client_id, product_id, invoice_id, status, title, period, price, currency,
		       config, activated_at, expires_at, next_due_date, suspended_at, suspension_reason, created_at, updated_at
		FROM client_orders
		WHERE client_id = $1
		ORDER BY id DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.pool.Query(ctx, query, clientID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(
			&o.ID, &o.ClientID, &o.ProductID, &o.InvoiceID, &o.Status, &o.Title, &o.Period, &o.Price, &o.Currency,
			&o.Config, &o.ActivatedAt, &o.ExpiresAt, &o.NextDueDate, &o.SuspendedAt, &o.SuspensionReason, &o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		orders = append(orders, &o)
	}

	return orders, total, nil
}

func (r *OrderRepository) List(ctx context.Context, limit, offset int) ([]*domain.Order, int, error) {
	countQuery := `SELECT COUNT(*) FROM client_orders`
	var total int
	if err := r.pool.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, client_id, product_id, invoice_id, status, title, period, price, currency,
		       config, activated_at, expires_at, next_due_date, suspended_at, suspension_reason, created_at, updated_at
		FROM client_orders
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(
			&o.ID, &o.ClientID, &o.ProductID, &o.InvoiceID, &o.Status, &o.Title, &o.Period, &o.Price, &o.Currency,
			&o.Config, &o.ActivatedAt, &o.ExpiresAt, &o.NextDueDate, &o.SuspendedAt, &o.SuspensionReason, &o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		orders = append(orders, &o)
	}

	return orders, total, nil
}

func (r *OrderRepository) ListDueOrders(ctx context.Context, dueBefore time.Time) ([]*domain.Order, error) {
	query := `
		SELECT id, client_id, product_id, invoice_id, status, title, period, price, currency,
		       config, activated_at, expires_at, next_due_date, suspended_at, suspension_reason, created_at, updated_at
		FROM client_orders
		WHERE status = 'active' AND next_due_date <= $1
	`
	rows, err := r.pool.Query(ctx, query, dueBefore)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(
			&o.ID, &o.ClientID, &o.ProductID, &o.InvoiceID, &o.Status, &o.Title, &o.Period, &o.Price, &o.Currency,
			&o.Config, &o.ActivatedAt, &o.ExpiresAt, &o.NextDueDate, &o.SuspendedAt, &o.SuspensionReason, &o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, &o)
	}
	return orders, nil
}

func (r *OrderRepository) ListOverdueSuspensions(ctx context.Context, overdueDays int) ([]*domain.Order, error) {
	query := `
		SELECT id, client_id, product_id, invoice_id, status, title, period, price, currency,
		       config, activated_at, expires_at, next_due_date, suspended_at, suspension_reason, created_at, updated_at
		FROM client_orders
		WHERE status = 'active' AND expires_at < (CURRENT_TIMESTAMP - ($1 || ' days')::INTERVAL)
	`
	rows, err := r.pool.Query(ctx, query, overdueDays)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(
			&o.ID, &o.ClientID, &o.ProductID, &o.InvoiceID, &o.Status, &o.Title, &o.Period, &o.Price, &o.Currency,
			&o.Config, &o.ActivatedAt, &o.ExpiresAt, &o.NextDueDate, &o.SuspendedAt, &o.SuspensionReason, &o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, &o)
	}
	return orders, nil
}

func (r *OrderRepository) Create(ctx context.Context, o *domain.Order) error {
	query := `
		INSERT INTO client_orders (
			client_id, product_id, invoice_id, status, title, period, price, currency,
			config, activated_at, expires_at, next_due_date, suspended_at, suspension_reason, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at
	`
	if o.Status == "" {
		o.Status = domain.OrderStatusPendingSetup
	}
	if o.Currency == "" {
		o.Currency = "USD"
	}

	return r.pool.QueryRow(ctx, query,
		o.ClientID, o.ProductID, o.InvoiceID, o.Status, o.Title, o.Period, o.Price, o.Currency,
		o.Config, o.ActivatedAt, o.ExpiresAt, o.NextDueDate, o.SuspendedAt, o.SuspensionReason,
	).Scan(&o.ID, &o.CreatedAt, &o.UpdatedAt)
}

func (r *OrderRepository) Update(ctx context.Context, o *domain.Order) error {
	query := `
		UPDATE client_orders SET
			invoice_id = $1, status = $2, title = $3, period = $4, price = $5, currency = $6,
			config = $7, activated_at = $8, expires_at = $9, next_due_date = $10,
			suspended_at = $11, suspension_reason = $12, updated_at = CURRENT_TIMESTAMP
		WHERE id = $13
	`
	tag, err := r.pool.Exec(ctx, query,
		o.InvoiceID, o.Status, o.Title, o.Period, o.Price, o.Currency,
		o.Config, o.ActivatedAt, o.ExpiresAt, o.NextDueDate,
		o.SuspendedAt, o.SuspensionReason, o.ID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return appErrors.ErrNotFound
	}
	return nil
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id int64, status domain.OrderStatus, reason *string) error {
	var query string
	var err error

	if status == domain.OrderStatusSuspended {
		query = `UPDATE client_orders SET status = $1, suspended_at = CURRENT_TIMESTAMP, suspension_reason = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`
		_, err = r.pool.Exec(ctx, query, status, reason, id)
	} else if status == domain.OrderStatusActive {
		query = `UPDATE client_orders SET status = $1, suspended_at = NULL, suspension_reason = NULL, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
		_, err = r.pool.Exec(ctx, query, status, id)
	} else {
		query = `UPDATE client_orders SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
		_, err = r.pool.Exec(ctx, query, status, id)
	}
	return err
}
