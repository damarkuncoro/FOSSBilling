package listener

import (
	"context"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/notification"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
)

type ClientListener struct {
	es *notification.EmailService; cr domain.ClientRepository
}

func NewClientListener(es *notification.EmailService, cr domain.ClientRepository) *ClientListener {
	return &ClientListener{es, cr}
}

func (l *ClientListener) HandleClientRegistered(ctx context.Context, e events.Event) error {
	var id int64
	if p, ok := e.Payload.(domain.ClientRegisteredPayload); ok { id = p.ClientID } else if p, ok := e.Payload.(*domain.ClientRegisteredPayload); ok && p != nil { id = p.ClientID }
	if id == 0 { return nil }

	c, err := l.cr.GetByID(ctx, id); if err != nil || c == nil { return err }
	return l.es.SendWelcomeEmail(ctx, c)
}
