package domain

import (
	"context"
	"time"
)

type BlockedIP struct {
	ID        int64     `json:"id"`
	IP        string    `json:"ip"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}

type AntispamConfig struct {
	CaptchaEnabled        bool         `json:"captcha_enabled"`
	CaptchaProvider       string       `json:"captcha_provider"` // "turnstile", "recaptcha"
	CaptchaSiteKey        string       `json:"captcha_site_key"`
	CaptchaSecret         string       `json:"captcha_secret"`
	StopForumSpamEnabled  bool         `json:"sfs_enabled"`
	SFSMinConfidence      float64      `json:"sfs_min_confidence"` // e.g. 80.0%
	TempEmailBlockEnabled bool         `json:"temp_email_block_enabled"`
	CustomBlockedDomains  []string     `json:"custom_blocked_domains"`
	HoneypotEnabled       bool         `json:"honeypot_enabled"`
	HoneypotFieldName     string       `json:"honeypot_field_name"`
	BlockedIPs            []*BlockedIP `json:"blocked_ips,omitempty"`
}

type AntispamRepository interface {
	GetConfig(ctx context.Context) (*AntispamConfig, error)
	UpdateConfig(ctx context.Context, config *AntispamConfig) error
	ListBlockedIPs(ctx context.Context) ([]*BlockedIP, error)
	AddBlockedIP(ctx context.Context, ip, reason string) (*BlockedIP, error)
	DeleteBlockedIP(ctx context.Context, ip string) error
	IsIPBlocked(ctx context.Context, ip string) (bool, error)
}
