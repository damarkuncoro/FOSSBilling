package postgres

import (
	"context"
	"errors"

	"github.com/fossbilling/backend-go/internal/domain"
	appErrors "github.com/fossbilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MassMailRepository struct {
	pool *pgxpool.Pool
}

func NewMassMailRepository(pool *pgxpool.Pool) *MassMailRepository {
	return &MassMailRepository{pool: pool}
}

func (r *MassMailRepository) GetByID(ctx context.Context, id int64) (*domain.MassMailCampaign, error) {
	query := `
		SELECT id, admin_id, subject, content, status, sent_count, created_at, sent_at
		FROM mass_mail_campaigns
		WHERE id = $1
	`
	var c domain.MassMailCampaign
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.AdminID, &c.Subject, &c.Content, &c.Status, &c.SentCount, &c.CreatedAt, &c.SentAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *MassMailRepository) List(ctx context.Context, limit, offset int) ([]*domain.MassMailCampaign, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM mass_mail_campaigns`).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, admin_id, subject, content, status, sent_count, created_at, sent_at
		FROM mass_mail_campaigns
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*domain.MassMailCampaign
	for rows.Next() {
		var c domain.MassMailCampaign
		if err := rows.Scan(&c.ID, &c.AdminID, &c.Subject, &c.Content, &c.Status, &c.SentCount, &c.CreatedAt, &c.SentAt); err != nil {
			return nil, 0, err
		}
		list = append(list, &c)
	}
	return list, total, rows.Err()
}

func (r *MassMailRepository) Create(ctx context.Context, c *domain.MassMailCampaign) error {
	query := `
		INSERT INTO mass_mail_campaigns (admin_id, subject, content, status, sent_count, created_at, sent_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, $6)
		RETURNING id, created_at
	`
	return r.pool.QueryRow(ctx, query,
		c.AdminID, c.Subject, c.Content, c.Status, c.SentCount, c.SentAt,
	).Scan(&c.ID, &c.CreatedAt)
}

func (r *MassMailRepository) Update(ctx context.Context, c *domain.MassMailCampaign) error {
	query := `
		UPDATE mass_mail_campaigns
		SET subject = $2, content = $3, status = $4, sent_count = $5, sent_at = $6
		WHERE id = $1
	`
	cmd, err := r.pool.Exec(ctx, query, c.ID, c.Subject, c.Content, c.Status, c.SentCount, c.SentAt)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return appErrors.ErrNotFound
	}
	return nil
}
