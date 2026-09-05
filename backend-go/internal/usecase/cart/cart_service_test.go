package cart

import (
	"context"
	"testing"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/usecase/billing"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

func setupCartService() (*CartService, *memory.MockPromoRepository, *memory.MockOrderRepository, *memory.MockClientRepository) {
	promoRepo := memory.NewMockPromoRepository()
	orderRepo := memory.NewMockOrderRepository()
	invRepo := memory.NewMockInvoiceRepository()
	clientRepo := memory.NewMockClientRepository()

	taxCalc := billing.NewTaxCalculator(nil)
	invService := billing.NewInvoiceService(invRepo, clientRepo, taxCalc)
	promoCalc := NewPromoCalculator(promoRepo)

	cartService := NewCartService(promoCalc, promoRepo, orderRepo, invService)
	return cartService, promoRepo, orderRepo, clientRepo
}

func TestPromoCalculator_Discounts(t *testing.T) {
	promoRepo := memory.NewMockPromoRepository()
	calc := NewPromoCalculator(promoRepo)

	// 1. Percentage Discount (20% off $100.00 -> $20.00)
	pctPromo := &domain.Promo{
		Code:   "SAVE20",
		Type:   domain.PromoTypePercentage,
		Value:  decimal.FromFloat(20.00), // 20%
		Active: true,
	}
	subtotal100 := decimal.FromFloat(100.00)
	discountPct := calc.CalculateDiscount(subtotal100, pctPromo)
	if discountPct.String() != "20.00" {
		t.Errorf("20%% discount of $100 = %s; want 20.00", discountPct.String())
	}

	// 2. Absolute Fixed Discount ($15.00 off $50.00 -> $15.00)
	absPromo := &domain.Promo{
		Code:   "FLAT15",
		Type:   domain.PromoTypeAbsolute,
		Value:  decimal.FromFloat(15.00),
		Active: true,
	}
	subtotal50 := decimal.FromFloat(50.00)
	discountAbs := calc.CalculateDiscount(subtotal50, absPromo)
	if discountAbs.String() != "15.00" {
		t.Errorf("Flat discount of $15 on $50 = %s; want 15.00", discountAbs.String())
	}

	// 3. Absolute Fixed Discount exceeds subtotal ($100 off $30 -> capped at $30)
	bigPromo := &domain.Promo{
		Code:   "BIG100",
		Type:   domain.PromoTypeAbsolute,
		Value:  decimal.FromFloat(100.00),
		Active: true,
	}
	discountCapped := calc.CalculateDiscount(decimal.FromFloat(30.00), bigPromo)
	if discountCapped.String() != "30.00" {
		t.Errorf("Capped discount = %s; want 30.00", discountCapped.String())
	}
}

func TestPromoCalculator_ValidationRules(t *testing.T) {
	ctx := context.Background()
	promoRepo := memory.NewMockPromoRepository()
	calc := NewPromoCalculator(promoRepo)

	now := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)
	pastDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	futureDate := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)

	// 1. Inactive Promo
	inactivePromo := &domain.Promo{Code: "INACTIVE", Active: false}
	if err := calc.ValidatePromo(ctx, inactivePromo, 1, now); err != ErrPromoInactive {
		t.Errorf("Expected ErrPromoInactive, got: %v", err)
	}

	// 2. Expired Promo
	expiredPromo := &domain.Promo{Code: "EXPIRED", Active: true, EndDate: &pastDate}
	if err := calc.ValidatePromo(ctx, expiredPromo, 1, now); err != ErrPromoExpired {
		t.Errorf("Expected ErrPromoExpired, got: %v", err)
	}

	// 3. Max Uses Limit Exceeded
	maxUsesPromo := &domain.Promo{Code: "MAXED", Active: true, MaxUses: 5, UsedCount: 5}
	if err := calc.ValidatePromo(ctx, maxUsesPromo, 1, now); err != ErrPromoMaxUses {
		t.Errorf("Expected ErrPromoMaxUses, got: %v", err)
	}

	// 4. Once Per Client Rule
	oncePromo := &domain.Promo{
		ID:            10,
		Code:          "ONCE",
		Active:        true,
		OncePerClient: true,
		StartDate:     &pastDate,
		EndDate:       &futureDate,
	}
	_ = promoRepo.Create(ctx, oncePromo)

	// First time redemption -> valid
	if err := calc.ValidatePromo(ctx, oncePromo, 42, now); err != nil {
		t.Errorf("First use should be valid, got: %v", err)
	}

	// Record redemption for client 42
	_ = promoRepo.IncrementUsed(ctx, oncePromo.ID, 42, nil)

	// Second time redemption for client 42 -> blocked
	if err := calc.ValidatePromo(ctx, oncePromo, 42, now); err != ErrPromoAlreadyUsed {
		t.Errorf("Second use should return ErrPromoAlreadyUsed, got: %v", err)
	}
}

func TestCartService_CheckoutFlow(t *testing.T) {
	ctx := context.Background()
	service, promoRepo, _, clientRepo := setupCartService()

	// 1. Create Client
	client := &domain.Client{
		Email:     "shopper@example.com",
		FirstName: "Shopper",
		Country:   "US",
		Currency:  "USD",
	}
	_ = clientRepo.Create(ctx, client)

	// 2. Create Promo Coupon (20% off)
	promo := &domain.Promo{
		Code:   "HOST20",
		Type:   domain.PromoTypePercentage,
		Value:  decimal.FromFloat(20.00),
		Active: true,
	}
	_ = promoRepo.Create(ctx, promo)

	// 3. Create Cart with 2 items and Promo
	cart := &Cart{
		ClientID:  client.ID,
		PromoCode: "HOST20",
		Items: []CartItem{
			{ProductID: 1, Title: "Hosting Plan A", Period: "1M", Price: decimal.FromFloat(50.00), Quantity: 1},
			{ProductID: 2, Title: "Domain Registration", Period: "1Y", Price: decimal.FromFloat(15.00), Quantity: 1},
		},
	}

	// 4. Execute Checkout
	result, err := service.Checkout(ctx, cart)
	if err != nil {
		t.Fatalf("Checkout failed: %v", err)
	}

	// Subtotal = $65.00, Discount 20% = $13.00, Total = $52.00
	if len(result.Orders) != 2 {
		t.Fatalf("Expected 2 orders, got: %d", len(result.Orders))
	}
	if result.Orders[0].Status != domain.OrderStatusPendingSetup {
		t.Errorf("Order[0] status = %s; want pending_setup", result.Orders[0].Status)
	}

	invoice := result.Invoice
	if invoice.Subtotal.String() != "52.00" { // Subtotal after discount line item
		t.Errorf("Invoice Subtotal = %s; want 52.00", invoice.Subtotal.String())
	}
	if invoice.Status != domain.InvoiceStatusUnpaid {
		t.Errorf("Invoice status = %s; want unpaid", invoice.Status)
	}

	// Verify promo usage incremented
	updatedPromo, _ := promoRepo.GetByCode(ctx, "HOST20")
	if updatedPromo.UsedCount != 1 {
		t.Errorf("Promo used count = %d; want 1", updatedPromo.UsedCount)
	}
}
