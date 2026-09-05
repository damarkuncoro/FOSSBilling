package apikey

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/domain"
)

var (
	ErrEmptyKeyName = errors.New("api key name cannot be empty")
)

type APIKeyService struct {
	apiKeyRepo domain.APIKeyRepository
}

func NewAPIKeyService(apiKeyRepo domain.APIKeyRepository) *APIKeyService {
	return &APIKeyService{apiKeyRepo: apiKeyRepo}
}

func (s *APIKeyService) ListKeys(ctx context.Context, clientID int64) ([]*domain.APIKey, error) {
	return s.apiKeyRepo.ListByClientID(ctx, clientID)
}

func (s *APIKeyService) GenerateKey(ctx context.Context, clientID int64, name string, expireDays int) (*domain.APIKey, error) {
	if name == "" {
		return nil, ErrEmptyKeyName
	}

	keyBytes := make([]byte, 16)
	_, _ = rand.Read(keyBytes)
	key := fmt.Sprintf("fb_%s", hex.EncodeToString(keyBytes))

	secretBytes := make([]byte, 32)
	_, _ = rand.Read(secretBytes)
	secret := hex.EncodeToString(secretBytes)

	var expiresAt *time.Time
	if expireDays > 0 {
		exp := time.Now().UTC().AddDate(0, 0, expireDays)
		expiresAt = &exp
	}

	apiKey := &domain.APIKey{
		ClientID:  clientID,
		Name:      name,
		Key:       key,
		Secret:    secret,
		ExpiresAt: expiresAt,
	}

	if err := s.apiKeyRepo.Create(ctx, apiKey); err != nil {
		return nil, err
	}

	return apiKey, nil
}

func (s *APIKeyService) RevokeKey(ctx context.Context, id, clientID int64) error {
	return s.apiKeyRepo.Delete(ctx, id, clientID)
}
