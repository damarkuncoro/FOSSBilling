package http_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	clientHandler "github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/http/client"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/downloadable"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockDownloadRepoForHandler struct {
	files map[int64]*domain.DownloadableFile
}

func (m *mockDownloadRepoForHandler) GetByID(ctx context.Context, id int64) (*domain.DownloadableFile, error) {
	if f, ok := m.files[id]; ok {
		return f, nil
	}
	return nil, appErrors.ErrNotFound
}
func (m *mockDownloadRepoForHandler) GetByProductID(ctx context.Context, productID int64) (*domain.DownloadableFile, error) {
	return nil, nil
}
func (m *mockDownloadRepoForHandler) ListByClientID(ctx context.Context, clientID int64) ([]*domain.DownloadableFile, error) {
	var list []*domain.DownloadableFile
	if clientID == 1 {
		for _, f := range m.files {
			list = append(list, f)
		}
	}
	return list, nil
}
func (m *mockDownloadRepoForHandler) Create(ctx context.Context, file *domain.DownloadableFile) error {
	return nil
}
func (m *mockDownloadRepoForHandler) IncrementDownloads(ctx context.Context, id int64) error {
	return nil
}
func (m *mockDownloadRepoForHandler) Delete(ctx context.Context, id int64) error {
	return nil
}

func TestClientDownloadHandler_ListDownloads(t *testing.T) {
	mockRepo := &mockDownloadRepoForHandler{
		files: map[int64]*domain.DownloadableFile{
			1: {
				ID:          1,
				ProductID:   10,
				Filename:    "fossbilling-agent.zip",
				FileSize:    10485760, // 10 MB
				ContentType: "application/zip",
				Version:     "2.0.1",
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
			},
		},
	}
	mockOrderRepo := &mockOrderRepoForLicense{orders: map[int64]*domain.Order{}}

	downloadSvc := downloadable.NewDownloadableService(mockRepo, mockOrderRepo, "test-secret")
	h := clientHandler.NewDownloadHandler(downloadSvc)

	t.Run("List downloads for client with active digital product", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/client/downloads", nil)
		ctx := middleware.WithClientID(req.Context(), 1)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		h.ListDownloads(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		var res struct {
			Success bool                             `json:"success"`
			Data    []downloadable.ClientDownloadDTO `json:"data"`
		}
		err := json.NewDecoder(rec.Body).Decode(&res)
		require.NoError(t, err)
		assert.True(t, res.Success)
		require.Len(t, res.Data, 1)
		assert.Equal(t, "fossbilling-agent.zip", res.Data[0].Title)
		assert.Equal(t, "10.0 MB", res.Data[0].FileSize)
		assert.Contains(t, res.Data[0].DownloadURL, "/api/v1/client/downloads/1/file")
	})

	t.Run("Client with no downloads receives empty array", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/client/downloads", nil)
		ctx := middleware.WithClientID(req.Context(), 999)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		h.ListDownloads(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		var res struct {
			Success bool                             `json:"success"`
			Data    []downloadable.ClientDownloadDTO `json:"data"`
		}
		err := json.NewDecoder(rec.Body).Decode(&res)
		require.NoError(t, err)
		assert.True(t, res.Success)
		assert.Empty(t, res.Data)
	})
}
