package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const clientCols = `id, aid, group_id, email, password_hash, type, first_name, last_name, gender, birthday, company, address_1, address_2, city, state, postcode, country, phone_cc, phone, notes, currency, billing_email, referred_by, tax_exempt, status, two_factor_enabled, two_factor_secret, custom_1, custom_2, custom_3, custom_4, custom_5, custom_6, custom_7, custom_8, custom_9, custom_10, custom_11, custom_12, custom_13, custom_14, custom_15, custom_16, custom_17, custom_18, custom_19, custom_20, created_at, updated_at`

type ClientRepository struct{ pool *pgxpool.Pool }

func NewClientRepository(p *pgxpool.Pool) *ClientRepository { return &ClientRepository{p} }

func scanClient(r pgx.Row) (*domain.Client, error) {
	var c domain.Client
	err := r.Scan(&c.ID, &c.AID, &c.GroupID, &c.Email, &c.PasswordHash, &c.Type, &c.FirstName, &c.LastName, &c.Gender, &c.Birthday, &c.Company, &c.Address1, &c.Address2, &c.City, &c.State, &c.Postcode, &c.Country, &c.PhoneCC, &c.Phone, &c.Notes, &c.Currency, &c.BillingEmail, &c.ReferredBy, &c.TaxExempt, &c.Status, &c.TwoFactorEnabled, &c.TwoFactorSecret, &c.Custom1, &c.Custom2, &c.Custom3, &c.Custom4, &c.Custom5, &c.Custom6, &c.Custom7, &c.Custom8, &c.Custom9, &c.Custom10, &c.Custom11, &c.Custom12, &c.Custom13, &c.Custom14, &c.Custom15, &c.Custom16, &c.Custom17, &c.Custom18, &c.Custom19, &c.Custom20, &c.CreatedAt, &c.UpdatedAt)
	return &c, err
}

func (r *ClientRepository) GetByID(ctx context.Context, id int64) (*domain.Client, error) {
	c, err := scanClient(r.pool.QueryRow(ctx, "SELECT "+clientCols+" FROM clients WHERE id = $1", id))
	if err != nil && errors.Is(err, pgx.ErrNoRows) { return nil, appErrors.ErrNotFound }; return c, err
}

func (r *ClientRepository) GetByEmail(ctx context.Context, email string) (*domain.Client, error) {
	c, err := scanClient(r.pool.QueryRow(ctx, "SELECT "+clientCols+" FROM clients WHERE LOWER(email) = LOWER($1)", email))
	if err != nil && errors.Is(err, pgx.ErrNoRows) { return nil, appErrors.ErrNotFound }; return c, err
}

func (r *ClientRepository) List(ctx context.Context, l, o int) ([]*domain.Client, int, error) {
	res, err := list(ctx, r.pool, "SELECT "+clientCols+" FROM clients ORDER BY id DESC LIMIT $1 OFFSET $2", scanClient, l, o)
	return res, total(ctx, r.pool, "SELECT COUNT(*) FROM clients"), err
}

