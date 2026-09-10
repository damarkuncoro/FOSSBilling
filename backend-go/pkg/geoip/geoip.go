package geoip

import (
	"net"
	"strings"
)

type CountryInfo struct { ISOCode, Name, Flag, Currency string; Locales []string }

var db = map[string]CountryInfo{
	"US": {ISOCode: "US", Name: "United States", Currency: "USD", Locales: []string{"en_US"}},
	"ID": {ISOCode: "ID", Name: "Indonesia", Currency: "IDR", Locales: []string{"id_ID"}},
	"GB": {ISOCode: "GB", Name: "United Kingdom", Currency: "GBP", Locales: []string{"en_GB"}},
}

func CountryFlagEmoji(c string) string {
	c = strings.ToUpper(strings.TrimSpace(c)); if len(c) != 2 { return "🌐" }
	return string([]rune{rune(c[0]) + 127397, rune(c[1]) + 127397})
}

func LookupCountry(c string) CountryInfo {
	c = strings.ToUpper(strings.TrimSpace(c)); i, ok := db[c]
	if !ok { i = CountryInfo{ISOCode: c, Name: c, Currency: "USD", Locales: []string{"en_US"}} }
	i.Flag = CountryFlagEmoji(i.ISOCode); return i
}

func IsPrivateIP(s string) bool {
	ip := net.ParseIP(strings.TrimSpace(s)); if ip == nil { return false }
	return ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsUnspecified()
}
