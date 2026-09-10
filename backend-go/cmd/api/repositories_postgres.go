package main

import (
	"log"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresRepositories(p *pgxpool.Pool) *Repositories {
	log.Println("📦 Initializing PostgreSQL repositories...")
	return &Repositories{
		Client: postgres.NewClientRepository(p), Order: postgres.NewOrderRepository(p), Invoice: postgres.NewInvoiceRepository(p),
		Transaction: postgres.NewTransactionRepository(p), Promo: postgres.NewPromoRepository(p), Support: postgres.NewSupportRepository(p),
		Staff: postgres.NewStaffRepository(p), Currency: postgres.NewCurrencyRepository(p), News: postgres.NewNewsRepository(p),
		Downloadable: postgres.NewDownloadableRepository(p), APIKey: postgres.NewAPIKeyRepository(p), MassMail: postgres.NewMassMailRepository(p),
		Company: postgres.NewCompanyRepository(p), Product: postgres.NewProductRepository(p), Catalog: postgres.NewCatalogRepository(p),
		System: postgres.NewSystemRepository(p), Activity: postgres.NewActivityRepository(p), Notification: postgres.NewNotificationRepository(p),
		AdminNotification: postgres.NewNotificationRepository(p), Page: postgres.NewPageRepository(p), KB: postgres.NewKBRepository(p),
		Antispam: postgres.NewAntispamRepository(p), Formbuilder: postgres.NewFormbuilderRepository(p), Extension: postgres.NewExtensionRepository(p),
		Redirect: postgres.NewRedirectRepository(p), Theme: memory.NewThemeRepository(), Tax: postgres.NewTaxRepository(p),
		EmailTemplate: postgres.NewEmailTemplateRepository(p),
	}
}
