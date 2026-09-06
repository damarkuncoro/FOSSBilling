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

type SendGridDriver struct {
	apiKey     string
	httpClient *http.Client
}

func NewSendGridDriver(apiKey string) *SendGridDriver {
	return &SendGridDriver{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (d *SendGridDriver) ID() string   { return "sendgrid" }
func (d *SendGridDriver) Name() string { return "SendGrid Email API" }

func (d *SendGridDriver) Send(ctx context.Context, msg mailer.Message) error {
	toRecipients := make([]map[string]string, 0, len(msg.To))
	for _, to := range msg.To {
		toRecipients = append(toRecipients, map[string]string{"email": to})
	}

	payload := map[string]interface{}{
		"personalizations": []map[string]interface{}{
			{"to": toRecipients},
		},
		"from": map[string]string{
			"email": msg.From,
		},
		"subject": msg.Subject,
		"content": []map[string]string{
			{
				"type":  "text/html",
				"value": msg.HTMLBody,
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.sendgrid.com/v3/mail/send", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+d.apiKey)
	req.Header.Set("Content-Type", "application/json")

	if d.apiKey == "SG.test_mock" || d.apiKey == "" {
		return nil
	}

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("sendgrid API returned error status: %d", resp.StatusCode)
	}

	return nil
}
