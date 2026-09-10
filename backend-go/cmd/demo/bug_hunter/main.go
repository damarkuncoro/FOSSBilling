package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/http/admin"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/provisioning"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/order"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
)

func main() {
	ctx := context.Background()
	fmt.Println("==================================================================")
	fmt.Println("🕷️  FOSSBilling Next-Gen - BUG HUNTER SESSION 7 (PRIVILEGE ESCALATION)")
	fmt.Println("==================================================================")

	// Setup Infrastructure
	staffRepo := memory.NewMockStaffRepository()
	orderRepo := memory.NewMockOrderRepository()
	productRepo := memory.NewMockProductRepository()
	eventBus := events.NewEventBus()

	staffUc := staff.NewStaffService(staffRepo, "secret")
	orderUc := order.NewOrderService(orderRepo, productRepo, provisioning.NewProvisionerRegistry(), provisioning.NewRegistrarRegistry(), eventBus)

	handler := admin.NewStaffManagementHandler(staffUc, nil, orderRepo, orderUc, nil)

	// 1. Create a "Support" staff with NO 'write' permissions for orders
	_ = staffRepo.CreateGroup(ctx, &domain.AdminGroup{
		ID: 2, Name: "Support",
		Permissions: map[string][]string{"orders": {"read"}}, // READ ONLY
	})

	supportStaff := &domain.Staff{
		ID: 10, GroupID: 2, Name: "Support Guy", Email: "support@test.com", Role: domain.StaffRoleSupport, Status: "active",
	}
	_ = staffRepo.Create(ctx, supportStaff)

	// 2. Create a Suspended Order
	testOrder := &domain.Order{
		ID: 1, ClientID: 1, Status: domain.OrderStatusSuspended, Title: "Victim VPS", Price: decimal.FromFloat(10.0), Period: "1M",
	}
	_ = orderRepo.Create(ctx, testOrder)

	// --- BUG TEST 16: PRIVILEGE ESCALATION (UNSUSPEND) ---
	fmt.Println("\n[BUG TEST 16] 🔑 Privilege Escalation: Staff with READ-ONLY access unsuspending order")

	// Simulate request from Support Staff
	req, _ := http.NewRequest("POST", "/api/v1/admin/orders/1/unsuspend", nil)
	// Manually set path value
	req.SetPathValue("id", "1")

	// Inject StaffID into context
	req = req.WithContext(middleware.WithClientID(req.Context(), supportStaff.ID))

	rr := httptest.NewRecorder()
	handler.UnsuspendOrder(rr, req)

	if rr.Code == http.StatusOK {
		updatedOrder, _ := orderRepo.GetByID(ctx, 1)
		fmt.Printf("   ❌ BUG FOUND: Staff without write permissions successfully unsuspended the order!\n")
		fmt.Printf("      Order Status: %s\n", updatedOrder.Status)
	} else {
		fmt.Printf("   ✅ Success: System blocked unauthorized action: %d\n", rr.Code)
	}

	// --- BUG TEST 17: MISSING AUTH ON SOME ENDPOINTS ---
	fmt.Println("\n[BUG TEST 17] 🕵️  Missing Permission Check: ActivateOrder")
	rr2 := httptest.NewRecorder()
	handler.ActivateOrder(rr2, req)

	if rr2.Code == http.StatusOK {
		fmt.Println("   ❌ BUG FOUND: Missing permission check on ActivateOrder endpoint.")
	} else {
		fmt.Printf("   ✅ Success: Blocked. Code: %d\n", rr2.Code)
	}

	fmt.Println("\n🎉 BUG HUNTING SESSION 7 COMPLETED")
}
