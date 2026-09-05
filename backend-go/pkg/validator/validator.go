package validator

import (
	"net/mail"
	"strings"
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
		v.Add(field, field+" must be at least "+string(rune('0'+minLen))+" characters")
	}
}
