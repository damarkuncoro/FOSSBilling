package seed

import (
	"context"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
)

type MemoryRepositories struct {
	ClientRepo   *memory.MockClientRepository; StaffRepo *memory.MockStaffRepository; PromoRepo *memory.MockPromoRepository
	NewsRepo *memory.MockNewsRepository; CurrencyRepo *memory.MockCurrencyRepository; InvoiceRepo *memory.MockInvoiceRepository
}

func SeedAll(ctx context.Context, r MemoryRepositories) {
	SeedClients(ctx, r.ClientRepo); SeedInvoices(ctx, r.InvoiceRepo); SeedStaff(ctx, r.StaffRepo)
	SeedPromos(ctx, r.PromoRepo); SeedCurrencies(ctx, r.CurrencyRepo); SeedNews(ctx, r.NewsRepo)
}
