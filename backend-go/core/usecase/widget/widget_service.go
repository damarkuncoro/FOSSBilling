package widget

import (
	"context"
	"sort"
	"sync"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type WidgetService struct {
	mu sync.RWMutex; ws map[string]*domain.Widget
}

func NewWidgetService() *WidgetService {
	s := &WidgetService{ws: make(map[string]*domain.Widget)}
	s.RegisterWidget(&domain.Widget{ID: "branding_footer", Slot: "client.footer", Title: "FOSSBilling", Component: "BrandingFooter", Priority: 100})
	s.RegisterWidget(&domain.Widget{ID: "cookie_banner", Slot: "client.banner", Title: "GDPR", Component: "CookieBanner", Priority: 10})
	s.RegisterWidget(&domain.Widget{ID: "dashboard_stats", Slot: "admin.dashboard.top", Title: "Metrics", Component: "Stats", Priority: 1})
	return s
}

func (s *WidgetService) RegisterWidget(w *domain.Widget) { s.mu.Lock(); defer s.mu.Unlock(); if w.Priority <= 0 { w.Priority = 10 }; s.ws[w.ID] = w }
func (s *WidgetService) UnregisterWidget(id string) { s.mu.Lock(); defer s.mu.Unlock(); delete(s.ws, id) }

func (s *WidgetService) GetWidgetsForSlot(ctx context.Context, slot string) []*domain.Widget {
	s.mu.RLock(); defer s.mu.RUnlock(); var res []*domain.Widget
	for _, w := range s.ws { if w.Slot == slot { c := *w; res = append(res, &c) } }
	sort.Slice(res, func(i, j int) bool { return res[i].Priority < res[j].Priority }); return res
}

func (s *WidgetService) GetRegistry(ctx context.Context) map[string][]*domain.Widget {
	s.mu.RLock(); defer s.mu.RUnlock(); reg := make(map[string][]*domain.Widget)
	for _, w := range s.ws { c := *w; reg[w.Slot] = append(reg[w.Slot], &c) }
	for sl := range reg { sort.Slice(reg[sl], func(i, j int) bool { return reg[sl][i].Priority < reg[sl][j].Priority }) }
	return reg
}
