package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const clientCols = `id, aid, group_id, email, password_hash, type, first_name, last_name, gender, birthday, company, address_1, address_2, city, state, postcode, country, phone_cc, phone, notes, currency, billing_email, referrer_id, referred_by, tax_exempt, status, oauth_provider, oauth_id, two_factor_enabled, two_factor_secret, custom_1, custom_2, custom_3, custom_4, custom_5, custom_6, custom_7, custom_8, custom_9, custom_10, custom_11, custom_12, custom_13, custom_14, custom_15, custom_16, custom_17, custom_18, custom_19, custom_20, created_at, updated_at`

type ClientRepository struct{ pool *pgxpool.Pool }

func NewClientRepository(p *pgxpool.Pool) *ClientRepository { return &ClientRepository{p} }

func scanClient(r pgx.Row) (*domain.Client, error) {
	var c domain.Client
	var aid, gender, company, companyVat, companyNumber, address1, address2, city, state, postcode, country, phoneCc, phone, notes, billingEmail, referredBy, oauthProvider, oauthId sql.NullString
	var c1, c2, c3, c4, c5, c6, c7, c8, c9, c10, c11, c12, c13, c14, c15, c16, c17, c18, c19, c20 sql.NullString

	err := r.Scan(
		&c.ID, &aid, &c.GroupID, &c.Email, &c.PasswordHash, &c.Type, &c.FirstName, &c.LastName,
		&gender, &c.Birthday, &company, &address1, &address2, &city, &state, &postcode,
		&country, &phoneCc, &phone, &notes, &c.Currency, &billingEmail, &c.ReferrerID,
		&referredBy, &c.TaxExempt, &c.Status, &oauthProvider, &oauthId, &c.TwoFactorEnabled,
		&c.TwoFactorSecret, &c1, &c2, &c3, &c4, &c5, &c6, &c7, &c8, &c9, &c10, &c11, &c12, &c13, &c14, &c15, &c16, &c17, &c18, &c19, &c20,
		&c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if aid.Valid { c.AID = &aid.String }
	c.Gender = gender.String; c.Company = company.String; c.Address1 = address1.String
	c.CompanyVat = companyVat.String; c.CompanyNumber = companyNumber.String
	c.Address2 = address2.String; c.City = city.String; c.State = state.String
	c.Postcode = postcode.String; c.Country = country.String; c.PhoneCC = phoneCc.String
	c.Phone = phone.String; c.Notes = notes.String; c.BillingEmail = billingEmail.String
	c.ReferredBy = referredBy.String; c.OAuthProvider = oauthProvider.String; c.OAuthID = oauthId.String
	c.Custom1 = c1.String; c.Custom2 = c2.String; c.Custom3 = c3.String; c.Custom4 = c4.String
	c.Custom5 = c5.String; c.Custom6 = c6.String; c.Custom7 = c7.String; c.Custom8 = c8.String
	c.Custom9 = c9.String; c.Custom10 = c10.String; c.Custom11 = c11.String; c.Custom12 = c12.String
	c.Custom13 = c13.String; c.Custom14 = c14.String; c.Custom15 = c15.String; c.Custom16 = c16.String
	c.Custom17 = c17.String; c.Custom18 = c18.String; c.Custom19 = c19.String; c.Custom20 = c20.String

	return &c, nil
}

func (r *ClientRepository) GetByID(ctx context.Context, id int64) (*domain.Client, error) {
	c, err := scanClient(r.pool.QueryRow(ctx, "SELECT "+clientCols+" FROM clients WHERE id = $1", id))
	if err != nil && errors.Is(err, pgx.ErrNoRows) { return nil, appErrors.ErrNotFound }; return c, err
}

func (r *ClientRepository) GetByEmail(ctx context.Context, email string) (*domain.Client, error) {
	c, err := scanClient(r.pool.QueryRow(ctx, "SELECT "+clientCols+" FROM clients WHERE LOWER(email) = LOWER($1)", email))
	if err != nil && errors.Is(err, pgx.ErrNoRows) { return nil, appErrors.ErrNotFound }; return c, err
}

func (r *ClientRepository) GetByOAuth(ctx context.Context, provider, oauthID string) (*domain.Client, error) {
	c, err := scanClient(r.pool.QueryRow(ctx, "SELECT "+clientCols+" FROM clients WHERE oauth_provider = $1 AND oauth_id = $2", provider, oauthID))
	if err != nil && errors.Is(err, pgx.ErrNoRows) { return nil, appErrors.ErrNotFound }; return c, err
}

func (r *ClientRepository) List(ctx context.Context, l, o int) ([]*domain.Client, int, error) {
	res, err := list(ctx, r.pool, "SELECT "+clientCols+" FROM clients ORDER BY id DESC LIMIT $1 OFFSET $2", scanClient, l, o)
	return res, total(ctx, r.pool, "SELECT COUNT(*) FROM clients"), err
}

func (r *ClientRepository) Create(ctx context.Context, c *domain.Client) error {
	q := `INSERT INTO clients (
		aid, group_id, email, password_hash, type, first_name, last_name, gender, birthday, company,
		address_1, address_2, city, state, postcode, country, phone_cc, phone, notes, currency,
		billing_email, referrer_id, referred_by, tax_exempt, status, oauth_provider, oauth_id,
		two_factor_enabled, two_factor_secret, custom_1, custom_2, custom_3, custom_4, custom_5,
		custom_6, custom_7, custom_8, custom_9, custom_10, custom_11, custom_12, custom_13,
		custom_14, custom_15, custom_16, custom_17, custom_18, custom_19, custom_20,
		created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
		$11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
		$21, $22, $23, $24, $25, $26, $27,
		$28, $29, $30, $31, $32, $33, $34, $35,
		$36, $37, $38, $39, $40, $41, $42, $43,
		$44, $45, $46, $47, $48, $49, $50, $51
	) RETURNING id`
	c.CreatedAt, c.UpdatedAt = time.Now().UTC(), time.Now().UTC()
	if c.Status == "" { c.Status = domain.ClientStatusActive }
	if c.Currency == "" { c.Currency = "USD" }
	if c.Country == "" { c.Country = "US" }
	if c.Type == "" { c.Type = "individual" }

	return r.pool.QueryRow(ctx, q,
		c.AID, c.GroupID, c.Email, c.PasswordHash, c.Type, c.FirstName, c.LastName, c.Gender, c.Birthday, c.Company,
		c.Address1, c.Address2, c.City, c.State, c.Postcode, c.Country, c.PhoneCC, c.Phone, c.Notes, c.Currency,
		c.BillingEmail, c.ReferrerID, c.ReferredBy, c.TaxExempt, c.Status, c.OAuthProvider, c.OAuthID,
		c.TwoFactorEnabled, c.TwoFactorSecret, c.Custom1, c.Custom2, c.Custom3, c.Custom4, c.Custom5,
		c.Custom6, c.Custom7, c.Custom8, c.Custom9, c.Custom10, c.Custom11, c.Custom12, c.Custom13,
		c.Custom14, c.Custom15, c.Custom16, c.Custom17, c.Custom18, c.Custom19, c.Custom20,
		c.CreatedAt, c.UpdatedAt,
	).Scan(&c.ID)
}

func (r *ClientRepository) Update(ctx context.Context, c *domain.Client) error {
	q := `UPDATE clients SET first_name = $1, last_name = $2, company = $3, address_1 = $4, address_2 = $5, city = $6, state = $7, postcode = $8, country = $9, phone_cc = $10, phone = $11, currency = $12, tax_exempt = $13, status = $14, oauth_provider = $15, oauth_id = $16, two_factor_enabled = $17, two_factor_secret = $18, gender = $19, birthday = $20, type = $21, notes = $22, billing_email = $23, referrer_id = $24, referred_by = $25, custom_1 = $26, custom_2 = $27, custom_3 = $28, custom_4 = $29, custom_5 = $30, custom_6 = $31, custom_7 = $32, custom_8 = $33, custom_9 = $34, custom_10 = $35, custom_11 = $36, custom_12 = $37, custom_13 = $38, custom_14 = $39, custom_15 = $40, custom_16 = $41, custom_17 = $42, custom_18 = $43, custom_19 = $44, custom_20 = $45, updated_at = $46 WHERE id = $47`
	c.UpdatedAt = time.Now().UTC()
	_, err := r.pool.Exec(ctx, q, c.FirstName, c.LastName, c.Company, c.Address1, c.Address2, c.City, c.State, c.Postcode, c.Country, c.PhoneCC, c.Phone, c.Currency, c.TaxExempt, c.Status, c.OAuthProvider, c.OAuthID, c.TwoFactorEnabled, c.TwoFactorSecret, c.Gender, c.Birthday, c.Type, c.Notes, c.BillingEmail, c.ReferrerID, c.ReferredBy, c.Custom1, c.Custom2, c.Custom3, c.Custom4, c.Custom5, c.Custom6, c.Custom7, c.Custom8, c.Custom9, c.Custom10, c.Custom11, c.Custom12, c.Custom13, c.Custom14, c.Custom15, c.Custom16, c.Custom17, c.Custom18, c.Custom19, c.Custom20, c.UpdatedAt, c.ID)
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
