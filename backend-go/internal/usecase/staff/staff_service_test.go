package staff

import (
	"context"
	"testing"

	"github.com/fossbilling/backend-go/internal/domain"
	"github.com/fossbilling/backend-go/internal/repository/memory"
	"github.com/fossbilling/backend-go/pkg/auth"
)

func setupStaffService() (*StaffService, *memory.MockStaffRepository) {
	repo := memory.NewMockStaffRepository()
	service := NewStaffService(repo, "staff-jwt-secret-key-32-bytes")
	return service, repo
}

func TestStaffService_LoginAndAuditLog(t *testing.T) {
	ctx := context.Background()
	service, repo := setupStaffService()

	hashedPassword, _ := auth.HashPassword("AdminPass!123")
	staff := &domain.Staff{
		GroupID:      1,
		Email:        "admin@fossbilling.org",
		PasswordHash: hashedPassword,
		Name:         "Super Admin",
		Role:         domain.StaffRoleSuperAdmin,
		Status:       "active",
	}
	_ = repo.Create(ctx, staff)

	res, err := service.Login(ctx, StaffLoginDTO{
		Email:    "admin@fossbilling.org",
		Password: "AdminPass!123",
	})
	if err != nil {
		t.Fatalf("Staff login failed: %v", err)
	}

	if res.Token == "" {
		t.Error("Expected JWT token from staff login")
	}
	if res.Staff.Name != "Super Admin" {
		t.Errorf("Staff.Name = %s; want Super Admin", res.Staff.Name)
	}

	// Verify Audit Log was recorded
	logs := repo.GetAuditLogs()
	if len(logs) != 1 {
		t.Fatalf("Expected 1 audit log, got: %d", len(logs))
	}
	if logs[0].Module != "staff" || logs[0].Action != "login" {
		t.Errorf("Log = %s:%s; want staff:login", logs[0].Module, logs[0].Action)
	}
}

func TestStaffService_RBACPermissions(t *testing.T) {
	ctx := context.Background()
	service, repo := setupStaffService()

	// 1. Create Support Group (Permissions: support: read, write; clients: read only)
	supportGroup := &domain.AdminGroup{
		Name: "Support Team",
		Permissions: map[string][]string{
			"support": {"read", "write"},
			"clients": {"read"},
		},
	}
	_ = repo.CreateGroup(ctx, supportGroup)

	// 2. Create Support Staff Member
	supportStaff := &domain.Staff{
		GroupID: supportGroup.ID,
		Email:   "support@fossbilling.org",
		Role:    domain.StaffRoleSupport,
		Status:  "active",
	}
	_ = repo.Create(ctx, supportStaff)

	// 3. Create Superadmin
	superAdmin := &domain.Staff{
		Email:  "root@fossbilling.org",
		Role:   domain.StaffRoleSuperAdmin,
		Status: "active",
	}
	_ = repo.Create(ctx, superAdmin)

	// Test Superadmin has all permissions
	hasPerm, _ := service.HasPermission(ctx, superAdmin.ID, "billing", "delete")
	if !hasPerm {
		t.Error("Superadmin should have permission on any module")
	}

	// Test Support Staff permissions:
	// a. Support write -> allowed
	canWriteTickets, _ := service.HasPermission(ctx, supportStaff.ID, "support", "write")
	if !canWriteTickets {
		t.Error("Support staff should have 'support:write' permission")
	}

	// b. Clients read -> allowed
	canReadClients, _ := service.HasPermission(ctx, supportStaff.ID, "clients", "read")
	if !canReadClients {
		t.Error("Support staff should have 'clients:read' permission")
	}

	// c. Clients delete -> denied
	canDeleteClients, _ := service.HasPermission(ctx, supportStaff.ID, "clients", "delete")
	if canDeleteClients {
		t.Error("Support staff should NOT have 'clients:delete' permission")
	}

	// d. Billing module -> denied
	canAccessBilling, _ := service.HasPermission(ctx, supportStaff.ID, "billing", "read")
	if canAccessBilling {
		t.Error("Support staff should NOT have 'billing:read' permission")
	}
}
