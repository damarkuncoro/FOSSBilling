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
	ErrIPBlocked   = errors.New("IP blacklisted")
	ErrDispEmail   = errors.New("disposable email not allowed")
	ErrSpam        = errors.New("spam detected")
	ErrBot         = errors.New("bot activity detected")
	ErrCaptcha     = errors.New("captcha failed")
)

type AntispamService struct {
	repo domain.AntispamRepository; email *security.DisposableEmailChecker; sfs *security.StopForumSpamChecker; cap security.CaptchaVerifier
}

func NewAntispamService(r domain.AntispamRepository, cv security.CaptchaVerifier, e *security.DisposableEmailChecker, s *security.StopForumSpamChecker) *AntispamService {
	if e == nil { e = security.NewDisposableEmailChecker() }; if s == nil { s = security.NewStopForumSpamChecker() }
	return &AntispamService{r, e, s, cv}
}

func (s *AntispamService) ValidateSignup(ctx context.Context, em, ip, hp, token string) error {
	cfg, _ := s.repo.GetConfig(ctx)
	if cfg == nil { cfg = &domain.AntispamConfig{TempEmailBlockEnabled: true, StopForumSpamEnabled: true, HoneypotEnabled: true} }

	if ip != "" { if b, _ := s.repo.IsIPBlocked(ctx, ip); b { return ErrIPBlocked } }
	if cfg.HoneypotEnabled && strings.TrimSpace(hp) != "" { return ErrBot }
	if cfg.CaptchaEnabled && s.cap != nil { if p, err := s.cap.Verify(ctx, token, ip); err != nil || !p { return ErrCaptcha } }
	if cfg.TempEmailBlockEnabled && em != "" {
		if len(cfg.CustomBlockedDomains) > 0 { s.email.SetCustomBlocks(cfg.CustomBlockedDomains) }
		if s.email.IsDisposable(em) { return ErrDispEmail }
	}
	if cfg.StopForumSpamEnabled && s.sfs != nil {
		conf := cfg.SFSMinConfidence; if conf <= 0 { conf = 80.0 }
		if sp, r, _ := s.sfs.CheckSpam(ctx, em, ip, conf); sp { return fmt.Errorf("%w: %s", ErrSpam, r) }
	}
	return nil
}

func (s *AntispamService) ValidateTicketSubmission(ctx context.Context, em, ip, token string) error {
	cfg, _ := s.repo.GetConfig(ctx); if cfg == nil { cfg = &domain.AntispamConfig{TempEmailBlockEnabled: true, StopForumSpamEnabled: true} }
	if ip != "" { if b, _ := s.repo.IsIPBlocked(ctx, ip); b { return ErrIPBlocked } }
	if cfg.CaptchaEnabled && s.cap != nil { if p, err := s.cap.Verify(ctx, token, ip); err != nil || !p { return ErrCaptcha } }
	if cfg.TempEmailBlockEnabled && em != "" && s.email.IsDisposable(em) { return ErrDispEmail }
	return nil
}

func (s *AntispamService) IsIPBlocked(ctx context.Context, ip string) (bool, error) { return s.repo.IsIPBlocked(ctx, ip) }
func (s *AntispamService) GetConfig(ctx context.Context) (*domain.AntispamConfig, error) { return s.repo.GetConfig(ctx) }
func (s *AntispamService) UpdateConfig(ctx context.Context, cfg *domain.AntispamConfig) error { return s.repo.UpdateConfig(ctx, cfg) }
func (s *AntispamService) ListBlockedIPs(ctx context.Context) ([]*domain.BlockedIP, error) { return s.repo.ListBlockedIPs(ctx) }

func (s *AntispamService) BlockIP(ctx context.Context, ip, reason string) (*domain.BlockedIP, error) {
	ip = strings.TrimSpace(ip)
	if net.ParseIP(ip) == nil { if _, _, err := net.ParseCIDR(ip); err != nil { return nil, fmt.Errorf("invalid IP/CIDR: %s", ip) } }
	return s.repo.AddBlockedIP(ctx, ip, reason)
}

func (s *AntispamService) UnblockIP(ctx context.Context, ip string) error { return s.repo.DeleteBlockedIP(ctx, strings.TrimSpace(ip)) }
