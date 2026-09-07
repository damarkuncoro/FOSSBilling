package system

import (
	"context"
	"encoding/json"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type SystemService struct {
	systemRepo domain.SystemRepository
}

func NewSystemService(systemRepo domain.SystemRepository) *SystemService {
	return &SystemService{
		systemRepo: systemRepo,
	}
}

func (s *SystemService) GetSecuritySettings(ctx context.Context) (map[string]interface{}, error) {
	settings, err := s.systemRepo.ListSettings(ctx, "security")
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	for _, setting := range settings {
		var val interface{}
		_ = json.Unmarshal(setting.Value, &val)
		result[setting.Key] = val
	}

	// Defaults if empty
	if len(result) == 0 {
		return map[string]interface{}{
			"recaptcha_enabled": false,
			"max_login_attempts": 5,
			"force_ssl": true,
		}, nil
	}

	return result, nil
}

func (s *SystemService) UpdateSecuritySetting(ctx context.Context, key string, value interface{}) error {
	valJSON, _ := json.Marshal(value)
	return s.systemRepo.UpdateSetting(ctx, "security", key, valJSON)
}

func (s *SystemService) GetSystemStatus(ctx context.Context) *domain.SystemStatus {
	// In a real app, this would query OS stats and DB stats
	return &domain.SystemStatus{
		EngineVersion:  "v0.7.0-NextGen (Go 1.22)",
		DatabaseType:   "PostgreSQL 16",
		DatabaseSize:   "24.5 MB",
		ActiveSessions: 8,
		CronLastRun:    "5 minutes ago",
		CronStatus:     "healthy",
		SystemLoad:     "0.18, 0.22, 0.15",
		MemoryUsage:    "142 MB / 8 GB",
		Uptime:         "14 days",
	}
}
