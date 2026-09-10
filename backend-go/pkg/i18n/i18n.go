package i18n

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/geoip"
)

type Locale struct { Code, Name, Native, Flag, Dir string }
var SupportedLocales = []Locale{
	{Code: "en_US", Name: "English", Native: "English", Flag: "🇺🇸", Dir: "ltr"},
	{Code: "id_ID", Name: "Indonesian", Native: "Bahasa Indonesia", Flag: "🇮🇩", Dir: "ltr"},
	{Code: "de_DE", Name: "German", Native: "Deutsch", Flag: "🇩🇪", Dir: "ltr"},
	{Code: "fr_FR", Name: "French", Native: "Français", Flag: "🇫🇷", Dir: "ltr"},
}

type I18n struct { mu sync.RWMutex; cats map[string]map[string]string; def string }
var defMgr = &I18n{cats: map[string]map[string]string{
	"en_US": {"welcome": "Welcome to :app", "invalid_credentials": "Invalid email or password", "order_created": "Order #:id created successfully", "items_count_one": ":count item in cart", "items_count_other": ":count items in cart"},
	"id_ID": {"welcome": "Selamat datang di :app", "invalid_credentials": "Email atau sandi salah", "order_created": "Pesanan #:id berhasil dibuat", "items_count_one": ":count item di keranjang", "items_count_other": ":count item di keranjang"},
	"de_DE": {"welcome": "Willkommen bei :app"},
	"fr_FR": {"welcome": "Bienvenue sur :app"},
}, def: "en_US"}

func (m *I18n) Translate(l, k string, p ...map[string]string) string {
	m.mu.RLock(); defer m.mu.RUnlock(); msg := k
	if c, ok := m.cats[l]; ok { if v, ex := c[k]; ex { msg = v } } else if c, ok = m.cats[m.def]; ok { if v, ex := c[k]; ex { msg = v } }
	if len(p) > 0 && p[0] != nil { for k, v := range p[0] { msg = strings.ReplaceAll(msg, ":"+k, v) } }
	return msg
}

func T(l, k string, p ...map[string]string) string { return defMgr.Translate(l, k, p...) }

func TPlural(l, s, pl string, c int, p ...map[string]string) string {
	k := pl; if c == 1 { k = s }; pm := map[string]string{}; if len(p) > 0 { pm = p[0] }
	pm["count"] = fmt.Sprintf("%d", c); return T(l, k, pm)
}

func ParseAcceptLanguage(h string) string {
	if h == "" { return "en_US" }
	for _, t := range strings.Split(h, ",") {
		lang := strings.ReplaceAll(strings.Split(strings.TrimSpace(t), ";")[0], "-", "_")
		for _, s := range SupportedLocales {
			if strings.EqualFold(s.Code, lang) || strings.HasPrefix(strings.ToLower(lang), strings.ToLower(s.Code[:2])) { return s.Code }
		}
	}
	return "en_US"
}

func LocaleFromContext(ctx context.Context) string { if l, ok := ctx.Value("app_locale").(string); ok && l != "" { return l }; return "en_US" }

func LocaleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := r.URL.Query().Get("locale")
		if l == "" { if c, err := r.Cookie("BBLANG"); err == nil { l = c.Value } }
		if l == "" {
			h := r.Header.Get("Accept-Language"); if h != "" {
				for _, t := range strings.Split(h, ",") {
					lang := strings.ReplaceAll(strings.Split(strings.TrimSpace(t), ";")[0], "-", "_")
					for _, s := range SupportedLocales { if strings.EqualFold(s.Code, lang) { l = s.Code; break } }
					if l != "" { break }
				}
			}
		}
		if l == "" { l = "en_US" }
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "app_locale", l)))
	})
}

func GetLocaleFlag(l string) string {
	p := strings.Split(l, "_"); if len(p) == 2 { return geoip.CountryFlagEmoji(p[1]) }; return "🌐"
}
