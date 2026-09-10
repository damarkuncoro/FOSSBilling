package security

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	scriptRegex  = regexp.MustCompile(`(?i)<script\b[^>]*>[\s\S]*?</script>`)
	styleRegex   = regexp.MustCompile(`(?i)<style\b[^>]*>[\s\S]*?</style>`)
	iframeRegex  = regexp.MustCompile(`(?i)<iframe\b[^>]*>[\s\S]*?</iframe>`)
	handlerRegex = regexp.MustCompile(`(?i)\s+on\w+\s*=\s*(?:'[^']*'|"[^"]*"|[^\s>]+)`)
	jsRegex      = regexp.MustCompile(`(?i)javascript:\s*`)
	slugRegex    = regexp.MustCompile(`[^a-z0-9\-]+`)
)

func SanitizeHTML(s string) string {
	if s == "" { return "" }
	s = strings.ReplaceAll(s, "\x00", "")
	s = scriptRegex.ReplaceAllString(s, "")
	s = styleRegex.ReplaceAllString(s, "")
	s = iframeRegex.ReplaceAllString(s, "")
	s = handlerRegex.ReplaceAllString(s, "")
	return strings.TrimSpace(jsRegex.ReplaceAllString(s, ""))
}

func SanitizeSlug(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = slugRegex.ReplaceAllString(s, "-")
	return strings.Trim(regexp.MustCompile(`-+`).ReplaceAllString(s, "-"), "-")
}

func SanitizeAlphaNumeric(s string) string {
	var b strings.Builder
	for _, r := range s { if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' || unicode.IsSpace(r) { b.WriteRune(r) } }
	return b.String()
}
