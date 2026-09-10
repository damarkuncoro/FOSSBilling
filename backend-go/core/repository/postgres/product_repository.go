package postgres

import (
	"context"
	"errors"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const prodCols = `id, category_id, form_id, type, name, slug, description, status, setup_type, config, price_monthly, price_annually, setup_fee, stock, created_at, updated_at`

type ProductRepository struct{ pool *pgxpool.Pool }

func NewProductRepository(p *pgxpool.Pool) *ProductRepository { return &ProductRepository{p} }

func scanProd(r pgx.Row) (*domain.Product, error) {
	var p domain.Product
	err := r.Scan(&p.ID, &p.CategoryID, &p.FormID, &p.Type, &p.Name, &p.Slug, &p.Description, &p.Status, &p.SetupType, &p.Config, &p.PriceMonthly, &p.PriceAnnually, &p.SetupFee, &p.Stock, &p.CreatedAt, &p.UpdatedAt)
	p.Title, p.IsActive = p.Name, p.Status == "enabled"; return &p, err
}

func (r *ProductRepository) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	p, err := scanProd(r.pool.QueryRow(ctx, "SELECT "+prodCols+" FROM products WHERE id = $1", id))
	if err != nil && errors.Is(err, pgx.ErrNoRows) { return nil, appErrors.ErrNotFound }; return p, err
}

func (r *ProductRepository) GetBySlug(ctx context.Context, sl string) (*domain.Product, error) {
	p, err := scanProd(r.pool.QueryRow(ctx, "SELECT "+prodCols+" FROM products WHERE slug = $1", sl))
	if err != nil && errors.Is(err, pgx.ErrNoRows) { return nil, appErrors.ErrNotFound }; return p, err
}

func (r *ProductRepository) List(ctx context.Context, l, o int) ([]*domain.Product, int, error) {
	res, err := list(ctx, r.pool, "SELECT "+prodCols+" FROM products ORDER BY id ASC LIMIT $1 OFFSET $2", scanProd, l, o)
	return res, total(ctx, r.pool, "SELECT COUNT(*) FROM products"), err
}

func (r *ProductRepository) Create(ctx context.Context, p *domain.Product) error {
	q := `INSERT INTO products (category_id, form_id, type, name, slug, description, status, setup_type, config, price_monthly, price_annually, setup_fee, stock, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id, created_at, updated_at`
	if p.Name == "" { p.Name = p.Title }; if p.Status == "" { p.Status = "enabled" }
	return r.pool.QueryRow(ctx, q, p.CategoryID, p.FormID, p.Type, p.Name, p.Slug, p.Description, p.Status, p.SetupType, p.Config, p.PriceMonthly, p.PriceAnnually, p.SetupFee, p.Stock).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *ProductRepository) Update(ctx context.Context, p *domain.Product) error {
	q := `UPDATE products SET category_id = $1, form_id = $2, type = $3, name = $4, slug = $5, description = $6, status = $7, setup_type = $8, config = $9, price_monthly = $10, price_annually = $11, setup_fee = $12, stock = $13, updated_at = CURRENT_TIMESTAMP WHERE id = $14`
	if p.Name == "" { p.Name = p.Title }; if p.Status == "" { p.Status = "enabled" }
	_, err := r.pool.Exec(ctx, q, p.CategoryID, p.FormID, p.Type, p.Name, p.Slug, p.Description, p.Status, p.SetupType, p.Config, p.PriceMonthly, p.PriceAnnually, p.SetupFee, p.Stock, p.ID)
	return err
}

func (r *ProductRepository) DecrementStock(ctx context.Context, id int64, qty int) error {
	t, err := r.pool.Exec(ctx, `UPDATE products SET stock = stock - $1 WHERE id = $2 AND stock >= $1`, qty, id)
	if err == nil && t.RowsAffected() == 0 { return errors.New("insufficient stock") }; return err
}

func (r *ProductRepository) Delete(ctx context.Context, id int64) error { _, err := r.pool.Exec(ctx, `DELETE FROM products WHERE id = $1`, id); return err }
