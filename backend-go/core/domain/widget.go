package domain

type Widget struct {
	ID        string                 `json:"id"`
	Module    string                 `json:"module"`
	Slot      string                 `json:"slot"` // e.g. "client.footer", "admin.dashboard.top", "client.cart.summary"
	Title     string                 `json:"title"`
	Component string                 `json:"component"` // UI Component identifier or template name
	Priority  int                    `json:"priority"`  // Lower number = higher priority
	Config    map[string]interface{} `json:"config,omitempty"`
}
