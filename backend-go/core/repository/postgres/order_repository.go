package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const orderCols = `id, client_id, product_id, invoice_id, status, title, period, price, currency, config, activated_at, expires_at, next_due_date, suspended_at, suspension_reason, created_at, updated_at`

type OrderRepository struct{ pool *pgxpool.Pool }

func NewOrderRepository(p *pgxpool.Pool) *OrderRepository { return &OrderRepository{p} }

func scanOrder(r pgx.Row) (*domain.Order, error) {
	var o domain.Order; var pid *int64
	err := r.Scan(&o.ID, &o.ClientID, &pid, &o.InvoiceID, &o.Status, &o.Title, &o.Period, &o.Price, &o.Currency, &o.Config, &o.ActivatedAt, &o.ExpiresAt, &o.NextDueDate, &o.SuspendedAt, &o.SuspensionReason, &o.CreatedAt, &o.UpdatedAt)
	if pid != nil { o.ProductID = *pid }; return &o, err
}

func (r *OrderRepository) GetByID(ctx context.Context, id int64) (*domain.Order, error) {
	o, err := scanOrder(r.pool.QueryRow(ctx, "SELECT "+orderCols+" FROM client_orders WHERE id = $1", id))
	if err != nil && errors.Is(err, pgx.ErrNoRows) { return nil, appErrors.ErrNotFound }; return o, err
}

func (r *OrderRepository) ListByClientID(ctx context.Context, cID int64, l, o int) ([]*domain.Order, int, error) {
	res, err := list(ctx, r.pool, "SELECT "+orderCols+" FROM client_orders WHERE client_id = $1 ORDER BY id DESC LIMIT $2 OFFSET $3", scanOrder, cID, l, o)
	return res, total(ctx, r.pool, "SELECT COUNT(*) FROM client_orders WHERE client_id = $1", cID), err
}

func (r *OrderRepository) List(ctx context.Context, l, o int) ([]*domain.Order, int, error) {
	res, err := list(ctx, r.pool, "SELECT "+orderCols+" FROM client_orders ORDER BY id DESC LIMIT $1 OFFSET $2", scanOrder, l, o)
	return res, total(ctx, r.pool, "SELECT COUNT(*) FROM client_orders"), err
}

func (r *OrderRepository) ListByInvoiceID(ctx context.Context, invID int64) ([]*domain.Order, error) {
	return list(ctx, r.pool, "SELECT "+orderCols+" FROM client_orders WHERE invoice_id = $1", scanOrder, invID)
}

func (r *OrderRepository) ListDueOrders(ctx context.Context, due time.Time) ([]*domain.Order, error) {
	return list(ctx, r.pool, "SELECT "+orderCols+" FROM client_orders WHERE status = 'active' AND next_due_date <= $1", scanOrder, due)
}

func (r *OrderRepository) ListOverdueSuspensions(ctx context.Context, days int) ([]*domain.Order, error) {
	return list(ctx, r.pool, "SELECT "+orderCols+" FROM client_orders WHERE status = 'active' AND expires_at < (CURRENT_TIMESTAMP - ($1 * INTERVAL '1 day'))", scanOrder, days)
}

func (r *OrderRepository) ListPendingProvisioning(ctx context.Context) ([]*domain.Order, error) {
	return list(ctx, r.pool, "SELECT o."+orderCols+" FROM client_orders o JOIN invoices i ON o.invoice_id = i.id WHERE o.status = 'pending_setup' AND i.status = 'paid'", scanOrder)
}

func (r *OrderRepository) Create(ctx context.Context, o *domain.Order) error {
	q := `INSERT INTO client_orders (client_id, product_id, invoice_id, status, title, period, price, currency, config, activated_at, expires_at, next_due_date, suspended_at, suspension_reason, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id, created_at, updated_at`
	if o.Status == "" { o.Status = domain.OrderStatusPendingSetup }; if o.Currency == "" { o.Currency = "USD" }
	var pid *int64; if o.ProductID > 0 { pid = &o.ProductID }
	return r.pool.QueryRow(ctx, q, o.ClientID, pid, o.InvoiceID, o.Status, o.Title, o.Period, o.Price, o.Currency, o.Config, o.ActivatedAt, o.ExpiresAt, o.NextDueDate, o.SuspendedAt, o.SuspensionReason).Scan(&o.ID, &o.CreatedAt, &o.UpdatedAt)
}

func (r *OrderRepository) Update(ctx context.Context, o *domain.Order) error {
	q := `UPDATE client_orders SET invoice_id = $1, status = $2, title = $3, period = $4, price = $5, currency = $6, config = $7, activated_at = $8, expires_at = $9, next_due_date = $10, suspended_at = $11, suspension_reason = $12, updated_at = CURRENT_TIMESTAMP WHERE id = $13`
	t, err := r.pool.Exec(ctx, q, o.InvoiceID, o.Status, o.Title, o.Period, o.Price, o.Currency, o.Config, o.ActivatedAt, o.ExpiresAt, o.NextDueDate, o.SuspendedAt, o.SuspensionReason, o.ID)
	if err == nil && t.RowsAffected() == 0 { return appErrors.ErrNotFound }; return err
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id int64, st domain.OrderStatus, res *string) error {
	if st == domain.OrderStatusSuspended { return r.pool.QueryRow(ctx, "UPDATE client_orders SET status = $1, suspended_at = CURRENT_TIMESTAMP, suspension_reason = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3", st, res, id).Scan() }
	if st == domain.OrderStatusActive { return r.pool.QueryRow(ctx, "UPDATE client_orders SET status = $1, suspended_at = NULL, suspension_reason = NULL, updated_at = CURRENT_TIMESTAMP WHERE id = $2", st, id).Scan() }
	_, err := r.pool.Exec(ctx, "UPDATE client_orders SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2", st, id); return err
}
