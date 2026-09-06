package downloadable

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
)

var (
	ErrDownloadExpired   = errors.New("download link has expired")
	ErrInvalidSignature  = errors.New("invalid download signature")
	ErrProductNotOrdered = errors.New("client does not own an active order for this digital product")
)

type DownloadLink struct {
	FileID    int64     `json:"file_id"`
	Filename  string    `json:"filename"`
	URL       string    `json:"download_url"`
	ExpiresAt time.Time `json:"expires_at"`
}

type DownloadableService struct {
	downloadRepo domain.DownloadableRepository
	orderRepo    domain.OrderRepository
	jwtSecret    string
}

func NewDownloadableService(
	downloadRepo domain.DownloadableRepository,
	orderRepo domain.OrderRepository,
	jwtSecret string,
) *DownloadableService {
	return &DownloadableService{
		downloadRepo: downloadRepo,
		orderRepo:    orderRepo,
		jwtSecret:    jwtSecret,
	}
}

// GenerateDownloadLink verifies client has active order for product and creates secure HMAC signed link
func (s *DownloadableService) GenerateDownloadLink(ctx context.Context, clientID, fileID int64, ttl time.Duration) (*DownloadLink, error) {
	file, err := s.downloadRepo.GetByID(ctx, fileID)
	if err != nil {
		return nil, err
	}

	// Verify client owns active order for this product
	orders, _, err := s.orderRepo.ListByClientID(ctx, clientID, 100, 0)
	if err != nil {
		return nil, err
	}

	hasActiveOrder := false
	for _, o := range orders {
		if o.ProductID == file.ProductID && o.Status == domain.OrderStatusActive {
			hasActiveOrder = true
			break
		}
	}
	if !hasActiveOrder {
		return nil, ErrProductNotOrdered
	}

	expiresAt := time.Now().UTC().Add(ttl)
	expUnix := expiresAt.Unix()

	// HMAC-SHA256 signature
	raw := fmt.Sprintf("%d:%d:%d", clientID, fileID, expUnix)
	h := hmac.New(sha256.New, []byte(s.jwtSecret))
	h.Write([]byte(raw))
	signature := hex.EncodeToString(h.Sum(nil))

	downloadURL := fmt.Sprintf("/api/v1/client/downloads/%d/file?client_id=%d&expires=%d&sig=%s",
		fileID, clientID, expUnix, signature)

	return &DownloadLink{
		FileID:    file.ID,
		Filename:  file.Filename,
		URL:       downloadURL,
		ExpiresAt: expiresAt,
	}, nil
}

// VerifyAndGetFile validates signed URL and increments download counter
func (s *DownloadableService) VerifyAndGetFile(ctx context.Context, clientID, fileID, expiresUnix int64, signature string) (*domain.DownloadableFile, error) {
	if time.Now().UTC().Unix() > expiresUnix {
		return nil, ErrDownloadExpired
	}

	raw := fmt.Sprintf("%d:%d:%d", clientID, fileID, expiresUnix)
	h := hmac.New(sha256.New, []byte(s.jwtSecret))
	h.Write([]byte(raw))
	expectedSig := hex.EncodeToString(h.Sum(nil))

	if !hmac.Equal([]byte(signature), []byte(expectedSig)) {
		return nil, ErrInvalidSignature
	}

	file, err := s.downloadRepo.GetByID(ctx, fileID)
	if err != nil {
		return nil, appErrors.ErrNotFound
	}

	_ = s.downloadRepo.IncrementDownloads(ctx, fileID)
	file.Downloads++
	return file, nil
}

type ClientDownloadDTO struct {
	ID                    int64     `json:"id"`
	Title                 string    `json:"title"`
	Category              string    `json:"category"`
	Version               string    `json:"version"`
	FileSize              string    `json:"file_size"`
	Description           string    `json:"description"`
	RequiresActiveService bool      `json:"requires_active_service"`
	DownloadURL           string    `json:"download_url"`
	UpdatedAt             time.Time `json:"updated_at"`
}

func formatBytes(bytes int64) string {
	if bytes >= 1024*1024*1024 {
		return fmt.Sprintf("%.1f GB", float64(bytes)/(1024*1024*1024))
	}
	if bytes >= 1024*1024 {
		return fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
	}
	if bytes >= 1024 {
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024)
	}
	return fmt.Sprintf("%d B", bytes)
}

func (s *DownloadableService) ListClientDownloads(ctx context.Context, clientID int64) ([]ClientDownloadDTO, error) {
	files, err := s.downloadRepo.ListByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	result := make([]ClientDownloadDTO, 0, len(files))
	for _, f := range files {
		sizeStr := formatBytes(f.FileSize)
		if f.FileSize <= 0 {
			sizeStr = "N/A"
		}

		expiresAt := time.Now().UTC().Add(24 * time.Hour)
		expUnix := expiresAt.Unix()
		raw := fmt.Sprintf("%d:%d:%d", clientID, f.ID, expUnix)
		h := hmac.New(sha256.New, []byte(s.jwtSecret))
		h.Write([]byte(raw))
		signature := hex.EncodeToString(h.Sum(nil))
		downloadURL := fmt.Sprintf("/api/v1/client/downloads/%d/file?client_id=%d&expires=%d&sig=%s",
			f.ID, clientID, expUnix, signature)

		result = append(result, ClientDownloadDTO{
			ID:                    f.ID,
			Title:                 f.Filename,
			Category:              "Digital Goods",
			Version:               f.Version,
			FileSize:              sizeStr,
			Description:           "Product download file for " + f.Filename,
			RequiresActiveService: true,
			DownloadURL:           downloadURL,
			UpdatedAt:             f.UpdatedAt,
		})
	}

	return result, nil
}


