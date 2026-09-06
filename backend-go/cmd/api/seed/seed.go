package seed

import (
	"context"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
)

type MemoryRepositories struct {
	ClientRepo   *memory.MockClientRepository
	StaffRepo    *memory.MockStaffRepository
	PromoRepo    *memory.MockPromoRepository
	NewsRepo     *memory.MockNewsRepository
	CurrencyRepo *memory.MockCurrencyRepository
	InvoiceRepo  *memory.MockInvoiceRepository
}

// SeedAll populates memory repositories with demo fixtures
func SeedAll(ctx context.Context, repos MemoryRepositories) {
	SeedClients(ctx, repos.ClientRepo)
	SeedInvoices(ctx, repos.InvoiceRepo)
	SeedStaff(ctx, repos.StaffRepo)
	SeedPromos(ctx, repos.PromoRepo)
	SeedCurrencies(ctx, repos.CurrencyRepo)
	SeedNews(ctx, repos.NewsRepo)
}
