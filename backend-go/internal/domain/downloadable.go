package domain

import (
	"context"
	"time"
)

type DownloadableFile struct {
	ID          int64     `json:"id"`
	ProductID   int64     `json:"product_id"`
	Filename    string    `json:"filename"`
	FilePath    string    `json:"-"`
	FileSize    int64     `json:"file_size"`
	ContentType string    `json:"content_type"`
	Version     string    `json:"version"`
	Downloads   int       `json:"downloads"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DownloadableRepository interface {
	GetByID(ctx context.Context, id int64) (*DownloadableFile, error)
	GetByProductID(ctx context.Context, productID int64) (*DownloadableFile, error)
	ListByClientID(ctx context.Context, clientID int64) ([]*DownloadableFile, error)
	Create(ctx context.Context, file *DownloadableFile) error
	IncrementDownloads(ctx context.Context, id int64) error
	Delete(ctx context.Context, id int64) error
}
