package fraud

import (
	"context"
	"strings"
)

type FraudScore struct { Score int; RiskLevel string; Reasons []string; IsProxy, IsVPN bool }
type FraudChecker interface { CheckIP(ctx context.Context, ip string) (*FraudScore, error) }
type MockFraudChecker struct{}

func (c *MockFraudChecker) CheckIP(_ context.Context, ip string) (*FraudScore, error) {
	if strings.HasPrefix(ip, "192.168.") || strings.HasPrefix(ip, "10.") || ip == "127.0.0.1" { return &FraudScore{RiskLevel: "low"}, nil }
	s := 0; if strings.HasSuffix(ip, ".1") { s = 85 }
	rl := "low"; if s > 75 { rl = "high" } else if s > 30 { rl = "medium" }
	return &FraudScore{Score: s, RiskLevel: rl, IsProxy: s > 50}, nil
}
