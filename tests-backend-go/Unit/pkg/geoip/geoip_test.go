package geoip_test

import (
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/geoip"
	"github.com/stretchr/testify/assert"
)

func TestLookupCountry(t *testing.T) {
	us := geoip.LookupCountry("US")
	assert.Equal(t, "US", us.ISOCode)
	assert.Equal(t, "United States", us.Name)
	assert.Equal(t, "USD", us.Currency)
	assert.NotEmpty(t, us.Flag)

	id := geoip.LookupCountry("ID")
	assert.Equal(t, "ID", id.ISOCode)
	assert.Equal(t, "Indonesia", id.Name)
	assert.Equal(t, "IDR", id.Currency)
}

func TestCountryFlagEmoji(t *testing.T) {
	assert.Equal(t, "🇺🇸", geoip.CountryFlagEmoji("US"))
	assert.Equal(t, "🇮🇩", geoip.CountryFlagEmoji("ID"))
	assert.Equal(t, "🇩🇪", geoip.CountryFlagEmoji("DE"))
	assert.Equal(t, "🌐", geoip.CountryFlagEmoji("INVALID"))
}

func TestIsPrivateIP(t *testing.T) {
	assert.True(t, geoip.IsPrivateIP("127.0.0.1"))
	assert.True(t, geoip.IsPrivateIP("10.0.0.1"))
	assert.True(t, geoip.IsPrivateIP("192.168.1.1"))
	assert.True(t, geoip.IsPrivateIP("::1"))

	assert.False(t, geoip.IsPrivateIP("8.8.8.8"))
	assert.False(t, geoip.IsPrivateIP("1.1.1.1"))
	assert.False(t, geoip.IsPrivateIP("not-an-ip"))
}
