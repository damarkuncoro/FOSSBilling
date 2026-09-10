package system

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
}

type HealthUsecase struct {
	db    *pgxpool.Pool
	cache interface{}
}

func NewHealthUsecase(db *pgxpool.Pool, cache interface{}) *HealthUsecase {
	return &HealthUsecase{db: db, cache: cache}
}

func (u *HealthUsecase) Check(ctx context.Context) HealthStatus {
	res := HealthStatus{
		Status:    "ok",
		Timestamp: time.Now(),
		Services:  make(map[string]string),
	}

	// 1. Check Database
	if u.db != nil {
		if err := u.db.Ping(ctx); err != nil {
			res.Status = "degraded"
			res.Services["database"] = "unreachable: " + err.Error()
		} else {
			res.Services["database"] = "healthy"
		}
	} else {
		res.Services["database"] = "not_configured"
	}

	// 2. Check Cache
	if u.cache != nil {
		res.Services["cache"] = "active"
	} else {
		res.Services["cache"] = "in-memory"
	}

	// 3. Check App State
	res.Services["api"] = "running"

	return res
}
