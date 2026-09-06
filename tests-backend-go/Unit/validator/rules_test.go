package validator_test

import (
	"testing"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/validator"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

func TestClientRules(t *testing.T) {
	// Valid
	validClient := &domain.Client{
		Email:     "valid.client@example.com",
		FirstName: "Jane",
	}
	if err := validator.ValidateClientRegistration(validClient, "SecurePass123!"); err != nil {
		t.Fatalf("expected valid client to pass, got: %v", err)
	}

	// Short password
	if err := validator.ValidateClientRegistration(validClient, "short"); err == nil {
		t.Fatalf("expected error for short password")
	}

	// Invalid email
	invalidEmailClient := &domain.Client{
		Email:     "not-an-email",
		FirstName: "Jane",
	}
	if err := validator.ValidateClientRegistration(invalidEmailClient, "SecurePass123!"); err == nil {
		t.Fatalf("expected error for invalid email")
	}
}

func TestInvoiceRules(t *testing.T) {
	now := time.Now().UTC().Add(48 * time.Hour)
	inv := &domain.Invoice{
		ClientID: 1,
		TaxRate:  11.0,
		DueAt:    now,
	}
	items := []domain.InvoiceItem{
		{Title: "Hosting", Price: decimal.FromFloat(10.0), Quantity: 1},
	}

	if err := validator.ValidateInvoiceCreation(inv, items); err != nil {
		t.Fatalf("expected valid invoice to pass, got: %v", err)
	}

	// No items
	if err := validator.ValidateInvoiceCreation(inv, []domain.InvoiceItem{}); err == nil {
		t.Fatalf("expected error for empty items")
	}

	// Negative tax rate
	inv.TaxRate = -5.0
	if err := validator.ValidateInvoiceCreation(inv, items); err == nil {
		t.Fatalf("expected error for negative tax rate")
	}
}

func TestOrderRules(t *testing.T) {
	validOrder := &domain.Order{
		ClientID:  1,
		ProductID: 101,
		Title:     "cPanel Web Hosting",
		Period:    "1M",
	}

	if err := validator.ValidateOrderPlacement(validOrder); err != nil {
		t.Fatalf("expected valid order to pass, got: %v", err)
	}

	// Invalid period
	validOrder.Period = "99Y"
	if err := validator.ValidateOrderPlacement(validOrder); err == nil {
		t.Fatalf("expected error for invalid period")
	}
}
