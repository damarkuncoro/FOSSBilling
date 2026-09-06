package builder

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

type OrderBuilder struct {
	clientID      int64
	productID     int64
	invoiceID     *int64
	title         string
	period        string
	price         decimal.Money
	currency      string
	status        domain.OrderStatus
	config        map[string]interface{}
	activatedDays int
}

func NewOrderBuilder() *OrderBuilder {
	return &OrderBuilder{
		period:        "1M",
		currency:      "USD",
		status:        domain.OrderStatusPendingSetup,
		config:        make(map[string]interface{}),
		activatedDays: 30,
	}
}

func (b *OrderBuilder) ForClient(clientID int64) *OrderBuilder {
	b.clientID = clientID
	return b
}

func (b *OrderBuilder) ForProduct(productID int64, title string) *OrderBuilder {
	b.productID = productID
	b.title = title
	return b
}

func (b *OrderBuilder) WithInvoice(invoiceID int64) *OrderBuilder {
	b.invoiceID = &invoiceID
	return b
}

func (b *OrderBuilder) WithPeriod(period string) *OrderBuilder {
	b.period = period
	return b
}

func (b *OrderBuilder) WithPrice(price decimal.Money, currency string) *OrderBuilder {
	b.price = price
	b.currency = currency
	return b
}

func (b *OrderBuilder) WithConfig(key string, value interface{}) *OrderBuilder {
	b.config[key] = value
	return b
}

func (b *OrderBuilder) AsActive() *OrderBuilder {
	b.status = domain.OrderStatusActive
	return b
}

func (b *OrderBuilder) Build() (*domain.Order, error) {
	if b.clientID == 0 {
		return nil, errors.New("order requires a valid client ID")
	}
	if b.productID == 0 {
		return nil, errors.New("order requires a valid product ID")
	}
	if b.title == "" {
		return nil, errors.New("order requires a non-empty title")
	}

	configJSON, _ := json.Marshal(b.config)
	now := time.Now().UTC()
	var activatedAt *time.Time
	var expiresAt *time.Time
	var nextDueDate *time.Time

	if b.status == domain.OrderStatusActive {
		activatedAt = &now
		exp := now.Add(time.Duration(b.activatedDays) * 24 * time.Hour)
		expiresAt = &exp
		nextDueDate = &exp
	}

	return &domain.Order{
		ClientID:    b.clientID,
		ProductID:   b.productID,
		InvoiceID:   b.invoiceID,
		Title:       b.title,
		Period:      b.period,
		Price:       b.price,
		Currency:    b.currency,
		Status:      b.status,
		Config:      configJSON,
		ActivatedAt: activatedAt,
		ExpiresAt:   expiresAt,
		NextDueDate: nextDueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}
