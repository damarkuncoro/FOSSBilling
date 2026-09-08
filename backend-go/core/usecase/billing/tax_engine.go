package billing

import (
	"context"
	"math"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

type TaxCalculator struct {
	repo domain.TaxRepository
}

func NewTaxCalculator(repo domain.TaxRepository) *TaxCalculator {
	return &TaxCalculator{repo: repo}
}

func (tc *TaxCalculator) GetTaxRateForClient(ctx context.Context, client *domain.Client) (rate float64, taxName string) {
	if client == nil || client.TaxExempt {
		return 0.0, "Tax Exempt"
	}

	if tc.repo == nil {
		return 0.0, ""
	}

	rule, err := tc.repo.GetByLocation(ctx, client.Country, client.State)
	if err != nil || rule == nil {
		return 0.0, ""
	}

	if rule.TaxExempt {
		return 0.0, rule.Name
	}

	return rule.Rate, rule.Name
}

func (tc *TaxCalculator) ListRules(ctx context.Context) ([]*domain.TaxRule, error) {
	return tc.repo.List(ctx)
}

func (tc *TaxCalculator) CreateRule(ctx context.Context, rule *domain.TaxRule) error {
	return tc.repo.Create(ctx, rule)
}

func (tc *TaxCalculator) UpdateRule(ctx context.Context, rule *domain.TaxRule) error {
	return tc.repo.Update(ctx, rule)
}

func (tc *TaxCalculator) DeleteRule(ctx context.Context, id int64) error {
	return tc.repo.Delete(ctx, id)
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
