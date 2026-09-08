package notification

import (
	"context"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type AdminNotificationService struct {
	repo domain.AdminNotificationRepository
}

func NewAdminNotificationService(repo domain.AdminNotificationRepository) *AdminNotificationService {
	return &AdminNotificationService{repo: repo}
}

func (s *AdminNotificationService) CreateAlert(ctx context.Context, title, message, alertType, module string) error {
	notif := &domain.AdminNotification{
		Title:   title,
		Message: message,
		Type:    alertType,
		Module:  module,
		IsRead:  false,
	}
	return s.repo.Create(ctx, notif)
}

func (s *AdminNotificationService) ListAlerts(ctx context.Context, limit, offset int, unreadOnly bool) ([]*domain.AdminNotification, int, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.repo.List(ctx, limit, offset, unreadOnly)
}

func (s *AdminNotificationService) MarkRead(ctx context.Context, id int64) error {
	return s.repo.MarkAsRead(ctx, id)
}

func (s *AdminNotificationService) MarkAllRead(ctx context.Context) error {
	return s.repo.MarkAllAsRead(ctx)
}
