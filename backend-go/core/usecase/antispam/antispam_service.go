package antispam

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/security"
)

var (
	ErrIPBlocked         = errors.New("access denied: your IP address is blacklisted")
	ErrDisposableEmail   = errors.New("registration with disposable or temporary email domains is not allowed")
	ErrSpamDetected      = errors.New("submission flagged by anti-spam system")
	ErrHoneypotTriggered = errors.New("bot activity detected")
	ErrCaptchaFailed     = errors.New("captcha verification failed, please try again")
)

type AntispamService struct {
	repo            domain.AntispamRepository
	emailChecker    *security.DisposableEmailChecker
	sfsChecker      *security.StopForumSpamChecker
	captchaVerifier security.CaptchaVerifier
}

func NewAntispamService(
	repo domain.AntispamRepository,
	captchaVerifier security.CaptchaVerifier,
	emailChecker *security.DisposableEmailChecker,
	sfsChecker *security.StopForumSpamChecker,
) *AntispamService {
	if emailChecker == nil {
		emailChecker = security.NewDisposableEmailChecker()
	}
	if sfsChecker == nil {
		sfsChecker = security.NewStopForumSpamChecker()
	}
	return &AntispamService{
		repo:            repo,
		captchaVerifier: captchaVerifier,
		emailChecker:    emailChecker,
		sfsChecker:      sfsChecker,
	}
}

// ValidateSignup performs complete multi-layer spam & bot validation on user registrations
func (s *AntispamService) ValidateSignup(ctx context.Context, email, remoteIP, honeypotVal, captchaToken string) error {
	config, _ := s.repo.GetConfig(ctx)
	if config == nil {
		config = &domain.AntispamConfig{
			TempEmailBlockEnabled: true,
			StopForumSpamEnabled:  true,
			HoneypotEnabled:       true,
			HoneypotFieldName:     "website_hp",
		}
	}

	// 1. IP Blacklist check
	if remoteIP != "" {
		blocked, err := s.repo.IsIPBlocked(ctx, remoteIP)
		if err == nil && blocked {
			return ErrIPBlocked
		}
	}

	// 2. Honeypot check (hidden bot field should remain empty)
	if config.HoneypotEnabled && strings.TrimSpace(honeypotVal) != "" {
		return ErrHoneypotTriggered
	}

	// 3. CAPTCHA verification (Cloudflare Turnstile / reCAPTCHA)
	if config.CaptchaEnabled && s.captchaVerifier != nil {
		passed, err := s.captchaVerifier.Verify(ctx, captchaToken, remoteIP)
		if err != nil || !passed {
			return ErrCaptchaFailed
		}
	}

	// 4. Disposable/Temporary Email Check
	if config.TempEmailBlockEnabled && email != "" {
		if len(config.CustomBlockedDomains) > 0 {
			s.emailChecker.SetCustomBlocks(config.CustomBlockedDomains)
		}
		if s.emailChecker.IsDisposable(email) {
			return ErrDisposableEmail
		}
	}

	// 5. StopForumSpam External Lookup
	if config.StopForumSpamEnabled && s.sfsChecker != nil {
		minConfidence := config.SFSMinConfidence
		if minConfidence <= 0 {
			minConfidence = 80.0
		}
		isSpam, reason, _ := s.sfsChecker.CheckSpam(ctx, email, remoteIP, minConfidence)
		if isSpam {
			return fmt.Errorf("%w: %s", ErrSpamDetected, reason)
		}
	}

	return nil
}

// ValidateTicketSubmission validates guest/client ticket opening against spam
func (s *AntispamService) ValidateTicketSubmission(ctx context.Context, email, remoteIP, captchaToken string) error {
	config, _ := s.repo.GetConfig(ctx)
	if config == nil {
		config = &domain.AntispamConfig{
			TempEmailBlockEnabled: true,
			StopForumSpamEnabled:  true,
		}
	}

	if remoteIP != "" {
		blocked, err := s.repo.IsIPBlocked(ctx, remoteIP)
		if err == nil && blocked {
			return ErrIPBlocked
		}
	}

	if config.CaptchaEnabled && s.captchaVerifier != nil {
		passed, err := s.captchaVerifier.Verify(ctx, captchaToken, remoteIP)
		if err != nil || !passed {
			return ErrCaptchaFailed
		}
	}

	if config.TempEmailBlockEnabled && email != "" && s.emailChecker.IsDisposable(email) {
		return ErrDisposableEmail
	}

	return nil
}

func (s *AntispamService) IsIPBlocked(ctx context.Context, ip string) (bool, error) {
	return s.repo.IsIPBlocked(ctx, ip)
}

func (s *AntispamService) GetConfig(ctx context.Context) (*domain.AntispamConfig, error) {
	return s.repo.GetConfig(ctx)
}

func (s *AntispamService) UpdateConfig(ctx context.Context, cfg *domain.AntispamConfig) error {
	return s.repo.UpdateConfig(ctx, cfg)
}

func (s *AntispamService) ListBlockedIPs(ctx context.Context) ([]*domain.BlockedIP, error) {
	return s.repo.ListBlockedIPs(ctx)
}

func (s *AntispamService) BlockIP(ctx context.Context, ip, reason string) (*domain.BlockedIP, error) {
	parsed := net.ParseIP(strings.TrimSpace(ip))
	if parsed == nil {
		// Also support CIDR format
		_, _, err := net.ParseCIDR(strings.TrimSpace(ip))
		if err != nil {
			return nil, fmt.Errorf("invalid IP address or CIDR notation: %s", ip)
		}
	}
	return s.repo.AddBlockedIP(ctx, strings.TrimSpace(ip), reason)
}

func (s *AntispamService) UnblockIP(ctx context.Context, ip string) error {
	return s.repo.DeleteBlockedIP(ctx, strings.TrimSpace(ip))
}
