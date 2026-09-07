package activity

import (
	"context"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type ActivityService struct {
	activityRepo domain.ActivityRepository
}

func NewActivityService(activityRepo domain.ActivityRepository) *ActivityService {
	return &ActivityService{activityRepo: activityRepo}
}

func (s *ActivityService) ListLogs(ctx context.Context, limit, offset int) ([]*domain.Activity, int, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.activityRepo.List(ctx, limit, offset)
}

func (s *ActivityService) ListClientLogs(ctx context.Context, clientID int64, limit, offset int) ([]*domain.Activity, int, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.activityRepo.ListByClientID(ctx, clientID, limit, offset)
}

func (s *ActivityService) LogSystemEvent(ctx context.Context, event, message string) error {
	return s.activityRepo.Log(ctx, &domain.Activity{
		Type:    "info",
		Event:   event,
		Message: message,
	})
}

func (s *ActivityService) LogClientEvent(ctx context.Context, clientID int64, event, message, ip string) error {
	return s.activityRepo.Log(ctx, &domain.Activity{
		ClientID:  &clientID,
		Type:      "info",
		Event:     event,
		Message:   message,
		IPAddress: ip,
	})
}

func (s *ActivityService) LogAdminEvent(ctx context.Context, adminID int64, event, message, ip string) error {
	return s.activityRepo.Log(ctx, &domain.Activity{
		AdminID:   &adminID,
		Type:      "info",
		Event:     event,
		Message:   message,
		IPAddress: ip,
	})
}
