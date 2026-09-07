package memory

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type MockAntispamRepository struct {
	mu         sync.RWMutex
	config     *domain.AntispamConfig
	blockedIPs map[string]*domain.BlockedIP
	nextID     int64
}

func NewMockAntispamRepository() *MockAntispamRepository {
	return &MockAntispamRepository{
		config: &domain.AntispamConfig{
			CaptchaEnabled:        false,
			CaptchaProvider:       "turnstile",
			StopForumSpamEnabled:  true,
			SFSMinConfidence:      80.0,
			TempEmailBlockEnabled: true,
			HoneypotEnabled:       true,
			HoneypotFieldName:     "website_hp",
		},
		blockedIPs: make(map[string]*domain.BlockedIP),
		nextID:     1,
	}
}

func (m *MockAntispamRepository) GetConfig(ctx context.Context) (*domain.AntispamConfig, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	cfg := *m.config
	return &cfg, nil
}

func (m *MockAntispamRepository) UpdateConfig(ctx context.Context, config *domain.AntispamConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.config = config
	return nil
}

func (m *MockAntispamRepository) ListBlockedIPs(ctx context.Context) ([]*domain.BlockedIP, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var list []*domain.BlockedIP
	for _, b := range m.blockedIPs {
		list = append(list, b)
	}
	return list, nil
}

func (m *MockAntispamRepository) AddBlockedIP(ctx context.Context, ip, reason string) (*domain.BlockedIP, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	cleanIP := strings.TrimSpace(ip)
	b := &domain.BlockedIP{
		ID:        m.nextID,
		IP:        cleanIP,
		Reason:    reason,
		CreatedAt: time.Now().UTC(),
	}
	m.nextID++
	m.blockedIPs[cleanIP] = b
	return b, nil
}

func (m *MockAntispamRepository) DeleteBlockedIP(ctx context.Context, ip string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.blockedIPs, strings.TrimSpace(ip))
	return nil
}

func (m *MockAntispamRepository) IsIPBlocked(ctx context.Context, ip string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.blockedIPs[strings.TrimSpace(ip)]
	return ok, nil
}
