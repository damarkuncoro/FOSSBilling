package notification

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/mailer"
)

type EmailService struct {
	mailer    mailer.Mailer
	fromEmail string
	appName   string
}

func NewEmailService(m mailer.Mailer, fromEmail, appName string) *EmailService {
	if fromEmail == "" {
		fromEmail = "no-reply@fossbilling.org"
	}
	if appName == "" {
		appName = "FOSSBilling"
	}
	return &EmailService{mailer: m, fromEmail: fromEmail, appName: appName}
}

func (s *EmailService) SendWelcomeEmail(ctx context.Context, client *domain.Client) error {
	tmpl, err := template.New("welcome").Parse(welcomeTemplate)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, map[string]interface{}{
		"AppName":   s.appName,
		"FirstName": client.FirstName,
		"Email":     client.Email,
	}); err != nil {
		return err
	}

	return s.mailer.Send(ctx, mailer.Message{
		From:     fmt.Sprintf("%s <%s>", s.appName, s.fromEmail),
		To:       []string{client.Email},
		Subject:  fmt.Sprintf("Selamat Datang di %s!", s.appName),
		HTMLBody: buf.String(),
	})
}

func (s *EmailService) SendInvoiceCreated(ctx context.Context, client *domain.Client, invoice *domain.Invoice) error {
	tmpl, err := template.New("invoiceCreated").Parse(invoiceCreatedTemplate)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, map[string]interface{}{
		"AppName":   s.appName,
		"FirstName": client.FirstName,
		"InvoiceNr": invoice.Nr,
		"Currency":  invoice.Currency,
		"Total":     invoice.Total.String(),
		"DueAt":     invoice.DueAt.Format("02 Jan 2006"),
	}); err != nil {
		return err
	}

	return s.mailer.Send(ctx, mailer.Message{
		From:     fmt.Sprintf("%s Billing <%s>", s.appName, s.fromEmail),
		To:       []string{client.Email},
		Subject:  fmt.Sprintf("Tagihan Baru #%s Diterbitkan", invoice.Nr),
		HTMLBody: buf.String(),
	})
}

func (s *EmailService) SendPaymentReceipt(ctx context.Context, client *domain.Client, invoice *domain.Invoice, txn *domain.Transaction) error {
	tmpl, err := template.New("paymentReceipt").Parse(paymentReceiptTemplate)
	if err != nil {
		return err
	}

	gateway, txnID := "Online Payment", "-"
	if txn != nil {
		gateway = txn.GatewayID
		txnID = txn.TxnID
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, map[string]interface{}{
		"AppName":   s.appName,
		"FirstName": client.FirstName,
		"InvoiceNr": invoice.Nr,
		"Currency":  invoice.Currency,
		"Total":     invoice.Total.String(),
		"GatewayID": gateway,
		"TxnID":     txnID,
	}); err != nil {
		return err
	}

	return s.mailer.Send(ctx, mailer.Message{
		From:     fmt.Sprintf("%s Billing <%s>", s.appName, s.fromEmail),
		To:       []string{client.Email},
		Subject:  fmt.Sprintf("Bukti Pembayaran Lunas - Invoice #%s", invoice.Nr),
		HTMLBody: buf.String(),
	})
}

func (s *EmailService) SendTicketReply(ctx context.Context, client *domain.Client, ticket *domain.Ticket, messageContent string) error {
	tmpl, err := template.New("ticketReply").Parse(ticketReplyTemplate)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, map[string]interface{}{
		"AppName":   s.appName,
		"FirstName": client.FirstName,
		"TicketID":  ticket.ID,
		"Subject":   ticket.Subject,
		"Message":   messageContent,
	}); err != nil {
		return err
	}

	return s.mailer.Send(ctx, mailer.Message{
		From:     fmt.Sprintf("%s Support <%s>", s.appName, s.fromEmail),
		To:       []string{client.Email},
		Subject:  fmt.Sprintf("[#%d] Balasan Baru: %s", ticket.ID, ticket.Subject),
		HTMLBody: buf.String(),
	})
}

func (s *EmailService) SendServiceActivatedEmail(ctx context.Context, client *domain.Client, order *domain.Order) error {
	tmpl, err := template.New("serviceActivated").Parse(serviceActivatedTemplate)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, map[string]interface{}{
		"AppName":   s.appName,
		"FirstName": client.FirstName,
		"OrderID":   order.ID,
		"Title":     order.Title,
	}); err != nil {
		return err
	}

	return s.mailer.Send(ctx, mailer.Message{
		From:     fmt.Sprintf("%s Support <%s>", s.appName, s.fromEmail),
		To:       []string{client.Email},
		Subject:  fmt.Sprintf("Layanan Anda Telah Aktif: %s", order.Title),
		HTMLBody: buf.String(),
	})
}

