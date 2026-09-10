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
	ErrExpired = errors.New("link expired")
	ErrSig     = errors.New("invalid signature")
	ErrNoSvc   = errors.New("no active service")
)

type DownloadLink struct { FileID int64 `json:"file_id"`; Filename string `json:"filename"`; URL string `json:"download_url"`; ExpiresAt time.Time `json:"expires_at"` }

type DownloadableService struct {
	repo domain.DownloadableRepository; orderRepo domain.OrderRepository; secret string
}

func NewDownloadableService(r domain.DownloadableRepository, or domain.OrderRepository, s string) *DownloadableService {
	return &DownloadableService{r, or, s}
}

func (s *DownloadableService) sign(cid, fid, exp int64) string {
	h := hmac.New(sha256.New, []byte(s.secret)); h.Write([]byte(fmt.Sprintf("%d:%d:%d", cid, fid, exp)))
	return hex.EncodeToString(h.Sum(nil))
}

func (s *DownloadableService) GenerateDownloadLink(ctx context.Context, cid, fid int64, ttl time.Duration) (*DownloadLink, error) {
	f, err := s.repo.GetByID(ctx, fid); if err != nil { return nil, err }
	os, _, _ := s.orderRepo.ListByClientID(ctx, cid, 100, 0); active := false
	for _, o := range os { if o.ProductID == f.ProductID && o.Status == domain.OrderStatusActive { active = true; break } }
	if !active { return nil, ErrNoSvc }

	exp := time.Now().UTC().Add(ttl); expU := exp.Unix()
	url := fmt.Sprintf("/api/v1/client/downloads/%d/file?client_id=%d&expires=%d&sig=%s", fid, cid, expU, s.sign(cid, fid, expU))
	return &DownloadLink{FileID: f.ID, Filename: f.Filename, URL: url, ExpiresAt: exp}, nil
}

func (s *DownloadableService) VerifyAndGetFile(ctx context.Context, cid, fid, exp int64, sig string) (*domain.DownloadableFile, error) {
	if time.Now().UTC().Unix() > exp { return nil, ErrExpired }
	if !hmac.Equal([]byte(sig), []byte(s.sign(cid, fid, exp))) { return nil, ErrSig }
	f, err := s.repo.GetByID(ctx, fid); if err != nil { return nil, appErrors.ErrNotFound }
	_ = s.repo.IncrementDownloads(ctx, fid); f.Downloads++; return f, nil
}

type ClientDownloadDTO struct {
	ID                    int64     `json:"id"`
	Title                 string    `json:"title"`
	Category              string    `json:"category"`
	Version               string    `json:"version"`
	FileSize              string    `json:"file_size"`
	Description           string    `json:"description"`
	DownloadURL           string    `json:"download_url"`
	RequiresActiveService bool      `json:"requires_active_service"`
	UpdatedAt             time.Time `json:"updated_at"`
}

func (s *DownloadableService) ListClientDownloads(ctx context.Context, cid int64) ([]ClientDownloadDTO, error) {
	fs, err := s.repo.ListByClientID(ctx, cid); if err != nil { return nil, err }
	res := make([]ClientDownloadDTO, 0, len(fs))
	for _, f := range fs {
		sz := "N/A"; if f.FileSize > 0 {
			v := float64(f.FileSize)
			if f.FileSize >= 1024*1024*1024 { sz = fmt.Sprintf("%.1f GB", v/(1024*1024*1024)) } else if f.FileSize >= 1024*1024 { sz = fmt.Sprintf("%.1f MB", v/(1024*1024)) } else { sz = fmt.Sprintf("%.1f KB", v/1024) }
		}
		expU := time.Now().UTC().Add(24 * time.Hour).Unix()
		url := fmt.Sprintf("/api/v1/client/downloads/%d/file?client_id=%d&expires=%d&sig=%s", f.ID, cid, expU, s.sign(cid, f.ID, expU))
		res = append(res, ClientDownloadDTO{ID: f.ID, Title: f.Filename, Category: "Digital Goods", Version: f.Version, FileSize: sz, Description: "Download " + f.Filename, RequiresActiveService: true, DownloadURL: url, UpdatedAt: f.UpdatedAt})
	}
	return res, nil
}
