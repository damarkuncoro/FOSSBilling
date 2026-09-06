package main

import (
	"context"
	"log"

	"github.com/damarkuncoro/FOSSBilling/backend-go/cmd/api/seed"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
)

// NewMockRepositories instantiates in-memory mock repositories and seeds mock demonstration data
func NewMockRepositories(ctx context.Context) *Repositories {
	log.Println("🧪 [Mock Data] Initializing in-memory mock repositories & seeding mock data...")

	memClient := memory.NewMockClientRepository()
	memStaff := memory.NewMockStaffRepository()
	memPromo := memory.NewMockPromoRepository()
	memNews := memory.NewMockNewsRepository()
	memCurr := memory.NewMockCurrencyRepository()
	memInv := memory.NewMockInvoiceRepository()

	seed.SeedAll(ctx, seed.MemoryRepositories{
		ClientRepo:   memClient,
		StaffRepo:    memStaff,
		PromoRepo:    memPromo,
		NewsRepo:     memNews,
		CurrencyRepo: memCurr,
		InvoiceRepo:  memInv,
	})

	return &Repositories{
		Client:       memClient,
		Order:        memory.NewMockOrderRepository(),
		Invoice:      memInv,
		Transaction:  memory.NewMockTransactionRepository(),
		Promo:        memPromo,
		Support:      memory.NewMockSupportRepository(),
		Staff:        memStaff,
		Currency:     memCurr,
		News:         memNews,
		Downloadable: memory.NewMockDownloadableRepository(),
		APIKey:       memory.NewMockAPIKeyRepository(),
		MassMail:     memory.NewMockMassMailRepository(),
		Company:      memory.NewMockCompanyRepository(),
	}
}
