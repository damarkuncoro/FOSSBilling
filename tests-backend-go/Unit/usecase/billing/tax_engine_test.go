package billing_test

import (
	"context"
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/billing"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

type mockTaxRepo struct {
	rules []*domain.TaxRule
}

func (m *mockTaxRepo) GetByID(ctx context.Context, id int64) (*domain.TaxRule, error) { return nil, nil }
func (m *mockTaxRepo) List(ctx context.Context) ([]*domain.TaxRule, error)           { return m.rules, nil }
func (m *mockTaxRepo) GetByLocation(ctx context.Context, country, state string) (*domain.TaxRule, error) {
	for _, r := range m.rules {
		c := ""
		if r.Country != nil {
			c = *r.Country
		}
		s := ""
		if r.State != nil {
			s = *r.State
		}
		if c == country && s == state {
			return r, nil
		}
		if c == country && s == "" {
			return r, nil
		}
		if c == "" && s == "" {
			return r, nil
		}
	}
	return nil, nil
}
func (m *mockTaxRepo) Create(ctx context.Context, rule *domain.TaxRule) error { return nil }
func (m *mockTaxRepo) Update(ctx context.Context, rule *domain.TaxRule) error { return nil }
func (m *mockTaxRepo) Delete(ctx context.Context, id int64) error             { return nil }

func TestTaxEngine_GetTaxRateForClient(t *testing.T) {
	ca := "CA"
	us := "US"
	nl := "NL"
	id := "ID"

	mock := &mockTaxRepo{
		rules: []*domain.TaxRule{
			{Name: "California Tax", Country: &us, State: &ca, Rate: 9.5, IsActive: true},
			{Name: "US Federal", Country: &us, Rate: 5.0, IsActive: true},
			{Name: "EU VAT", Country: &nl, Rate: 21.0, IsActive: true},
			{Name: "Indonesian PPN", Country: &id, Rate: 11.0, IsActive: true},
		},
	}

	tc := billing.NewTaxCalculator(mock)
	ctx := context.Background()

	// Scenario 1: State & Country Match (California)
	clientCA := &domain.Client{Country: "US", State: "CA", TaxExempt: false}
	rate, name := tc.GetTaxRateForClient(ctx, clientCA)
	if rate != 9.5 || name != "California Tax" {
		t.Errorf("CA rate = %v, %s; want 9.5, California Tax", rate, name)
	}

	// Scenario 2: Country Fallback Match (New York -> US Federal)
	clientNY := &domain.Client{Country: "US", State: "NY", TaxExempt: false}
	rate, name = tc.GetTaxRateForClient(ctx, clientNY)
	if rate != 5.0 || name != "US Federal" {
		t.Errorf("NY rate = %v, %s; want 5.0, US Federal", rate, name)
	}

	// Scenario 3: Tax Exempt Client
	clientExempt := &domain.Client{Country: "NL", State: "", TaxExempt: true}
	rate, name = tc.GetTaxRateForClient(ctx, clientExempt)
	if rate != 0.0 || name != "Tax Exempt" {
		t.Errorf("Exempt rate = %v, %s; want 0.0, Tax Exempt", rate, name)
	}

	// Scenario 4: Country with no specific rule
	clientSG := &domain.Client{Country: "SG", State: "", TaxExempt: false}
	rate, name = tc.GetTaxRateForClient(ctx, clientSG)
	if rate != 0.0 {
		t.Errorf("SG rate = %v; want 0.0", rate)
	}
}

func TestTaxEngine_CalculateInvoiceTotals(t *testing.T) {
	tc := billing.NewTaxCalculator(nil)

	// Scenario: Subtotal $100.00 with 21% VAT
	subtotal := decimal.FromFloat(100.00)
	tax, total := tc.CalculateInvoiceTotals(subtotal, 21.0)

	if tax.String() != "21.00" {
		t.Errorf("tax = %s; want 21.00", tax.String())
	}
	if total.String() != "121.00" {
		t.Errorf("total = %s; want 121.00", total.String())
	}

	// Scenario: Subtotal $49.99 with 11% PPN
	subtotal2 := decimal.FromFloat(49.99)
	tax2, total2 := tc.CalculateInvoiceTotals(subtotal2, 11.0)
	if tax2.String() != "5.50" {
		t.Errorf("tax2 = %s; want 5.50", tax2.String())
	}
	if total2.String() != "55.49" {
		t.Errorf("total2 = %s; want 55.49", total2.String())
	}

	// Scenario: 0% tax rate
	tax3, total3 := tc.CalculateInvoiceTotals(subtotal, 0.0)
	if tax3 != 0 || total3 != subtotal {
		t.Errorf("0%% tax: tax = %v, total = %v; want 0, %v", tax3, total3, subtotal)
	}
}
