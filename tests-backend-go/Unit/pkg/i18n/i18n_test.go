package i18n_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/i18n"
	"github.com/stretchr/testify/assert"
)

func TestI18n_Translations(t *testing.T) {
	// English
	assert.Equal(t, "Welcome to FOSSBilling", i18n.T("en_US", "welcome", map[string]string{"app": "FOSSBilling"}))
	assert.Equal(t, "Order #101 created successfully", i18n.T("en_US", "order_created", map[string]string{"id": "101"}))

	// Indonesian
	assert.Equal(t, "Selamat datang di FOSSBilling", i18n.T("id_ID", "welcome", map[string]string{"app": "FOSSBilling"}))
	assert.Equal(t, "Pesanan #101 berhasil dibuat", i18n.T("id_ID", "order_created", map[string]string{"id": "101"}))

	// German
	assert.Equal(t, "Willkommen bei FOSSBilling", i18n.T("de_DE", "welcome", map[string]string{"app": "FOSSBilling"}))

	// Fallback to default
	assert.Equal(t, "Welcome to FOSSBilling", i18n.T("unknown_LOCALE", "welcome", map[string]string{"app": "FOSSBilling"}))
}

func TestI18n_Pluralization(t *testing.T) {
	assert.Equal(t, "1 item in cart", i18n.TPlural("en_US", "items_count_one", "items_count_other", 1))
	assert.Equal(t, "5 items in cart", i18n.TPlural("en_US", "items_count_one", "items_count_other", 5))

	assert.Equal(t, "1 item di keranjang", i18n.TPlural("id_ID", "items_count_one", "items_count_other", 1))
	assert.Equal(t, "3 item di keranjang", i18n.TPlural("id_ID", "items_count_one", "items_count_other", 3))
}

func TestI18n_ParseAcceptLanguage(t *testing.T) {
	assert.Equal(t, "id_ID", i18n.ParseAcceptLanguage("id-ID,id;q=0.9,en-US;q=0.8,en;q=0.7"))
	assert.Equal(t, "de_DE", i18n.ParseAcceptLanguage("de,de-DE;q=0.9,en;q=0.8"))
	assert.Equal(t, "fr_FR", i18n.ParseAcceptLanguage("fr-FR,fr;q=0.9"))
	assert.Equal(t, "en_US", i18n.ParseAcceptLanguage(""))
}

func TestI18n_LocaleMiddleware(t *testing.T) {
	handler := i18n.LocaleMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loc := i18n.LocaleFromContext(r.Context())
		w.Write([]byte(loc))
	}))

	// Via query param
	req1 := httptest.NewRequest("GET", "/api/test?locale=id_ID", nil)
	rr1 := httptest.NewRecorder()
	handler.ServeHTTP(rr1, req1)
	assert.Equal(t, "id_ID", rr1.Body.String())

	// Via Accept-Language header
	req2 := httptest.NewRequest("GET", "/api/test", nil)
	req2.Header.Set("Accept-Language", "de-DE,de;q=0.9")
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)
	assert.Equal(t, "de_DE", rr2.Body.String())
}
