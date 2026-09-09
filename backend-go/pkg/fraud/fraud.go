package fraud

import (
	"context"
	"strings"
)

type FraudScore struct {
	Score     int      `json:"score"`      // 0-100, higher is riskier
	RiskLevel string   `json:"risk_level"` // neutral, low, medium, high
	Reasons   []string `json:"reasons"`
	IsProxy   bool     `json:"is_proxy"`
	IsVPN     bool     `json:"is_vpn"`
}

type FraudChecker interface {
	CheckIP(ctx context.Context, ip string) (*FraudScore, error)
}

type MockFraudChecker struct {
	BlockedASNs []string
}

func (c *MockFraudChecker) CheckIP(ctx context.Context, ip string) (*FraudScore, error) {
	score := 0
	reasons := []string{}

	// Simulate proxy detection for data center IPs (mock)
	if strings.HasPrefix(ip, "192.168.") || strings.HasPrefix(ip, "10.") || ip == "127.0.0.1" {
		return &FraudScore{Score: 0, RiskLevel: "low", Reasons: []string{"Private/Local IP"}}, nil
	}

	// Mock logic: IPs ending in .1 are considered risky bots
	if strings.HasSuffix(ip, ".1") {
		score = 85
		reasons = append(reasons, "Known botnet pattern detected")
	}

	risk := "low"
	if score > 75 {
		risk = "high"
	} else if score > 30 {
		risk = "medium"
	}

	return &FraudScore{
		Score:     score,
		RiskLevel: risk,
		Reasons:   reasons,
		IsProxy:   score > 50,
	}, nil
}
