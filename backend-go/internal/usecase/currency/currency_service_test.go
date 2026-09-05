package currency_test

import (
	"context"
	"testing"

	"github.com/fossbilling/backend-go/internal/repository/memory"
	"github.com/fossbilling/backend-go/internal/usecase/currency"
)

func TestCurrencyService_CRUD(t *testing.T) {
	repo := memory.NewMockCurrencyRepository()
	service := currency.NewCurrencyService(repo)
	ctx := context.Background()

	// 1. List Currencies
	list, err := service.ListCurrencies(ctx)
	if err != nil || len(list) < 3 {
		t.Fatalf("expected at least 3 currencies, got %d (err: %v)", len(list), err)
	}

	// 2. Create Currency (SGD)
	c, err := service.CreateCurrency(ctx, currency.CreateCurrencyDTO{
		Code:           "SGD",
		Title:          "Singapore Dollar",
		ConversionRate: 1.34,
		Format:         "S$ {{price}}",
	})
	if err != nil {
		t.Fatalf("CreateCurrency failed: %v", err)
	}
	if c.Code != "SGD" {
		t.Errorf("expected code SGD, got %s", c.Code)
	}

	// 3. Update Currency Rate
	updated, err := service.UpdateCurrency(ctx, "SGD", currency.UpdateCurrencyDTO{
		ConversionRate: 1.35,
	})
	if err != nil {
		t.Fatalf("UpdateCurrency failed: %v", err)
	}
	if updated.ConversionRate != 1.35 {
		t.Errorf("expected rate 1.35, got %f", updated.ConversionRate)
	}

	// 4. Set Default Currency
	err = service.SetDefault(ctx, "SGD")
	if err != nil {
		t.Fatalf("SetDefault failed: %v", err)
	}

	// 5. Delete non-default (IDR)
	err = service.DeleteCurrency(ctx, "IDR")
	if err != nil {
		t.Fatalf("DeleteCurrency IDR failed: %v", err)
	}

	// 6. Delete default should fail
	err = service.DeleteCurrency(ctx, "SGD")
	if err != currency.ErrCannotDeleteDefault {
		t.Errorf("expected ErrCannotDeleteDefault, got %v", err)
	}
}
