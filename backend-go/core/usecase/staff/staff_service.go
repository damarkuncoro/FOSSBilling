package staff

import (
	"context"
	"errors"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/auth"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/security"
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
	Token             string             `json:"token,omitempty"`
	TwoFactorRequired bool               `json:"two_factor_required,omitempty"`
	Staff             *domain.Staff      `json:"staff,omitempty"`
	Group             *domain.AdminGroup `json:"group,omitempty"`
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

	if staff.TwoFactorEnabled {
		return &StaffAuthResponse{
			TwoFactorRequired: true,
			Staff:             staff,
		}, nil
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
		Staff: staff,
		Group: group,
	}, nil
}

func (s *StaffService) VerifyTwoFactor(ctx context.Context, email, code string) (*StaffAuthResponse, error) {
	staff, err := s.staffRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, appErrors.ErrUnauthorized
	}

	if staff.TwoFactorSecret == nil || !security.VerifyTOTP(*staff.TwoFactorSecret, code) {
		return nil, errors.New("invalid two-factor code")
	}

	token, err := auth.GenerateToken(s.jwtSecret, staff.ID, staff.Email, string(staff.Role), 12*time.Hour)
	if err != nil {
		return nil, err
	}

	group, _ := s.staffRepo.GetGroupByID(ctx, staff.GroupID)

	return &StaffAuthResponse{
		Token: token,
		Staff: staff,
		Group: group,
	}, nil
}

func (s *StaffService) SetupTwoFactor(ctx context.Context, staffID int64) (string, string, error) {
	staff, err := s.staffRepo.GetByID(ctx, staffID)
	if err != nil {
		return "", "", err
	}

	secret := security.GenerateTOTPSecret()
	qrURL := security.GenerateTOTPURL(staff.Email, "FOSSBilling-Admin", secret)

	staff.TwoFactorSecret = &secret
	err = s.staffRepo.Update(ctx, staff) // Need Update in StaffRepository
	return secret, qrURL, err
}

func (s *StaffService) EnableTwoFactor(ctx context.Context, staffID int64, code string) error {
	staff, err := s.staffRepo.GetByID(ctx, staffID)
	if err != nil {
		return err
	}

	if staff.TwoFactorSecret == nil || !security.VerifyTOTP(*staff.TwoFactorSecret, code) {
		return errors.New("invalid verification code")
	}

	staff.TwoFactorEnabled = true
	return s.staffRepo.Update(ctx, staff)
}

func (s *StaffService) DisableTwoFactor(ctx context.Context, staffID int64) error {
	staff, err := s.staffRepo.GetByID(ctx, staffID)
	if err != nil {
		return err
	}

	staff.TwoFactorEnabled = false
	staff.TwoFactorSecret = nil
	return s.staffRepo.Update(ctx, staff)
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
