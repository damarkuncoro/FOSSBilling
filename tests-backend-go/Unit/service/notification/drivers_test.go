package notification_test

import (
	"context"
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/notification"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/notification/drivers"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/mailer"
)

func TestNotificationDriverRegistry(t *testing.T) {
	registry := notification.NewDriverRegistry()

	mockDrv := drivers.NewMockDriver()
	smtpDrv := drivers.NewSMTPDriver("smtp.example.com", "587", "user", "pass", "no-reply@example.com")
	resendDrv := drivers.NewResendDriver("re_test_mock")
	sendgridDrv := drivers.NewSendGridDriver("SG.test_mock")

	registry.Register(mockDrv)
	registry.Register(smtpDrv)
	registry.Register(resendDrv)
	registry.Register(sendgridDrv)

	list := registry.List()
	if len(list) != 4 {
		t.Fatalf("expected 4 drivers, got %d", len(list))
	}

	drv, err := registry.Get("resend")
	if err != nil || drv.ID() != "resend" {
		t.Fatalf("failed to retrieve resend driver: %v", err)
	}

	_, err = registry.Get("nonexistent")
	if err == nil {
		t.Fatalf("expected error for nonexistent driver")
	}

	// Test sending via drivers
	ctx := context.Background()
	msg := mailer.Message{
		From:     "no-reply@fossbilling.org",
		To:       []string{"user@example.com"},
		Subject:  "Test Email",
		HTMLBody: "<p>Hello World</p>",
	}

	if err := mockDrv.Send(ctx, msg); err != nil {
		t.Fatalf("MockDriver send failed: %v", err)
	}
	if mockDrv.Count() != 1 {
		t.Fatalf("expected 1 sent message in mock, got %d", mockDrv.Count())
	}

	if err := resendDrv.Send(ctx, msg); err != nil {
		t.Fatalf("ResendDriver test send failed: %v", err)
	}

	if err := sendgridDrv.Send(ctx, msg); err != nil {
		t.Fatalf("SendGridDriver test send failed: %v", err)
	}
}
