package postgres

import (
	"context"
	"errors"

	"github.com/fossbilling/backend-go/internal/domain"
	appErrors "github.com/fossbilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DownloadableRepository struct {
	pool *pgxpool.Pool
}

func NewDownloadableRepository(pool *pgxpool.Pool) *DownloadableRepository {
	return &DownloadableRepository{pool: pool}
}

func (r *DownloadableRepository) GetByID(ctx context.Context, id int64) (*domain.DownloadableFile, error) {
	query := `
		SELECT id, product_id, filename, file_path, file_size, content_type, version, downloads, created_at, updated_at
		FROM downloadable_files
		WHERE id = $1
	`
	var f domain.DownloadableFile
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&f.ID, &f.ProductID, &f.Filename, &f.FilePath, &f.FileSize, &f.ContentType, &f.Version, &f.Downloads, &f.CreatedAt, &f.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, err
	}
	return &f, nil
}

func (r *DownloadableRepository) GetByProductID(ctx context.Context, productID int64) (*domain.DownloadableFile, error) {
	query := `
		SELECT id, product_id, filename, file_path, file_size, content_type, version, downloads, created_at, updated_at
		FROM downloadable_files
		WHERE product_id = $1
		LIMIT 1
	`
	var f domain.DownloadableFile
	err := r.pool.QueryRow(ctx, query, productID).Scan(
		&f.ID, &f.ProductID, &f.Filename, &f.FilePath, &f.FileSize, &f.ContentType, &f.Version, &f.Downloads, &f.CreatedAt, &f.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, err
	}
	return &f, nil
}

func (r *DownloadableRepository) ListByClientID(ctx context.Context, clientID int64) ([]*domain.DownloadableFile, error) {
	query := `
		SELECT DISTINCT df.id, df.product_id, df.filename, df.file_path, df.file_size, df.content_type, df.version, df.downloads, df.created_at, df.updated_at
		FROM downloadable_files df
		INNER JOIN orders o ON o.product_id = df.product_id
		WHERE o.client_id = $1 AND o.status = 'active'
		ORDER BY df.id ASC
	`
	rows, err := r.pool.Query(ctx, query, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.DownloadableFile
	for rows.Next() {
		var f domain.DownloadableFile
		if err := rows.Scan(&f.ID, &f.ProductID, &f.Filename, &f.FilePath, &f.FileSize, &f.ContentType, &f.Version, &f.Downloads, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, &f)
	}
	return list, rows.Err()
}

func (r *DownloadableRepository) Create(ctx context.Context, f *domain.DownloadableFile) error {
	query := `
		INSERT INTO downloadable_files (product_id, filename, file_path, file_size, content_type, version, downloads, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at
	`
	return r.pool.QueryRow(ctx, query,
		f.ProductID, f.Filename, f.FilePath, f.FileSize, f.ContentType, f.Version, f.Downloads,
	).Scan(&f.ID, &f.CreatedAt, &f.UpdatedAt)
}

func (r *DownloadableRepository) IncrementDownloads(ctx context.Context, id int64) error {
	cmd, err := r.pool.Exec(ctx, `UPDATE downloadable_files SET downloads = downloads + 1, updated_at = CURRENT_TIMESTAMP WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return appErrors.ErrNotFound
	}
	return nil
}

func (r *DownloadableRepository) Delete(ctx context.Context, id int64) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM downloadable_files WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return appErrors.ErrNotFound
	}
	return nil
}
