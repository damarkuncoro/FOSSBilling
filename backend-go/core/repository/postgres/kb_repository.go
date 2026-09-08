package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type KBRepository struct {
	pool *pgxpool.Pool
}

func NewKBRepository(pool *pgxpool.Pool) *KBRepository {
	return &KBRepository{pool: pool}
}

func (r *KBRepository) ListCategories(ctx context.Context) ([]*domain.KBCategory, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, title, slug, description, icon, created_at, updated_at FROM kb_categories ORDER BY title ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.KBCategory
	for rows.Next() {
		cat := &domain.KBCategory{}
		if err := rows.Scan(&cat.ID, &cat.Title, &cat.Slug, &cat.Description, &cat.Icon, &cat.CreatedAt, &cat.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, cat)
	}
	return list, nil
}

func (r *KBRepository) GetCategoryBySlug(ctx context.Context, slug string) (*domain.KBCategory, error) {
	cat := &domain.KBCategory{}
	err := r.pool.QueryRow(ctx, "SELECT id, title, slug, description, icon, created_at, updated_at FROM kb_categories WHERE slug = $1", slug).
		Scan(&cat.ID, &cat.Title, &cat.Slug, &cat.Description, &cat.Icon, &cat.CreatedAt, &cat.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, err
	}
	return cat, nil
}

func (r *KBRepository) CreateCategory(ctx context.Context, cat *domain.KBCategory) error {
	return r.pool.QueryRow(ctx, "INSERT INTO kb_categories (title, slug, description, icon) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at",
		cat.Title, cat.Slug, cat.Description, cat.Icon).Scan(&cat.ID, &cat.CreatedAt, &cat.UpdatedAt)
}

func (r *KBRepository) ListArticles(ctx context.Context, categoryID int64, limit, offset int) ([]*domain.KBArticle, int, error) {
	where := "WHERE 1=1"
	args := []interface{}{limit, offset}
	if categoryID > 0 {
		where += " AND category_id = $3"
		args = append(args, categoryID)
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM kb_articles %s", where)
	var total int
	err := r.pool.QueryRow(ctx, countQuery, args[2:]...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT id, category_id, title, slug, content, status, views, created_at, updated_at
		FROM kb_articles
		%s
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, where)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*domain.KBArticle
	for rows.Next() {
		art := &domain.KBArticle{}
		if err := rows.Scan(&art.ID, &art.CategoryID, &art.Title, &art.Slug, &art.Content, &art.Status, &art.Views, &art.CreatedAt, &art.UpdatedAt); err != nil {
			return nil, 0, err
		}
		list = append(list, art)
	}
	return list, total, nil
}

func (r *KBRepository) GetArticleBySlug(ctx context.Context, slug string) (*domain.KBArticle, error) {
	art := &domain.KBArticle{}
	err := r.pool.QueryRow(ctx, "SELECT id, category_id, title, slug, content, status, views, created_at, updated_at FROM kb_articles WHERE slug = $1", slug).
		Scan(&art.ID, &art.CategoryID, &art.Title, &art.Slug, &art.Content, &art.Status, &art.Views, &art.CreatedAt, &art.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, err
	}
	return art, nil
}

func (r *KBRepository) CreateArticle(ctx context.Context, art *domain.KBArticle) error {
	return r.pool.QueryRow(ctx, "INSERT INTO kb_articles (category_id, title, slug, content, status) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at",
		art.CategoryID, art.Title, art.Slug, art.Content, art.Status).Scan(&art.ID, &art.CreatedAt, &art.UpdatedAt)
}

func (r *KBRepository) UpdateArticle(ctx context.Context, art *domain.KBArticle) error {
	_, err := r.pool.Exec(ctx, "UPDATE kb_articles SET category_id = $1, title = $2, slug = $3, content = $4, status = $5, updated_at = NOW() WHERE id = $6",
		art.CategoryID, art.Title, art.Slug, art.Content, art.Status, art.ID)
	return err
}

func (r *KBRepository) DeleteArticle(ctx context.Context, id int64) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM kb_articles WHERE id = $1", id)
	return err
}

func (r *KBRepository) IncrementViews(ctx context.Context, id int64) error {
	_, err := r.pool.Exec(ctx, "UPDATE kb_articles SET views = views + 1 WHERE id = $1", id)
	return err
}
