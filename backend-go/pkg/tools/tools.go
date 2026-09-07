package tools

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

const (
	lowerChars   = "abcdefghijklmnopqrstuvwxyz"
	upperChars   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitChars   = "0123456789"
	specialChars = "!@#$%^&*()-_=+"
)

// GeneratePassword generates a cryptographically secure random password
func GeneratePassword(length int, includeSpecial bool) (string, error) {
	if length < 6 {
		length = 8
	}

	charSet := lowerChars + upperChars + digitChars
	if includeSpecial {
		charSet += specialChars
	}

	// Guarantee at least one from required categories
	req := []string{lowerChars, upperChars, digitChars}
	if includeSpecial {
		req = append(req, specialChars)
	}

	res := make([]byte, length)
	for i, cat := range req {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(cat))))
		if err != nil {
			return "", err
		}
		res[i] = cat[n.Int64()]
	}

	// Fill remaining characters
	for i := len(req); i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charSet))))
		if err != nil {
			return "", err
		}
		res[i] = charSet[n.Int64()]
	}

	// Fisher-Yates shuffle
	for i := len(res) - 1; i > 0; i-- {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return "", err
		}
		j := n.Int64()
		res[i], res[j] = res[j], res[i]
	}

	return string(res), nil
}

// FormatBytes formats byte counts into human-readable strings (e.g. 1024 -> 1.00 KB, 1073741824 -> 1.00 GB)
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	suffixes := []string{"KB", "MB", "GB", "TB", "PB"}
	return fmt.Sprintf("%.2f %s", float64(bytes)/float64(div), suffixes[exp])
}

// NormalizeEmail normalizes email addresses by trimming whitespace and lowercasing
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
