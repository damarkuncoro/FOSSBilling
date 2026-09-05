package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InvoiceRepository struct {
	pool *pgxpool.Pool
}

func NewInvoiceRepository(pool *pgxpool.Pool) *InvoiceRepository {
	return &InvoiceRepository{pool: pool}
}

func (r *InvoiceRepository) GetByID(ctx context.Context, id int64) (*domain.Invoice, error) {
	query := `
		SELECT id, serie, nr, client_id, status, currency, currency_rate, subtotal, tax, total, tax_rate, due_at, paid_at, created_at, updated_at
		FROM invoices
		WHERE id = $1
	`
	var inv domain.Invoice
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&inv.ID, &inv.Serie, &inv.Nr, &inv.ClientID, &inv.Status, &inv.Currency, &inv.CurrencyRate,
		&inv.Subtotal, &inv.Tax, &inv.Total, &inv.TaxRate, &inv.DueAt, &inv.PaidAt, &inv.CreatedAt, &inv.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, err
	}

	itemsQuery := `
		SELECT id, invoice_id, order_id, title, period, price, quantity, unit, taxable, created_at
		FROM invoice_items
		WHERE invoice_id = $1
		ORDER BY id ASC
	`
	rows, err := r.pool.Query(ctx, itemsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var it domain.InvoiceItem
		if err := rows.Scan(
			&it.ID, &it.InvoiceID, &it.OrderID, &it.Title, &it.Period, &it.Price, &it.Quantity, &it.Unit, &it.Taxable, &it.CreatedAt,
		); err != nil {
			return nil, err
		}
		inv.Items = append(inv.Items, it)
	}

	return &inv, nil
}

func (r *InvoiceRepository) ListByClientID(ctx context.Context, clientID int64, limit, offset int) ([]*domain.Invoice, int, error) {
	countQuery := `SELECT COUNT(*) FROM invoices WHERE client_id = $1`
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, clientID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, serie, nr, client_id, status, currency, currency_rate, subtotal, tax, total, tax_rate, due_at, paid_at, created_at, updated_at
		FROM invoices
		WHERE client_id = $1
		ORDER BY id DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.pool.Query(ctx, query, clientID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var invoices []*domain.Invoice
	for rows.Next() {
		var inv domain.Invoice
		if err := rows.Scan(
			&inv.ID, &inv.Serie, &inv.Nr, &inv.ClientID, &inv.Status, &inv.Currency, &inv.CurrencyRate,
			&inv.Subtotal, &inv.Tax, &inv.Total, &inv.TaxRate, &inv.DueAt, &inv.PaidAt, &inv.CreatedAt, &inv.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		invoices = append(invoices, &inv)
	}

	return invoices, total, nil
}

func (r *InvoiceRepository) List(ctx context.Context, limit, offset int) ([]*domain.Invoice, int, error) {
	countQuery := `SELECT COUNT(*) FROM invoices`
	var total int
	if err := r.pool.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, serie, nr, client_id, status, currency, currency_rate, subtotal, tax, total, tax_rate, due_at, paid_at, created_at, updated_at
		FROM invoices
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var invoices []*domain.Invoice
	for rows.Next() {
		var inv domain.Invoice
		if err := rows.Scan(
			&inv.ID, &inv.Serie, &inv.Nr, &inv.ClientID, &inv.Status, &inv.Currency, &inv.CurrencyRate,
			&inv.Subtotal, &inv.Tax, &inv.Total, &inv.TaxRate, &inv.DueAt, &inv.PaidAt, &inv.CreatedAt, &inv.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		invoices = append(invoices, &inv)
	}

	return invoices, total, nil
}

func (r *InvoiceRepository) Create(ctx context.Context, inv *domain.Invoice, items []domain.InvoiceItem) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if inv.Status == "" {
		inv.Status = domain.InvoiceStatusUnpaid
	}
	if inv.Currency == "" {
		inv.Currency = "USD"
	}
	if inv.CurrencyRate == 0 {
		inv.CurrencyRate = 1.0
	}
	if inv.Serie == "" {
		inv.Serie = "INV"
	}

	invoiceQuery := `
		INSERT INTO invoices (
			serie, nr, client_id, status, currency, currency_rate, subtotal, tax, total, tax_rate, due_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at
	`
	if inv.Nr == "" {
		inv.Nr = "PENDING"
	}

	if err := tx.QueryRow(ctx, invoiceQuery,
		inv.Serie, inv.Nr, inv.ClientID, inv.Status, inv.Currency, inv.CurrencyRate,
		inv.Subtotal, inv.Tax, inv.Total, inv.TaxRate, inv.DueAt,
	).Scan(&inv.ID, &inv.CreatedAt, &inv.UpdatedAt); err != nil {
		return err
	}

	// Format Nr with padded ID
	inv.Nr = fmt.Sprintf("%05d", inv.ID)
	_, _ = tx.Exec(ctx, `UPDATE invoices SET nr = $1 WHERE id = $2`, inv.Nr, inv.ID)

	itemQuery := `
		INSERT INTO invoice_items (invoice_id, order_id, title, period, price, quantity, unit, taxable, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP)
		RETURNING id, created_at
	`
	for i := range items {
		items[i].InvoiceID = inv.ID
		if items[i].Quantity <= 0 {
			items[i].Quantity = 1
		}
		if items[i].Unit == "" {
			items[i].Unit = "unit"
		}

		if err := tx.QueryRow(ctx, itemQuery,
			items[i].InvoiceID, items[i].OrderID, items[i].Title, items[i].Period,
			items[i].Price, items[i].Quantity, items[i].Unit, items[i].Taxable,
		).Scan(&items[i].ID, &items[i].CreatedAt); err != nil {
			return err
		}
	}
	inv.Items = items

	return tx.Commit(ctx)
}

func (r *InvoiceRepository) MarkAsPaid(ctx context.Context, id int64, paidAt time.Time) error {
	query := `
		UPDATE invoices
		SET status = 'paid', paid_at = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	tag, err := r.pool.Exec(ctx, query, paidAt, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return appErrors.ErrNotFound
	}
	return nil
}

func (r *InvoiceRepository) Update(ctx context.Context, inv *domain.Invoice) error {
	query := `
		UPDATE invoices SET
			serie = $1, nr = $2, status = $3, currency = $4, currency_rate = $5,
			subtotal = $6, tax = $7, total = $8, tax_rate = $9, due_at = $10,
			paid_at = $11, updated_at = CURRENT_TIMESTAMP
		WHERE id = $12
	`
	tag, err := r.pool.Exec(ctx, query,
		inv.Serie, inv.Nr, inv.Status, inv.Currency, inv.CurrencyRate,
		inv.Subtotal, inv.Tax, inv.Total, inv.TaxRate, inv.DueAt,
		inv.PaidAt, inv.ID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return appErrors.ErrNotFound
	}
	return nil
}
