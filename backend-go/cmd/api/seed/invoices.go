package seed

import (
	"context"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

// SeedInvoices populates sample paid and unpaid invoices
func SeedInvoices(ctx context.Context, repo *memory.MockInvoiceRepository) {
	if repo == nil {
		return
	}

	now := time.Now().UTC()
	dueDate := now.Add(14 * 24 * time.Hour)
	paidDate := now.Add(-2 * 24 * time.Hour)

	_ = repo.Create(ctx, &domain.Invoice{
		ID:           1,
		Serie:        "INV",
		Nr:           "2026-0001",
		ClientID:     1,
		Status:       domain.InvoiceStatusPaid,
		Currency:     "USD",
		CurrencyRate: 1.0,
		Subtotal:     decimal.FromFloat(100.0),
		Tax:          decimal.FromFloat(11.0),
		Total:        decimal.FromFloat(111.0),
		TaxRate:      11.0,
		DueAt:        dueDate,
		PaidAt:       &paidDate,
		CreatedAt:    now.Add(-3 * 24 * time.Hour),
		UpdatedAt:    paidDate,
	}, []domain.InvoiceItem{
		{
			ID:        1,
			InvoiceID: 1,
			Title:     "Cloud VPS Premium (2 vCPU, 4GB RAM, 80GB NVMe)",
			Price:     decimal.FromFloat(100.0),
			Quantity:  1,
			Taxable:   true,
			CreatedAt: now.Add(-3 * 24 * time.Hour),
		},
	})

	_ = repo.Create(ctx, &domain.Invoice{
		ID:           2,
		Serie:        "INV",
		Nr:           "2026-0002",
		ClientID:     1,
		Status:       domain.InvoiceStatusUnpaid,
		Currency:     "USD",
		CurrencyRate: 1.0,
		Subtotal:     decimal.FromFloat(49.0),
		Tax:          decimal.FromFloat(5.39),
		Total:        decimal.FromFloat(54.39),
		TaxRate:      11.0,
		DueAt:        dueDate,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, []domain.InvoiceItem{
		{
			ID:        2,
			InvoiceID: 2,
			Title:     "cPanel Web Hosting Professional - 1 Year Renewal",
			Price:     decimal.FromFloat(49.0),
			Quantity:  1,
			Taxable:   true,
			CreatedAt: now,
		},
	})
}
