package order

import (
	"context"
	"testing"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

func setupOrderService() (*OrderService, *memory.MockOrderRepository) {
	repo := memory.NewMockOrderRepository()
	service := NewOrderService(repo)
	return service, repo
}

func TestOrderService_ActivationFlow(t *testing.T) {
	ctx := context.Background()
	service, repo := setupOrderService()

	order := &domain.Order{
		ClientID:  1,
		ProductID: 10,
		Period:    "1M",
		Price:     decimal.FromFloat(15.00),
		Currency:  "USD",
		Status:    domain.OrderStatusPendingSetup,
	}
	_ = repo.Create(ctx, order)

	baseDate := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	activated, err := service.Activate(ctx, order.ID, baseDate)
	if err != nil {
		t.Fatalf("Activate failed: %v", err)
	}

	if activated.Status != domain.OrderStatusActive {
		t.Errorf("Status = %s; want active", activated.Status)
	}

	expectedExpiry := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	if activated.ExpiresAt == nil || !activated.ExpiresAt.Equal(expectedExpiry) {
		t.Errorf("ExpiresAt = %v; want %v", activated.ExpiresAt, expectedExpiry)
	}
}

func TestOrderService_SuspensionAndUnsuspension(t *testing.T) {
	ctx := context.Background()
	service, repo := setupOrderService()

	order := &domain.Order{
		ClientID:  1,
		ProductID: 10,
		Period:    "1M",
		Status:    domain.OrderStatusPendingSetup,
	}
	_ = repo.Create(ctx, order)
	_, _ = service.Activate(ctx, order.ID, time.Now().UTC())

	// 1. Suspend active order
	suspended, err := service.Suspend(ctx, order.ID, "Payment overdue 7 days")
	if err != nil {
		t.Fatalf("Suspend failed: %v", err)
	}
	if suspended.Status != domain.OrderStatusSuspended {
		t.Errorf("Status = %s; want suspended", suspended.Status)
	}
	if suspended.SuspensionReason == nil || *suspended.SuspensionReason != "Payment overdue 7 days" {
		t.Errorf("SuspensionReason = %v; want 'Payment overdue 7 days'", suspended.SuspensionReason)
	}

	// 2. Unsuspend back to active
	unsuspended, err := service.Unsuspend(ctx, order.ID)
	if err != nil {
		t.Fatalf("Unsuspend failed: %v", err)
	}
	if unsuspended.Status != domain.OrderStatusActive {
		t.Errorf("Status = %s; want active", unsuspended.Status)
	}
	if unsuspended.SuspensionReason != nil {
		t.Errorf("Expected nil reason after unsuspend, got: %v", *unsuspended.SuspensionReason)
	}
}

func TestOrderService_Renew(t *testing.T) {
	ctx := context.Background()
	service, repo := setupOrderService()

	initExpiry := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	order := &domain.Order{
		ClientID:  1,
		ProductID: 10,
		Period:    "1Y",
		Status:    domain.OrderStatusActive,
		ExpiresAt: &initExpiry,
	}
	_ = repo.Create(ctx, order)

	renewed, err := service.Renew(ctx, order.ID)
	if err != nil {
		t.Fatalf("Renew failed: %v", err)
	}

	expectedNextExpiry := time.Date(2027, 5, 1, 0, 0, 0, 0, time.UTC)
	if renewed.ExpiresAt == nil || !renewed.ExpiresAt.Equal(expectedNextExpiry) {
		t.Errorf("Renewed ExpiresAt = %v; want %v", renewed.ExpiresAt, expectedNextExpiry)
	}
}

func TestOrderService_CheckGracePeriodOverdue(t *testing.T) {
	service, _ := setupOrderService()

	expiry := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
	order := &domain.Order{
		Status:    domain.OrderStatusActive,
		ExpiresAt: &expiry,
	}

	// 5 days grace period -> deadline is Jan 15
	gracePeriodDays := 5

	// Date: Jan 12 (Within grace period)
	nowWithin := time.Date(2026, 1, 12, 0, 0, 0, 0, time.UTC)
	if isOverdue := service.CheckGracePeriodOverdue(order, gracePeriodDays, nowWithin); isOverdue {
		t.Errorf("CheckGracePeriodOverdue(Jan 12) = true; want false")
	}

	// Date: Jan 16 (Past grace period)
	nowPast := time.Date(2026, 1, 16, 0, 0, 0, 0, time.UTC)
	if isOverdue := service.CheckGracePeriodOverdue(order, gracePeriodDays, nowPast); !isOverdue {
		t.Errorf("CheckGracePeriodOverdue(Jan 16) = false; want true")
	}
}
