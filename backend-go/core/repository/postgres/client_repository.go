package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const clientCols = `id, group_id, email, password_hash, first_name, last_name, 
	company, address_1, address_2, city, state, postcode, country, phone_cc, 
	phone, currency, tax_exempt, status, created_at, updated_at`

type ClientRepository struct {
	pool *pgxpool.Pool
}

func NewClientRepository(pool *pgxpool.Pool) *ClientRepository {
	return &ClientRepository{pool: pool}
}

func scanClient(scanner interface{ Scan(...any) error }) (*domain.Client, error) {
	var c domain.Client
	err := scanner.Scan(
		&c.ID, &c.GroupID, &c.Email, &c.PasswordHash, &c.FirstName, &c.LastName,
		&c.Company, &c.Address1, &c.Address2, &c.City, &c.State, &c.Postcode,
		&c.Country, &c.PhoneCC, &c.Phone, &c.Currency, &c.TaxExempt, &c.Status,
		&c.CreatedAt, &c.UpdatedAt,
	)
	return &c, err
}

func (r *ClientRepository) GetByID(ctx context.Context, id int64) (*domain.Client, error) {
	row := r.pool.QueryRow(ctx, "SELECT "+clientCols+" FROM clients WHERE id = $1", id)
	c, err := scanClient(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get client by id: %w", err)
	}
	return c, nil
}

func (r *ClientRepository) GetByEmail(ctx context.Context, email string) (*domain.Client, error) {
	row := r.pool.QueryRow(ctx, "SELECT "+clientCols+" FROM clients WHERE LOWER(email) = LOWER($1)", email)
	c, err := scanClient(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get client by email: %w", err)
	}
	return c, nil
}

func (r *ClientRepository) Create(ctx context.Context, c *domain.Client) error {
	query := `INSERT INTO clients (group_id, email, password_hash, first_name, last_name, company, address_1, address_2, city, state, postcode, country, phone_cc, phone, currency, tax_exempt, status, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19) RETURNING id, created_at, updated_at`
	now := time.Now().UTC()
	c.CreatedAt, c.UpdatedAt = now, now
	if c.Status == "" { c.Status = domain.ClientStatusActive }
	if c.Currency == "" { c.Currency = "USD" }
	if c.Country == "" { c.Country = "US" }

	return r.pool.QueryRow(ctx, query,
		c.GroupID, c.Email, c.PasswordHash, c.FirstName, c.LastName, c.Company,
		c.Address1, c.Address2, c.City, c.State, c.Postcode, c.Country,
		c.PhoneCC, c.Phone, c.Currency, c.TaxExempt, c.Status, c.CreatedAt, c.UpdatedAt,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *ClientRepository) Update(ctx context.Context, c *domain.Client) error {
	query := `UPDATE clients SET first_name = $1, last_name = $2, company = $3, address_1 = $4, address_2 = $5, city = $6, state = $7, postcode = $8, country = $9, phone_cc = $10, phone = $11, currency = $12, tax_exempt = $13, status = $14, updated_at = $15 WHERE id = $16`
	c.UpdatedAt = time.Now().UTC()
	_, err := r.pool.Exec(ctx, query,
		c.FirstName, c.LastName, c.Company, c.Address1, c.Address2, c.City, c.State,
		c.Postcode, c.Country, c.PhoneCC, c.Phone, c.Currency, c.TaxExempt, c.Status, c.UpdatedAt, c.ID,
	)
	return err
}

func (r *ClientRepository) GetBalance(ctx context.Context, clientID int64) (decimal.Money, error) {
	query := `SELECT COALESCE(SUM(CASE WHEN type = 'credit' THEN amount ELSE -amount END), 0) FROM client_balances WHERE client_id = $1`
	var balance int64
	err := r.pool.QueryRow(ctx, query, clientID).Scan(&balance)
	return decimal.Money(balance), err
}

func (r *ClientRepository) AddBalanceTransaction(ctx context.Context, b *domain.ClientBalance) error {
	query := `INSERT INTO client_balances (client_id, type, amount, description, rel_id, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`
	b.CreatedAt = time.Now().UTC()
	return r.pool.QueryRow(ctx, query, b.ClientID, b.Type, int64(b.Amount), b.Description, b.RelID, b.CreatedAt).Scan(&b.ID, &b.CreatedAt)
}

func (r *ClientRepository) List(ctx context.Context, limit, offset int) ([]*domain.Client, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM clients`).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.pool.Query(ctx, "SELECT "+clientCols+" FROM clients ORDER BY id DESC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var clients []*domain.Client
	for rows.Next() {
		c, err := scanClient(rows)
		if err != nil {
			return nil, 0, err
		}
		clients = append(clients, c)
	}
	return clients, total, nil
}

func (r *ClientRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.pool.Exec(ctx, "DELETE FROM clients WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return appErrors.ErrNotFound
	}
	return nil
}

