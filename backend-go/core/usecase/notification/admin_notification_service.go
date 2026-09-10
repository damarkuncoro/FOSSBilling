package notification

import (
	"context"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type AdminNotificationService struct{ repo domain.AdminNotificationRepository }

func NewAdminNotificationService(r domain.AdminNotificationRepository) *AdminNotificationService { return &AdminNotificationService{r} }

func (s *AdminNotificationService) CreateAlert(ctx context.Context, title, msg, tp, mod string) error {
	return s.repo.CreateAdmin(ctx, &domain.AdminNotification{Title: title, Message: msg, Type: tp, Module: mod, IsRead: false})
}

func (s *AdminNotificationService) ListAlerts(ctx context.Context, l, o int, unread bool) ([]*domain.AdminNotification, int, error) {
	if l <= 0 { l = 20 }; return s.repo.ListAdmin(ctx, l, o, unread)
}

func (s *AdminNotificationService) MarkRead(ctx context.Context, id int64) error { return s.repo.MarkAdminAsRead(ctx, id) }
func (s *AdminNotificationService) MarkAllRead(ctx context.Context) error { return s.repo.MarkAdminAllAsRead(ctx) }
