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
	sldRegex  = regexp.MustCompile(`^[a-z0-9]+([a-z0-9\-]*[a-z0-9]+)?$`)
	tldRegex  = regexp.MustCompile(`^[a-z]{2,24}$`)
	slugRegex = regexp.MustCompile(`^[a-z0-9\-]+$`)
)

type ValidationErrors map[string]string

func New() ValidationErrors {
	return make(ValidationErrors)
}

func (v ValidationErrors) Add(field, message string) {
	if _, exists := v[field]; !exists {
		v[field] = message
	}
}

func (v ValidationErrors) IsValid() bool {
	return len(v) == 0
}

func (v ValidationErrors) CheckEmail(field, email string) {
	if strings.TrimSpace(email) == "" {
		v.Add(field, "Email cannot be empty")
		return
	}
	if _, err := mail.ParseAddress(email); err != nil {
		v.Add(field, "Invalid email address format")
	}
}

func (v ValidationErrors) CheckRequired(field, value string) {
	if strings.TrimSpace(value) == "" {
		v.Add(field, field+" is required")
	}
}

func (v ValidationErrors) CheckMinLength(field, value string, minLen int) {
	if len(strings.TrimSpace(value)) < minLen {
		v.Add(field, fmt.Sprintf("%s must be at least %d characters", field, minLen))
	}
}

func (v ValidationErrors) CheckURL(field, rawURL string) {
	if strings.TrimSpace(rawURL) == "" {
		return
	}
	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		v.Add(field, field+" must be a valid HTTP or HTTPS URL")
	}
}

func (v ValidationErrors) CheckIP(field, ipStr string) {
	if strings.TrimSpace(ipStr) == "" {
		return
	}
	if net.ParseIP(strings.TrimSpace(ipStr)) == nil {
		v.Add(field, field+" must be a valid IPv4 or IPv6 address")
	}
}

func (v ValidationErrors) CheckSLD(field, sld string) {
	cleaned := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(sld), "."))
	if cleaned == "" || len(cleaned) > 63 || !sldRegex.MatchString(cleaned) || strings.HasPrefix(cleaned, "--") {
		v.Add(field, field+" must be a valid second-level domain name")
	}
}

func (v ValidationErrors) CheckTLD(field, tld string) {
	cleaned := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(tld), "."))
	if cleaned == "" || !tldRegex.MatchString(cleaned) {
		v.Add(field, field+" must be a valid top-level domain (e.g. com, org, net)")
	}
}

func (v ValidationErrors) CheckSlug(field, slug string) {
	if strings.TrimSpace(slug) == "" || !slugRegex.MatchString(slug) {
		v.Add(field, field+" must contain only lowercase letters, numbers, and hyphens")
	}
}

func (v ValidationErrors) CheckPasswordStrength(field, password string, minLength int) {
	if len(password) < minLength {
		v.Add(field, fmt.Sprintf("Password must be at least %d characters long", minLength))
		return
	}
	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}
	if !hasUpper || !hasLower || !hasDigit {
		v.Add(field, "Password must contain uppercase, lowercase, and numeric characters")
	}
	_ = hasSpecial
}
