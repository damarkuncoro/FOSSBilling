package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FormbuilderRepository struct {
	pool *pgxpool.Pool
}

func NewFormbuilderRepository(pool *pgxpool.Pool) *FormbuilderRepository {
	return &FormbuilderRepository{pool: pool}
}

func (r *FormbuilderRepository) GetFormByID(ctx context.Context, id int64) (*domain.Form, error) {
	query := `SELECT id, name, style, created_at, updated_at FROM forms WHERE id = $1`
	var form domain.Form
	var styleRaw []byte

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&form.ID,
		&form.Name,
		&styleRaw,
		&form.CreatedAt,
		&form.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, err
	}

	if len(styleRaw) > 0 {
		_ = json.Unmarshal(styleRaw, &form.Style)
	}

	// Fetch fields
	fields, err := r.GetFieldsByFormID(ctx, id)
	if err == nil {
		form.Fields = fields
	}

	return &form, nil
}

func (r *FormbuilderRepository) ListForms(ctx context.Context, limit, offset int) ([]*domain.Form, int, error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM forms`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, name, style, created_at, updated_at
		FROM forms
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*domain.Form
	for rows.Next() {
		var form domain.Form
		var styleRaw []byte
		if err := rows.Scan(&form.ID, &form.Name, &styleRaw, &form.CreatedAt, &form.UpdatedAt); err != nil {
			return nil, 0, err
		}
		if len(styleRaw) > 0 {
			_ = json.Unmarshal(styleRaw, &form.Style)
		}
		list = append(list, &form)
	}

	// Attach fields to listed forms
	for _, form := range list {
		fields, err := r.GetFieldsByFormID(ctx, form.ID)
		if err == nil {
			form.Fields = fields
		}
	}

	return list, total, nil
}

func (r *FormbuilderRepository) CreateForm(ctx context.Context, form *domain.Form) error {
	styleData, err := json.Marshal(form.Style)
	if err != nil {
		styleData = []byte("{}")
	}

	query := `
		INSERT INTO forms (name, style, created_at, updated_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at
	`
	return r.pool.QueryRow(ctx, query, form.Name, styleData).Scan(&form.ID, &form.CreatedAt, &form.UpdatedAt)
}

func (r *FormbuilderRepository) UpdateForm(ctx context.Context, form *domain.Form) error {
	styleData, err := json.Marshal(form.Style)
	if err != nil {
		styleData = []byte("{}")
	}

	query := `
		UPDATE forms
		SET name = $1, style = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		RETURNING updated_at
	`
	err = r.pool.QueryRow(ctx, query, form.Name, styleData, form.ID).Scan(&form.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return appErrors.ErrNotFound
		}
		return err
	}
	return nil
}

func (r *FormbuilderRepository) DeleteForm(ctx context.Context, id int64) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM forms WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return appErrors.ErrNotFound
	}
	return nil
}

func (r *FormbuilderRepository) GetFieldByID(ctx context.Context, fieldID int64) (*domain.FormField, error) {
	query := `
		SELECT id, form_id, name, label, hide_label, description, type,
		       default_value, required, hidden, readonly, options, prefix,
		       suffix, text_size, created_at, updated_at
		FROM form_fields
		WHERE id = $1
	`
	var f domain.FormField
	var optionsRaw []byte
	var desc, defVal, pfx, sfx *string

	err := r.pool.QueryRow(ctx, query, fieldID).Scan(
		&f.ID,
		&f.FormID,
		&f.Name,
		&f.Label,
		&f.HideLabel,
		&desc,
		&f.Type,
		&defVal,
		&f.Required,
		&f.Hidden,
		&f.Readonly,
		&optionsRaw,
		&pfx,
		&sfx,
		&f.TextSize,
		&f.CreatedAt,
		&f.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, err
	}

	if desc != nil {
		f.Description = *desc
	}
	if defVal != nil {
		f.DefaultValue = *defVal
	}
	if pfx != nil {
		f.Prefix = *pfx
	}
	if sfx != nil {
		f.Suffix = *sfx
	}
	if len(optionsRaw) > 0 {
		_ = json.Unmarshal(optionsRaw, &f.Options)
	}

	return &f, nil
}

func (r *FormbuilderRepository) GetFieldsByFormID(ctx context.Context, formID int64) ([]*domain.FormField, error) {
	query := `
		SELECT id, form_id, name, label, hide_label, description, type,
		       default_value, required, hidden, readonly, options, prefix,
		       suffix, text_size, created_at, updated_at
		FROM form_fields
		WHERE form_id = $1
		ORDER BY id ASC
	`
	rows, err := r.pool.Query(ctx, query, formID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fields []*domain.FormField
	for rows.Next() {
		var f domain.FormField
		var optionsRaw []byte
		var desc, defVal, pfx, sfx *string

		err := rows.Scan(
			&f.ID,
			&f.FormID,
			&f.Name,
			&f.Label,
			&f.HideLabel,
			&desc,
			&f.Type,
			&defVal,
			&f.Required,
			&f.Hidden,
			&f.Readonly,
			&optionsRaw,
			&pfx,
			&sfx,
			&f.TextSize,
			&f.CreatedAt,
			&f.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if desc != nil {
			f.Description = *desc
		}
		if defVal != nil {
			f.DefaultValue = *defVal
		}
		if pfx != nil {
			f.Prefix = *pfx
		}
		if sfx != nil {
			f.Suffix = *sfx
		}
		if len(optionsRaw) > 0 {
			_ = json.Unmarshal(optionsRaw, &f.Options)
		}

		fields = append(fields, &f)
	}

	return fields, nil
}

func (r *FormbuilderRepository) AddField(ctx context.Context, field *domain.FormField) error {
	optionsData, err := json.Marshal(field.Options)
	if err != nil || field.Options == nil {
		optionsData = []byte("{}")
	}

	query := `
		INSERT INTO form_fields (
			form_id, name, label, hide_label, description, type,
			default_value, required, hidden, readonly, options, prefix,
			suffix, text_size, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		)
		RETURNING id, created_at, updated_at
	`
	return r.pool.QueryRow(
		ctx,
		query,
		field.FormID,
		field.Name,
		field.Label,
		field.HideLabel,
		field.Description,
		field.Type,
		field.DefaultValue,
		field.Required,
		field.Hidden,
		field.Readonly,
		optionsData,
		field.Prefix,
		field.Suffix,
		field.TextSize,
	).Scan(&field.ID, &field.CreatedAt, &field.UpdatedAt)
}

func (r *FormbuilderRepository) UpdateField(ctx context.Context, field *domain.FormField) error {
	optionsData, err := json.Marshal(field.Options)
	if err != nil || field.Options == nil {
		optionsData = []byte("{}")
	}

	query := `
		UPDATE form_fields
		SET name = $1, label = $2, hide_label = $3, description = $4,
		    type = $5, default_value = $6, required = $7, hidden = $8,
		    readonly = $9, options = $10, prefix = $11, suffix = $12,
		    text_size = $13, updated_at = CURRENT_TIMESTAMP
		WHERE id = $14
		RETURNING updated_at
	`
	var updatedAt time.Time
	err = r.pool.QueryRow(
		ctx,
		query,
		field.Name,
		field.Label,
		field.HideLabel,
		field.Description,
		field.Type,
		field.DefaultValue,
		field.Required,
		field.Hidden,
		field.Readonly,
		optionsData,
		field.Prefix,
		field.Suffix,
		field.TextSize,
		field.ID,
	).Scan(&updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return appErrors.ErrNotFound
		}
		return err
	}
	field.UpdatedAt = updatedAt
	return nil
}

func (r *FormbuilderRepository) DeleteField(ctx context.Context, fieldID int64) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM form_fields WHERE id = $1`, fieldID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return appErrors.ErrNotFound
	}
	return nil
}

func (r *FormbuilderRepository) CountFields(ctx context.Context, formID int64) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM form_fields WHERE form_id = $1`, formID).Scan(&count)
	return count, err
}
