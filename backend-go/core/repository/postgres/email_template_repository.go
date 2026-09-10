package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EmailTemplateRepository struct {
	pool *pgxpool.Pool
}

func NewEmailTemplateRepository(pool *pgxpool.Pool) *EmailTemplateRepository {
	return &EmailTemplateRepository{pool: pool}
}

func (r *EmailTemplateRepository) GetByCode(ctx context.Context, code string) (*domain.EmailTemplate, error) {
	q := `SELECT id, code, subject, content, description, created_at, updated_at FROM email_templates WHERE code = $1 LIMIT 1`
	var t domain.EmailTemplate
	err := r.pool.QueryRow(ctx, q, code).Scan(&t.ID, &t.Code, &t.Subject, &t.Content, &t.Description, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *EmailTemplateRepository) List(ctx context.Context) ([]*domain.EmailTemplate, error) {
	q := `SELECT id, code, subject, content, description, created_at, updated_at FROM email_templates ORDER BY code ASC`
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []*domain.EmailTemplate
	for rows.Next() {
		t := &domain.EmailTemplate{}
		if err := rows.Scan(&t.ID, &t.Code, &t.Subject, &t.Content, &t.Description, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		templates = append(templates, t)
	}
	return templates, nil
}

func (r *EmailTemplateRepository) Update(ctx context.Context, t *domain.EmailTemplate) error {
	q := `UPDATE email_templates SET subject = $1, content = $2, description = $3, updated_at = $4 WHERE id = $5`
	_, err := r.pool.Exec(ctx, q, t.Subject, t.Content, t.Description, time.Now().UTC(), t.ID)
	if err != nil {
		return fmt.Errorf("failed to update email template: %w", err)
	}
	return nil
}
