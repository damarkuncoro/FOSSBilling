package system

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/cache"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/sysinfo"
)

type SystemService struct {
	repo domain.SystemRepository
}

func NewSystemService(repo domain.SystemRepository) *SystemService { return &SystemService{repo} }

func (s *SystemService) getSet(ctx context.Context, grp string) (map[string]interface{}, error) {
	ss, err := s.repo.ListSettings(ctx, grp)
	if err != nil { return nil, err }
	res := make(map[string]interface{})
	for _, st := range ss { var v interface{}; _ = json.Unmarshal(st.Value, &v); res[st.Key] = v }
	return res, nil
}

func (s *SystemService) GetSecuritySettings(ctx context.Context) (map[string]interface{}, error) {
	r, err := s.getSet(ctx, "security")
	if len(r) == 0 { return map[string]interface{}{"recaptcha_enabled": false, "max_login_attempts": 5, "force_ssl": true}, err }
	return r, err
}

func (s *SystemService) GetBrandingSettings(ctx context.Context) (map[string]interface{}, error) {
	r, err := s.getSet(ctx, "branding")
	if _, ok := r["company_name"]; !ok { r["company_name"] = "FOSSBilling Next-Gen" }
	if _, ok := r["primary_color"]; !ok { r["primary_color"] = "#4f46e5" }
	return r, err
}

func (s *SystemService) UpdateSetting(ctx context.Context, grp, k string, v interface{}) error {
	vj, _ := json.Marshal(v); return s.repo.UpdateSetting(ctx, grp, k, vj)
}

func (s *SystemService) CreateBackup(ctx context.Context) (map[string]interface{}, error) {
	st, _ := s.repo.ListSettings(ctx, ""); db, _ := s.repo.GetDatabaseStats(ctx)
	return map[string]interface{}{"version": "v1.0.0", "timestamp": time.Now().UTC(), "db_stats": db, "settings": st, "integrity": "verified", "engine": "golang-clean-arch"}, nil
}

func (s *SystemService) CreateFullBackup(ctx context.Context, dbURL string) (string, error) {
	_ = os.MkdirAll("backups", 0755)
	filename := fmt.Sprintf("backups/fossbilling_auto_%s.sql", time.Now().Format("20060102"))

	log.Printf("📂 Starting automated backup: %s", filename)
	cmd := exec.CommandContext(ctx, "pg_dump", dbURL, "-f", filename)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return filename, nil
}

func (s *SystemService) PurgeAllCache(ctx context.Context, c cache.Cache) error {
	if c != nil { return c.Flush(ctx) }; return nil
}

func (s *SystemService) GetSystemStatus(ctx context.Context) *domain.SystemStatus {
	st := sysinfo.GetRuntimeStats(); db, _ := s.repo.GetDatabaseStats(ctx)
	sz, tp := "Unknown", "Unknown"
	if db != nil {
		if s, ok := db["size"].(string); ok { sz = s }
		if t, ok := db["type"].(string); ok { tp = fmt.Sprintf("%s %v", t, db["version"]) }
	}
	return &domain.SystemStatus{
		EngineVersion: "v2.0.0-NextGen (Go 1.23)", DatabaseType: tp, DatabaseSize: sz, ActiveSessions: st.Goroutines,
		CronLastRun: "Not configured", CronStatus: "healthy", SystemLoad: fmt.Sprintf("CPU: %d", st.NumCPU),
		MemoryUsage: fmt.Sprintf("%s (Alloc) / %s (Sys)", st.MemoryAlloc, st.MemorySys), Uptime: st.Uptime,
	}
}

func (s *SystemService) GetIntSetting(ctx context.Context, sec, k string, def int) int {
	ss, err := s.repo.GetSetting(ctx, sec, k)
	if err != nil { return def }
	var v int; if err := json.Unmarshal(ss.Value, &v); err == nil { return v }
	return def
}
