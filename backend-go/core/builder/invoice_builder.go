package builder

import (
	"errors"
	"fmt"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

type InvoiceBuilder struct {
	clientID     int64
	serie        string
	nr           string
	currency     string
	currencyRate float64
	dueDays      int
	taxRate      float64
	promo        *domain.Promo
	items        []domain.InvoiceItem
}

func NewInvoiceBuilder() *InvoiceBuilder {
	return &InvoiceBuilder{
		serie:        "INV",
		currency:     "USD",
		currencyRate: 1.0,
		dueDays:      14,
		items:        make([]domain.InvoiceItem, 0),
	}
}

func (b *InvoiceBuilder) ForClient(clientID int64) *InvoiceBuilder {
	b.clientID = clientID
	return b
}

func (b *InvoiceBuilder) WithSerieAndNr(serie, nr string) *InvoiceBuilder {
	b.serie = serie
	b.nr = nr
	return b
}

func (b *InvoiceBuilder) WithCurrency(currency string, rate float64) *InvoiceBuilder {
	b.currency = currency
	b.currencyRate = rate
	return b
}

func (b *InvoiceBuilder) WithDueDays(days int) *InvoiceBuilder {
	b.dueDays = days
	return b
}

func (b *InvoiceBuilder) WithTaxRate(taxRate float64) *InvoiceBuilder {
	b.taxRate = taxRate
	return b
}

func (b *InvoiceBuilder) ApplyPromo(promo *domain.Promo) *InvoiceBuilder {
	b.promo = promo
	return b
}

func (b *InvoiceBuilder) AddItem(title string, price decimal.Money, quantity int, taxable bool) *InvoiceBuilder {
	b.items = append(b.items, domain.InvoiceItem{
		Title:    title,
		Price:    price,
		Quantity: quantity,
		Taxable:  taxable,
	})
	return b
}

func (b *InvoiceBuilder) Build() (*domain.Invoice, []domain.InvoiceItem, error) {
	if b.clientID == 0 {
		return nil, nil, errors.New("invoice requires a valid client ID")
	}
	if len(b.items) == 0 {
		return nil, nil, errors.New("invoice must have at least one line item")
	}

	var subtotal decimal.Money
	var taxableSubtotal decimal.Money

	for _, item := range b.items {
		itemTotal := item.Price * decimal.Money(item.Quantity)
		subtotal += itemTotal
		if item.Taxable {
			taxableSubtotal += itemTotal
		}
	}

	// Apply discount if promo voucher present
	var discount decimal.Money
	if b.promo != nil && b.promo.Active {
		if b.promo.Type == domain.PromoTypePercentage {
			discount = decimal.FromFloat(subtotal.ToFloat() * (float64(b.promo.Value) / 1000000.0))
		} else {
			discount = decimal.FromFloat(float64(b.promo.Value) / 10000.0)
		}
		if discount > subtotal {
			discount = subtotal
		}
	}

	afterDiscount := subtotal - discount
	taxAmount := decimal.FromFloat(taxableSubtotal.ToFloat() * (b.taxRate / 100.0))
	total := afterDiscount + taxAmount

	now := time.Now().UTC()
	dueDate := now.Add(time.Duration(b.dueDays) * 24 * time.Hour)

	if b.nr == "" {
		b.nr = fmt.Sprintf("%d-%04d", now.Year(), now.Unix()%10000)
	}

	invoice := &domain.Invoice{
		Serie:        b.serie,
		Nr:           b.nr,
		ClientID:     b.clientID,
		Status:       domain.InvoiceStatusUnpaid,
		Currency:     b.currency,
		CurrencyRate: b.currencyRate,
		Subtotal:     subtotal,
		Tax:          taxAmount,
		Total:        total,
		TaxRate:      b.taxRate,
		DueAt:        dueDate,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	return invoice, b.items, nil
}
