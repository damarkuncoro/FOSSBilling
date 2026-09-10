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
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/fraud"
)

var ErrEmptyCart = errors.New("cart is empty")

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
	fraudChecker    fraud.FraudChecker
	eventBus        *events.EventBus
}

func NewCartService(pc *PromoCalculator, pr domain.PromoRepository, or domain.OrderRepository, pdr domain.ProductRepository, cr domain.ClientRepository, fs *formbuilder.FormbuilderService, tc *billing.TaxCalculator, is *billing.InvoiceService, fc fraud.FraudChecker, eb *events.EventBus) *CartService {
	return &CartService{pc, pr, or, pdr, cr, fs, tc, is, fc, eb}
}

func (s *CartService) CalculateTotals(ctx context.Context, cart *Cart) error {
	cart.Subtotal, cart.Discount, cart.Tax = 0, 0, 0
	for i, it := range cart.Items {
		if it.ProductID <= 0 || it.Quantity <= 0 {
			return fmt.Errorf("invalid item: product_id=%d, qty=%d", it.ProductID, it.Quantity)
		}
		prod, err := s.productRepo.GetByID(ctx, it.ProductID)
		if err != nil || prod == nil {
			return fmt.Errorf("product %d not found", it.ProductID)
		}
		cart.Items[i].Price = map[string]decimal.Money{"1Y": prod.PriceAnnually}[it.Period]
		if cart.Items[i].Price == 0 {
			cart.Items[i].Price = prod.PriceMonthly
		}
		if cart.Items[i].Title == "" {
			cart.Items[i].Title = prod.Name
		}
		cart.Subtotal += cart.Items[i].Price * decimal.Money(it.Quantity)
	}

	if cart.PromoCode != "" {
		if promo, err := s.promoRepo.GetByCode(ctx, cart.PromoCode); err == nil && s.promoCalculator.ValidatePromo(ctx, promo, cart.ClientID, time.Now().UTC()) == nil {
			cart.Discount = s.promoCalculator.CalculateDiscount(cart.Subtotal, promo)
		}
	}

	cart.Total = cart.Subtotal - cart.Discount
	if client, err := s.clientRepo.GetByID(ctx, cart.ClientID); err == nil && client != nil {
		rate, _ := s.taxCalculator.GetTaxRateForClient(ctx, client)
		cart.Tax, cart.Total = s.taxCalculator.CalculateInvoiceTotals(cart.Total, rate)
	}
	return nil
}

func (s *CartService) Checkout(ctx context.Context, cart *Cart, remoteIP string) (*CheckoutResult, error) {
	if len(cart.Items) == 0 {
		return nil, ErrEmptyCart
	}
	if err := s.CalculateTotals(ctx, cart); err != nil {
		return nil, err
	}

	status := domain.OrderStatusPendingSetup
	if s.fraudChecker != nil {
		if score, _ := s.fraudChecker.CheckIP(ctx, remoteIP); score != nil && score.RiskLevel == "high" {
			status = "manual_review"
		}
	}

	currency := "USD"
	if client, err := s.clientRepo.GetByID(ctx, cart.ClientID); err == nil && client != nil {
		currency = client.Currency
	}

	var orders []*domain.Order
	var invItems []billing.CreateInvoiceItemDTO
	for i, it := range cart.Items {
		if prod, _ := s.productRepo.GetByID(ctx, it.ProductID); prod != nil && prod.FormID != nil && s.formService != nil {
			var sub map[string]interface{}
			_ = json.Unmarshal(it.Config, &sub)
			if clean, err := s.formService.ValidateFormSubmission(ctx, *prod.FormID, sub); err == nil {
				it.Config, _ = json.Marshal(clean)
				cart.Items[i].Config = it.Config
			}
		}

		if err := s.productRepo.DecrementStock(ctx, it.ProductID, it.Quantity); err != nil {
			return nil, err
		}

		ord := &domain.Order{ClientID: cart.ClientID, ProductID: it.ProductID, Status: status, Title: it.Title, Period: it.Period, Price: it.Price, Currency: currency, Config: it.Config}
		if err := s.orderRepo.Create(ctx, ord); err != nil {
			return nil, err
		}
		orders = append(orders, ord)
		invItems = append(invItems, billing.CreateInvoiceItemDTO{OrderID: &ord.ID, Title: it.Title, Period: &it.Period, Price: it.Price, Quantity: it.Quantity, Taxable: true})
	}

	if cart.Discount > 0 {
		invItems = append(invItems, billing.CreateInvoiceItemDTO{Title: "Discount (" + cart.PromoCode + ")", Price: -cart.Discount, Quantity: 1, Taxable: true})
	}

	inv, err := s.invoiceService.CreateInvoice(ctx, billing.CreateInvoiceDTO{ClientID: cart.ClientID, Currency: currency, DueDays: 14, Items: invItems})
	if err != nil {
		return nil, err
	}

	for _, o := range orders {
		o.InvoiceID = &inv.ID
		_ = s.orderRepo.Update(ctx, o)
	}

	if cart.PromoCode != "" {
		if promo, _ := s.promoRepo.GetByCode(ctx, cart.PromoCode); promo != nil {
			_ = s.promoRepo.IncrementUsed(ctx, promo.ID, cart.ClientID, &orders[0].ID)
		}
	}

	return &CheckoutResult{Orders: orders, Invoice: inv}, nil
}
