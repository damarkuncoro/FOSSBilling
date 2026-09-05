package billing

import (
	"math"

	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

type TaxRule struct {
	Name      string
	Country   string // ISO 2-letter
	State     string
	Rate      float64 // e.g. 21.0 for 21%
	TaxExempt bool
}

type TaxCalculator struct {
	rules []TaxRule
}

func NewTaxCalculator(rules []TaxRule) *TaxCalculator {
	return &TaxCalculator{rules: rules}
}

func (tc *TaxCalculator) GetTaxRateForClient(client *domain.Client) (rate float64, taxName string) {
	if client.TaxExempt {
		return 0.0, "Tax Exempt"
	}

	// 1. Match State & Country
	for _, rule := range tc.rules {
		if rule.Country == client.Country && rule.State != "" && rule.State == client.State {
			return rule.Rate, rule.Name
		}
	}

	// 2. Match Country
	for _, rule := range tc.rules {
		if rule.Country == client.Country && rule.State == "" {
			return rule.Rate, rule.Name
		}
	}

	// 3. Fallback to Global / Default rule
	for _, rule := range tc.rules {
		if rule.Country == "" && rule.State == "" {
			return rule.Rate, rule.Name
		}
	}

	return 0.0, ""
}

func (tc *TaxCalculator) CalculateInvoiceTotals(subtotal decimal.Money, taxRate float64) (tax decimal.Money, total decimal.Money) {
	if taxRate <= 0 {
		return 0, subtotal
	}

	taxFloat := subtotal.ToFloat() * (taxRate / 100.0)
	tax = decimal.FromFloat(math.Round(taxFloat*100) / 100)
	total = subtotal + tax
	return tax, total
}
