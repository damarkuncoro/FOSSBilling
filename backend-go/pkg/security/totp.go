package security

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"math"
	"net/url"
	"strings"
	"time"
)

// GenerateTOTPSecret creates a random base32 encoded secret for TOTP (160-bit)
func GenerateTOTPSecret() string {
	// For simplicity in this implementation, we use a robust random string generator
	// and encode it to base32.
	secret := GenerateRandomString(20)
	return base32.StdEncoding.EncodeToString([]byte(secret))
}

// GetTOTPCode generates the current 6-digit TOTP code for a given secret
func GetTOTPCode(secret string) (string, error) {
	key, err := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	if err != nil {
		return "", err
	}

	// Calculate counter (30-second steps)
	counter := uint64(time.Now().Unix() / 30)
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, counter)

	// HMAC-SHA1
	mac := hmac.New(sha1.New, key)
	mac.Write(buf)
	sum := mac.Sum(nil)

	// Dynamic truncation
	offset := sum[len(sum)-1] & 0xf
	binCode := int64(sum[offset]&0x7f)<<24 |
		int64(sum[offset+1]&0xff)<<16 |
		int64(sum[offset+2]&0xff)<<8 |
		int64(sum[offset+3]&0xff)

	code := binCode % int64(math.Pow10(6))
	return fmt.Sprintf("%06d", code), nil
}

// VerifyTOTP checks if the provided code is valid for the secret
func VerifyTOTP(secret, code string) bool {
	expected, err := GetTOTPCode(secret)
	if err != nil {
		return false
	}
	return expected == strings.TrimSpace(code)
}

// GenerateTOTPURL creates a standard otpauth:// URL for QR code generation
func GenerateTOTPURL(account, issuer, secret string) string {
	u := url.URL{
		Scheme: "otpauth",
		Host:   "totp",
		Path:   fmt.Sprintf("/%s:%s", issuer, account),
	}
	q := u.Query()
	q.Set("secret", secret)
	q.Set("issuer", issuer)
	q.Set("digits", "6")
	q.Set("period", "30")
	u.RawQuery = q.Encode()

	return u.String()
}
