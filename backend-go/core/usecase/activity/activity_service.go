package activity

import (
	"context"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type ActivityService struct{ repo domain.ActivityRepository }

func NewActivityService(r domain.ActivityRepository) *ActivityService { return &ActivityService{r} }

func (s *ActivityService) ListLogs(ctx context.Context, l, o int) ([]*domain.Activity, int, error) {
	if l <= 0 { l = 50 }; return s.repo.List(ctx, l, o)
}

func (s *ActivityService) ListClientLogs(ctx context.Context, cid int64, l, o int) ([]*domain.Activity, int, error) {
	if l <= 0 { l = 20 }; return s.repo.ListByClientID(ctx, cid, l, o)
}

func (s *ActivityService) GetTrend(ctx context.Context, days int) (map[string]int, error) {
	if days <= 0 { days = 7 }
	return s.repo.GetTrend(ctx, days)
}

func (s *ActivityService) log(ctx context.Context, cid, aid *int64, ev, msg, ip string) error {
	return s.repo.Log(ctx, &domain.Activity{ClientID: cid, AdminID: aid, Type: "info", Event: ev, Message: msg, IPAddress: ip})
}

func (s *ActivityService) LogSystemEvent(ctx context.Context, ev, msg string) error { return s.log(ctx, nil, nil, ev, msg, "") }
func (s *ActivityService) LogClientEvent(ctx context.Context, cid int64, ev, msg, ip string) error { return s.log(ctx, &cid, nil, ev, msg, ip) }
func (s *ActivityService) LogAdminEvent(ctx context.Context, aid int64, ev, msg, ip string) error { return s.log(ctx, nil, &aid, ev, msg, ip) }
