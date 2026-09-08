package domain

import (
	"context"
	"time"
)

type Notification struct {
	ID        int64     `json:"id"`
	ClientID  int64     `json:"client_id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Type      string    `json:"type"` // info, success, warning, error
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

type NotificationRepository interface {
	Create(ctx context.Context, n *Notification) error
	ListByClientID(ctx context.Context, clientID int64, limit, offset int) ([]*Notification, int, error)
	MarkAsRead(ctx context.Context, id int64) error
	MarkAllAsRead(ctx context.Context, clientID int64) error
	DeleteOld(ctx context.Context, days int) error
}

type AdminNotification struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Type      string    `json:"type"` // info, success, warning, danger
	Module    string    `json:"module,omitempty"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

type AdminNotificationRepository interface {
	CreateAdmin(ctx context.Context, n *AdminNotification) error
	ListAdmin(ctx context.Context, limit, offset int, unreadOnly bool) ([]*AdminNotification, int, error)
	MarkAdminAsRead(ctx context.Context, id int64) error
	MarkAdminAllAsRead(ctx context.Context) error
	DeleteAdminOld(ctx context.Context, days int) error
}
