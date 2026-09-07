package security_test

import (
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/security"
	"github.com/stretchr/testify/assert"
)

func TestSanitizeHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Strip script tags",
			input:    "<div>Hello <script>alert('xss');</script>World</div>",
			expected: "<div>Hello World</div>",
		},
		{
			name:     "Strip style tags and iframe",
			input:    "<p>Text</p><iframe src='evil.com'></iframe><style>body{display:none;}</style>",
			expected: "<p>Text</p>",
		},
		{
			name:     "Strip inline event handlers",
			input:    `<img src="image.jpg" onerror="alert(1)" onload="evil()" alt="test" />`,
			expected: `<img src="image.jpg" alt="test" />`,
		},
		{
			name:     "Strip javascript pseudo protocols",
			input:    `<a href="javascript:alert(1)">Click Me</a>`,
			expected: `<a href="alert(1)">Click Me</a>`,
		},
		{
			name:     "Clean safe HTML",
			input:    `<h1>Welcome to FOSSBilling</h1><p>Enjoy our cloud hosting.</p>`,
			expected: `<h1>Welcome to FOSSBilling</h1><p>Enjoy our cloud hosting.</p>`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := security.SanitizeHTML(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestSanitizeSlug(t *testing.T) {
	assert.Equal(t, "my-cloud-vps-1", security.SanitizeSlug("My Cloud VPS #1"))
	assert.Equal(t, "fossbilling-enterprise-plan", security.SanitizeSlug("FOSSBilling - Enterprise Plan!"))
	assert.Equal(t, "domain-names", security.SanitizeSlug("Domain  ---  Names"))
}

func TestSanitizeAlphaNumeric(t *testing.T) {
	assert.Equal(t, "valid_user-123", security.SanitizeAlphaNumeric("valid_user-123!@#$%^&*()"))
}
