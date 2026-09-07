package domain

// CookieConsentConfig represents the cookie compliance banner settings.
type CookieConsentConfig struct {
	Enabled          bool   `json:"enabled"`
	Message          string `json:"message"`
	ButtonText       string `json:"button_text"`
	Position         string `json:"position"` // "bottom", "top", "bottom-left", "bottom-right"
	PrivacyPolicyURL string `json:"privacy_policy_url"`
	Theme            string `json:"theme"` // "classic", "dark", "light"
}
