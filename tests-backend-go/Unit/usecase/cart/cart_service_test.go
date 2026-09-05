package cart_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/billing"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/cart"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

func setupCartService() (*cart.CartService, *memory.MockPromoRepository, *memory.MockOrderRepository, *memory.MockClientRepository) {
	promoRepo := memory.NewMockPromoRepository()
	orderRepo := memory.NewMockOrderRepository()
	invRepo := memory.NewMockInvoiceRepository()
	clientRepo := memory.NewMockClientRepository()

	taxCalc := billing.NewTaxCalculator(nil)
	invService := billing.NewInvoiceService(invRepo, clientRepo, taxCalc)
	promoCalc := cart.NewPromoCalculator(promoRepo)

	cartService := cart.NewCartService(promoCalc, promoRepo, orderRepo, invService)
	return cartService, promoRepo, orderRepo, clientRepo
}

func TestPromoCalculator_Discounts(t *testing.T) {
	promoRepo := memory.NewMockPromoRepository()
	calc := cart.NewPromoCalculator(promoRepo)

	// 1. Percentage discount
	p1 := &domain.Promo{
		Code:   "DISC10",
		Type:   domain.PromoTypePercentage,
		Value:  decimal.FromFloat(10.0),
		Active: true,
	}

	subtotal := decimal.FromFloat(100.00)
	discount := calc.CalculateDiscount(subtotal, p1)
	if discount.String() != "10.00" {
		t.Errorf("got discount %s, want 10.00", discount.String())
	}

	// 2. Absolute discount
	p2 := &domain.Promo{
		Code:   "FLAT15",
		Type:   domain.PromoTypeAbsolute,
		Value:  decimal.FromFloat(15.0),
		Active: true,
	}

	discount = calc.CalculateDiscount(subtotal, p2)
	if discount.String() != "15.00" {
		t.Errorf("got discount %s, want 15.00", discount.String())
	}
}

func TestPromoCalculator_ValidationRules(t *testing.T) {
	promoRepo := memory.NewMockPromoRepository()
	calc := cart.NewPromoCalculator(promoRepo)
	ctx := context.Background()
	now := time.Now().UTC()

	// Inactive promo
	pInactive := &domain.Promo{Code: "INACTIVE", Active: false}
	err := calc.ValidatePromo(ctx, pInactive, 1, now)
	if !errors.Is(err, cart.ErrPromoInactive) {
		t.Errorf("got error %v, want ErrPromoInactive", err)
	}

	// Expired promo
	past := now.Add(-24 * time.Hour)
	pExpired := &domain.Promo{Code: "EXPIRED", Active: true, EndDate: &past}
	err = calc.ValidatePromo(ctx, pExpired, 1, now)
	if !errors.Is(err, cart.ErrPromoExpired) {
		t.Errorf("got error %v, want ErrPromoExpired", err)
	}

	// Max uses reached
	pMax := &domain.Promo{Code: "MAXED", Active: true, MaxUses: 5, UsedCount: 5}
	err = calc.ValidatePromo(ctx, pMax, 1, now)
	if !errors.Is(err, cart.ErrPromoMaxUses) {
		t.Errorf("got error %v, want ErrPromoMaxUses", err)
	}

	// Once per client already used
	pOnce := &domain.Promo{Code: "ONCE", Active: true, OncePerClient: true}
	_ = promoRepo.Create(ctx, pOnce)
	_ = promoRepo.IncrementUsed(ctx, pOnce.ID, 1, nil) // Client 1 used it
	err = calc.ValidatePromo(ctx, pOnce, 1, now)
	if !errors.Is(err, cart.ErrPromoAlreadyUsed) {
		t.Errorf("got error %v, want ErrPromoAlreadyUsed", err)
	}

	// Client 2 can still use it
	err = calc.ValidatePromo(ctx, pOnce, 2, now)
	if err != nil {
		t.Errorf("client 2 should be able to use promo, got %v", err)
	}
}

func TestCartService_CheckoutFlow(t *testing.T) {
	cartService, promoRepo, orderRepo, clientRepo := setupCartService()
	ctx := context.Background()

	// Create test client
	client := &domain.Client{
		Email:     "buyer@example.com",
		FirstName: "Alice",
		LastName:  "Smith",
		Country:   "US",
		Currency:  "USD",
		Status:    "active",
	}
	_ = clientRepo.Create(ctx, client)

	// Create Promo
	promo := &domain.Promo{
		Code:   "SAVE20",
		Type:   domain.PromoTypePercentage,
		Value:  decimal.FromFloat(20.0),
		Active: true,
	}
	_ = promoRepo.Create(ctx, promo)

	shoppingCart := &cart.Cart{
		ClientID:  client.ID,
		PromoCode: "SAVE20",
		Items: []cart.CartItem{
			{ProductID: 10, Title: "cPanel Shared Hosting", Period: "1M", Price: decimal.FromFloat(50.00), Quantity: 1},
			{ProductID: 20, Title: "Domain Registration", Period: "1Y", Price: decimal.FromFloat(15.00), Quantity: 1},
		},
	}

	res, err := cartService.Checkout(ctx, shoppingCart)
	if err != nil {
		t.Fatalf("Checkout failed: %v", err)
	}

	// Verify Invoice
	if res.Invoice == nil {
		t.Fatal("expected invoice in checkout response")
	}
	// Subtotal = 50 + 15 = 65. Discount 20% = 13. Total = 52.
	if res.Invoice.Subtotal.String() != "52.00" {
		t.Errorf("got subtotal %s, want 52.00", res.Invoice.Subtotal.String())
	}
	if len(res.Invoice.Items) != 3 {
		t.Errorf("got %d invoice items, want 3 (2 products + 1 promo discount line item)", len(res.Invoice.Items))
	}

	// Verify Orders created in pending_setup
	if len(res.Orders) != 2 {
		t.Errorf("got %d orders, want 2", len(res.Orders))
	}
	for _, ord := range res.Orders {
		if ord.Status != domain.OrderStatusPendingSetup {
			t.Errorf("order %d status = %s, want pending_setup", ord.ID, ord.Status)
		}
		if ord.InvoiceID == nil || *ord.InvoiceID != res.Invoice.ID {
			t.Errorf("order %d invoice_id mismatch", ord.ID)
		}
	}

	// Verify Orders in Repository
	orders, _, _ := orderRepo.ListByClientID(ctx, client.ID, 10, 0)
	if len(orders) != 2 {
		t.Errorf("repository has %d orders, want 2", len(orders))
	}
}
