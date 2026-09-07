package notification

import (
	"context"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type NotificationService struct {
	repo domain.NotificationRepository
}

func NewNotificationService(repo domain.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) CreateNotification(ctx context.Context, clientID int64, title, message, nType string) error {
	return s.repo.Create(ctx, &domain.Notification{
		ClientID: clientID,
		Title:    title,
		Message:  message,
		Type:     nType,
		IsRead:   false,
	})
}

func (s *NotificationService) ListMyNotifications(ctx context.Context, clientID int64, limit, offset int) ([]*domain.Notification, int, error) {
	if limit <= 0 { limit = 20 }
	return s.repo.ListByClientID(ctx, clientID, limit, offset)
}

func (s *NotificationService) MarkAsRead(ctx context.Context, id int64) error {
	return s.repo.MarkAsRead(ctx, id)
}

func (s *NotificationService) MarkAllAsRead(ctx context.Context, clientID int64) error {
	return s.repo.MarkAllAsRead(ctx, clientID)
}
