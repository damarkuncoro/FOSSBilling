package downloadable_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/fossbilling/backend-go/internal/domain"
	"github.com/fossbilling/backend-go/internal/repository/memory"
	"github.com/fossbilling/backend-go/internal/usecase/downloadable"
)

func TestDownloadableService_SignedDownloadFlow(t *testing.T) {
	downloadRepo := memory.NewMockDownloadableRepository()
	orderRepo := memory.NewMockOrderRepository()
	secret := "digital-goods-secret-key-12345"

	service := downloadable.NewDownloadableService(downloadRepo, orderRepo, secret)
	ctx := context.Background()

	// 1. Create Digital File
	file := &domain.DownloadableFile{
		ProductID:   100,
		Filename:    "fossbilling-v1.0.0.zip",
		FilePath:    "/data/files/fossbilling.zip",
		FileSize:    10485760,
		ContentType: "application/zip",
		Version:     "1.0.0",
	}
	_ = downloadRepo.Create(ctx, file)

	// 2. Client without active order should fail
	_, err := service.GenerateDownloadLink(ctx, 1, file.ID, 1*time.Hour)
	if err != downloadable.ErrProductNotOrdered {
		t.Fatalf("expected ErrProductNotOrdered, got: %v", err)
	}

	// 3. Client buys and activates order
	_ = orderRepo.Create(ctx, &domain.Order{
		ClientID:  1,
		ProductID: 100,
		Status:    domain.OrderStatusActive,
	})

	// 4. Generate signed download link
	link, err := service.GenerateDownloadLink(ctx, 1, file.ID, 1*time.Hour)
	if err != nil {
		t.Fatalf("GenerateDownloadLink failed: %v", err)
	}
	if !strings.Contains(link.URL, "sig=") {
		t.Errorf("expected signature in download url: %s", link.URL)
	}

	// 5. Extract signature & verify
	expUnix := link.ExpiresAt.Unix()
	parts := strings.Split(link.URL, "sig=")
	sig := parts[1]

	downloaded, err := service.VerifyAndGetFile(ctx, 1, file.ID, expUnix, sig)
	if err != nil || downloaded.ID != file.ID {
		t.Fatalf("VerifyAndGetFile failed: %v", err)
	}
	if downloaded.Downloads != 1 {
		t.Errorf("expected downloads incremented to 1, got %d", downloaded.Downloads)
	}

	// 6. Invalid signature should fail
	_, err = service.VerifyAndGetFile(ctx, 1, file.ID, expUnix, "invalidsig1234")
	if err != downloadable.ErrInvalidSignature {
		t.Errorf("expected ErrInvalidSignature, got: %v", err)
	}
}
