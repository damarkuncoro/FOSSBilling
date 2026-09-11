package billing

import (
	"context"
	"fmt"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/metrics"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/pdf"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/plugins"
)

type InvoiceService struct {
	invoiceRepo   domain.InvoiceRepository
	clientRepo    domain.ClientRepository
	companyRepo   domain.CompanyRepository
	taxCalculator *TaxCalculator
	hookManager   *plugins.HookManager
	eventBus      *events.EventBus
}

func NewInvoiceService(ir domain.InvoiceRepository, cr domain.ClientRepository, cor domain.CompanyRepository, tc *TaxCalculator, hm *plugins.HookManager, eb ...*events.EventBus) *InvoiceService {
	var bus *events.EventBus
	if len(eb) > 0 {
		bus = eb[0]
	}
	return &InvoiceService{ir, cr, cor, tc, hm, bus}
}

type CreateInvoiceDTO struct {
	ClientID int64                  `json:"client_id"`
	Currency string                 `json:"currency"`
	DueDays  int                    `json:"due_days"`
	Items    []CreateInvoiceItemDTO `json:"items"`
}

type CreateInvoiceItemDTO struct {
	OrderID  *int64        `json:"order_id"`
	Title    string        `json:"title"`
	Period   *string       `json:"period"`
	Price    decimal.Money `json:"price"`
	Quantity int           `json:"quantity"`
	Taxable  bool          `json:"taxable"`
}

func (s *InvoiceService) GetInvoice(ctx context.Context, id int64) (*domain.Invoice, error) {
	return s.invoiceRepo.GetByID(ctx, id)
}

func (s *InvoiceService) CreateInvoice(ctx context.Context, dto CreateInvoiceDTO) (*domain.Invoice, error) {
	c, err := s.clientRepo.GetByID(ctx, dto.ClientID)
	if err != nil {
		return nil, err
	}
	if dto.Currency == "" {
		dto.Currency = c.Currency
	}
	if dto.DueDays <= 0 {
		dto.DueDays = 14
	}

	var sub, taxSub decimal.Money
	var items []domain.InvoiceItem
	for _, it := range dto.Items {
		if it.Quantity <= 0 {
			it.Quantity = 1
		}
		title := it.Title
		if s.hookManager != nil {
			if f, err := s.hookManager.Apply(ctx, "filter_invoice_item_title", title); err == nil {
				title = f.(string)
			}
		}
		lt := it.Price * decimal.Money(it.Quantity)
		sub += lt
		if it.Taxable {
			taxSub += lt
		}
		items = append(items, domain.InvoiceItem{OrderID: it.OrderID, Title: title, Period: it.Period, Price: it.Price, Quantity: it.Quantity, Unit: "unit", Taxable: it.Taxable})
	}

	var rate float64; var tax decimal.Money
	if s.taxCalculator != nil {
		rate, _ = s.taxCalculator.GetTaxRateForClient(ctx, c)
		tax, _ = s.taxCalculator.CalculateInvoiceTotals(taxSub, rate)
	}

	now, total := time.Now().UTC(), sub+tax
	status, paidAt := domain.InvoiceStatusUnpaid, (*time.Time)(nil)
	if total == 0 {
		status, paidAt = domain.InvoiceStatusPaid, &now
	}

	inv := &domain.Invoice{Serie: "INV", ClientID: c.ID, Status: status, Currency: dto.Currency, CurrencyRate: 1.0, Subtotal: sub, Tax: tax, Total: total, TaxRate: rate, DueAt: now.AddDate(0, 0, dto.DueDays), PaidAt: paidAt}
	if err := s.invoiceRepo.Create(ctx, inv, items); err != nil {
		return nil, err
	}

	if total == 0 && s.eventBus != nil {
		s.eventBus.PublishAsync(ctx, events.Event{Type: events.EventInvoicePaid, Payload: domain.InvoicePaidPayload{InvoiceID: inv.ID, ClientID: inv.ClientID, Amount: 0, Currency: inv.Currency, PaidAt: now}})
	}
	return s.invoiceRepo.GetByID(ctx, inv.ID)
}

func (s *InvoiceService) PayWithBalance(ctx context.Context, cID, iID int64) (*domain.Invoice, error) {
	inv, err := s.invoiceRepo.GetByID(ctx, iID)
	if err != nil || inv.ClientID != cID {
		return nil, appErrors.ErrNotFound
	}
	if inv.Status == domain.InvoiceStatusPaid {
		return inv, nil
	}
	bal, _ := s.clientRepo.GetBalance(ctx, inv.ClientID)
	if bal < inv.Total {
		return nil, appErrors.ErrNoFunds
	}

	now := time.Now().UTC()
	if err := s.clientRepo.AddBalanceTransaction(ctx, &domain.ClientBalance{ClientID: inv.ClientID, Type: domain.BalanceTypeDebit, Amount: inv.Total, Description: "Payment for invoice #" + inv.Serie + inv.Nr, RelID: &inv.ID}); err != nil {
		return nil, err
	}
	if err := s.invoiceRepo.MarkAsPaid(ctx, inv.ID, now); err != nil {
		return nil, err
	}

	metrics.InvoicesPaidTotal.Inc()
	metrics.RevenueTotal.WithLabelValues(inv.Currency).Add(inv.Total.ToFloat())

	if s.eventBus != nil {
		_ = s.eventBus.Publish(ctx, events.Event{Type: events.EventInvoicePaid, Payload: domain.InvoicePaidPayload{InvoiceID: inv.ID, ClientID: inv.ClientID, Amount: inv.Total, Currency: inv.Currency, PaidAt: now}})
	}
	return s.invoiceRepo.GetByID(ctx, inv.ID)
}

func (s *InvoiceService) RefundInvoice(ctx context.Context, id int64) error {
	inv, err := s.invoiceRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if ok, _ := s.invoiceRepo.UpdateStatusAtomic(ctx, inv.ID, domain.InvoiceStatusRefunded, domain.InvoiceStatusPaid); !ok {
		return fmt.Errorf("invalid refund state")
	}
	return s.clientRepo.AddBalanceTransaction(ctx, &domain.ClientBalance{ClientID: inv.ClientID, Type: domain.BalanceTypeCredit, Amount: inv.Total, Description: "Refund for invoice #" + inv.Serie + inv.Nr, RelID: &inv.ID})
}

func (s *InvoiceService) GeneratePDF(ctx context.Context, id int64) ([]byte, string, error) {
	inv, err := s.invoiceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, "", err
	}
	client, err := s.clientRepo.GetByID(ctx, inv.ClientID)
	if err != nil {
		return nil, "", err
	}

	compName, compAddr, compTax := "FOSSBilling", "", ""
	if s.companyRepo != nil {
		settings, err := s.companyRepo.Get(ctx)
		if err == nil && settings != nil {
			compName = settings.Name
			compAddr = fmt.Sprintf("%s, %s, %s %s, %s", settings.Address1, settings.City, settings.State, settings.Postcode, settings.Country)
			compTax = settings.VatNumber
		}
	}

	data, err := pdf.GenerateInvoicePDF(inv, client, compName, compAddr, compTax)
	filename := fmt.Sprintf("Invoice_%s_%s.pdf", inv.Serie, inv.Nr)
	return data, filename, err
}
