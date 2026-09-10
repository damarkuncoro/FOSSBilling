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
	Product      domain.ProductRepository
	Catalog      domain.CatalogRepository
	System       domain.SystemRepository
	Activity          domain.ActivityRepository
	Notification      domain.NotificationRepository
	AdminNotification domain.AdminNotificationRepository
	Page              domain.PageRepository
	KB           domain.KBRepository
	Antispam     domain.AntispamRepository
	Formbuilder  domain.FormbuilderRepository
	Extension    domain.ExtensionRepository
	Redirect          domain.RedirectRepository
	Theme             domain.ThemeRepository
	Tax               domain.TaxRepository
	EmailTemplate     domain.EmailTemplateRepository
}

// InitRepositories factory that determines whether to instantiate real PostgreSQL or mock repositories
func InitRepositories(ctx context.Context, cfg *config.Config, pool *pgxpool.Pool) *Repositories {
	if pool != nil {
		return NewPostgresRepositories(pool)
	}

	return NewMockRepositories(ctx)
}
