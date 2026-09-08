package listener

import (
	"context"
	"log"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/notification"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
)

type SystemListener struct {
	emailService *notification.EmailService
	adminEmail   string
}

func NewSystemListener(emailService *notification.EmailService, adminEmail string) *SystemListener {
	if adminEmail == "" {
		adminEmail = "admin@fossbilling.org"
	}
	return &SystemListener{
		emailService: emailService,
		adminEmail:   adminEmail,
	}
}

func (l *SystemListener) HandleLowStock(ctx context.Context, e events.Event) error {
	prod := e.Payload.(*domain.Product)
	log.Printf("⚠️ [System Alert] Low stock detected for product %s (ID: %d): %d remaining", prod.Name, prod.ID, prod.Stock)

	return l.emailService.SendLowStockWarning(ctx, l.adminEmail, prod.ID, prod.Name, prod.Stock)
}
