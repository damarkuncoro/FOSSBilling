package security_test

import (
	"context"
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/security"
)

func TestDisposableEmailChecker(t *testing.T) {
	checker := security.NewDisposableEmailChecker()

	tests := []struct {
		email    string
		expected bool
	}{
		{"user@mailinator.com", true},
		{"test@10minutemail.com", true},
		{"spam@guerrillamail.com", true},
		{"bot@tempmail.com", true},
		{"client@trashmail.com", true},
		{"legit@gmail.com", false},
		{"corporate@fossbilling.org", false},
		{"admin@yahoo.com", false},
		{"someone@outlook.com", false},
		{"invalid-email-format", false},
	}

	for _, tt := range tests {
		got := checker.IsDisposable(tt.email)
		if got != tt.expected {
			t.Errorf("IsDisposable(%q) = %v; want %v", tt.email, got, tt.expected)
		}
	}

	// Test custom blocked domains
	checker.SetCustomBlocks([]string{"fakecompany.xyz", "spammer.net"})
	if !checker.IsDisposable("john@fakecompany.xyz") {
		t.Errorf("Expected fakecompany.xyz to be detected as disposable")
	}
	if !checker.IsDisposable("bot@spammer.net") {
		t.Errorf("Expected spammer.net to be detected as disposable")
	}
}

func TestStopForumSpamChecker(t *testing.T) {
	checker := security.NewStopForumSpamChecker()
	ctx := context.Background()

	// Clean inputs should not flag
	isSpam, reason, err := checker.CheckSpam(ctx, "safe_user@example.com", "127.0.0.1", 80.0)
	if err != nil {
		t.Logf("StopForumSpam network check returned (fallback safe): %v", err)
	}
	if isSpam {
		t.Errorf("Expected clean user to not be spam, got spam reason: %s", reason)
	}
}

func TestTurnstileVerifier(t *testing.T) {
	verifier := security.NewTurnstileVerifier("1x0000000000000000000000000000000AA") // Cloudflare dummy test secret

	// Empty token should fail
	passed, err := verifier.Verify(context.Background(), "", "127.0.0.1")
	if passed || err == nil {
		t.Error("Expected empty token to fail verification")
	}
}
