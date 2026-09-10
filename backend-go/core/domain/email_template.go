package domain

import (
	"context"
	"time"
)

type EmailTemplate struct {
	ID          int64     `json:"id" db:"id"`
	Code        string    `json:"code" db:"code"`
	Subject     string    `json:"subject" db:"subject"`
	Content     string    `json:"content" db:"content"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type EmailTemplateRepository interface {
	GetByCode(ctx context.Context, code string) (*EmailTemplate, error)
	List(ctx context.Context) ([]*EmailTemplate, error)
	Update(ctx context.Context, t *EmailTemplate) error
}
