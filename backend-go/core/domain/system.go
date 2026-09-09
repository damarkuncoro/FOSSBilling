package domain

import (
	"context"
	"encoding/json"
	"time"
)

type SystemSetting struct {
	ID        int64           `json:"id"`
	Section   string          `json:"section"` // security, branding, localization, etc.
	Key       string          `json:"key"`
	Value     json.RawMessage `json:"value"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type SystemRepository interface {
	GetSetting(ctx context.Context, section, key string) (*SystemSetting, error)
	ListSettings(ctx context.Context, section string) ([]*SystemSetting, error)
	UpdateSetting(ctx context.Context, section, key string, value json.RawMessage) error
	GetDatabaseStats(ctx context.Context) (map[string]interface{}, error)
}

type SystemStatus struct {
	EngineVersion  string `json:"engine_version"`
	DatabaseType   string `json:"database_type"`
	DatabaseSize   string `json:"database_size"`
	ActiveSessions int    `json:"active_sessions"`
	CronLastRun    string `json:"cron_last_run"`
	CronStatus     string `json:"cron_status"`
	SystemLoad     string `json:"system_load"`
	MemoryUsage    string `json:"memory_usage"`
	Uptime         string `json:"uptime"`
}
