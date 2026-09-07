package postgres

import (
	"context"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CatalogRepository struct {
	pool *pgxpool.Pool
}

func NewCatalogRepository(pool *pgxpool.Pool) *CatalogRepository {
	return &CatalogRepository{pool: pool}
}

func (r *CatalogRepository) ListCategories(ctx context.Context) ([]*domain.ProductCategory, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, title, slug, description, created_at, updated_at FROM product_categories ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*domain.ProductCategory
	for rows.Next() {
		c := &domain.ProductCategory{}
		if err := rows.Scan(&c.ID, &c.Title, &c.Slug, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *CatalogRepository) GetCategoryByID(ctx context.Context, id int64) (*domain.ProductCategory, error) {
	c := &domain.ProductCategory{}
	err := r.pool.QueryRow(ctx, `SELECT id, title, slug, description, created_at, updated_at FROM product_categories WHERE id = $1`, id).
		Scan(&c.ID, &c.Title, &c.Slug, &c.Description, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *CatalogRepository) ListServers(ctx context.Context) ([]*domain.Server, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name, hostname, ip, manager, status, is_default, max_accounts, created_at, updated_at FROM servers ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var servers []*domain.Server
	for rows.Next() {
		s := &domain.Server{}
		if err := rows.Scan(&s.ID, &s.Name, &s.Hostname, &s.IP, &s.Manager, &s.Status, &s.IsDefault, &s.MaxAccounts, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		servers = append(servers, s)
	}
	return servers, nil
}

func (r *CatalogRepository) GetServerByID(ctx context.Context, id int64) (*domain.Server, error) {
	s := &domain.Server{}
	err := r.pool.QueryRow(ctx, `SELECT id, name, hostname, ip, manager, status, is_default, max_accounts, created_at, updated_at FROM servers WHERE id = $1`, id).
		Scan(&s.ID, &s.Name, &s.Hostname, &s.IP, &s.Manager, &s.Status, &s.IsDefault, &s.MaxAccounts, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *CatalogRepository) ListTlds(ctx context.Context) ([]*domain.TLD, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, tld, registrar_id, price_registration, price_renewal, price_transfer, min_years, is_active FROM tlds ORDER BY tld ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tlds []*domain.TLD
	for rows.Next() {
		t := &domain.TLD{}
		if err := rows.Scan(&t.ID, &t.Tld, &t.RegistrarID, &t.PriceRegistration, &t.PriceRenewal, &t.PriceTransfer, &t.MinYears, &t.IsActive); err != nil {
			return nil, err
		}
		tlds = append(tlds, t)
	}
	return tlds, nil
}
