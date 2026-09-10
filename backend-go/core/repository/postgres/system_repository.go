package postgres

import (
	"context"
	"encoding/json"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SystemRepository struct{ pool *pgxpool.Pool }

func NewSystemRepository(p *pgxpool.Pool) *SystemRepository { return &SystemRepository{p} }

func scanSet(r pgx.Row) (*domain.SystemSetting, error) {
	var s domain.SystemSetting; err := r.Scan(&s.ID, &s.Section, &s.Key, &s.Value, &s.UpdatedAt); return &s, err
}

func (r *SystemRepository) GetSetting(ctx context.Context, sec, k string) (*domain.SystemSetting, error) {
	return scanSet(r.pool.QueryRow(ctx, "SELECT id, section, key, value, updated_at FROM system_settings WHERE section = $1 AND key = $2", sec, k))
}

func (r *SystemRepository) ListSettings(ctx context.Context, sec string) ([]*domain.SystemSetting, error) {
	q := "SELECT id, section, key, value, updated_at FROM system_settings"; var args []any
	if sec != "" { q += " WHERE section = $1"; args = append(args, sec) }
	return list(ctx, r.pool, q+" ORDER BY section, key", scanSet, args...)
}

func (r *SystemRepository) UpdateSetting(ctx context.Context, sec, k string, v json.RawMessage) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO system_settings (section, key, value, updated_at) VALUES ($1, $2, $3, CURRENT_TIMESTAMP) ON CONFLICT (section, key) DO UPDATE SET value = EXCLUDED.value, updated_at = CURRENT_TIMESTAMP`, sec, k, v); return err
}

func (r *SystemRepository) GetDatabaseStats(ctx context.Context) (map[string]interface{}, error) {
	var sz, v string
	_ = r.pool.QueryRow(ctx, "SELECT pg_size_pretty(pg_database_size(current_database()))").Scan(&sz)
	_ = r.pool.QueryRow(ctx, "SHOW server_version").Scan(&v)
	return map[string]interface{}{"size": sz, "version": v, "type": "PostgreSQL"}, nil
}
