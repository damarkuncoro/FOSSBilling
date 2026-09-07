package security

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
)

// RequestFingerprint generates a normalized client device & environment fingerprint
type RequestFingerprint struct {
	UserAgent string `json:"user_agent"`
	IP        string `json:"ip"`
	Language  string `json:"language"`
	Encoding  string `json:"encoding"`
	Hash      string `json:"hash"`
}

func GenerateFingerprint(r *http.Request) RequestFingerprint {
	ua := r.UserAgent()
	ip := extractClientIP(r)
	lang := r.Header.Get("Accept-Language")
	enc := r.Header.Get("Accept-Encoding")

	// Calculate deterministic SHA-256 fingerprint hash
	raw := fmt.Sprintf("%s|%s|%s|%s", ua, ip, lang, enc)
	h := sha256.Sum256([]byte(raw))
	hashHex := hex.EncodeToString(h[:])

	return RequestFingerprint{
		UserAgent: ua,
		IP:        ip,
		Language:  lang,
		Encoding:  enc,
		Hash:      hashHex,
	}
}

func extractClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}
	parts := strings.Split(r.RemoteAddr, ":")
	if len(parts) > 0 {
		return parts[0]
	}
	return r.RemoteAddr
}
