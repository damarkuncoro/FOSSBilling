package staff

import (
	"context"
	"errors"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/auth"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
)

type StaffService struct {
	staffRepo domain.StaffRepository
	jwtSecret string
}

func NewStaffService(staffRepo domain.StaffRepository, jwtSecret string) *StaffService {
	return &StaffService{
		staffRepo: staffRepo,
		jwtSecret: jwtSecret,
	}
}

type StaffLoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type StaffAuthResponse struct {
	Token string        `json:"token"`
	Staff domain.Staff  `json:"staff"`
	Group *domain.AdminGroup `json:"group,omitempty"`
}

func (s *StaffService) Login(ctx context.Context, dto StaffLoginDTO) (*StaffAuthResponse, error) {
	staff, err := s.staffRepo.GetByEmail(ctx, dto.Email)
	if err != nil {
		return nil, appErrors.ErrUnauthorized
	}

	if !auth.CheckPassword(dto.Password, staff.PasswordHash) {
		return nil, appErrors.ErrUnauthorized
	}

	if staff.Status != "active" {
		return nil, errors.New("staff account is inactive")
	}

	token, err := auth.GenerateToken(s.jwtSecret, staff.ID, staff.Email, string(staff.Role), 12*time.Hour)
	if err != nil {
		return nil, err
	}

	group, _ := s.staffRepo.GetGroupByID(ctx, staff.GroupID)

	// Log audit trail
	_ = s.staffRepo.AddAuditLog(ctx, &domain.AuditLog{
		StaffID: &staff.ID,
		Module:  "staff",
		Action:  "login",
		Details: "Staff logged into admin dashboard",
	})

	return &StaffAuthResponse{
		Token: token,
		Staff: *staff,
		Group: group,
	}, nil
}

// HasPermission checks if staff has access to perform an action on a module
func (s *StaffService) HasPermission(ctx context.Context, staffID int64, module, action string) (bool, error) {
	staff, err := s.staffRepo.GetByID(ctx, staffID)
	if err != nil {
		return false, err
	}

	if staff.Role == domain.StaffRoleSuperAdmin {
		return true, nil
	}

	group, err := s.staffRepo.GetGroupByID(ctx, staff.GroupID)
	if err != nil {
		return false, err
	}

	actions, ok := group.Permissions[module]
	if !ok {
		return false, nil
	}

	for _, a := range actions {
		if a == "*" || a == action {
			return true, nil
		}
	}

	return false, nil
}

func (s *StaffService) ListAuditLogs(ctx context.Context, limit, offset int) ([]*domain.AuditLog, int, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.staffRepo.ListAuditLogs(ctx, limit, offset)
}

