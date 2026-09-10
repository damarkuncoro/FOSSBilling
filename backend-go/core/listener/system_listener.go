package listener

import (
	"context"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/notification"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
)

type SystemListener struct {
	es *notification.EmailService; adm string
}

func NewSystemListener(es *notification.EmailService, adm string) *SystemListener {
	return &SystemListener{es, adm}
}

func (l *SystemListener) HandleLowStock(ctx context.Context, e events.Event) error {
	p, ok := e.Payload.(*domain.Product); if !ok { return nil }
	return l.es.SendLowStockWarning(ctx, l.adm, p.ID, p.Name, p.Stock)
}
