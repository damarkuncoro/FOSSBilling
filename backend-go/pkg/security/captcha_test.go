package security_test

import (
	"context"
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/security"
)

func TestTurnstileVerifier_MockPassed(t *testing.T) {
	verifier := security.NewTurnstileVerifier("0x4AAAAAAAB...")
	ctx := context.Background()

	valid, err := verifier.Verify(ctx, "MOCK_PASSED_TOKEN", "127.0.0.1")
	if err != nil {
		t.Fatalf("unexpected verify error: %v", err)
	}
	if !valid {
		t.Errorf("expected MOCK_PASSED_TOKEN to return true")
	}

	// Empty secret key passes by default in testing/dev
	emptyVerifier := security.NewTurnstileVerifier("")
	valid2, err := emptyVerifier.Verify(ctx, "any-token", "127.0.0.1")
	if err != nil || !valid2 {
		t.Errorf("expected empty secret to pass")
	}
}
