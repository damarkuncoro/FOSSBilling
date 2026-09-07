package provisioning

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/notification"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/mailer"
)

// EmailRegistrarDriver sends notifications to an admin email for manual domain processing.
// It can optionally use RDAP for availability checks.
type EmailRegistrarDriver struct {
	emailService *notification.EmailService
	adminEmail   string
	rdapDriver   *RDAPRegistrarDriver
}

// NewEmailRegistrarDriver creates a new email-based registrar driver
func NewEmailRegistrarDriver(emailService *notification.EmailService, adminEmail string, rdapDriver ...*RDAPRegistrarDriver) *EmailRegistrarDriver {
	var rdap *RDAPRegistrarDriver
	if len(rdapDriver) > 0 {
		rdap = rdapDriver[0]
	}
	return &EmailRegistrarDriver{
		emailService: emailService,
		adminEmail:   adminEmail,
		rdapDriver:   rdap,
	}
}

func (d *EmailRegistrarDriver) CheckAvailability(ctx context.Context, domainName string) (*DomainAvailability, error) {
	if d.rdapDriver != nil {
		return d.rdapDriver.CheckAvailability(ctx, domainName)
	}

	// Fallback to mock behavior if RDAP is not available
	return NewMockRegistrarDriver().CheckAvailability(ctx, domainName)
}

func (d *EmailRegistrarDriver) RegisterDomain(ctx context.Context, req DomainRegistrationRequest) (*DomainRegistrationResult, error) {
	subject := fmt.Sprintf("Register domain: %s", req.DomainName)
	content := fmt.Sprintf("A request to register domain %s for %d years has been received.\n\nNameservers:\n%s",
		req.DomainName, req.Years, strings.Join(req.Nameservers, "\n"))

	err := d.sendEmail(ctx, subject, content)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	return &DomainRegistrationResult{
		DomainName:    req.DomainName,
		Status:        "pending",
		RegisteredAt:  now,
		ExpiresAt:     now.AddDate(req.Years, 0, 0),
		Nameservers:   req.Nameservers,
		TransactionID: fmt.Sprintf("EMAIL-REG-%d", time.Now().UnixNano()),
	}, nil
}

func (d *EmailRegistrarDriver) RenewDomain(ctx context.Context, domainName string, years int) (*DomainRegistrationResult, error) {
	subject := fmt.Sprintf("Renew domain: %s", domainName)
	content := fmt.Sprintf("A request to renew domain %s for %d years has been received.", domainName, years)

	err := d.sendEmail(ctx, subject, content)
	if err != nil {
		return nil, err
	}

	return &DomainRegistrationResult{
		DomainName:    domainName,
		Status:        "pending",
		TransactionID: fmt.Sprintf("EMAIL-REN-%d", time.Now().UnixNano()),
	}, nil
}

func (d *EmailRegistrarDriver) sendEmail(ctx context.Context, subject, content string) error {
	if d.emailService == nil {
		return fmt.Errorf("email service not configured for email registrar")
	}

	// We use the underlying mailer directly to send to admin
	return d.emailService.GetMailer().Send(ctx, mailer.Message{
		To:      []string{d.adminEmail},
		Subject: subject,
		TextBody: content,
	})
}
