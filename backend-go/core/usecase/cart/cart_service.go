package cart

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/billing"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

var (
	ErrEmptyCart = errors.New("cart is empty")
)

type CartItem struct {
	ProductID int64           `json:"product_id"`
	Title     string          `json:"title"`
	Period    string          `json:"period"`
	Price     decimal.Money   `json:"price"`
	Quantity  int             `json:"quantity"`
	Config    json.RawMessage `json:"config,omitempty"`
}

type Cart struct {
	ClientID  int64         `json:"client_id"`
	Items     []CartItem    `json:"items"`
	PromoCode string        `json:"promo_code,omitempty"`
	Subtotal  decimal.Money `json:"subtotal"`
	Discount  decimal.Money `json:"discount"`
	Total     decimal.Money `json:"total"`
}

type CheckoutResult struct {
	Orders  []*domain.Order `json:"orders"`
	Invoice *domain.Invoice `json:"invoice"`
}

type CartService struct {
	promoCalculator *PromoCalculator
	promoRepo       domain.PromoRepository
	orderRepo       domain.OrderRepository
	invoiceService  *billing.InvoiceService
}

func NewCartService(
	promoCalculator *PromoCalculator,
	promoRepo domain.PromoRepository,
	orderRepo domain.OrderRepository,
	invoiceService *billing.InvoiceService,
) *CartService {
	return &CartService{
		promoCalculator: promoCalculator,
		promoRepo:       promoRepo,
		orderRepo:       orderRepo,
		invoiceService:  invoiceService,
	}
}

// CalculateTotals calculates subtotal, applied promo discount, and total
func (s *CartService) CalculateTotals(ctx context.Context, cart *Cart) error {
	var subtotal decimal.Money
	for _, it := range cart.Items {
		qty := it.Quantity
		if qty <= 0 {
			qty = 1
		}
		subtotal += it.Price * decimal.Money(qty)
	}

	cart.Subtotal = subtotal
	cart.Discount = 0
	cart.Total = subtotal

	if cart.PromoCode != "" {
		promo, err := s.promoRepo.GetByCode(ctx, cart.PromoCode)
		if err == nil && promo != nil {
			if err := s.promoCalculator.ValidatePromo(ctx, promo, cart.ClientID, time.Now().UTC()); err == nil {
				discount := s.promoCalculator.CalculateDiscount(subtotal, promo)
				cart.Discount = discount
				cart.Total = subtotal - discount
			}
		}
	}

	return nil
}

// Checkout converts cart items into Orders (pending_setup) and generates an Invoice
func (s *CartService) Checkout(ctx context.Context, cart *Cart) (*CheckoutResult, error) {
	if len(cart.Items) == 0 {
		return nil, ErrEmptyCart
	}

	_ = s.CalculateTotals(ctx, cart)

	var createdOrders []*domain.Order
	var invoiceItems []billing.CreateInvoiceItemDTO

	for _, item := range cart.Items {
		qty := item.Quantity
		if qty <= 0 {
			qty = 1
		}

		order := &domain.Order{
			ClientID:  cart.ClientID,
			ProductID: item.ProductID,
			Status:    domain.OrderStatusPendingSetup,
			Title:     item.Title,
			Period:    item.Period,
			Price:     item.Price,
			Currency:  "USD",
			Config:    item.Config,
		}

		if err := s.orderRepo.Create(ctx, order); err != nil {
			return nil, err
		}
		createdOrders = append(createdOrders, order)

		invoiceItems = append(invoiceItems, billing.CreateInvoiceItemDTO{
			OrderID:  &order.ID,
			Title:    item.Title,
			Period:   &item.Period,
			Price:    item.Price,
			Quantity: qty,
			Taxable:  true,
		})
	}

	// Apply discount line item if promo exists
	if cart.Discount > 0 {
		negDiscount := -cart.Discount
		discountTitle := "Discount (" + cart.PromoCode + ")"
		invoiceItems = append(invoiceItems, billing.CreateInvoiceItemDTO{
			Title:    discountTitle,
			Price:    negDiscount,
			Quantity: 1,
			Taxable:  false,
		})
	}

	// Generate Invoice
	invoice, err := s.invoiceService.CreateInvoice(ctx, billing.CreateInvoiceDTO{
		ClientID: cart.ClientID,
		Currency: "USD",
		DueDays:  14,
		Items:    invoiceItems,
	})
	if err != nil {
		return nil, err
	}

	// Update orders with invoice ID
	for _, o := range createdOrders {
		o.InvoiceID = &invoice.ID
		_ = s.orderRepo.Update(ctx, o)
	}

	// Record promo redemption if applicable
	if cart.PromoCode != "" {
		if promo, err := s.promoRepo.GetByCode(ctx, cart.PromoCode); err == nil && promo != nil {
			_ = s.promoRepo.IncrementUsed(ctx, promo.ID, cart.ClientID, &createdOrders[0].ID)
		}
	}

	return &CheckoutResult{
		Orders:  createdOrders,
		Invoice: invoice,
	}, nil
}
