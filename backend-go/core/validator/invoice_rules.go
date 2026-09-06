package validator

import (
	"errors"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

// ValidateInvoiceCreation checks invoice domain constraints
func ValidateInvoiceCreation(inv *domain.Invoice, items []domain.InvoiceItem) error {
	if inv == nil {
		return errors.New("invoice cannot be nil")
	}
	if inv.ClientID == 0 {
		return errors.New("invoice must be associated with a valid client ID")
	}
	if len(items) == 0 {
		return errors.New("invoice must contain at least one line item")
	}
	if inv.TaxRate < 0.0 || inv.TaxRate > 100.0 {
		return errors.New("tax rate must be between 0 and 100 percent")
	}
	if inv.DueAt.Before(time.Now().UTC().Add(-24 * time.Hour)) {
		return errors.New("invoice due date cannot be in the past")
	}
	return nil
}
