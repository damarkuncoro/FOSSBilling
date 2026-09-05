package postgres

import (
	"context"
	"errors"

	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StaffRepository struct {
	pool *pgxpool.Pool
}

func NewStaffRepository(pool *pgxpool.Pool) *StaffRepository {
	return &StaffRepository{pool: pool}
}

func (r *StaffRepository) GetByID(ctx context.Context, id int64) (*domain.Staff, error) {
	query := `
		SELECT id, group_id, email, password_hash, name, role, status, created_at, updated_at
		FROM staff
		WHERE id = $1
	`
	var s domain.Staff
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.GroupID, &s.Email, &s.PasswordHash, &s.Name, &s.Role, &s.Status, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *StaffRepository) GetByEmail(ctx context.Context, email string) (*domain.Staff, error) {
	query := `
		SELECT id, group_id, email, password_hash, name, role, status, created_at, updated_at
		FROM staff
		WHERE LOWER(email) = LOWER($1)
	`
	var s domain.Staff
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&s.ID, &s.GroupID, &s.Email, &s.PasswordHash, &s.Name, &s.Role, &s.Status, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *StaffRepository) Create(ctx context.Context, s *domain.Staff) error {
	query := `
		INSERT INTO staff (group_id, email, password_hash, name, role, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at
	`
	if s.Role == "" {
		s.Role = domain.StaffRoleAdmin
	}
	if s.Status == "" {
		s.Status = "active"
	}

	return r.pool.QueryRow(ctx, query,
		s.GroupID, s.Email, s.PasswordHash, s.Name, s.Role, s.Status,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

func (r *StaffRepository) GetGroupByID(ctx context.Context, groupID int64) (*domain.AdminGroup, error) {
	query := `
		SELECT id, name, permissions, created_at, updated_at
		FROM admin_groups
		WHERE id = $1
	`
	var g domain.AdminGroup
	err := r.pool.QueryRow(ctx, query, groupID).Scan(
		&g.ID, &g.Name, &g.Permissions, &g.CreatedAt, &g.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, err
	}
	return &g, nil
}

func (r *StaffRepository) CreateGroup(ctx context.Context, group *domain.AdminGroup) error {
	query := `
		INSERT INTO admin_groups (name, permissions, created_at, updated_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at
	`
	return r.pool.QueryRow(ctx, query, group.Name, group.Permissions).Scan(&group.ID, &group.CreatedAt, &group.UpdatedAt)
}

func (r *StaffRepository) AddAuditLog(ctx context.Context, log *domain.AuditLog) error {
	query := `
		INSERT INTO audit_logs (staff_id, client_id, module, action, details, ip_address, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP)
		RETURNING id, created_at
	`
	return r.pool.QueryRow(ctx, query,
		log.StaffID, log.ClientID, log.Module, log.Action, log.Details, log.IPAddress,
	).Scan(&log.ID, &log.CreatedAt)
}

func (r *StaffRepository) ListAuditLogs(ctx context.Context, limit, offset int) ([]*domain.AuditLog, int, error) {
	countQuery := `SELECT COUNT(*) FROM audit_logs`
	var total int
	if err := r.pool.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, staff_id, client_id, module, action, details, ip_address, created_at
		FROM audit_logs
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []*domain.AuditLog
	for rows.Next() {
		var l domain.AuditLog
		if err := rows.Scan(
			&l.ID, &l.StaffID, &l.ClientID, &l.Module, &l.Action, &l.Details, &l.IPAddress, &l.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		logs = append(logs, &l)
	}

	return logs, total, nil
}
