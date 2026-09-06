package validator

import (
	"errors"
	"strings"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

var allowedPeriods = map[string]bool{
	"1W": true, "2W": true, "1M": true, "3M": true,
	"6M": true, "1Y": true, "2Y": true, "3Y": true,
	"FREE": true, "ONETIME": true,
}

// ValidateOrderPlacement checks order domain constraints
func ValidateOrderPlacement(order *domain.Order) error {
	if order == nil {
		return errors.New("order cannot be nil")
	}
	if order.ClientID == 0 {
		return errors.New("order must belong to a client")
	}
	if order.ProductID == 0 {
		return errors.New("order must reference a product ID")
	}
	if strings.TrimSpace(order.Title) == "" {
		return errors.New("order title is required")
	}
	if !allowedPeriods[order.Period] {
		return errors.New("invalid billing period cycle")
	}
	return nil
}
