package postgres

import (
	"context"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CatalogRepository struct{ pool *pgxpool.Pool }

func NewCatalogRepository(p *pgxpool.Pool) *CatalogRepository { return &CatalogRepository{p} }

func (r *CatalogRepository) ListCategories(ctx context.Context) ([]*domain.ProductCategory, error) {
	return list(ctx, r.pool, `SELECT id, title, slug, description, created_at, updated_at FROM product_categories ORDER BY id ASC`, func(r pgx.Row) (*domain.ProductCategory, error) {
		var c domain.ProductCategory; err := r.Scan(&c.ID, &c.Title, &c.Slug, &c.Description, &c.CreatedAt, &c.UpdatedAt); return &c, err
	})
}

func (r *CatalogRepository) GetCategoryByID(ctx context.Context, id int64) (*domain.ProductCategory, error) {
	var c domain.ProductCategory; err := r.pool.QueryRow(ctx, `SELECT id, title, slug, description, created_at, updated_at FROM product_categories WHERE id = $1`, id).Scan(&c.ID, &c.Title, &c.Slug, &c.Description, &c.CreatedAt, &c.UpdatedAt); return &c, err
}

func scanSrv(r pgx.Row) (*domain.Server, error) {
	var s domain.Server; err := r.Scan(&s.ID, &s.Name, &s.Hostname, &s.IP, &s.Manager, &s.Status, &s.IsDefault, &s.MaxAccounts, &s.CreatedAt, &s.UpdatedAt); return &s, err
}

func (r *CatalogRepository) ListServers(ctx context.Context) ([]*domain.Server, error) {
	return list(ctx, r.pool, `SELECT id, name, hostname, ip, manager, status, is_default, max_accounts, created_at, updated_at FROM servers ORDER BY id ASC`, scanSrv)
}

func (r *CatalogRepository) GetServerByID(ctx context.Context, id int64) (*domain.Server, error) {
	return scanSrv(r.pool.QueryRow(ctx, `SELECT id, name, hostname, ip, manager, status, is_default, max_accounts, created_at, updated_at FROM servers WHERE id = $1`, id))
}

func (r *CatalogRepository) CreateServer(ctx context.Context, s *domain.Server) error {
	return r.pool.QueryRow(ctx, `INSERT INTO servers (name, hostname, ip, manager, status, is_default, max_accounts, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW()) RETURNING id, created_at, updated_at`, s.Name, s.Hostname, s.IP, s.Manager, s.Status, s.IsDefault, s.MaxAccounts).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

func (r *CatalogRepository) UpdateServer(ctx context.Context, s *domain.Server) error {
	_, err := r.pool.Exec(ctx, `UPDATE servers SET name = $1, hostname = $2, ip = $3, manager = $4, status = $5, is_default = $6, max_accounts = $7, updated_at = NOW() WHERE id = $8`, s.Name, s.Hostname, s.IP, s.Manager, s.Status, s.IsDefault, s.MaxAccounts, s.ID); return err
}

func (r *CatalogRepository) DeleteServer(ctx context.Context, id int64) error { _, err := r.pool.Exec(ctx, `DELETE FROM servers WHERE id = $1`, id); return err }

func (r *CatalogRepository) ListTlds(ctx context.Context) ([]*domain.TLD, error) {
	return list(ctx, r.pool, `SELECT id, tld, registrar_id, price_registration, price_renewal, price_transfer, min_years, is_active FROM tlds ORDER BY tld ASC`, func(r pgx.Row) (*domain.TLD, error) {
		var t domain.TLD; err := r.Scan(&t.ID, &t.Tld, &t.RegistrarID, &t.PriceRegistration, &t.PriceRenewal, &t.PriceTransfer, &t.MinYears, &t.IsActive); return &t, err
	})
}
