package seed

import (
	"context"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/auth"
)

func SeedStaff(ctx context.Context, r *memory.MockStaffRepository) {
	if r == nil { return }
	g := &domain.AdminGroup{ID: 1, Name: "Super Administrators", Permissions: map[string][]string{"clients": {"*"}, "orders": {"*"}, "support": {"*"}, "system": {"*"}, "billing": {"*"}, "staff": {"*"}, "news": {"*"}, "currencies": {"*"}}}
	_ = r.CreateGroup(ctx, g); ap, _ := auth.HashPassword("SuperSecretAdmin123!")
	_ = r.Create(ctx, &domain.Staff{ID: 1, GroupID: g.ID, Email: "admin@fossbilling.org", PasswordHash: ap, Name: "Super Administrator", Role: domain.StaffRoleSuperAdmin, Status: "active"})
}
