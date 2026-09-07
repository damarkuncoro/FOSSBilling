package memory

import (
	"context"
	"encoding/json"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type MockSystemRepository struct{}

func NewMockSystemRepository() *MockSystemRepository {
	return &MockSystemRepository{}
}

func (m *MockSystemRepository) GetSetting(ctx context.Context, section, key string) (*domain.SystemSetting, error) {
	return nil, nil
}

func (m *MockSystemRepository) ListSettings(ctx context.Context, section string) ([]*domain.SystemSetting, error) {
	if section == "security" {
		return []*domain.SystemSetting{
			{Key: "force_ssl", Value: json.RawMessage("true")},
		}, nil
	}
	return nil, nil
}

func (m *MockSystemRepository) UpdateSetting(ctx context.Context, section, key string, value json.RawMessage) error {
	return nil
}
