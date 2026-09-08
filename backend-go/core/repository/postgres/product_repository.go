package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{pool: pool}
}

func (r *ProductRepository) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	query := `
		SELECT id, category_id, form_id, type, name, slug, description, status, setup_type, config,
		       price_monthly, price_annually, setup_fee, stock, created_at, updated_at
		FROM products
		WHERE id = $1
	`
	p := &domain.Product{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.CategoryID, &p.FormID, &p.Type, &p.Name, &p.Slug, &p.Description,
		&p.Status, &p.SetupType, &p.Config, &p.PriceMonthly, &p.PriceAnnually,
		&p.SetupFee, &p.Stock, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get product by id: %w", err)
	}
	p.Title = p.Name
	p.IsActive = p.Status == "enabled"
	return p, nil
}

func (r *ProductRepository) GetBySlug(ctx context.Context, slug string) (*domain.Product, error) {
	query := `
		SELECT id, category_id, form_id, type, name, slug, description, status, setup_type, config,
		       price_monthly, price_annually, setup_fee, stock, created_at, updated_at
		FROM products
		WHERE slug = $1
	`
	p := &domain.Product{}
	err := r.pool.QueryRow(ctx, query, slug).Scan(
		&p.ID, &p.CategoryID, &p.FormID, &p.Type, &p.Name, &p.Slug, &p.Description,
		&p.Status, &p.SetupType, &p.Config, &p.PriceMonthly, &p.PriceAnnually,
		&p.SetupFee, &p.Stock, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get product by slug: %w", err)
	}
	p.Title = p.Name
	p.IsActive = p.Status == "enabled"
	return p, nil
}

func (r *ProductRepository) List(ctx context.Context, limit, offset int) ([]*domain.Product, int, error) {
	countQuery := `SELECT COUNT(*) FROM products`
	var total int
	if err := r.pool.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	query := `
		SELECT id, category_id, form_id, type, name, slug, description, status, setup_type, config,
		       price_monthly, price_annually, setup_fee, stock, created_at, updated_at
		FROM products
		ORDER BY id ASC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		p := &domain.Product{}
		if err := rows.Scan(
			&p.ID, &p.CategoryID, &p.FormID, &p.Type, &p.Name, &p.Slug, &p.Description,
			&p.Status, &p.SetupType, &p.Config, &p.PriceMonthly, &p.PriceAnnually,
			&p.SetupFee, &p.Stock, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan product: %w", err)
		}
		p.Title = p.Name
		p.IsActive = p.Status == "enabled"
		products = append(products, p)
	}

	return products, total, nil
}

func (r *ProductRepository) Create(ctx context.Context, p *domain.Product) error {
	query := `
		INSERT INTO products (category_id, form_id, type, name, slug, description, status, setup_type, config,
		                      price_monthly, price_annually, setup_fee, stock, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id, created_at, updated_at
	`
	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now

	// Handle compatibility mappings
	if p.Name == "" && p.Title != "" {
		p.Name = p.Title
	}
	if p.Status == "" {
		if p.IsActive {
			p.Status = "enabled"
		} else {
			p.Status = "disabled"
		}
	}

	err := r.pool.QueryRow(ctx, query,
		p.CategoryID, p.FormID, p.Type, p.Name, p.Slug, p.Description,
		p.Status, p.SetupType, p.Config, p.PriceMonthly, p.PriceAnnually,
		p.SetupFee, p.Stock, p.CreatedAt, p.UpdatedAt,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

func (r *ProductRepository) Update(ctx context.Context, p *domain.Product) error {
	query := `
		UPDATE products SET
			category_id = $1, form_id = $2, type = $3, name = $4, slug = $5,
			description = $6, status = $7, setup_type = $8, config = $9,
			price_monthly = $10, price_annually = $11, setup_fee = $12, stock = $13, updated_at = $14
		WHERE id = $15
	`
	p.UpdatedAt = time.Now().UTC()

	if p.Name == "" && p.Title != "" {
		p.Name = p.Title
	}
	if p.Status == "" {
		if p.IsActive {
			p.Status = "enabled"
		} else {
			p.Status = "disabled"
		}
	}

	_, err := r.pool.Exec(ctx, query,
		p.CategoryID, p.FormID, p.Type, p.Name, p.Slug,
		p.Description, p.Status, p.SetupType, p.Config,
		p.PriceMonthly, p.PriceAnnually, p.SetupFee, p.Stock, p.UpdatedAt, p.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}
