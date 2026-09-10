package postgres

import (
	"context"
	"errors"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const kbCatCols = `id, title, slug, description, icon, created_at, updated_at`
const kbArtCols = `id, category_id, title, slug, content, status, views, created_at, updated_at`

type KBRepository struct{ pool *pgxpool.Pool }

func NewKBRepository(p *pgxpool.Pool) *KBRepository { return &KBRepository{p} }

func scanKBCat(r pgx.Row) (*domain.KBCategory, error) {
	var c domain.KBCategory; err := r.Scan(&c.ID, &c.Title, &c.Slug, &c.Description, &c.Icon, &c.CreatedAt, &c.UpdatedAt); return &c, err
}

func scanKBArt(r pgx.Row) (*domain.KBArticle, error) {
	var a domain.KBArticle; err := r.Scan(&a.ID, &a.CategoryID, &a.Title, &a.Slug, &a.Content, &a.Status, &a.Views, &a.CreatedAt, &a.UpdatedAt); return &a, err
}

func (r *KBRepository) ListCategories(ctx context.Context) ([]*domain.KBCategory, error) {
	return list(ctx, r.pool, "SELECT "+kbCatCols+" FROM kb_categories ORDER BY title ASC", scanKBCat)
}

func (r *KBRepository) CreateCategory(ctx context.Context, c *domain.KBCategory) error {
	return r.pool.QueryRow(ctx, `INSERT INTO kb_categories (title, slug, description, icon) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`, c.Title, c.Slug, c.Description, c.Icon).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *KBRepository) GetCategoryBySlug(ctx context.Context, sl string) (*domain.KBCategory, error) {
	c, err := scanKBCat(r.pool.QueryRow(ctx, "SELECT "+kbCatCols+" FROM kb_categories WHERE slug = $1", sl))
	if err != nil && errors.Is(err, pgx.ErrNoRows) { return nil, appErrors.ErrNotFound }; return c, err
}

func (r *KBRepository) ListArticles(ctx context.Context, cid int64, l, o int) ([]*domain.KBArticle, int, error) {
	q := "SELECT "+kbArtCols+" FROM kb_articles"; qc := "SELECT COUNT(*) FROM kb_articles"; var args []any
	if cid > 0 { q += " WHERE category_id = $3"; qc += " WHERE category_id = $1"; args = append(args, cid) }
	res, err := list(ctx, r.pool, q+" ORDER BY created_at DESC LIMIT $1 OFFSET $2", scanKBArt, append([]any{l, o}, args...)...)
	return res, total(ctx, r.pool, qc, args...), err
}

func (r *KBRepository) GetArticleBySlug(ctx context.Context, sl string) (*domain.KBArticle, error) {
	a, err := scanKBArt(r.pool.QueryRow(ctx, "SELECT "+kbArtCols+" FROM kb_articles WHERE slug = $1", sl))
	if err != nil && errors.Is(err, pgx.ErrNoRows) { return nil, appErrors.ErrNotFound }; return a, err
}

func (r *KBRepository) CreateArticle(ctx context.Context, a *domain.KBArticle) error {
	return r.pool.QueryRow(ctx, `INSERT INTO kb_articles (category_id, title, slug, content, status) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`, a.CategoryID, a.Title, a.Slug, a.Content, a.Status).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
}

func (r *KBRepository) UpdateArticle(ctx context.Context, a *domain.KBArticle) error {
	_, err := r.pool.Exec(ctx, `UPDATE kb_articles SET category_id = $1, title = $2, slug = $3, content = $4, status = $5, updated_at = NOW() WHERE id = $6`, a.CategoryID, a.Title, a.Slug, a.Content, a.Status, a.ID); return err
}

func (r *KBRepository) DeleteArticle(ctx context.Context, id int64) error { _, err := r.pool.Exec(ctx, `DELETE FROM kb_articles WHERE id = $1`, id); return err }
func (r *KBRepository) IncrementViews(ctx context.Context, id int64) error { _, err := r.pool.Exec(ctx, `UPDATE kb_articles SET views = views + 1 WHERE id = $1`, id); return err }
