package seed

import (
	"context"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
)

func SeedInvoices(ctx context.Context, r *memory.MockInvoiceRepository) {
	if r == nil { return }
	n := time.Now().UTC(); da, pa := n.AddDate(0, 0, 14), n.AddDate(0, 0, -2); ca := n.AddDate(0, 0, -3)
	_ = r.Create(ctx, &domain.Invoice{ID: 1, Serie: "INV", Nr: "2026-0001", ClientID: 1, Status: domain.InvoiceStatusPaid, Currency: "USD", CurrencyRate: 1.0, Subtotal: 1000000, Tax: 110000, Total: 1110000, TaxRate: 11.0, DueAt: da, PaidAt: &pa, CreatedAt: ca, UpdatedAt: pa}, []domain.InvoiceItem{{ID: 1, InvoiceID: 1, Title: "VPS Premium", Price: 1000000, Quantity: 1, Taxable: true, CreatedAt: ca}})
	_ = r.Create(ctx, &domain.Invoice{ID: 2, Serie: "INV", Nr: "2026-0002", ClientID: 1, Status: domain.InvoiceStatusUnpaid, Currency: "USD", CurrencyRate: 1.0, Subtotal: 490000, Tax: 53900, Total: 543900, TaxRate: 11.0, DueAt: da, CreatedAt: n, UpdatedAt: n}, []domain.InvoiceItem{{ID: 2, InvoiceID: 2, Title: "Web Hosting", Price: 490000, Quantity: 1, Taxable: true, CreatedAt: n}})
}
