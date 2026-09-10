package notification

import (
	"context"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type NotificationService struct{ repo domain.NotificationRepository }

func NewNotificationService(r domain.NotificationRepository) *NotificationService { return &NotificationService{r} }

func (s *NotificationService) CreateNotification(ctx context.Context, cid int64, title, msg, tp string) error {
	return s.repo.Create(ctx, &domain.Notification{ClientID: cid, Title: title, Message: msg, Type: tp, IsRead: false})
}

func (s *NotificationService) ListMyNotifications(ctx context.Context, cid int64, l, o int) ([]*domain.Notification, int, error) {
	if l <= 0 { l = 20 }; return s.repo.ListByClientID(ctx, cid, l, o)
}

func (s *NotificationService) MarkAsRead(ctx context.Context, id int64) error { return s.repo.MarkAsRead(ctx, id) }
func (s *NotificationService) MarkAllAsRead(ctx context.Context, cid int64) error { return s.repo.MarkAllAsRead(ctx, cid) }
