package billing

import (
	"context"
	"fmt"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
)

type InvoiceService struct {
	invoiceRepo   domain.InvoiceRepository
	clientRepo    domain.ClientRepository
	taxCalculator *TaxCalculator
	eventBus      *events.EventBus
}

func NewInvoiceService(
	invoiceRepo domain.InvoiceRepository,
	clientRepo domain.ClientRepository,
	taxCalculator *TaxCalculator,
	eventBus ...*events.EventBus,
) *InvoiceService {
	var bus *events.EventBus
	if len(eventBus) > 0 {
		bus = eventBus[0]
	}
	return &InvoiceService{
		invoiceRepo:   invoiceRepo,
		clientRepo:    clientRepo,
		taxCalculator: taxCalculator,
		eventBus:      bus,
	}
}

type CreateInvoiceDTO struct {
	ClientID int64
	Currency string
	DueDays  int
	Items    []CreateInvoiceItemDTO
}

type CreateInvoiceItemDTO struct {
	OrderID  *int64
	Title    string
	Period   *string
	Price    decimal.Money
	Quantity int
	Taxable  bool
}

// CreateInvoice generates a new invoice with subtotal, tax calculation, and line items
func (s *InvoiceService) CreateInvoice(ctx context.Context, dto CreateInvoiceDTO) (*domain.Invoice, error) {
	client, err := s.clientRepo.GetByID(ctx, dto.ClientID)
	if err != nil {
		return nil, err
	}

	if dto.Currency == "" {
		dto.Currency = client.Currency
	}
	if dto.DueDays <= 0 {
		dto.DueDays = 14
	}

	var subtotal decimal.Money
	var taxableSubtotal decimal.Money
	var items []domain.InvoiceItem

	for _, it := range dto.Items {
		if it.Quantity <= 0 {
			it.Quantity = 1
		}
		lineTotal := it.Price * decimal.Money(it.Quantity)
		subtotal += lineTotal
		if it.Taxable {
			taxableSubtotal += lineTotal
		}

		items = append(items, domain.InvoiceItem{
			OrderID:  it.OrderID,
			Title:    it.Title,
			Period:   it.Period,
			Price:    it.Price,
			Quantity: it.Quantity,
			Unit:     "unit",
			Taxable:  it.Taxable,
		})
	}

	// Calculate Tax
	var taxRate float64
	var tax decimal.Money
	var total decimal.Money = subtotal

	if s.taxCalculator != nil {
		rate, _ := s.taxCalculator.GetTaxRateForClient(client)
		taxRate = rate
		tax, _ = s.taxCalculator.CalculateInvoiceTotals(taxableSubtotal, taxRate)
		total = subtotal + tax
	}

	now := time.Now().UTC()
	dueAt := now.AddDate(0, 0, dto.DueDays)

	invoice := &domain.Invoice{
		Serie:        "INV",
		ClientID:     client.ID,
		Status:       domain.InvoiceStatusUnpaid,
		Currency:     dto.Currency,
		CurrencyRate: 1.0,
		Subtotal:     subtotal,
		Tax:          tax,
		Total:        total,
		TaxRate:      taxRate,
		DueAt:        dueAt,
	}

	if err := s.invoiceRepo.Create(ctx, invoice, items); err != nil {
		return nil, err
	}

	return s.invoiceRepo.GetByID(ctx, invoice.ID)
}

// PayWithBalance deducts from client balance to mark an invoice as paid
func (s *InvoiceService) PayWithBalance(ctx context.Context, invoiceID int64) (*domain.Invoice, error) {
	inv, err := s.invoiceRepo.GetByID(ctx, invoiceID)
	if err != nil {
		return nil, err
	}

	if inv.Status == domain.InvoiceStatusPaid {
		return inv, nil
	}

	balance, err := s.clientRepo.GetBalance(ctx, inv.ClientID)
	if err != nil {
		return nil, err
	}

	if balance < inv.Total {
		return nil, appErrors.ErrInsufficientFunds
	}

	// Deduct balance
	deduct := &domain.ClientBalance{
		ClientID:    inv.ClientID,
		Type:        domain.BalanceTypeDebit,
		Amount:      inv.Total,
		Description: fmt.Sprintf("Payment for invoice #%s%s", inv.Serie, inv.Nr),
		RelID:       &inv.ID,
	}
	if err := s.clientRepo.AddBalanceTransaction(ctx, deduct); err != nil {
		return nil, err
	}

	// Mark invoice as paid
	now := time.Now().UTC()
	if err := s.invoiceRepo.MarkAsPaid(ctx, inv.ID, now); err != nil {
		return nil, err
	}

	// Publish Event
	if s.eventBus != nil {
		_ = s.eventBus.Publish(ctx, events.Event{
			Type: events.EventInvoicePaid,
			Payload: domain.InvoicePaidPayload{
				InvoiceID: inv.ID,
				ClientID:  inv.ClientID,
				Amount:    inv.Total,
				Currency:  inv.Currency,
				PaidAt:    now,
			},
		})
	}

	return s.invoiceRepo.GetByID(ctx, inv.ID)
}
