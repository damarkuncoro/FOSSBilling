package listener

import (
	"context"
	"log"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/notification"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
)

type ClientListener struct {
	emailService *notification.EmailService
	clientRepo   domain.ClientRepository
}

func NewClientListener(emailService *notification.EmailService, clientRepo domain.ClientRepository) *ClientListener {
	return &ClientListener{
		emailService: emailService,
		clientRepo:   clientRepo,
	}
}

func (l *ClientListener) HandleClientRegistered(ctx context.Context, e events.Event) error {
	var clientID int64

	switch p := e.Payload.(type) {
	case domain.ClientRegisteredPayload:
		clientID = p.ClientID
	case *domain.ClientRegisteredPayload:
		if p != nil {
			clientID = p.ClientID
		}
	case map[string]interface{}:
		if id, ok := p["client_id"].(int64); ok {
			clientID = id
		}
	}

	if clientID == 0 {
		return nil
	}

	log.Printf("📢 [Event Listener] Client #%d registered. Dispatching welcome email...", clientID)

	client, err := l.clientRepo.GetByID(ctx, clientID)
	if err != nil || client == nil {
		return err
	}

	return l.emailService.SendWelcomeEmail(ctx, client)
}
