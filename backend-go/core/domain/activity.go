package domain

import (
	"context"
	"time"
)

type Activity struct {
	ID        int64     `json:"id"`
	ClientID  *int64    `json:"client_id,omitempty"`
	AdminID   *int64    `json:"admin_id,omitempty"`
	Type      string    `json:"type"` // info, warning, danger, success
	Event     string    `json:"event"`
	Message   string    `json:"message"`
	IPAddress string    `json:"ip_address,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type ActivityRepository interface {
	Log(ctx context.Context, activity *Activity) error
	List(ctx context.Context, limit, offset int) ([]*Activity, int, error)
	ListByClientID(ctx context.Context, clientID int64, limit, offset int) ([]*Activity, int, error)
	GetTrend(ctx context.Context, days int) (map[string]int, error)
	DeleteOld(ctx context.Context, days int) error
}
