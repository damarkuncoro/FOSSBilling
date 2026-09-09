package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SystemRepository struct {
	pool *pgxpool.Pool
}

func NewSystemRepository(pool *pgxpool.Pool) *SystemRepository {
	return &SystemRepository{pool: pool}
}

func (r *SystemRepository) GetSetting(ctx context.Context, section, key string) (*domain.SystemSetting, error) {
	query := `SELECT id, section, key, value, updated_at FROM system_settings WHERE section = $1 AND key = $2`
	s := &domain.SystemSetting{}
	err := r.pool.QueryRow(ctx, query, section, key).Scan(&s.ID, &s.Section, &s.Key, &s.Value, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *SystemRepository) ListSettings(ctx context.Context, section string) ([]*domain.SystemSetting, error) {
	query := `SELECT id, section, key, value, updated_at FROM system_settings`
	var args []interface{}
	if section != "" {
		query += ` WHERE section = $1`
		args = append(args, section)
	}
	query += ` ORDER BY section, key`

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []*domain.SystemSetting
	for rows.Next() {
		s := &domain.SystemSetting{}
		if err := rows.Scan(&s.ID, &s.Section, &s.Key, &s.Value, &s.UpdatedAt); err != nil {
			return nil, err
		}
		settings = append(settings, s)
	}
	return settings, nil
}

func (r *SystemRepository) UpdateSetting(ctx context.Context, section, key string, value json.RawMessage) error {
	query := `
		INSERT INTO system_settings (section, key, value, updated_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
		ON CONFLICT (section, key) DO UPDATE SET
			value = EXCLUDED.value,
			updated_at = CURRENT_TIMESTAMP
	`
	_, err := r.pool.Exec(ctx, query, section, key, value)
	if err != nil {
		return fmt.Errorf("failed to update system setting: %w", err)
	}
	return nil
}

func (r *SystemRepository) GetDatabaseStats(ctx context.Context) (map[string]interface{}, error) {
	var dbSize string
	err := r.pool.QueryRow(ctx, "SELECT pg_size_pretty(pg_database_size(current_database()))").Scan(&dbSize)
	if err != nil {
		return nil, err
	}

	var version string
	err = r.pool.QueryRow(ctx, "SHOW server_version").Scan(&version)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"size":    dbSize,
		"version": version,
		"type":    "PostgreSQL",
	}, nil
}
