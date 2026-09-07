package geoip

import (
	"net"
	"strings"
)

type CountryInfo struct {
	ISOCode  string   `json:"iso_code"`
	Name     string   `json:"name"`
	Flag     string   `json:"flag"`
	Currency string   `json:"currency"`
	Locales  []string `json:"locales"`
}

var countryDatabase = map[string]CountryInfo{
	"US": {ISOCode: "US", Name: "United States", Currency: "USD", Locales: []string{"en_US"}},
	"ID": {ISOCode: "ID", Name: "Indonesia", Currency: "IDR", Locales: []string{"id_ID", "en_US"}},
	"GB": {ISOCode: "GB", Name: "United Kingdom", Currency: "GBP", Locales: []string{"en_GB"}},
	"DE": {ISOCode: "DE", Name: "Germany", Currency: "EUR", Locales: []string{"de_DE", "en_US"}},
	"FR": {ISOCode: "FR", Name: "France", Currency: "EUR", Locales: []string{"fr_FR", "en_US"}},
	"SG": {ISOCode: "SG", Name: "Singapore", Currency: "SGD", Locales: []string{"en_SG", "zh_CN"}},
	"JP": {ISOCode: "JP", Name: "Japan", Currency: "JPY", Locales: []string{"ja_JP"}},
	"AU": {ISOCode: "AU", Name: "Australia", Currency: "AUD", Locales: []string{"en_AU"}},
	"CA": {ISOCode: "CA", Name: "Canada", Currency: "CAD", Locales: []string{"en_CA", "fr_CA"}},
	"NL": {ISOCode: "NL", Name: "Netherlands", Currency: "EUR", Locales: []string{"nl_NL", "en_US"}},
	"IN": {ISOCode: "IN", Name: "India", Currency: "INR", Locales: []string{"en_IN", "hi_IN"}},
	"BR": {ISOCode: "BR", Name: "Brazil", Currency: "BRL", Locales: []string{"pt_BR"}},
}

// CountryFlagEmoji returns Unicode regional indicator flag emoji for 2-letter ISO code
func CountryFlagEmoji(isoCode string) string {
	code := strings.ToUpper(strings.TrimSpace(isoCode))
	if len(code) != 2 {
		return "🌐"
	}
	// Unicode regional indicator symbol letter A is 0x1F1E6
	const regionalIndicatorBase = 0x1F1E6 - 'A'
	r1 := rune(code[0]) + regionalIndicatorBase
	r2 := rune(code[1]) + regionalIndicatorBase
	return string([]rune{r1, r2})
}

// LookupCountry returns enriched country info for a 2-letter ISO country code
func LookupCountry(isoCode string) CountryInfo {
	code := strings.ToUpper(strings.TrimSpace(isoCode))
	if info, ok := countryDatabase[code]; ok {
		info.Flag = CountryFlagEmoji(info.ISOCode)
		return info
	}

	return CountryInfo{
		ISOCode:  code,
		Name:     code,
		Flag:     CountryFlagEmoji(code),
		Currency: "USD",
		Locales:  []string{"en_US"},
	}
}

// IsPrivateIP checks if an IP is a loopback, private RFC1918, or local link address
func IsPrivateIP(ipStr string) bool {
	ip := net.ParseIP(strings.TrimSpace(ipStr))
	if ip == nil {
		return false
	}
	return ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsUnspecified()
}
