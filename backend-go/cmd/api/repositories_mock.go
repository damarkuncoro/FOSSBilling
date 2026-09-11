package main

import (
	"context"
	"log"

	"github.com/damarkuncoro/FOSSBilling/backend-go/cmd/api/seed"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
)

func NewMockRepositories(ctx context.Context) *Repositories {
	log.Println("🧪 Seeding mock data...")
	mc, ms, mp, mn, mcur, mi, mnot := memory.NewMockClientRepository(), memory.NewMockStaffRepository(), memory.NewMockPromoRepository(), memory.NewMockNewsRepository(), memory.NewMockCurrencyRepository(), memory.NewMockInvoiceRepository(), memory.NewMockNotificationRepository()
	seed.SeedAll(ctx, seed.MemoryRepositories{ClientRepo: mc, StaffRepo: ms, PromoRepo: mp, NewsRepo: mn, CurrencyRepo: mcur, InvoiceRepo: mi})

	return &Repositories{
		Client: mc, Staff: ms, Promo: mp, News: mn, Currency: mcur, Invoice: mi, Notification: mnot, AdminNotification: mnot,
		Order: memory.NewMockOrderRepository(), Transaction: memory.NewMockTransactionRepository(), Support: memory.NewMockSupportRepository(),
		Downloadable: memory.NewMockDownloadableRepository(), APIKey: memory.NewMockAPIKeyRepository(), MassMail: memory.NewMockMassMailRepository(),
		Company: memory.NewMockCompanyRepository(), Product: memory.NewMockProductRepository(), Catalog: memory.NewMockCatalogRepository(),
		System: memory.NewMockSystemRepository(), Activity: memory.NewMockActivityRepository(), Page: memory.NewMockPageRepository(),
		KB: memory.NewMockKBRepository(), Antispam: memory.NewMockAntispamRepository(), Formbuilder: memory.NewMockFormbuilderRepository(),
		Extension: memory.NewMockExtensionRepository(), Redirect: memory.NewRedirectRepository(), Theme: memory.NewThemeRepository(), Tax: memory.NewMockTaxRepository(),
		EmailTemplate: memory.NewMockEmailTemplateRepository(),
		Affiliate:     memory.NewMockAffiliateRepository(),
	}
}
