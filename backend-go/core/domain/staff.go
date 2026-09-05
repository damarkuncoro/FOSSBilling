package domain

import (
	"context"
	"time"
)

type StaffRole string

const (
	StaffRoleSuperAdmin StaffRole = "superadmin"
	StaffRoleAdmin      StaffRole = "admin"
	StaffRoleSupport    StaffRole = "support"
	StaffRoleBilling    StaffRole = "billing"
)

type Staff struct {
	ID           int64     `json:"id"`
	GroupID      int64     `json:"group_id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Name         string    `json:"name"`
	Role         StaffRole `json:"role"`
	Status       string    `json:"status"` // active, inactive
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type AdminGroup struct {
	ID          int64               `json:"id"`
	Name        string              `json:"name"`
	Permissions map[string][]string `json:"permissions"` // e.g. "clients": ["read", "write", "delete"]
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

type AuditLog struct {
	ID        int64     `json:"id"`
	StaffID   *int64    `json:"staff_id,omitempty"`
	ClientID  *int64    `json:"client_id,omitempty"`
	Module    string    `json:"module"`
	Action    string    `json:"action"`
	Details   string    `json:"details"`
	IPAddress string    `json:"ip_address,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type StaffRepository interface {
	GetByID(ctx context.Context, id int64) (*Staff, error)
	GetByEmail(ctx context.Context, email string) (*Staff, error)
	Create(ctx context.Context, staff *Staff) error
	GetGroupByID(ctx context.Context, groupID int64) (*AdminGroup, error)
	CreateGroup(ctx context.Context, group *AdminGroup) error
	AddAuditLog(ctx context.Context, log *AuditLog) error
	ListAuditLogs(ctx context.Context, limit, offset int) ([]*AuditLog, int, error)
}

