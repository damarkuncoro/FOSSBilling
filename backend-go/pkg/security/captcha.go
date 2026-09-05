package security

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)


type CaptchaVerifier interface {
	Verify(ctx context.Context, token, remoteIP string) (bool, error)
}

// TurnstileVerifier validates Cloudflare Turnstile tokens
type TurnstileVerifier struct {
	SecretKey  string
	HTTPClient *http.Client
}

func NewTurnstileVerifier(secretKey string) *TurnstileVerifier {
	return &TurnstileVerifier{
		SecretKey: secretKey,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type turnstileResponse struct {
	Success     bool     `json:"success"`
	ChallengeTs string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
}

func (v *TurnstileVerifier) Verify(ctx context.Context, token, remoteIP string) (bool, error) {
	// If secret key is empty or in mock/testing mode, pass
	if v.SecretKey == "" || token == "MOCK_PASSED_TOKEN" {
		return true, nil
	}

	data := url.Values{}
	data.Set("secret", v.SecretKey)
	data.Set("response", token)
	if remoteIP != "" {
		data.Set("remoteip", remoteIP)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://challenges.cloudflare.com/turnstile/v0/siteverify", strings.NewReader(data.Encode()))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := v.HTTPClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("turnstile request error: %w", err)
	}
	defer resp.Body.Close()


	var res turnstileResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return false, fmt.Errorf("failed to decode turnstile response: %w", err)
	}

	return res.Success, nil
}
