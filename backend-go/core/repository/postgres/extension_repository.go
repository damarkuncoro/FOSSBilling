package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ExtensionRepository struct {
	pool *pgxpool.Pool
}

func NewExtensionRepository(pool *pgxpool.Pool) *ExtensionRepository {
	return &ExtensionRepository{pool: pool}
}

func (r *ExtensionRepository) List(ctx context.Context, filter domain.ExtensionFilter) ([]*domain.Extension, error) {
	query := `
		SELECT id, name, type, version, description, author, author_url,
		       icon, status, has_settings, config, manifest, created_at, updated_at
		FROM extensions
		WHERE 1=1
	`
	var args []interface{}
	idx := 1

	if filter.Type != "" {
		query += fmt.Sprintf(" AND type = $%d", idx)
		args = append(args, filter.Type)
		idx++
	}
	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", idx)
		args = append(args, filter.Status)
		idx++
	}
	if filter.Active != nil {
		if *filter.Active {
			query += " AND status IN ('active', 'core')"
		} else {
			query += " AND status = 'inactive'"
		}
	}
	if filter.HasSettings != nil {
		query += fmt.Sprintf(" AND has_settings = $%d", idx)
		args = append(args, *filter.HasSettings)
		idx++
	}
	if filter.Search != "" {
		query += fmt.Sprintf(" AND (LOWER(name) LIKE $%d OR LOWER(description) LIKE $%d OR LOWER(id) LIKE $%d)", idx, idx, idx)
		term := "%" + strings.ToLower(filter.Search) + "%"
		args = append(args, term)
		idx++
	}

	query += " ORDER BY status ASC, name ASC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.Extension
	for rows.Next() {
		var ext domain.Extension
		var authorURL, icon *string
		var configRaw, manifestRaw []byte

		err := rows.Scan(
			&ext.ID,
			&ext.Name,
			&ext.Type,
			&ext.Version,
			&ext.Description,
			&ext.Author,
			&authorURL,
			&icon,
			&ext.Status,
			&ext.HasSettings,
			&configRaw,
			&manifestRaw,
			&ext.CreatedAt,
			&ext.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if authorURL != nil {
			ext.AuthorURL = *authorURL
		}
		if icon != nil {
			ext.Icon = *icon
		}
		if len(configRaw) > 0 {
			_ = json.Unmarshal(configRaw, &ext.Config)
		}
		if len(manifestRaw) > 0 {
			_ = json.Unmarshal(manifestRaw, &ext.Manifest)
		}

		list = append(list, &ext)
	}

	return list, nil
}

func (r *ExtensionRepository) GetByID(ctx context.Context, id string) (*domain.Extension, error) {
	query := `
		SELECT id, name, type, version, description, author, author_url,
		       icon, status, has_settings, config, manifest, created_at, updated_at
		FROM extensions
		WHERE id = $1
	`
	var ext domain.Extension
	var authorURL, icon *string
	var configRaw, manifestRaw []byte

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&ext.ID,
		&ext.Name,
		&ext.Type,
		&ext.Version,
		&ext.Description,
		&ext.Author,
		&authorURL,
		&icon,
		&ext.Status,
		&ext.HasSettings,
		&configRaw,
		&manifestRaw,
		&ext.CreatedAt,
		&ext.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, err
	}

	if authorURL != nil {
		ext.AuthorURL = *authorURL
	}
	if icon != nil {
		ext.Icon = *icon
	}
	if len(configRaw) > 0 {
		_ = json.Unmarshal(configRaw, &ext.Config)
	}
	if len(manifestRaw) > 0 {
		_ = json.Unmarshal(manifestRaw, &ext.Manifest)
	}

	return &ext, nil
}

func (r *ExtensionRepository) Create(ctx context.Context, ext *domain.Extension) error {
	configData, _ := json.Marshal(ext.Config)
	if ext.Config == nil {
		configData = []byte("{}")
	}
	manifestData, _ := json.Marshal(ext.Manifest)
	if ext.Manifest == nil {
		manifestData = []byte("{}")
	}

	query := `
		INSERT INTO extensions (
			id, name, type, version, description, author, author_url,
			icon, status, has_settings, config, manifest, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,
			CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		)
		RETURNING created_at, updated_at
	`
	return r.pool.QueryRow(
		ctx,
		query,
		ext.ID,
		ext.Name,
		ext.Type,
		ext.Version,
		ext.Description,
		ext.Author,
		ext.AuthorURL,
		ext.Icon,
		ext.Status,
		ext.HasSettings,
		configData,
		manifestData,
	).Scan(&ext.CreatedAt, &ext.UpdatedAt)
}

func (r *ExtensionRepository) Update(ctx context.Context, ext *domain.Extension) error {
	configData, _ := json.Marshal(ext.Config)
	if ext.Config == nil {
		configData = []byte("{}")
	}
	manifestData, _ := json.Marshal(ext.Manifest)
	if ext.Manifest == nil {
		manifestData = []byte("{}")
	}

	query := `
		UPDATE extensions
		SET name = $1, type = $2, version = $3, description = $4, author = $5,
		    author_url = $6, icon = $7, status = $8, has_settings = $9,
		    config = $10, manifest = $11, updated_at = CURRENT_TIMESTAMP
		WHERE id = $12
		RETURNING updated_at
	`
	err := r.pool.QueryRow(
		ctx,
		query,
		ext.Name,
		ext.Type,
		ext.Version,
		ext.Description,
		ext.Author,
		ext.AuthorURL,
		ext.Icon,
		ext.Status,
		ext.HasSettings,
		configData,
		manifestData,
		ext.ID,
	).Scan(&ext.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return appErrors.ErrNotFound
		}
		return err
	}
	return nil
}

func (r *ExtensionRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM extensions WHERE id = $1 AND status != 'core'`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return appErrors.ErrNotFound
	}
	return nil
}

func (r *ExtensionRepository) GetConfig(ctx context.Context, extID string) (map[string]interface{}, error) {
	query := `SELECT config FROM extensions WHERE id = $1`
	var configRaw []byte
	err := r.pool.QueryRow(ctx, query, extID).Scan(&configRaw)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, err
	}

	cfg := make(map[string]interface{})
	if len(configRaw) > 0 {
		_ = json.Unmarshal(configRaw, &cfg)
	}
	return cfg, nil
}

func (r *ExtensionRepository) UpdateConfig(ctx context.Context, extID string, config map[string]interface{}) error {
	configData, err := json.Marshal(config)
	if err != nil {
		return err
	}

	query := `
		UPDATE extensions
		SET config = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	tag, err := r.pool.Exec(ctx, query, configData, extID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return appErrors.ErrNotFound
	}
	return nil
}
