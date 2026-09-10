package postgres

import (
	"context"
	"errors"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const txnCols = `id, invoice_id, gateway_id, txn_id, type, amount, currency, status, raw_payload, created_at`

type TransactionRepository struct{ pool *pgxpool.Pool }

func NewTransactionRepository(p *pgxpool.Pool) *TransactionRepository { return &TransactionRepository{p} }

func scanTx(r pgx.Row) (*domain.Transaction, error) {
	var t domain.Transaction; err := r.Scan(&t.ID, &t.InvoiceID, &t.GatewayID, &t.TxnID, &t.Type, &t.Amount, &t.Currency, &t.Status, &t.RawPayload, &t.CreatedAt); return &t, err
}

func (r *TransactionRepository) GetByID(ctx context.Context, id int64) (*domain.Transaction, error) {
	tx, err := scanTx(r.pool.QueryRow(ctx, "SELECT "+txnCols+" FROM transactions WHERE id = $1", id))
	if err != nil && errors.Is(err, pgx.ErrNoRows) { return nil, appErrors.ErrNotFound }; return tx, err
}

func (r *TransactionRepository) GetByTxnID(ctx context.Context, gid, tid string) (*domain.Transaction, error) {
	tx, err := scanTx(r.pool.QueryRow(ctx, "SELECT "+txnCols+" FROM transactions WHERE gateway_id = $1 AND txn_id = $2", gid, tid))
	if err != nil && errors.Is(err, pgx.ErrNoRows) { return nil, appErrors.ErrNotFound }; return tx, err
}

func (r *TransactionRepository) Create(ctx context.Context, t *domain.Transaction) error {
	if t.Type == "" { t.Type = domain.TransactionTypePayment }; if t.Status == "" { t.Status = domain.TransactionStatusPending }; if t.Currency == "" { t.Currency = "USD" }
	raw := t.RawPayload; if len(raw) == 0 { raw = []byte("{}") }
	return r.pool.QueryRow(ctx, `INSERT INTO transactions (invoice_id, gateway_id, txn_id, type, amount, currency, status, raw_payload, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP) RETURNING id, created_at`, t.InvoiceID, t.GatewayID, t.TxnID, t.Type, t.Amount, t.Currency, t.Status, raw).Scan(&t.ID, &t.CreatedAt)
}

func (r *TransactionRepository) UpdateStatus(ctx context.Context, id int64, st domain.TransactionStatus) error {
	t, err := r.pool.Exec(ctx, `UPDATE transactions SET status = $1 WHERE id = $2`, st, id)
	if err == nil && t.RowsAffected() == 0 { return appErrors.ErrNotFound }; return err
}
