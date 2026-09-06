package drivers

import (
	"context"

	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/mailer"
)

type SMTPDriver struct {
	mailer *mailer.SMTPMailer
}

func NewSMTPDriver(host, port, username, password, from string) *SMTPDriver {
	return &SMTPDriver{
		mailer: mailer.NewSMTPMailer(host, port, username, password, from),
	}
}

func (d *SMTPDriver) ID() string   { return "smtp" }
func (d *SMTPDriver) Name() string { return "Standard SMTP Server" }

func (d *SMTPDriver) Send(ctx context.Context, msg mailer.Message) error {
	return d.mailer.Send(ctx, msg)
}
