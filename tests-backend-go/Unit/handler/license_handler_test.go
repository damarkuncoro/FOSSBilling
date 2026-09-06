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
	licenseUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/license"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockOrderRepoForLicense struct {
	orders map[int64]*domain.Order
}

func (m *mockOrderRepoForLicense) GetByID(ctx context.Context, id int64) (*domain.Order, error) {
	if o, ok := m.orders[id]; ok {
		return o, nil
	}
	return nil, assert.AnError
}

func (m *mockOrderRepoForLicense) ListByClientID(ctx context.Context, clientID int64, limit, offset int) ([]*domain.Order, int, error) {
	var list []*domain.Order
	for _, o := range m.orders {
		if o.ClientID == clientID {
			list = append(list, o)
		}
	}
	return list, len(list), nil
}

func (m *mockOrderRepoForLicense) List(ctx context.Context, limit, offset int) ([]*domain.Order, int, error) {
	return nil, 0, nil
}
func (m *mockOrderRepoForLicense) ListDueOrders(ctx context.Context, dueBefore time.Time) ([]*domain.Order, error) {
	return nil, nil
}
func (m *mockOrderRepoForLicense) ListOverdueSuspensions(ctx context.Context, overdueDays int) ([]*domain.Order, error) {
	return nil, nil
}
func (m *mockOrderRepoForLicense) Create(ctx context.Context, order *domain.Order) error {
	return nil
}
func (m *mockOrderRepoForLicense) Update(ctx context.Context, order *domain.Order) error {
	m.orders[order.ID] = order
	return nil
}
func (m *mockOrderRepoForLicense) UpdateStatus(ctx context.Context, id int64, status domain.OrderStatus, reason *string) error {
	return nil
}

func TestClientLicenseHandler(t *testing.T) {
	mockRepo := &mockOrderRepoForLicense{
		orders: map[int64]*domain.Order{
			20: {
				ID:       20,
				ClientID: 1,
				Title:    "FOSSBilling Enterprise Edition License",
				Status:   domain.OrderStatusActive,
				Period:   "1Y",
				Price:    decimal.FromFloat(199.00),
				Currency: "USD",
				Config:   []byte(`{"license_key":"FB-ENT-9941-8999","max_instances":10,"licensed_domain":"app.example.com","licensed_ip":"192.168.1.100"}`),
			},
		},
	}

	licenseService := licenseUsecase.NewLicenseService(mockRepo)
	h := clientHandler.NewLicenseHandler(licenseService)

	t.Run("List client licenses authenticated", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/client/licenses", nil)
		ctx := middleware.WithClientID(req.Context(), 1)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		h.ListLicenses(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		var res struct {
			Success bool                              `json:"success"`
			Data    []licenseUsecase.LicenseRecordDTO `json:"data"`
		}
		err := json.NewDecoder(rec.Body).Decode(&res)
		require.NoError(t, err)
		assert.True(t, res.Success)
		require.NotEmpty(t, res.Data)
		assert.Equal(t, "FOSSBilling Enterprise Edition License", res.Data[0].ProductName)
		assert.Equal(t, "FB-ENT-9941-8999", res.Data[0].LicenseKey)
		assert.Equal(t, "app.example.com", res.Data[0].LicensedDomain)
	})

	t.Run("Reset license lock", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/client/licenses/20/reset", nil)
		req.SetPathValue("id", "20")
		ctx := middleware.WithClientID(req.Context(), 1)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		h.ResetLicenseLock(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		updated, err := mockRepo.GetByID(context.Background(), 20)
		require.NoError(t, err)
		assert.Contains(t, string(updated.Config), `"licensed_domain":""`)
	})
}
