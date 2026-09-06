package validator

import (
	"errors"
	"regexp"
	"strings"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// ValidateClientRegistration checks registration input for domain compliance
func ValidateClientRegistration(client *domain.Client, rawPassword string) error {
	if client == nil {
		return errors.New("client cannot be nil")
	}

	email := strings.TrimSpace(client.Email)
	if email == "" {
		return errors.New("email address is required")
	}
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email address format")
	}

	if strings.TrimSpace(client.FirstName) == "" {
		return errors.New("first name is required")
	}

	if len(rawPassword) < 8 {
		return errors.New("password must be at least 8 characters in length")
	}

	return nil
}
