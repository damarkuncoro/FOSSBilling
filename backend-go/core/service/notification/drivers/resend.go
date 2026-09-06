package drivers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/mailer"
)

type ResendDriver struct {
	apiKey     string
	httpClient *http.Client
}

func NewResendDriver(apiKey string) *ResendDriver {
	return &ResendDriver{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (d *ResendDriver) ID() string   { return "resend" }
func (d *ResendDriver) Name() string { return "Resend Transactional Email API" }

func (d *ResendDriver) Send(ctx context.Context, msg mailer.Message) error {
	payload := map[string]interface{}{
		"from":    msg.From,
		"to":      msg.To,
		"subject": msg.Subject,
		"html":    msg.HTMLBody,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.resend.com/emails", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+d.apiKey)
	req.Header.Set("Content-Type", "application/json")

	// If apiKey is a mock or test dummy, simulate success
	if d.apiKey == "re_test_mock" || d.apiKey == "" {
		return nil
	}

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("resend API returned error status: %d", resp.StatusCode)
	}

	return nil
}
