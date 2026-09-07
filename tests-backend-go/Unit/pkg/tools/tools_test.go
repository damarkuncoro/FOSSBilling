package tools_test

import (
	"testing"
	"unicode"

	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/tools"
	"github.com/stretchr/testify/assert"
)

func TestGeneratePassword(t *testing.T) {
	pass, err := tools.GeneratePassword(12, true)
	assert.NoError(t, err)
	assert.Len(t, pass, 12)

	var hasUpper, hasLower, hasDigit bool
	for _, r := range pass {
		if unicode.IsUpper(r) {
			hasUpper = true
		}
		if unicode.IsLower(r) {
			hasLower = true
		}
		if unicode.IsDigit(r) {
			hasDigit = true
		}
	}
	assert.True(t, hasUpper, "password should contain uppercase")
	assert.True(t, hasLower, "password should contain lowercase")
	assert.True(t, hasDigit, "password should contain digit")
}

func TestFormatBytes(t *testing.T) {
	assert.Equal(t, "500 B", tools.FormatBytes(500))
	assert.Equal(t, "1.00 KB", tools.FormatBytes(1024))
	assert.Equal(t, "1.50 MB", tools.FormatBytes(1572864))
	assert.Equal(t, "1.00 GB", tools.FormatBytes(1073741824))
}

func TestNormalizeEmail(t *testing.T) {
	assert.Equal(t, "admin@fossbilling.org", tools.NormalizeEmail("  Admin@FOSSBilling.Org  "))
}
