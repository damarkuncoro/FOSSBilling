package notification

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"github.com/fossbilling/backend-go/internal/domain"
	"github.com/fossbilling/backend-go/pkg/mailer"
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
	return &EmailService{
		mailer:    m,
		fromEmail: fromEmail,
		appName:   appName,
	}
}

const welcomeTemplate = `
<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
  <div style="max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #e0e0e0; border-radius: 8px;">
    <h2 style="color: #206bc4;">Selamat Datang di {{.AppName}}, {{.FirstName}}!</h2>
    <p>Akun Anda telah berhasil dibuat. Anda sekarang dapat masuk ke Client Portal untuk mengelola layanan, tagihan, dan tiket bantuan Anda.</p>
    <p><strong>Email Login:</strong> {{.Email}}</p>
    <hr style="border: 0; border-top: 1px solid #eee; margin: 20px 0;">
    <p style="font-size: 12px; color: #888;">&copy; {{.AppName}}. All rights reserved.</p>
  </div>
</body>
</html>
`

const invoiceCreatedTemplate = `
<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
  <div style="max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #e0e0e0; border-radius: 8px;">
    <h2 style="color: #206bc4;">Tagihan Baru Terbit #{{.InvoiceNr}}</h2>
    <p>Halo {{.FirstName}},</p>
    <p>Tagihan baru Anda telah diterbitkan dengan rincian sebagai berikut:</p>
    <ul>
      <li><strong>Nomor Invoice:</strong> {{.InvoiceNr}}</li>
      <li><strong>Total Tagihan:</strong> {{.Currency}} {{.Total}}</li>
      <li><strong>Batas Pembayaran:</strong> {{.DueAt}}</li>
      <li><strong>Status:</strong> Menunggu Pembayaran (Unpaid)</li>
    </ul>
    <p>Silakan lakukan pembayaran sebelum tanggal jatuh tempo untuk memastikan layanan Anda tetap aktif.</p>
    <hr style="border: 0; border-top: 1px solid #eee; margin: 20px 0;">
    <p style="font-size: 12px; color: #888;">&copy; {{.AppName}}</p>
  </div>
</body>
</html>
`

const paymentReceiptTemplate = `
<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
  <div style="max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #e0e0e0; border-radius: 8px;">
    <h2 style="color: #2fb344;">Bukti Pembayaran Lunas #{{.InvoiceNr}}</h2>
    <p>Halo {{.FirstName}},</p>
    <p>Terima kasih! Pembayaran Anda telah kami terima dan diverifikasi:</p>
    <ul>
      <li><strong>Nomor Invoice:</strong> {{.InvoiceNr}}</li>
      <li><strong>Jumlah Dibayar:</strong> {{.Currency}} {{.Total}}</li>
      <li><strong>Metode/Gateway:</strong> {{.GatewayID}}</li>
      <li><strong>ID Transaksi:</strong> {{.TxnID}}</li>
    </ul>
    <p>Layanan Anda telah otomatis diaktifkan / diperpanjang.</p>
    <hr style="border: 0; border-top: 1px solid #eee; margin: 20px 0;">
    <p style="font-size: 12px; color: #888;">&copy; {{.AppName}}</p>
  </div>
</body>
</html>
`

const ticketReplyTemplate = `
<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
  <div style="max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #e0e0e0; border-radius: 8px;">
    <h2 style="color: #206bc4;">Balasan Baru pada Tiket Bantuan [#{{.TicketID}}]</h2>
    <p>Halo {{.FirstName}},</p>
    <p>Staf kami telah membalas tiket bantuan Anda <strong>"{{.Subject}}"</strong>:</p>
    <div style="background: #f8fafc; border-left: 4px solid #206bc4; padding: 12px; margin: 15px 0;">
      {{.Message}}
    </div>
    <p>Anda dapat membalas tiket ini secara langsung melalui Client Portal.</p>
    <hr style="border: 0; border-top: 1px solid #eee; margin: 20px 0;">
    <p style="font-size: 12px; color: #888;">&copy; {{.AppName}}</p>
  </div>
</body>
</html>
`

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

	gateway := "Online Payment"
	txnID := "-"
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
