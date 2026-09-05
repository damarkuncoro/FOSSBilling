package billing

import (
	"context"
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/repository/memory"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

func setupInvoiceService() (*InvoiceService, *memory.MockInvoiceRepository, *memory.MockClientRepository) {
	invRepo := memory.NewMockInvoiceRepository()
	clientRepo := memory.NewMockClientRepository()

	rules := []TaxRule{
		{Name: "PPN", Country: "ID", State: "", Rate: 11.0},
	}
	taxCalc := NewTaxCalculator(rules)

	service := NewInvoiceService(invRepo, clientRepo, taxCalc)
	return service, invRepo, clientRepo
}

func TestCalculateProrata(t *testing.T) {
	// Full price: $30.00 for 30 days -> 10 days active = $10.00
	fullPrice := decimal.FromFloat(30.00)
	prorated := CalculateProrata(fullPrice, 30, 10)
	if prorated.String() != "10.00" {
		t.Errorf("Prorated 10/30 of $30.00 = %s; want 10.00", prorated.String())
	}

	// 15 days out of 30 for $19.99 -> 9.995 -> rounds to $10.00
	price2 := decimal.FromFloat(19.99)
	prorated2 := CalculateProrata(price2, 30, 15)
	if prorated2.String() != "10.00" {
		t.Errorf("Prorated 15/30 of $19.99 = %s; want 10.00", prorated2.String())
	}

	// Zero active days -> $0.00
	if zero := CalculateProrata(fullPrice, 30, 0); zero != 0 {
		t.Errorf("Prorated 0 days = %v; want 0", zero)
	}
}

func TestConvertCurrency(t *testing.T) {
	// USD base rate 1.0, IDR rate 16000.0
	usdAmount := decimal.FromFloat(10.00) // $10.00
	idrAmount := ConvertCurrency(usdAmount, 1.0, 16000.0)
	if idrAmount.String() != "160000.00" {
		t.Errorf("ConvertCurrency($10, 1.0 -> 16000) = %s; want 160000.00", idrAmount.String())
	}
}

func TestInvoiceService_CreateInvoiceWithTax(t *testing.T) {
	ctx := context.Background()
	service, _, clientRepo := setupInvoiceService()

	client := &domain.Client{
		Email:     "buyer@example.com",
		FirstName: "Rudi",
		LastName:  "Hermawan",
		Country:   "ID",
		Currency:  "IDR",
	}
	_ = clientRepo.Create(ctx, client)

	dto := CreateInvoiceDTO{
		ClientID: client.ID,
		Currency: "IDR",
		DueDays:  7,
		Items: []CreateInvoiceItemDTO{
			{Title: "Cloud VPS Basic", Price: decimal.FromFloat(100000.00), Quantity: 1, Taxable: true},
			{Title: "Domain .com", Price: decimal.FromFloat(150000.00), Quantity: 1, Taxable: true},
		},
	}

	invoice, err := service.CreateInvoice(ctx, dto)
	if err != nil {
		t.Fatalf("CreateInvoice failed: %v", err)
	}

	// Subtotal = 250,000.00, Tax 11% = 27,500.00, Total = 277,500.00
	if invoice.Subtotal.String() != "250000.00" {
		t.Errorf("Subtotal = %s; want 250000.00", invoice.Subtotal.String())
	}
	if invoice.Tax.String() != "27500.00" {
		t.Errorf("Tax = %s; want 27500.00", invoice.Tax.String())
	}
	if invoice.Total.String() != "277500.00" {
		t.Errorf("Total = %s; want 277500.00", invoice.Total.String())
	}
	if len(invoice.Items) != 2 {
		t.Errorf("Items count = %d; want 2", len(invoice.Items))
	}
}

func TestInvoiceService_PayWithBalance(t *testing.T) {
	ctx := context.Background()
	service, _, clientRepo := setupInvoiceService()

	client := &domain.Client{
		Email:     "credit.buyer@example.com",
		FirstName: "Siti",
		Country:   "US",
		Currency:  "USD",
	}
	_ = clientRepo.Create(ctx, client)

	// Create invoice of $50.00
	invoice, _ := service.CreateInvoice(ctx, CreateInvoiceDTO{
		ClientID: client.ID,
		Currency: "USD",
		Items: []CreateInvoiceItemDTO{
			{Title: "Web Hosting Pro", Price: decimal.FromFloat(50.00), Quantity: 1},
		},
	})

	// 1. Pay with insufficient balance -> should fail
	_, err := service.PayWithBalance(ctx, invoice.ID)
	if err != appErrors.ErrInsufficientFunds {
		t.Errorf("Expected ErrInsufficientFunds, got: %v", err)
	}

	// 2. Deposit credit of $100.00
	_ = clientRepo.AddBalanceTransaction(ctx, &domain.ClientBalance{
		ClientID: client.ID,
		Type:     domain.BalanceTypeCredit,
		Amount:   decimal.FromFloat(100.00),
	})

	// 3. Pay invoice with balance -> should succeed
	paidInv, err := service.PayWithBalance(ctx, invoice.ID)
	if err != nil {
		t.Fatalf("PayWithBalance failed: %v", err)
	}
	if paidInv.Status != domain.InvoiceStatusPaid {
		t.Errorf("Status = %s; want paid", paidInv.Status)
	}
	if paidInv.PaidAt == nil {
		t.Error("Expected PaidAt timestamp to be set")
	}

	// 4. Verify remaining client balance ($100 - $50 = $50)
	remainingBalance, _ := clientRepo.GetBalance(ctx, client.ID)
	if remainingBalance.String() != "50.00" {
		t.Errorf("Remaining balance = %s; want 50.00", remainingBalance.String())
	}
}
