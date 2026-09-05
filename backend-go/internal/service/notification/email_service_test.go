package notification_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/service/notification"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/mailer"
)

func TestEmailService_TransactionalEmails(t *testing.T) {
	mockMailer := mailer.NewMockMailer()
	service := notification.NewEmailService(mockMailer, "noreply@fossbilling.org", "Nusantara Cloud")
	ctx := context.Background()

	client := &domain.Client{
		ID:        1,
		Email:     "client.budi@example.com",
		FirstName: "Budi",
		LastName:  "Santoso",
	}

	// 1. Welcome Email
	err := service.SendWelcomeEmail(ctx, client)
	if err != nil {
		t.Fatalf("SendWelcomeEmail error: %v", err)
	}
	if mockMailer.GetSentCount() != 1 {
		t.Fatalf("expected 1 sent email, got %d", mockMailer.GetSentCount())
	}
	last := mockMailer.LastMessage()
	if !strings.Contains(last.Subject, "Selamat Datang") || last.To[0] != client.Email {
		t.Errorf("unexpected email content: %v", last)
	}

	// 2. Invoice Created
	inv := &domain.Invoice{
		ID:       101,
		Nr:       "INV-101",
		Currency: "USD",
		Total:    450000,
		DueAt:    time.Now().AddDate(0, 0, 7),
	}
	err = service.SendInvoiceCreated(ctx, client, inv)
	if err != nil {
		t.Fatalf("SendInvoiceCreated error: %v", err)
	}
	if mockMailer.GetSentCount() != 2 {
		t.Fatalf("expected 2 sent emails, got %d", mockMailer.GetSentCount())
	}

	// 3. Payment Receipt
	txn := &domain.Transaction{
		GatewayID: "midtrans",
		TxnID:     "M-123456",
	}
	err = service.SendPaymentReceipt(ctx, client, inv, txn)
	if err != nil {
		t.Fatalf("SendPaymentReceipt error: %v", err)
	}
	if mockMailer.GetSentCount() != 3 {
		t.Fatalf("expected 3 sent emails, got %d", mockMailer.GetSentCount())
	}

	// 4. Ticket Reply
	ticket := &domain.Ticket{
		ID:      50,
		Subject: "Domain Connection Issue",
	}
	err = service.SendTicketReply(ctx, client, ticket, "Silakan periksa A record domain Anda.")
	if err != nil {
		t.Fatalf("SendTicketReply error: %v", err)
	}
	if mockMailer.GetSentCount() != 4 {
		t.Fatalf("expected 4 sent emails, got %d", mockMailer.GetSentCount())
	}
}