func (r *ClientRepository) Create(ctx context.Context, c *domain.Client) error {
	q := `INSERT INTO clients (aid, group_id, email, password_hash, type, first_name, last_name, gender, birthday, company, address_1, address_2, city, state, postcode, country, phone_cc, phone, notes, currency, billing_email, referred_by, tax_exempt, status, two_factor_enabled, two_factor_secret, custom_1, custom_2, custom_3, custom_4, custom_5, custom_6, custom_7, custom_8, custom_9, custom_10, custom_11, custom_12, custom_13, custom_14, custom_15, custom_16, custom_17, custom_18, custom_19, custom_20, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48) RETURNING id, created_at, updated_at`
	c.CreatedAt, c.UpdatedAt = time.Now().UTC(), time.Now().UTC()
	if c.Status == "" { c.Status = domain.ClientStatusActive }; if c.Currency == "" { c.Currency = "USD" }; if c.Country == "" { c.Country = "US" }; if c.Type == "" { c.Type = "individual" }
	return r.pool.QueryRow(ctx, q, c.AID, c.GroupID, c.Email, c.PasswordHash, c.Type, c.FirstName, c.LastName, c.Gender, c.Birthday, c.Company, c.Address1, c.Address2, c.City, c.State, c.Postcode, c.Country, c.PhoneCC, c.Phone, c.Notes, c.Currency, c.BillingEmail, c.ReferredBy, c.TaxExempt, c.Status, c.TwoFactorEnabled, c.TwoFactorSecret, c.Custom1, c.Custom2, c.Custom3, c.Custom4, c.Custom5, c.Custom6, c.Custom7, c.Custom8, c.Custom9, c.Custom10, c.Custom11, c.Custom12, c.Custom13, c.Custom14, c.Custom15, c.Custom16, c.Custom17, c.Custom18, c.Custom19, c.Custom20, c.CreatedAt, c.UpdatedAt).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *ClientRepository) Update(ctx context.Context, c *domain.Client) error {
	q := `UPDATE clients SET first_name = $1, last_name = $2, company = $3, address_1 = $4, address_2 = $5, city = $6, state = $7, postcode = $8, country = $9, phone_cc = $10, phone = $11, currency = $12, tax_exempt = $13, status = $14, two_factor_enabled = $15, two_factor_secret = $16, gender = $17, birthday = $18, type = $19, notes = $20, billing_email = $21, referred_by = $22, custom_1 = $23, custom_2 = $24, custom_3 = $25, custom_4 = $26, custom_5 = $27, custom_6 = $28, custom_7 = $29, custom_8 = $30, custom_9 = $31, custom_10 = $32, custom_11 = $33, custom_12 = $34, custom_13 = $35, custom_14 = $36, custom_15 = $37, custom_16 = $38, custom_17 = $39, custom_18 = $40, custom_19 = $41, custom_20 = $42, updated_at = $43 WHERE id = $44`
	c.UpdatedAt = time.Now().UTC()
	_, err := r.pool.Exec(ctx, q, c.FirstName, c.LastName, c.Company, c.Address1, c.Address2, c.City, c.State, c.Postcode, c.Country, c.PhoneCC, c.Phone, c.Currency, c.TaxExempt, c.Status, c.TwoFactorEnabled, c.TwoFactorSecret, c.Gender, c.Birthday, c.Type, c.Notes, c.BillingEmail, c.ReferredBy, c.Custom1, c.Custom2, c.Custom3, c.Custom4, c.Custom5, c.Custom6, c.Custom7, c.Custom8, c.Custom9, c.Custom10, c.Custom11, c.Custom12, c.Custom13, c.Custom14, c.Custom15, c.Custom16, c.Custom17, c.Custom18, c.Custom19, c.Custom20, c.UpdatedAt, c.ID)
	return err
}

func (r *ClientRepository) GetBalance(ctx context.Context, id int64) (decimal.Money, error) {
	var bal int64; err := r.pool.QueryRow(ctx, `SELECT COALESCE(SUM(CASE WHEN type = 'credit' THEN amount ELSE -amount END), 0) FROM client_balances WHERE client_id = $1`, id).Scan(&bal); return decimal.Money(bal), err
}

func (r *ClientRepository) AddBalanceTransaction(ctx context.Context, b *domain.ClientBalance) error {
	b.CreatedAt = time.Now().UTC()
	return r.pool.QueryRow(ctx, `INSERT INTO client_balances (client_id, type, amount, description, rel_id, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`, b.ClientID, b.Type, int64(b.Amount), b.Description, b.RelID, b.CreatedAt).Scan(&b.ID, &b.CreatedAt)
}

func (r *ClientRepository) Delete(ctx context.Context, id int64) error {
	t, err := r.pool.Exec(ctx, "DELETE FROM clients WHERE id = $1", id)
	if err == nil && t.RowsAffected() == 0 { return appErrors.ErrNotFound }; return err
}
