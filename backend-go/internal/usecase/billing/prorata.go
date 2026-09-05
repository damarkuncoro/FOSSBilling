package billing

import (
	"github.com/fossbilling/backend-go/pkg/decimal"
)

// CalculateProrata calculates the price for a partial period
// e.g. upgrading service with 10 days remaining in a 30-day month
func CalculateProrata(unitPrice decimal.Money, totalPeriodDays, activeDays int) decimal.Money {
	if totalPeriodDays <= 0 || activeDays <= 0 {
		return 0
	}
	if activeDays >= totalPeriodDays {
		return unitPrice
	}

	// Exact financial integer calculation with half-up rounding
	numerator := int64(unitPrice) * int64(activeDays)
	denominator := int64(totalPeriodDays)
	// Add half denominator for round half up
	result := (numerator + denominator/2) / denominator
	return decimal.Money(result)
}

// ConvertCurrency converts an amount from a base currency using exchange rates
func ConvertCurrency(amount decimal.Money, fromRate, toRate float64) decimal.Money {
	if fromRate <= 0 || toRate <= 0 {
		return amount
	}
	if fromRate == toRate {
		return amount
	}

	amountFloat := amount.ToFloat()
	baseAmount := amountFloat / fromRate
	targetAmount := baseAmount * toRate
	return decimal.FromFloat(targetAmount)
}
