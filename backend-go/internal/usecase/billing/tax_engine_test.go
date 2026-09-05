package billing

import (
	"testing"

	"github.com/fossbilling/backend-go/internal/domain"
	"github.com/fossbilling/backend-go/pkg/decimal"
)

func TestTaxEngine_GetTaxRateForClient(t *testing.T) {
	rules := []TaxRule{
		{Name: "California Tax", Country: "US", State: "CA", Rate: 9.5},
		{Name: "US Federal", Country: "US", State: "", Rate: 5.0},
		{Name: "EU VAT", Country: "NL", State: "", Rate: 21.0},
		{Name: "Indonesian PPN", Country: "ID", State: "", Rate: 11.0},
	}

	tc := NewTaxCalculator(rules)

	// Scenario 1: State & Country Match (California)
	clientCA := &domain.Client{Country: "US", State: "CA", TaxExempt: false}
	rate, name := tc.GetTaxRateForClient(clientCA)
	if rate != 9.5 || name != "California Tax" {
		t.Errorf("CA rate = %v, %s; want 9.5, California Tax", rate, name)
	}

	// Scenario 2: Country Fallback Match (New York -> US Federal)
	clientNY := &domain.Client{Country: "US", State: "NY", TaxExempt: false}
	rate, name = tc.GetTaxRateForClient(clientNY)
	if rate != 5.0 || name != "US Federal" {
		t.Errorf("NY rate = %v, %s; want 5.0, US Federal", rate, name)
	}

	// Scenario 3: Tax Exempt Client
	clientExempt := &domain.Client{Country: "NL", State: "", TaxExempt: true}
	rate, name = tc.GetTaxRateForClient(clientExempt)
	if rate != 0.0 || name != "Tax Exempt" {
		t.Errorf("Exempt rate = %v, %s; want 0.0, Tax Exempt", rate, name)
	}

	// Scenario 4: Country with no specific rule
	clientSG := &domain.Client{Country: "SG", State: "", TaxExempt: false}
	rate, name = tc.GetTaxRateForClient(clientSG)
	if rate != 0.0 {
		t.Errorf("SG rate = %v; want 0.0", rate)
	}
}

func TestTaxEngine_CalculateInvoiceTotals(t *testing.T) {
	tc := NewTaxCalculator(nil)

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
	// 49.99 * 0.11 = 5.4989 -> rounds to 5.50
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
