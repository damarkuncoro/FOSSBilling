package apikey

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type APIKeyService struct{ repo domain.APIKeyRepository }

func NewAPIKeyService(r domain.APIKeyRepository) *APIKeyService { return &APIKeyService{r} }

func (s *APIKeyService) ListKeys(ctx context.Context, cid int64) ([]*domain.APIKey, error) { return s.repo.ListByClientID(ctx, cid) }

func (s *APIKeyService) GenerateKey(ctx context.Context, cid int64, name string, days int) (*domain.APIKey, error) {
	if name == "" { return nil, errors.New("name required") }
	kb, sb := make([]byte, 16), make([]byte, 32)
	_, _ = rand.Read(kb); _, _ = rand.Read(sb)
	var exp *time.Time; if days > 0 { e := time.Now().UTC().AddDate(0, 0, days); exp = &e }
	ak := &domain.APIKey{ClientID: cid, Name: name, Key: fmt.Sprintf("fb_%s", hex.EncodeToString(kb)), Secret: hex.EncodeToString(sb), ExpiresAt: exp}
	return ak, s.repo.Create(ctx, ak)
}

func (s *APIKeyService) RevokeKey(ctx context.Context, id, cid int64) error { return s.repo.Delete(ctx, id, cid) }
