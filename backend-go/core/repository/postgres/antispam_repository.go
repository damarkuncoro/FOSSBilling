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

type AntispamRepository struct {
	pool *pgxpool.Pool
}

func NewAntispamRepository(pool *pgxpool.Pool) *AntispamRepository {
	return &AntispamRepository{pool: pool}
}

func (r *AntispamRepository) GetConfig(ctx context.Context) (*domain.AntispamConfig, error) {
	query := `SELECT value FROM system_settings WHERE section = 'antispam' AND key = 'config'`
	var raw string
	err := r.pool.QueryRow(ctx, query).Scan(&raw)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Return default config if not yet set
			return &domain.AntispamConfig{
				CaptchaEnabled:        false,
				CaptchaProvider:       "turnstile",
				StopForumSpamEnabled:  true,
				SFSMinConfidence:      80.0,
				TempEmailBlockEnabled: true,
				HoneypotEnabled:       true,
				HoneypotFieldName:     "website_hp",
			}, nil
		}
		return nil, err
	}

	var config domain.AntispamConfig
	if err := json.Unmarshal([]byte(raw), &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *AntispamRepository) UpdateConfig(ctx context.Context, config *domain.AntispamConfig) error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO system_settings (section, key, value, updated_at)
		VALUES ('antispam', 'config', $1, CURRENT_TIMESTAMP)
		ON CONFLICT (section, key) DO UPDATE
		SET value = EXCLUDED.value, updated_at = CURRENT_TIMESTAMP
	`
	_, err = r.pool.Exec(ctx, query, string(data))
	return err
}

func (r *AntispamRepository) ListBlockedIPs(ctx context.Context) ([]*domain.BlockedIP, error) {
	query := `SELECT id, ip, reason, created_at FROM blocked_ips ORDER BY id DESC`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.BlockedIP
	for rows.Next() {
		b := &domain.BlockedIP{}
		if err := rows.Scan(&b.ID, &b.IP, &b.Reason, &b.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, b)
	}
	return list, nil
}

func (r *AntispamRepository) AddBlockedIP(ctx context.Context, ip, reason string) (*domain.BlockedIP, error) {
	cleanIP := strings.TrimSpace(ip)
	query := `
		INSERT INTO blocked_ips (ip, reason, created_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP)
		ON CONFLICT (ip) DO UPDATE
		SET reason = EXCLUDED.reason
		RETURNING id, created_at
	`
	b := &domain.BlockedIP{
		IP:     cleanIP,
		Reason: reason,
	}
	err := r.pool.QueryRow(ctx, query, cleanIP, reason).Scan(&b.ID, &b.CreatedAt)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (r *AntispamRepository) DeleteBlockedIP(ctx context.Context, ip string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM blocked_ips WHERE ip = $1`, strings.TrimSpace(ip))
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return appErrors.ErrNotFound
	}
	return nil
}

func (r *AntispamRepository) IsIPBlocked(ctx context.Context, ip string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM blocked_ips WHERE ip = $1)`
	err := r.pool.QueryRow(ctx, query, strings.TrimSpace(ip)).Scan(&exists)
	return exists, err
}
