package security

import (
	"strings"
	"sync"
)

var defaultDisposableDomains = map[string]struct{}{
	"mailinator.com":         {},
	"guerrillamail.com":      {},
	"guerrillamailblock.com": {},
	"10minutemail.com":       {},
	"10minutemail.net":       {},
	"tempmail.com":           {},
	"temp-mail.org":          {},
	"yopmail.com":            {},
	"sharklasers.com":        {},
	"throwawaymail.com":      {},
	"trashmail.com":          {},
	"dispostable.com":        {},
	"getairmail.com":         {},
	"crazymailing.com":       {},
	"burnermail.io":          {},
	"fakeinbox.com":          {},
	"fakemailgenerator.com":  {},
	"emailondeck.com":        {},
	"mohmal.com":             {},
	"mytemp.email":           {},
}

type DisposableEmailChecker struct {
	mu           sync.RWMutex
	customBlocks map[string]struct{}
}

func NewDisposableEmailChecker(customBlocks ...[]string) *DisposableEmailChecker {
	c := &DisposableEmailChecker{
		customBlocks: make(map[string]struct{}),
	}
	if len(customBlocks) > 0 {
		for _, domain := range customBlocks[0] {
			c.customBlocks[strings.ToLower(strings.TrimSpace(domain))] = struct{}{}
		}
	}
	return c
}

func (c *DisposableEmailChecker) IsDisposable(email string) bool {
	parts := strings.Split(strings.TrimSpace(email), "@")
	if len(parts) != 2 {
		return false
	}
	domain := strings.ToLower(parts[1])

	if _, ok := defaultDisposableDomains[domain]; ok {
		return true
	}

	c.mu.RLock()
	defer c.mu.RUnlock()
	if _, ok := c.customBlocks[domain]; ok {
		return true
	}

	return false
}

func (c *DisposableEmailChecker) SetCustomBlocks(domains []string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.customBlocks = make(map[string]struct{})
	for _, domain := range domains {
		d := strings.ToLower(strings.TrimSpace(domain))
		if d != "" {
			c.customBlocks[d] = struct{}{}
		}
	}
}
