package security

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	scriptTagRegex    = regexp.MustCompile(`(?i)<script\b[^>]*>[\s\S]*?</script>`)
	styleTagRegex     = regexp.MustCompile(`(?i)<style\b[^>]*>[\s\S]*?</style>`)
	iframeTagRegex    = regexp.MustCompile(`(?i)<iframe\b[^>]*>[\s\S]*?</iframe>`)
	eventHandlerRegex = regexp.MustCompile(`(?i)\s+on\w+\s*=\s*(?:'[^']*'|"[^"]*"|[^\s>]+)`)
	jsProtocolRegex   = regexp.MustCompile(`(?i)javascript:\s*`)
	dataProtocolRegex = regexp.MustCompile(`(?i)data:\s*text/html`)
	slugNonAlphaRegex = regexp.MustCompile(`[^a-z0-9\-]+`)
	slugHyphensRegex  = regexp.MustCompile(`-+`)
)

// SanitizeHTML strips dangerous script tags, iframes, inline event handlers, and javascript pseudo-protocols
func SanitizeHTML(input string) string {
	if input == "" {
		return ""
	}

	// Remove null bytes
	res := strings.ReplaceAll(input, "\x00", "")

	// Remove script, style, and iframe blocks
	res = scriptTagRegex.ReplaceAllString(res, "")
	res = styleTagRegex.ReplaceAllString(res, "")
	res = iframeTagRegex.ReplaceAllString(res, "")

	// Remove inline event handlers (onclick, onerror, onload, etc.)
	res = eventHandlerRegex.ReplaceAllString(res, "")

	// Remove javascript: and malicious data: URLs
	res = jsProtocolRegex.ReplaceAllString(res, "")
	res = dataProtocolRegex.ReplaceAllString(res, "")

	return strings.TrimSpace(res)
}

// SanitizeSlug converts strings to URL-friendly lowercase slugs (e.g. "My Cloud VPS #1" -> "my-cloud-vps-1")
func SanitizeSlug(input string) string {
	slug := strings.ToLower(strings.TrimSpace(input))
	slug = slugNonAlphaRegex.ReplaceAllString(slug, "-")
	slug = slugHyphensRegex.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}

// SanitizeAlphaNumeric keeps only alphanumeric characters, underscores, and dashes
func SanitizeAlphaNumeric(input string) string {
	var sb strings.Builder
	for _, r := range input {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}
