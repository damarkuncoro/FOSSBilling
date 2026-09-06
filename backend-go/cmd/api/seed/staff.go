package seed

import (
	"context"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/auth"
)

// SeedStaff populates admin group and default administrator accounts
func SeedStaff(ctx context.Context, repo *memory.MockStaffRepository) {
	if repo == nil {
		return
	}

	group := &domain.AdminGroup{
		ID:   1,
		Name: "Super Administrators",
		Permissions: map[string][]string{
			"clients":    {"*"},
			"orders":     {"*"},
			"support":    {"*"},
			"system":     {"*"},
			"billing":    {"*"},
			"staff":      {"*"},
			"news":       {"*"},
			"currencies": {"*"},
		},
	}
	_ = repo.CreateGroup(ctx, group)

	adminPass, _ := auth.HashPassword("SuperSecretAdmin123!")
	_ = repo.Create(ctx, &domain.Staff{
		ID:           1,
		GroupID:      group.ID,
		Email:        "admin@fossbilling.org",
		PasswordHash: adminPass,
		Name:         "Super Administrator",
		Role:         domain.StaffRoleSuperAdmin,
		Status:       "active",
	})
}
