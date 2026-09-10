package postgres

import (
	"context"
	"errors"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const mmCols = `id, admin_id, subject, content, status, sent_count, created_at, sent_at`

type MassMailRepository struct{ pool *pgxpool.Pool }

func NewMassMailRepository(p *pgxpool.Pool) *MassMailRepository { return &MassMailRepository{p} }

func scanMM(r pgx.Row) (*domain.MassMailCampaign, error) {
	var c domain.MassMailCampaign; err := r.Scan(&c.ID, &c.AdminID, &c.Subject, &c.Content, &c.Status, &c.SentCount, &c.CreatedAt, &c.SentAt); return &c, err
}

func (r *MassMailRepository) GetByID(ctx context.Context, id int64) (*domain.MassMailCampaign, error) {
	c, err := scanMM(r.pool.QueryRow(ctx, "SELECT "+mmCols+" FROM mass_mail_campaigns WHERE id = $1", id))
	if err != nil && errors.Is(err, pgx.ErrNoRows) { return nil, appErrors.ErrNotFound }; return c, err
}

func (r *MassMailRepository) List(ctx context.Context, l, o int) ([]*domain.MassMailCampaign, int, error) {
	res, err := list(ctx, r.pool, "SELECT "+mmCols+" FROM mass_mail_campaigns ORDER BY created_at DESC LIMIT $1 OFFSET $2", scanMM, l, o)
	return res, total(ctx, r.pool, "SELECT COUNT(*) FROM mass_mail_campaigns"), err
}

func (r *MassMailRepository) Create(ctx context.Context, c *domain.MassMailCampaign) error {
	return r.pool.QueryRow(ctx, `INSERT INTO mass_mail_campaigns (admin_id, subject, content, status, sent_count, created_at, sent_at) VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, $6) RETURNING id, created_at`, c.AdminID, c.Subject, c.Content, c.Status, c.SentCount, c.SentAt).Scan(&c.ID, &c.CreatedAt)
}

func (r *MassMailRepository) Update(ctx context.Context, c *domain.MassMailCampaign) error {
	t, err := r.pool.Exec(ctx, `UPDATE mass_mail_campaigns SET subject = $2, content = $3, status = $4, sent_count = $5, sent_at = $6 WHERE id = $1`, c.ID, c.Subject, c.Content, c.Status, c.SentCount, c.SentAt)
	if err == nil && t.RowsAffected() == 0 { return appErrors.ErrNotFound }; return err
}
