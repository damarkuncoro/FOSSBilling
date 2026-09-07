package i18n

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/geoip"
)

type contextKey string

const (
	LocaleContextKey contextKey = "app_locale"
	DefaultLocale               = "en_US"
)

// LocaleInfo describes a supported system locale
type LocaleInfo struct {
	Code       string `json:"code"`
	Name       string `json:"name"`
	NativeName string `json:"native_name"`
	Flag       string `json:"flag"`
	Direction  string `json:"direction"` // "ltr" or "rtl"
}

var SupportedLocales = []LocaleInfo{
	{Code: "en_US", Name: "English (United States)", NativeName: "English", Flag: "🇺🇸", Direction: "ltr"},
	{Code: "id_ID", Name: "Indonesian", NativeName: "Bahasa Indonesia", Flag: "🇮🇩", Direction: "ltr"},
	{Code: "de_DE", Name: "German", NativeName: "Deutsch", Flag: "🇩🇪", Direction: "ltr"},
	{Code: "fr_FR", Name: "French", NativeName: "Français", Flag: "🇫🇷", Direction: "ltr"},
	{Code: "es_ES", Name: "Spanish", NativeName: "Español", Flag: "🇪🇸", Direction: "ltr"},
	{Code: "ja_JP", Name: "Japanese", NativeName: "日本語", Flag: "🇯🇵", Direction: "ltr"},
	{Code: "ar_SA", Name: "Arabic", NativeName: "العربية", Flag: "🇸🇦", Direction: "rtl"},
}

// I18nManager manages multi-language catalogs and message translation
type I18nManager struct {
	mu            sync.RWMutex
	catalogs      map[string]map[string]string
	defaultLocale string
}

var defaultManager = NewI18nManager()

// NewI18nManager initializes a new translation manager with default dictionaries
func NewI18nManager() *I18nManager {
	mgr := &I18nManager{
		catalogs:      make(map[string]map[string]string),
		defaultLocale: DefaultLocale,
	}
	mgr.loadDefaultCatalogs()
	return mgr
}

func (m *I18nManager) loadDefaultCatalogs() {
	// English (en_US)
	m.catalogs["en_US"] = map[string]string{
		"welcome":             "Welcome to :app",
		"login_success":       "Successfully authenticated",
		"invalid_credentials": "Invalid email or password",
		"order_created":       "Order #:id created successfully",
		"invoice_paid":        "Invoice #:nr has been marked as paid",
		"ticket_closed":       "Support ticket #:id has been closed",
		"items_count_one":     ":count item in cart",
		"items_count_other":   ":count items in cart",
	}

	// Indonesian (id_ID)
	m.catalogs["id_ID"] = map[string]string{
		"welcome":             "Selamat datang di :app",
		"login_success":       "Berhasil masuk",
		"invalid_credentials": "Email atau kata sandi tidak valid",
		"order_created":       "Pesanan #:id berhasil dibuat",
		"invoice_paid":        "Faktur #:nr telah ditandai lunas",
		"ticket_closed":       "Tiket bantuan #:id telah ditutup",
		"items_count_one":     ":count item di keranjang",
		"items_count_other":   ":count item di keranjang",
	}

	// German (de_DE)
	m.catalogs["de_DE"] = map[string]string{
		"welcome":             "Willkommen bei :app",
		"login_success":       "Erfolgreich authentifiziert",
		"invalid_credentials": "Ungültige E-Mail oder Passwort",
		"order_created":       "Bestellung #:id erfolgreich erstellt",
		"invoice_paid":        "Rechnung #:nr wurde als bezahlt markiert",
		"ticket_closed":       "Support-Ticket #:id wurde geschlossen",
		"items_count_one":     ":count Artikel im Warenkorb",
		"items_count_other":   ":count Artikel im Warenkorb",
	}

	// French (fr_FR)
	m.catalogs["fr_FR"] = map[string]string{
		"welcome":             "Bienvenue sur :app",
		"login_success":       "Authentification réussie",
		"invalid_credentials": "Email ou mot de passe invalide",
		"order_created":       "Commande #:id créée avec succès",
		"invoice_paid":        "La facture #:nr a été marquée comme payée",
		"ticket_closed":       "Le ticket de support #:id a été fermé",
		"items_count_one":     ":count article dans le panier",
		"items_count_other":   ":count articles dans le panier",
	}
}

// AddTranslation adds or overrides a translation key in a specific locale
func (m *I18nManager) AddTranslation(locale, key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.catalogs[locale]; !ok {
		m.catalogs[locale] = make(map[string]string)
	}
	m.catalogs[locale][key] = value
}

// Translate retrieves a translated string with parameter substitution
func (m *I18nManager) Translate(locale, key string, params ...map[string]string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	msg := key
	if cat, ok := m.catalogs[locale]; ok {
		if val, exists := cat[key]; exists {
			msg = val
		}
	} else if defCat, ok := m.catalogs[m.defaultLocale]; ok {
		if val, exists := defCat[key]; exists {
			msg = val
		}
	}

	if len(params) > 0 && params[0] != nil {
		for k, v := range params[0] {
			msg = strings.ReplaceAll(msg, ":"+k, v)
		}
	}

	return msg
}

// TranslatePlural returns plural-aware translated string
func (m *I18nManager) TranslatePlural(locale, singularKey, pluralKey string, count int, params ...map[string]string) string {
	targetKey := pluralKey
	if count == 1 {
		targetKey = singularKey
	}

	p := make(map[string]string)
	if len(params) > 0 && params[0] != nil {
		for k, v := range params[0] {
			p[k] = v
		}
	}
	p["count"] = fmt.Sprintf("%d", count)

	return m.Translate(locale, targetKey, p)
}

// ParseAcceptLanguage parses HTTP Accept-Language header and matches the best supported locale
func ParseAcceptLanguage(header string) string {
	if header == "" {
		return DefaultLocale
	}

	tags := strings.Split(header, ",")
	for _, tag := range tags {
		parts := strings.Split(strings.TrimSpace(tag), ";")
		lang := strings.ReplaceAll(parts[0], "-", "_")

		for _, supp := range SupportedLocales {
			if strings.EqualFold(supp.Code, lang) || strings.HasPrefix(strings.ToLower(lang), strings.ToLower(supp.Code[:2])) {
				return supp.Code
			}
		}
	}

	return DefaultLocale
}

// Global helper wrappers
func T(locale, key string, params ...map[string]string) string {
	return defaultManager.Translate(locale, key, params...)
}

func TPlural(locale, singularKey, pluralKey string, count int, params ...map[string]string) string {
	return defaultManager.TranslatePlural(locale, singularKey, pluralKey, count, params...)
}

// LocaleFromContext extracts locale from request context, defaulting to en_US
func LocaleFromContext(ctx context.Context) string {
	if loc, ok := ctx.Value(LocaleContextKey).(string); ok && loc != "" {
		return loc
	}
	return DefaultLocale
}

// LocaleMiddleware detects locale from cookie, query param, or Accept-Language header
func LocaleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		locale := r.URL.Query().Get("locale")
		if locale == "" {
			if c, err := r.Cookie("BBLANG"); err == nil && c.Value != "" {
				locale = c.Value
			}
		}
		if locale == "" {
			locale = ParseAcceptLanguage(r.Header.Get("Accept-Language"))
		}

		ctx := context.WithValue(r.Context(), LocaleContextKey, locale)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CountryFlag helper using geoip
func GetLocaleFlag(localeCode string) string {
	parts := strings.Split(localeCode, "_")
	if len(parts) == 2 {
		return geoip.CountryFlagEmoji(parts[1])
	}
	return "🌐"
}
