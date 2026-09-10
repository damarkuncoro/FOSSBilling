package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const invoiceCols = `id, serie, nr, client_id, status, currency, currency_rate, subtotal, tax, total, tax_rate, due_at, paid_at, created_at, updated_at`

type InvoiceRepository struct{ pool *pgxpool.Pool }

func NewInvoiceRepository(p *pgxpool.Pool) *InvoiceRepository { return &InvoiceRepository{p} }

func scanInvoice(r pgx.Row) (*domain.Invoice, error) {
	var inv domain.Invoice
	err := r.Scan(&inv.ID, &inv.Serie, &inv.Nr, &inv.ClientID, &inv.Status, &inv.Currency, &inv.CurrencyRate, &inv.Subtotal, &inv.Tax, &inv.Total, &inv.TaxRate, &inv.DueAt, &inv.PaidAt, &inv.CreatedAt, &inv.UpdatedAt)
	return &inv, err
}

func (r *InvoiceRepository) GetByID(ctx context.Context, id int64) (*domain.Invoice, error) {
	inv, err := scanInvoice(r.pool.QueryRow(ctx, "SELECT "+invoiceCols+" FROM invoices WHERE id = $1", id))
	if err != nil { return nil, err }
	rows, _ := r.pool.Query(ctx, `SELECT id, invoice_id, order_id, title, period, price, quantity, unit, taxable, created_at FROM invoice_items WHERE invoice_id = $1 ORDER BY id ASC`, id)
	defer rows.Close()
	for rows.Next() {
		var it domain.InvoiceItem; _ = rows.Scan(&it.ID, &it.InvoiceID, &it.OrderID, &it.Title, &it.Period, &it.Price, &it.Quantity, &it.Unit, &it.Taxable, &it.CreatedAt)
		inv.Items = append(inv.Items, it)
	}
	return inv, nil
}

func (r *InvoiceRepository) ListByClientID(ctx context.Context, cID int64, l, o int) ([]*domain.Invoice, int, error) {
	res, err := list(ctx, r.pool, "SELECT "+invoiceCols+" FROM invoices WHERE client_id = $1 ORDER BY id DESC LIMIT $2 OFFSET $3", scanInvoice, cID, l, o)
	return res, total(ctx, r.pool, "SELECT COUNT(*) FROM invoices WHERE client_id = $1", cID), err
}

func (r *InvoiceRepository) List(ctx context.Context, l, o int) ([]*domain.Invoice, int, error) {
	res, err := list(ctx, r.pool, "SELECT "+invoiceCols+" FROM invoices ORDER BY id DESC LIMIT $1 OFFSET $2", scanInvoice, l, o)
	return res, total(ctx, r.pool, "SELECT COUNT(*) FROM invoices"), err
}

func (r *InvoiceRepository) Create(ctx context.Context, inv *domain.Invoice, items []domain.InvoiceItem) error {
	tx, err := r.pool.Begin(ctx); if err != nil { return err }; defer tx.Rollback(ctx)
	if inv.Status == "" { inv.Status = domain.InvoiceStatusUnpaid }; if inv.Currency == "" { inv.Currency = "USD" }; if inv.Serie == "" { inv.Serie = "INV" }
	if err := tx.QueryRow(ctx, `INSERT INTO invoices (serie, nr, client_id, status, currency, currency_rate, subtotal, tax, total, tax_rate, due_at, created_at, updated_at) VALUES ($1, 'PENDING', $2, $3, $4, 1.0, $5, $6, $7, $8, $9, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id, created_at, updated_at`, inv.Serie, inv.ClientID, inv.Status, inv.Currency, inv.Subtotal, inv.Tax, inv.Total, inv.TaxRate, inv.DueAt).Scan(&inv.ID, &inv.CreatedAt, &inv.UpdatedAt); err != nil { return err }
	inv.Nr = fmt.Sprintf("%05d", inv.ID); _, _ = tx.Exec(ctx, `UPDATE invoices SET nr = $1 WHERE id = $2`, inv.Nr, inv.ID)
	for i := range items {
		items[i].InvoiceID = inv.ID; if items[i].Quantity <= 0 { items[i].Quantity = 1 }
		_ = tx.QueryRow(ctx, `INSERT INTO invoice_items (invoice_id, order_id, title, period, price, quantity, unit, taxable, created_at) VALUES ($1, $2, $3, $4, $5, $6, 'unit', $7, CURRENT_TIMESTAMP) RETURNING id, created_at`, inv.ID, items[i].OrderID, items[i].Title, items[i].Period, items[i].Price, items[i].Quantity, items[i].Taxable).Scan(&items[i].ID, &items[i].CreatedAt)
	}
	inv.Items = items; return tx.Commit(ctx)
}

func (r *InvoiceRepository) MarkAsPaid(ctx context.Context, id int64, paid time.Time) error {
	t, err := r.pool.Exec(ctx, `UPDATE invoices SET status = 'paid', paid_at = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`, paid, id)
	if err == nil && t.RowsAffected() == 0 { return appErrors.ErrNotFound }; return err
}

func (r *InvoiceRepository) UpdateStatus(ctx context.Context, id int64, st domain.InvoiceStatus) error {
	_, err := r.pool.Exec(ctx, `UPDATE invoices SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`, st, id); return err
}

func (r *InvoiceRepository) UpdateStatusAtomic(ctx context.Context, id int64, nSt, oSt domain.InvoiceStatus) (bool, error) {
	t, err := r.pool.Exec(ctx, `UPDATE invoices SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2 AND status = $3`, nSt, id, oSt)
	return t.RowsAffected() > 0, err
}

func (r *InvoiceRepository) Update(ctx context.Context, inv *domain.Invoice) error {
	_, err := r.pool.Exec(ctx, `UPDATE invoices SET serie = $1, nr = $2, status = $3, currency = $4, subtotal = $5, tax = $6, total = $7, tax_rate = $8, due_at = $9, paid_at = $10, updated_at = CURRENT_TIMESTAMP WHERE id = $11`, inv.Serie, inv.Nr, inv.Status, inv.Currency, inv.Subtotal, inv.Tax, inv.Total, inv.TaxRate, inv.DueAt, inv.PaidAt, inv.ID)
	return err
}

func (r *InvoiceRepository) Delete(ctx context.Context, id int64) error { _, err := r.pool.Exec(ctx, `DELETE FROM invoices WHERE id = $1`, id); return err }
