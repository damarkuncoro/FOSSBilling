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
	mailer       mailer.Mailer
	templateRepo domain.EmailTemplateRepository
	from, app    string
}

func NewEmailService(m mailer.Mailer, tr domain.EmailTemplateRepository, f, a string) *EmailService {
	return &EmailService{m, tr, f, a}
}

func (s *EmailService) GetMailer() mailer.Mailer { return s.mailer }

func (s *EmailService) send(ctx context.Context, to, sub, body string) error {
	return s.mailer.Send(ctx, mailer.Message{From: fmt.Sprintf("%s <%s>", s.app, s.from), To: []string{to}, Subject: sub, HTMLBody: body})
}

func (s *EmailService) exec(tmplStr string, data any) string {
	t, _ := template.New("e").Parse(tmplStr); var b bytes.Buffer; _ = t.Execute(&b, data); return b.String()
}

func (s *EmailService) getTmpl(ctx context.Context, code, defSub, defBody string) (string, string) {
	if s.templateRepo != nil {
		t, err := s.templateRepo.GetByCode(ctx, code)
		if err == nil && t != nil {
			return t.Subject, t.Content
		}
	}
	return defSub, defBody
}

func (s *EmailService) SendWelcomeEmail(ctx context.Context, c *domain.Client) error {
	sub, body := s.getTmpl(ctx, "welcome", "Selamat Datang di {{.AppName}}", welcomeTemplate)
	data := map[string]any{"AppName": s.app, "FirstName": c.FirstName, "Email": c.Email}
	return s.send(ctx, c.Email, s.exec(sub, data), s.exec(body, data))
}

func (s *EmailService) SendInvoiceCreated(ctx context.Context, c *domain.Client, i *domain.Invoice) error {
	sub, body := s.getTmpl(ctx, "invoice_created", "Tagihan Baru #{{.InvoiceNr}}", invoiceCreatedTemplate)
	data := map[string]any{"AppName": s.app, "FirstName": c.FirstName, "InvoiceNr": i.Nr, "Currency": i.Currency, "Total": i.Total.String(), "DueAt": i.DueAt.Format("02 Jan 2006")}
	return s.send(ctx, c.Email, s.exec(sub, data), s.exec(body, data))
}

func (s *EmailService) SendPaymentReceipt(ctx context.Context, c *domain.Client, i *domain.Invoice, t *domain.Transaction) error {
	sub, body := s.getTmpl(ctx, "payment_receipt", "Bukti Pembayaran - #{{.InvoiceNr}}", paymentReceiptTemplate)
	gw, tid := "Pembayaran Online", "-"; if t != nil { gw, tid = t.GatewayID, t.TxnID }
	data := map[string]any{"AppName": s.app, "FirstName": c.FirstName, "InvoiceNr": i.Nr, "Currency": i.Currency, "Total": i.Total.String(), "GatewayID": gw, "TxnID": tid}
	return s.send(ctx, c.Email, s.exec(sub, data), s.exec(body, data))
}

func (s *EmailService) SendTicketReply(ctx context.Context, c *domain.Client, t *domain.Ticket, msg string) error {
	sub, body := s.getTmpl(ctx, "ticket_reply", "[#{{.TicketID}}] New Reply: {{.Subject}}", ticketReplyTemplate)
	data := map[string]any{"AppName": s.app, "FirstName": c.FirstName, "TicketID": t.ID, "Subject": t.Subject, "Message": msg}
	return s.send(ctx, c.Email, s.exec(sub, data), s.exec(body, data))
}

func (s *EmailService) SendServiceActivatedEmail(ctx context.Context, c *domain.Client, o *domain.Order) error {
	sub, body := s.getTmpl(ctx, "service_activated", "Service Activated: {{.Title}}", serviceActivatedTemplate)
	data := map[string]any{"AppName": s.app, "FirstName": c.FirstName, "OrderID": o.ID, "Title": o.Title}
	return s.send(ctx, c.Email, s.exec(sub, data), s.exec(body, data))
}

func (s *EmailService) SendLowStockWarning(ctx context.Context, adm string, pid int64, nm string, st int) error {
	sub, body := s.getTmpl(ctx, "low_stock", "Low Stock Warning: {{.ProductName}}", lowStockTemplate)
	data := map[string]any{"AppName": s.app, "ProductID": pid, "ProductName": nm, "CurrentStock": st}
	return s.send(ctx, adm, s.exec(sub, data), s.exec(body, data))
}

func (s *EmailService) SendInvoiceReminderEmail(ctx context.Context, c *domain.Client, i *domain.Invoice) error {
	sub, body := s.getTmpl(ctx, "invoice_reminder", "Pengingat Tagihan: #{{.InvoiceNr}}", invoiceReminderTemplate)
	data := map[string]any{"AppName": s.app, "FirstName": c.FirstName, "InvoiceNr": i.Nr, "Total": i.Total.String(), "DueAt": i.DueAt.Format("02 Jan 2006")}
	return s.send(ctx, c.Email, s.exec(sub, data), s.exec(body, data))
}
