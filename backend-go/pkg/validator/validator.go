package validator

import (
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"unicode"
)

var (
	sldRegex = regexp.MustCompile(`^[a-z0-9]+([a-z0-9\-]*[a-z0-9]+)?$`)
	tldRegex = regexp.MustCompile(`^[a-z]{2,24}$`)
)

type ValidationErrors map[string]string

func New() ValidationErrors { return make(ValidationErrors) }
func (v ValidationErrors) Add(f, m string) { if _, ok := v[f]; !ok { v[f] = m } }
func (v ValidationErrors) IsValid() bool { return len(v) == 0 }

func (v ValidationErrors) CheckEmail(f, e string) {
	if strings.TrimSpace(e) == "" { v.Add(f, "Empty email"); return }
	if _, err := mail.ParseAddress(e); err != nil { v.Add(f, "Invalid email") }
}

func (v ValidationErrors) CheckRequired(f, val string) { if strings.TrimSpace(val) == "" { v.Add(f, f+" required") } }
func (v ValidationErrors) CheckMinLength(f, val string, min int) { if len(val) < min { v.Add(f, fmt.Sprintf("%s too short", f)) } }
func (v ValidationErrors) CheckSlug(f, val string) { if !regexp.MustCompile(`^[a-z0-9\-]+$`).MatchString(val) { v.Add(f, "Invalid slug") } }

func (v ValidationErrors) CheckURL(f, u string) {
	if u == "" { return }
	if p, err := url.ParseRequestURI(u); err != nil || (p.Scheme != "http" && p.Scheme != "https") { v.Add(f, "Invalid URL") }
}

func (v ValidationErrors) CheckIP(f, i string) { if i != "" && net.ParseIP(i) == nil { v.Add(f, "Invalid IP") } }

func (v ValidationErrors) CheckSLD(f, s string) {
	c := strings.ToLower(strings.Trim(s, ". ")); if c == "" || len(c) > 63 || !sldRegex.MatchString(c) { v.Add(f, "Invalid domain name") }
}

func (v ValidationErrors) CheckTLD(f, t string) {
	c := strings.ToLower(strings.Trim(t, ". ")); if c == "" || !tldRegex.MatchString(c) { v.Add(f, "Invalid TLD") }
}

func (v ValidationErrors) CheckPasswordStrength(f, p string, min int) {
	if len(p) < min { v.Add(f, fmt.Sprintf("Too short (min %d)", min)); return }
	var u, l, d bool
	for _, r := range p {
		switch {
		case unicode.IsUpper(r): u = true
		case unicode.IsLower(r): l = true
		case unicode.IsDigit(r): d = true
		}
	}
	if !u || !l || !d { v.Add(f, "Must have upper, lower and digit") }
}
