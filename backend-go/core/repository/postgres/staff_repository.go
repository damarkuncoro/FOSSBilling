package postgres

import (
	"context"
	"errors"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const stCols = `id, group_id, email, password_hash, name, role, status, two_factor_enabled, two_factor_secret, created_at, updated_at`

type StaffRepository struct{ pool *pgxpool.Pool }

func NewStaffRepository(p *pgxpool.Pool) *StaffRepository { return &StaffRepository{p} }

func scanStaff(r pgx.Row) (*domain.Staff, error) {
	var s domain.Staff; err := r.Scan(&s.ID, &s.GroupID, &s.Email, &s.PasswordHash, &s.Name, &s.Role, &s.Status, &s.TwoFactorEnabled, &s.TwoFactorSecret, &s.CreatedAt, &s.UpdatedAt); return &s, err
}

func (r *StaffRepository) GetByID(ctx context.Context, id int64) (*domain.Staff, error) {
	s, err := scanStaff(r.pool.QueryRow(ctx, "SELECT "+stCols+" FROM staff WHERE id = $1", id))
	if err != nil && errors.Is(err, pgx.ErrNoRows) { return nil, appErrors.ErrNotFound }; return s, err
}

func (r *StaffRepository) GetByEmail(ctx context.Context, em string) (*domain.Staff, error) {
	s, err := scanStaff(r.pool.QueryRow(ctx, "SELECT "+stCols+" FROM staff WHERE LOWER(email) = LOWER($1)", em))
	if err != nil && errors.Is(err, pgx.ErrNoRows) { return nil, appErrors.ErrNotFound }; return s, err
}

func (r *StaffRepository) List(ctx context.Context, l, o int) ([]*domain.Staff, int, error) {
	res, err := list(ctx, r.pool, "SELECT "+stCols+" FROM staff ORDER BY id ASC LIMIT $1 OFFSET $2", scanStaff, l, o)
	return res, total(ctx, r.pool, "SELECT COUNT(*) FROM staff"), err
}

func (r *StaffRepository) Create(ctx context.Context, s *domain.Staff) error {
	if s.Role == "" { s.Role = domain.StaffRoleAdmin }; if s.Status == "" { s.Status = "active" }
	return r.pool.QueryRow(ctx, `INSERT INTO staff (group_id, email, password_hash, name, role, status, two_factor_enabled, two_factor_secret, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id, created_at, updated_at`, s.GroupID, s.Email, s.PasswordHash, s.Name, s.Role, s.Status, s.TwoFactorEnabled, s.TwoFactorSecret).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

func (r *StaffRepository) Update(ctx context.Context, s *domain.Staff) error {
	_, err := r.pool.Exec(ctx, `UPDATE staff SET group_id = $1, email = $2, name = $3, role = $4, status = $5, two_factor_enabled = $6, two_factor_secret = $7, updated_at = CURRENT_TIMESTAMP WHERE id = $8`, s.GroupID, s.Email, s.Name, s.Role, s.Status, s.TwoFactorEnabled, s.TwoFactorSecret, s.ID)
	return err
}

func (r *StaffRepository) GetGroupByID(ctx context.Context, id int64) (*domain.AdminGroup, error) {
	var g domain.AdminGroup; err := r.pool.QueryRow(ctx, `SELECT id, name, permissions, created_at, updated_at FROM admin_groups WHERE id = $1`, id).Scan(&g.ID, &g.Name, &g.Permissions, &g.CreatedAt, &g.UpdatedAt)
	if err != nil && errors.Is(err, pgx.ErrNoRows) { return nil, appErrors.ErrNotFound }; return &g, err
}

func (r *StaffRepository) CreateGroup(ctx context.Context, g *domain.AdminGroup) error {
	return r.pool.QueryRow(ctx, `INSERT INTO admin_groups (name, permissions, created_at, updated_at) VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id, created_at, updated_at`, g.Name, g.Permissions).Scan(&g.ID, &g.CreatedAt, &g.UpdatedAt)
}

func (r *StaffRepository) AddAuditLog(ctx context.Context, l *domain.AuditLog) error {
	return r.pool.QueryRow(ctx, `INSERT INTO audit_logs (staff_id, client_id, module, action, details, ip_address, created_at) VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP) RETURNING id, created_at`, l.StaffID, l.ClientID, l.Module, l.Action, l.Details, l.IPAddress).Scan(&l.ID, &l.CreatedAt)
}

func (r *StaffRepository) ListAuditLogs(ctx context.Context, l, o int) ([]*domain.AuditLog, int, error) {
	res, err := list(ctx, r.pool, `SELECT a.id, a.staff_id, a.client_id, a.module, a.action, a.details, a.ip_address, a.created_at, s.name as staff_name FROM audit_logs a LEFT JOIN staff s ON a.staff_id = s.id ORDER BY a.id DESC LIMIT $1 OFFSET $2`, func(r pgx.Row) (*domain.AuditLog, error) {
		var l domain.AuditLog; err := r.Scan(&l.ID, &l.StaffID, &l.ClientID, &l.Module, &l.Action, &l.Details, &l.IPAddress, &l.CreatedAt, &l.StaffName); return &l, err
	}, l, o)
	return res, total(ctx, r.pool, "SELECT COUNT(*) FROM audit_logs"), err
}
