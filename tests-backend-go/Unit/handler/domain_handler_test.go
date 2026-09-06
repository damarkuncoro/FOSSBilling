package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	clientHandler "github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/http/client"
	guestHandler "github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/http/guest"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/provisioning"
	domainUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockOrderRepoForDomain struct {
	orders map[int64]*domain.Order
}

func (m *mockOrderRepoForDomain) GetByID(ctx context.Context, id int64) (*domain.Order, error) {
	if o, ok := m.orders[id]; ok {
		return o, nil
	}
	return nil, assert.AnError
}

func (m *mockOrderRepoForDomain) ListByClientID(ctx context.Context, clientID int64, limit, offset int) ([]*domain.Order, int, error) {
	var list []*domain.Order
	for _, o := range m.orders {
		if o.ClientID == clientID {
			list = append(list, o)
		}
	}
	return list, len(list), nil
}

func (m *mockOrderRepoForDomain) List(ctx context.Context, limit, offset int) ([]*domain.Order, int, error) {
	return nil, 0, nil
}
func (m *mockOrderRepoForDomain) ListDueOrders(ctx context.Context, dueBefore time.Time) ([]*domain.Order, error) {
	return nil, nil
}
func (m *mockOrderRepoForDomain) ListOverdueSuspensions(ctx context.Context, overdueDays int) ([]*domain.Order, error) {
	return nil, nil
}
func (m *mockOrderRepoForDomain) Create(ctx context.Context, order *domain.Order) error {
	return nil
}
func (m *mockOrderRepoForDomain) Update(ctx context.Context, order *domain.Order) error {
	m.orders[order.ID] = order
	return nil
}
func (m *mockOrderRepoForDomain) UpdateStatus(ctx context.Context, id int64, status domain.OrderStatus, reason *string) error {
	return nil
}

func TestGuestDomainHandler_CheckAvailability(t *testing.T) {
	registrar := provisioning.NewMockRegistrarDriver()
	mockRepo := &mockOrderRepoForDomain{orders: make(map[int64]*domain.Order)}
	domainService := domainUsecase.NewDomainService(mockRepo, registrar)
	h := guestHandler.NewDomainHandler(domainService)

	t.Run("Valid available domain", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/guest/domains/check?domain=myfreshstartup.com", nil)
		rec := httptest.NewRecorder()

		h.CheckAvailability(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		var res struct {
			Success bool `json:"success"`
			Data    struct {
				Domain    string  `json:"domain"`
				TLD       string  `json:"tld"`
				Available bool    `json:"available"`
				Price     float64 `json:"price"`
			} `json:"data"`
		}
		err := json.NewDecoder(rec.Body).Decode(&res)
		require.NoError(t, err)
		assert.True(t, res.Success)
		assert.Equal(t, "myfreshstartup.com", res.Data.Domain)
		assert.Equal(t, "com", res.Data.TLD)
		assert.True(t, res.Data.Available)
		assert.Equal(t, 12.99, res.Data.Price)
	})

	t.Run("Missing domain query parameter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/guest/domains/check", nil)
		rec := httptest.NewRecorder()

		h.CheckAvailability(rec, req)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Taken domain", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/guest/domains/check?domain=google.com", nil)
		rec := httptest.NewRecorder()

		h.CheckAvailability(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		var res struct {
			Success bool `json:"success"`
			Data    struct {
				Available bool `json:"available"`
			} `json:"data"`
		}
		_ = json.NewDecoder(rec.Body).Decode(&res)
		assert.False(t, res.Data.Available)
	})
}

func TestClientDomainHandler_CRUD(t *testing.T) {
	mockRepo := &mockOrderRepoForDomain{
		orders: map[int64]*domain.Order{
			10: {
				ID:       10,
				ClientID: 1,
				Title:    "Domain Registration: nusantara-cloud.com",
				Status:   domain.OrderStatusActive,
				Period:   "1Y",
				Price:    decimal.FromFloat(12.99),
				Currency: "USD",
				Config:   []byte(`{"domain_name":"nusantara-cloud.com","nameservers":["ns1.fossbilling.org"],"auto_renew":true}`),
			},
		},
	}

	registrar := provisioning.NewMockRegistrarDriver()
	domainService := domainUsecase.NewDomainService(mockRepo, registrar)
	h := clientHandler.NewDomainHandler(domainService)

	t.Run("List client domains authenticated", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/client/domains", nil)
		ctx := middleware.WithClientID(req.Context(), 1)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		h.ListDomains(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		var res struct {
			Success bool                             `json:"success"`
			Data    []domainUsecase.DomainRecordDTO `json:"data"`
		}
		err := json.NewDecoder(rec.Body).Decode(&res)
		require.NoError(t, err)
		assert.True(t, res.Success)
		require.NotEmpty(t, res.Data)
		assert.Equal(t, "nusantara-cloud.com", res.Data[0].DomainName)
		assert.Equal(t, ".com", res.Data[0].TLD)
		assert.True(t, res.Data[0].AutoRenew)
	})

	t.Run("New client with no domains returns empty list", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/client/domains", nil)
		ctx := middleware.WithClientID(req.Context(), 999)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		h.ListDomains(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		var res struct {
			Success bool                             `json:"success"`
			Data    []domainUsecase.DomainRecordDTO `json:"data"`
		}
		err := json.NewDecoder(rec.Body).Decode(&res)
		require.NoError(t, err)
		assert.True(t, res.Success)
		assert.Empty(t, res.Data)
		assert.Equal(t, 0, len(res.Data))
	})

	t.Run("Update nameservers", func(t *testing.T) {
		body := bytes.NewBufferString(`{"nameservers":["ns1.cloudflare.com","ns2.cloudflare.com"]}`)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/client/domains/10/nameservers", body)
		req.SetPathValue("id", "10")
		ctx := middleware.WithClientID(req.Context(), 1)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		h.UpdateNameservers(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		updated, err := mockRepo.GetByID(context.Background(), 10)
		require.NoError(t, err)
		assert.Contains(t, string(updated.Config), "ns1.cloudflare.com")
	})

	t.Run("Toggle auto renew", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/client/domains/10/toggle-autorenew", nil)
		req.SetPathValue("id", "10")
		ctx := middleware.WithClientID(req.Context(), 1)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		h.ToggleAutoRenew(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		updated, err := mockRepo.GetByID(context.Background(), 10)
		require.NoError(t, err)
		assert.Contains(t, string(updated.Config), `"auto_renew":false`)
	})
}
