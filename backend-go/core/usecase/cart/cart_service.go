package cart

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/billing"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/formbuilder"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
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
	Tax       decimal.Money `json:"tax"`
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
	productRepo     domain.ProductRepository
	clientRepo      domain.ClientRepository
	formService     *formbuilder.FormbuilderService
	taxCalculator   *billing.TaxCalculator
	invoiceService  *billing.InvoiceService
	eventBus        *events.EventBus
}

func NewCartService(
	promoCalculator *PromoCalculator,
	promoRepo domain.PromoRepository,
	orderRepo domain.OrderRepository,
	productRepo domain.ProductRepository,
	clientRepo domain.ClientRepository,
	formService *formbuilder.FormbuilderService,
	taxCalculator *billing.TaxCalculator,
	invoiceService *billing.InvoiceService,
	eventBus *events.EventBus,
) *CartService {
	return &CartService{
		promoCalculator: promoCalculator,
		promoRepo:       promoRepo,
		orderRepo:       orderRepo,
		productRepo:     productRepo,
		clientRepo:      clientRepo,
		formService:     formService,
		taxCalculator:   taxCalculator,
		invoiceService:  invoiceService,
		eventBus:        eventBus,
	}
}

// CalculateTotals calculates subtotal, applied promo discount, tax and total
func (s *CartService) CalculateTotals(ctx context.Context, cart *Cart) error {
	var subtotal decimal.Money
	for i, it := range cart.Items {
		if it.ProductID <= 0 {
			return fmt.Errorf("invalid product ID: %d", it.ProductID)
		}

		if it.Quantity <= 0 {
			return fmt.Errorf("invalid quantity for product %d: must be greater than zero", it.ProductID)
		}

		// Security: Verify price against database
		if s.productRepo != nil {
			dbProd, err := s.productRepo.GetByID(ctx, it.ProductID)
			if err != nil || dbProd == nil {
				return fmt.Errorf("product not found: %d", it.ProductID)
			}

			// Override with DB price for security
			switch it.Period {
			case "1Y":
				if dbProd.PriceAnnually > 0 {
					cart.Items[i].Price = dbProd.PriceAnnually
				}
			default:
				if dbProd.PriceMonthly > 0 {
					cart.Items[i].Price = dbProd.PriceMonthly
				}
			}
			if cart.Items[i].Title == "" {
				cart.Items[i].Title = dbProd.Name
			}
		}

		if cart.Items[i].Price < 0 {
			return errors.New("item price cannot be negative")
		}

		subtotal += cart.Items[i].Price * decimal.Money(it.Quantity)
	}

	cart.Subtotal = subtotal
	cart.Discount = 0

	if cart.PromoCode != "" {
		promo, err := s.promoRepo.GetByCode(ctx, cart.PromoCode)
		if err == nil && promo != nil {
			if err := s.promoCalculator.ValidatePromo(ctx, promo, cart.ClientID, time.Now().UTC()); err == nil {
				discount := s.promoCalculator.CalculateDiscount(subtotal, promo)
				cart.Discount = discount
			}
		}
	}

	afterDiscount := subtotal - cart.Discount
	cart.Tax = 0
	cart.Total = afterDiscount

	if s.taxCalculator != nil && cart.ClientID > 0 {
		if client, err := s.clientRepo.GetByID(ctx, cart.ClientID); err == nil && client != nil {
			rate, _ := s.taxCalculator.GetTaxRateForClient(ctx, client)
			tax, total := s.taxCalculator.CalculateInvoiceTotals(afterDiscount, rate)
			cart.Tax = tax
			cart.Total = total
		}
	}

	return nil
}

// Checkout converts cart items into Orders (pending_setup) and generates an Invoice
func (s *CartService) Checkout(ctx context.Context, cart *Cart) (*CheckoutResult, error) {
	if len(cart.Items) == 0 {
		return nil, ErrEmptyCart
	}

	if err := s.CalculateTotals(ctx, cart); err != nil {
		return nil, err
	}

	var createdOrders []*domain.Order
	var invoiceItems []billing.CreateInvoiceItemDTO

	currency := "USD"
	if cart.ClientID > 0 && s.clientRepo != nil {
		if client, err := s.clientRepo.GetByID(ctx, cart.ClientID); err == nil && client != nil {
			currency = client.Currency
		}
	}

	for i, item := range cart.Items {
		// Validation: If product has a custom form, validate it
		if item.ProductID > 0 && s.productRepo != nil && s.formService != nil {
			prod, err := s.productRepo.GetByID(ctx, item.ProductID)
			if err == nil && prod != nil && prod.FormID != nil {
				var submitted map[string]interface{}
				if len(item.Config) > 0 {
					_ = json.Unmarshal(item.Config, &submitted)
				}

				cleaned, err := s.formService.ValidateFormSubmission(ctx, *prod.FormID, submitted)
				if err != nil {
					return nil, fmt.Errorf("configuration error for '%s': %w", item.Title, err)
				}

				// Re-encode cleaned config
				item.Config, _ = json.Marshal(cleaned)
				cart.Items[i].Config = item.Config
			}
		}

		// BUG-24 FIX: Use atomic decrement for stock
		if item.ProductID > 0 && s.productRepo != nil {
			err := s.productRepo.DecrementStock(ctx, item.ProductID, item.Quantity)
			if err != nil {
				return nil, fmt.Errorf("could not purchase '%s': %w", item.Title, err)
			}

			// Check stock for event triggering (optional check after atomic update)
			if prod, err := s.productRepo.GetByID(ctx, item.ProductID); err == nil && prod != nil && prod.Stock <= 5 {
				if s.eventBus != nil {
					s.eventBus.PublishAsync(ctx, events.Event{
						Type:    events.EventLowStock,
						Payload: prod,
					})
				}
			}
		}

		order := &domain.Order{
			ClientID:  cart.ClientID,
			ProductID: item.ProductID,
			Status:    domain.OrderStatusPendingSetup,
			Title:     item.Title,
			Period:    item.Period,
			Price:     item.Price,
			Currency:  currency,
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
			Quantity: item.Quantity,
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
			Taxable:  true,
		})
	}

	// Generate Invoice
	invoice, err := s.invoiceService.CreateInvoice(ctx, billing.CreateInvoiceDTO{
		ClientID: cart.ClientID,
		Currency: currency,
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
			if err := s.promoRepo.IncrementUsed(ctx, promo.ID, cart.ClientID, &createdOrders[0].ID); err != nil {
				return nil, fmt.Errorf("failed to apply promo: %w", err)
			}
		}
	}

	return &CheckoutResult{
		Orders:  createdOrders,
		Invoice: invoice,
	}, nil
}
