package main

import (
	"log"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPostgresRepositories instantiates real database repositories backed by PostgreSQL
func NewPostgresRepositories(pool *pgxpool.Pool) *Repositories {
	log.Println("📦 [Real Data] Initializing PostgreSQL database repositories...")
	return &Repositories{
		Client:       postgres.NewClientRepository(pool),
		Order:        postgres.NewOrderRepository(pool),
		Invoice:      postgres.NewInvoiceRepository(pool),
		Transaction:  postgres.NewTransactionRepository(pool),
		Promo:        postgres.NewPromoRepository(pool),
		Support:      postgres.NewSupportRepository(pool),
		Staff:        postgres.NewStaffRepository(pool),
		Currency:     postgres.NewCurrencyRepository(pool),
		News:         postgres.NewNewsRepository(pool),
		Downloadable: postgres.NewDownloadableRepository(pool),
		APIKey:       postgres.NewAPIKeyRepository(pool),
		MassMail:     postgres.NewMassMailRepository(pool),
		Company:      postgres.NewCompanyRepository(pool),
	}
}
