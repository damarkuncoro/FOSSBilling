package widget

import (
	"context"
	"sort"
	"sync"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type WidgetService struct {
	mu      sync.RWMutex
	widgets map[string]*domain.Widget
}

func NewWidgetService() *WidgetService {
	s := &WidgetService{
		widgets: make(map[string]*domain.Widget),
	}

	// Register default core widgets
	s.RegisterWidget(&domain.Widget{
		ID:        "branding_footer",
		Module:    "branding",
		Slot:      "client.footer",
		Title:     "Powered by FOSSBilling",
		Component: "BrandingFooter",
		Priority:  100,
	})
	s.RegisterWidget(&domain.Widget{
		ID:        "cookie_consent_banner",
		Module:    "cookieconsent",
		Slot:      "client.banner",
		Title:     "GDPR Cookie Banner",
		Component: "CookieConsentBanner",
		Priority:  10,
	})
	s.RegisterWidget(&domain.Widget{
		ID:        "quick_stats_overview",
		Module:    "stats",
		Slot:      "admin.dashboard.top",
		Title:     "Real-Time Metrics",
		Component: "DashboardQuickStats",
		Priority:  1,
	})

	return s
}

func (s *WidgetService) RegisterWidget(w *domain.Widget) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if w.Priority <= 0 {
		w.Priority = 10
	}
	s.widgets[w.ID] = w
}

func (s *WidgetService) UnregisterWidget(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.widgets, id)
}

func (s *WidgetService) GetWidgetsForSlot(ctx context.Context, slot string) []*domain.Widget {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*domain.Widget
	for _, w := range s.widgets {
		if w.Slot == slot {
			cpy := *w
			result = append(result, &cpy)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Priority < result[j].Priority
	})

	return result
}

func (s *WidgetService) GetRegistry(ctx context.Context) map[string][]*domain.Widget {
	s.mu.RLock()
	defer s.mu.RUnlock()

	registry := make(map[string][]*domain.Widget)
	for _, w := range s.widgets {
		cpy := *w
		registry[w.Slot] = append(registry[w.Slot], &cpy)
	}

	for slot := range registry {
		sort.Slice(registry[slot], func(i, j int) bool {
			return registry[slot][i].Priority < registry[slot][j].Priority
		})
	}

	return registry
}
