package antispam_test

import (
	"context"
	"errors"
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/antispam"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/security"
)

type MockCaptchaVerifier struct {
	ShouldPass bool
}

func (m *MockCaptchaVerifier) Verify(ctx context.Context, responseToken, remoteIP string) (bool, error) {
	if !m.ShouldPass {
		return false, errors.New("invalid captcha token")
	}
	return true, nil
}

func setupAntispamService() (*antispam.AntispamService, *memory.MockAntispamRepository) {
	repo := memory.NewMockAntispamRepository()
	captcha := &MockCaptchaVerifier{ShouldPass: true}
	emailChecker := security.NewDisposableEmailChecker()
	sfsChecker := security.NewStopForumSpamChecker()
	svc := antispam.NewAntispamService(repo, captcha, emailChecker, sfsChecker)
	return svc, repo
}

func TestAntispamService_ValidateSignup(t *testing.T) {
	ctx := context.Background()
	svc, repo := setupAntispamService()

	// 1. Valid registration
	err := svc.ValidateSignup(ctx, "valid.user@gmail.com", "192.168.1.10", "", "valid_token")
	if err != nil {
		t.Fatalf("Expected valid signup to pass, got: %v", err)
	}

	// 2. Honeypot triggered
	err = svc.ValidateSignup(ctx, "bot@gmail.com", "192.168.1.10", "http://spamwebsite.com", "token")
	if !errors.Is(err, antispam.ErrHoneypotTriggered) {
		t.Errorf("Expected ErrHoneypotTriggered, got: %v", err)
	}

	// 3. Disposable email triggered
	err = svc.ValidateSignup(ctx, "spammer@mailinator.com", "192.168.1.10", "", "token")
	if !errors.Is(err, antispam.ErrDisposableEmail) {
		t.Errorf("Expected ErrDisposableEmail, got: %v", err)
	}

	// 4. IP Blacklist triggered
	_, _ = repo.AddBlockedIP(ctx, "203.0.113.5", "Abuse detected")
	err = svc.ValidateSignup(ctx, "legit@company.org", "203.0.113.5", "", "token")
	if !errors.Is(err, antispam.ErrIPBlocked) {
		t.Errorf("Expected ErrIPBlocked, got: %v", err)
	}
}

func TestAntispamService_IPManagement(t *testing.T) {
	ctx := context.Background()
	svc, _ := setupAntispamService()

	// Add IP
	blocked, err := svc.BlockIP(ctx, "198.51.100.42", "Port scanning")
	if err != nil {
		t.Fatalf("BlockIP failed: %v", err)
	}
	if blocked.IP != "198.51.100.42" {
		t.Errorf("Blocked IP = %s; want 198.51.100.42", blocked.IP)
	}

	// Verify is blocked
	isBlocked, err := svc.IsIPBlocked(ctx, "198.51.100.42")
	if err != nil || !isBlocked {
		t.Errorf("Expected IP to be blocked")
	}

	// List blocked IPs
	list, err := svc.ListBlockedIPs(ctx)
	if err != nil || len(list) != 1 {
		t.Errorf("Expected 1 blocked IP in list, got %d", len(list))
	}

	// Unblock IP
	err = svc.UnblockIP(ctx, "198.51.100.42")
	if err != nil {
		t.Fatalf("UnblockIP failed: %v", err)
	}

	isBlocked, _ = svc.IsIPBlocked(ctx, "198.51.100.42")
	if isBlocked {
		t.Errorf("Expected IP to no longer be blocked")
	}
}

func TestAntispamService_Config(t *testing.T) {
	ctx := context.Background()
	svc, _ := setupAntispamService()

	cfg, err := svc.GetConfig(ctx)
	if err != nil {
		t.Fatalf("GetConfig failed: %v", err)
	}
	if !cfg.TempEmailBlockEnabled {
		t.Errorf("Expected TempEmailBlockEnabled to be true by default")
	}

	newCfg := &domain.AntispamConfig{
		TempEmailBlockEnabled: true,
		StopForumSpamEnabled:  true,
		SFSMinConfidence:      90.0,
		CustomBlockedDomains:  []string{"badmail.xyz"},
	}
	err = svc.UpdateConfig(ctx, newCfg)
	if err != nil {
		t.Fatalf("UpdateConfig failed: %v", err)
	}

	updated, _ := svc.GetConfig(ctx)
	if updated.SFSMinConfidence != 90.0 {
		t.Errorf("SFSMinConfidence = %f; want 90.0", updated.SFSMinConfidence)
	}
	if len(updated.CustomBlockedDomains) != 1 || updated.CustomBlockedDomains[0] != "badmail.xyz" {
		t.Errorf("CustomBlockedDomains not updated properly")
	}
}
