package main

import (
	"context"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/config"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repositories holds all domain repository interfaces
type Repositories struct {
	Client       domain.ClientRepository
	Order        domain.OrderRepository
	Invoice      domain.InvoiceRepository
	Transaction  domain.TransactionRepository
	Promo        domain.PromoRepository
	Support      domain.SupportRepository
	Staff        domain.StaffRepository
	Currency     domain.CurrencyRepository
	News         domain.NewsRepository
	Downloadable domain.DownloadableRepository
	APIKey       domain.APIKeyRepository
	MassMail     domain.MassMailRepository
	Company      domain.CompanyRepository
}

// InitRepositories factory that determines whether to instantiate real PostgreSQL or mock repositories
func InitRepositories(ctx context.Context, cfg *config.Config, pool *pgxpool.Pool) *Repositories {
	if pool != nil {
		return NewPostgresRepositories(pool)
	}

	return NewMockRepositories(ctx)
}
