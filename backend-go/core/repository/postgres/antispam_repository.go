package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AntispamRepository struct{ pool *pgxpool.Pool }

func NewAntispamRepository(p *pgxpool.Pool) *AntispamRepository { return &AntispamRepository{p} }

func (r *AntispamRepository) GetConfig(ctx context.Context) (*domain.AntispamConfig, error) {
	var raw string; err := r.pool.QueryRow(ctx, `SELECT value FROM system_settings WHERE section = 'antispam' AND key = 'config'`).Scan(&raw)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) { return &domain.AntispamConfig{CaptchaProvider: "turnstile", StopForumSpamEnabled: true, SFSMinConfidence: 80.0, TempEmailBlockEnabled: true, HoneypotEnabled: true, HoneypotFieldName: "website_hp"}, nil }
		return nil, err
	}
	var c domain.AntispamConfig; return &c, json.Unmarshal([]byte(raw), &c)
}

func (r *AntispamRepository) UpdateConfig(ctx context.Context, c *domain.AntispamConfig) error {
	d, _ := json.Marshal(c); _, err := r.pool.Exec(ctx, `INSERT INTO system_settings (section, key, value, updated_at) VALUES ('antispam', 'config', $1, CURRENT_TIMESTAMP) ON CONFLICT (section, key) DO UPDATE SET value = EXCLUDED.value, updated_at = CURRENT_TIMESTAMP`, string(d)); return err
}

func (r *AntispamRepository) ListBlockedIPs(ctx context.Context) ([]*domain.BlockedIP, error) {
	return list(ctx, r.pool, `SELECT id, ip, reason, created_at FROM blocked_ips ORDER BY id DESC`, func(r pgx.Row) (*domain.BlockedIP, error) {
		var b domain.BlockedIP; err := r.Scan(&b.ID, &b.IP, &b.Reason, &b.CreatedAt); return &b, err
	})
}

func (r *AntispamRepository) AddBlockedIP(ctx context.Context, ip, res string) (*domain.BlockedIP, error) {
	b := &domain.BlockedIP{IP: strings.TrimSpace(ip), Reason: res}
	err := r.pool.QueryRow(ctx, `INSERT INTO blocked_ips (ip, reason, created_at) VALUES ($1, $2, CURRENT_TIMESTAMP) ON CONFLICT (ip) DO UPDATE SET reason = EXCLUDED.reason RETURNING id, created_at`, b.IP, b.Reason).Scan(&b.ID, &b.CreatedAt)
	return b, err
}

func (r *AntispamRepository) DeleteBlockedIP(ctx context.Context, ip string) error {
	t, err := r.pool.Exec(ctx, `DELETE FROM blocked_ips WHERE ip = $1`, strings.TrimSpace(ip))
	if err == nil && t.RowsAffected() == 0 { return appErrors.ErrNotFound }; return err
}

func (r *AntispamRepository) IsIPBlocked(ctx context.Context, ip string) (bool, error) {
	var ex bool; err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM blocked_ips WHERE ip = $1)`, strings.TrimSpace(ip)).Scan(&ex); return ex, err
}
