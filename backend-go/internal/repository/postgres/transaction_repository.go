package postgres

import (
	"context"
	"errors"

	"github.com/fossbilling/backend-go/internal/domain"
	appErrors "github.com/fossbilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	pool *pgxpool.Pool
}

func NewTransactionRepository(pool *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{pool: pool}
}

func (r *TransactionRepository) GetByID(ctx context.Context, id int64) (*domain.Transaction, error) {
	query := `
		SELECT id, invoice_id, gateway_id, txn_id, type, amount, currency, status, raw_payload, created_at
		FROM transactions
		WHERE id = $1
	`
	var t domain.Transaction
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.InvoiceID, &t.GatewayID, &t.TxnID, &t.Type, &t.Amount, &t.Currency, &t.Status, &t.RawPayload, &t.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *TransactionRepository) GetByTxnID(ctx context.Context, gatewayID, txnID string) (*domain.Transaction, error) {
	query := `
		SELECT id, invoice_id, gateway_id, txn_id, type, amount, currency, status, raw_payload, created_at
		FROM transactions
		WHERE gateway_id = $1 AND txn_id = $2
	`
	var t domain.Transaction
	err := r.pool.QueryRow(ctx, query, gatewayID, txnID).Scan(
		&t.ID, &t.InvoiceID, &t.GatewayID, &t.TxnID, &t.Type, &t.Amount, &t.Currency, &t.Status, &t.RawPayload, &t.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *TransactionRepository) Create(ctx context.Context, txn *domain.Transaction) error {
	query := `
		INSERT INTO transactions (invoice_id, gateway_id, txn_id, type, amount, currency, status, raw_payload, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP)
		RETURNING id, created_at
	`
	if txn.Type == "" {
		txn.Type = domain.TransactionTypePayment
	}
	if txn.Status == "" {
		txn.Status = domain.TransactionStatusPending
	}
	if txn.Currency == "" {
		txn.Currency = "USD"
	}

	return r.pool.QueryRow(ctx, query,
		txn.InvoiceID, txn.GatewayID, txn.TxnID, txn.Type, txn.Amount, txn.Currency, txn.Status, txn.RawPayload,
	).Scan(&txn.ID, &txn.CreatedAt)
}

func (r *TransactionRepository) UpdateStatus(ctx context.Context, id int64, status domain.TransactionStatus) error {
	query := `UPDATE transactions SET status = $1 WHERE id = $2`
	tag, err := r.pool.Exec(ctx, query, status, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return appErrors.ErrNotFound
	}
	return nil
}
