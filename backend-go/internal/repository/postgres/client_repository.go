package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClientRepository struct {
	pool *pgxpool.Pool
}

func NewClientRepository(pool *pgxpool.Pool) *ClientRepository {
	return &ClientRepository{pool: pool}
}

func (r *ClientRepository) GetByID(ctx context.Context, id int64) (*domain.Client, error) {
	query := `
		SELECT id, group_id, email, password_hash, first_name, last_name, 
		       company, address_1, address_2, city, state, postcode, 
		       country, phone_cc, phone, currency, tax_exempt, status, 
		       created_at, updated_at
		FROM clients
		WHERE id = $1
	`
	c := &domain.Client{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.GroupID, &c.Email, &c.PasswordHash, &c.FirstName, &c.LastName,
		&c.Company, &c.Address1, &c.Address2, &c.City, &c.State, &c.Postcode,
		&c.Country, &c.PhoneCC, &c.Phone, &c.Currency, &c.TaxExempt, &c.Status,
		&c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get client by id: %w", err)
	}
	return c, nil
}

func (r *ClientRepository) GetByEmail(ctx context.Context, email string) (*domain.Client, error) {
	query := `
		SELECT id, group_id, email, password_hash, first_name, last_name, 
		       company, address_1, address_2, city, state, postcode, 
		       country, phone_cc, phone, currency, tax_exempt, status, 
		       created_at, updated_at
		FROM clients
		WHERE LOWER(email) = LOWER($1)
	`
	c := &domain.Client{}
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&c.ID, &c.GroupID, &c.Email, &c.PasswordHash, &c.FirstName, &c.LastName,
		&c.Company, &c.Address1, &c.Address2, &c.City, &c.State, &c.Postcode,
		&c.Country, &c.PhoneCC, &c.Phone, &c.Currency, &c.TaxExempt, &c.Status,
		&c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get client by email: %w", err)
	}
	return c, nil
}

func (r *ClientRepository) Create(ctx context.Context, c *domain.Client) error {
	query := `
		INSERT INTO clients (
			group_id, email, password_hash, first_name, last_name, company,
			address_1, address_2, city, state, postcode, country,
			phone_cc, phone, currency, tax_exempt, status, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11, $12,
			$13, $14, $15, $16, $17, $18, $19
		) RETURNING id, created_at, updated_at
	`
	now := time.Now().UTC()
	c.CreatedAt = now
	c.UpdatedAt = now
	if c.Status == "" {
		c.Status = domain.ClientStatusActive
	}
	if c.Currency == "" {
		c.Currency = "USD"
	}
	if c.Country == "" {
		c.Country = "US"
	}

	err := r.pool.QueryRow(ctx, query,
		c.GroupID, c.Email, c.PasswordHash, c.FirstName, c.LastName, c.Company,
		c.Address1, c.Address2, c.City, c.State, c.Postcode, c.Country,
		c.PhoneCC, c.Phone, c.Currency, c.TaxExempt, c.Status, c.CreatedAt, c.UpdatedAt,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert client: %w", err)
	}
	return nil
}

func (r *ClientRepository) Update(ctx context.Context, c *domain.Client) error {
	query := `
		UPDATE clients SET
			first_name = $1, last_name = $2, company = $3,
			address_1 = $4, address_2 = $5, city = $6, state = $7, postcode = $8,
			country = $9, phone_cc = $10, phone = $11, currency = $12,
			tax_exempt = $13, status = $14, updated_at = $15
		WHERE id = $16
	`
	c.UpdatedAt = time.Now().UTC()
	_, err := r.pool.Exec(ctx, query,
		c.FirstName, c.LastName, c.Company,
		c.Address1, c.Address2, c.City, c.State, c.Postcode,
		c.Country, c.PhoneCC, c.Phone, c.Currency,
		c.TaxExempt, c.Status, c.UpdatedAt, c.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update client: %w", err)
	}
	return nil
}

func (r *ClientRepository) GetBalance(ctx context.Context, clientID int64) (decimal.Money, error) {
	query := `
		SELECT 
			COALESCE(SUM(CASE WHEN type = 'credit' THEN amount ELSE -amount END), 0) as balance
		FROM client_balances
		WHERE client_id = $1
	`
	var balance int64
	err := r.pool.QueryRow(ctx, query, clientID).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate client balance: %w", err)
	}
	return decimal.Money(balance), nil
}

func (r *ClientRepository) AddBalanceTransaction(ctx context.Context, b *domain.ClientBalance) error {
	query := `
		INSERT INTO client_balances (client_id, type, amount, description, rel_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	b.CreatedAt = time.Now().UTC()
	err := r.pool.QueryRow(ctx, query,
		b.ClientID, b.Type, int64(b.Amount), b.Description, b.RelID, b.CreatedAt,
	).Scan(&b.ID, &b.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to add balance record: %w", err)
	}
	return nil
}

func (r *ClientRepository) List(ctx context.Context, limit, offset int) ([]*domain.Client, int, error) {
	countQuery := `SELECT COUNT(*) FROM clients`
	var total int
	if err := r.pool.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			id, group_id, email, password_hash, first_name, last_name, company,
			address_1, address_2, city, state, postcode, country,
			phone_cc, phone, currency, tax_exempt, status, created_at, updated_at
		FROM clients
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var clients []*domain.Client
	for rows.Next() {
		var c domain.Client
		if err := rows.Scan(
			&c.ID, &c.GroupID, &c.Email, &c.PasswordHash, &c.FirstName, &c.LastName, &c.Company,
			&c.Address1, &c.Address2, &c.City, &c.State, &c.Postcode, &c.Country,
			&c.PhoneCC, &c.Phone, &c.Currency, &c.TaxExempt, &c.Status, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		clients = append(clients, &c)
	}

	return clients, total, nil
}

