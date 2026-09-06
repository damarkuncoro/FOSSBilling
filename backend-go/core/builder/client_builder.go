package builder

import (
	"errors"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/auth"
)

type ClientBuilder struct {
	email     string
	password  string
	firstName string
	lastName  string
	company   string
	country   string
	currency  string
	status    domain.ClientStatus
}

func NewClientBuilder() *ClientBuilder {
	return &ClientBuilder{
		country:  "ID",
		currency: "USD",
		status:   domain.ClientStatusActive,
	}
}

func (b *ClientBuilder) WithEmail(email string) *ClientBuilder {
	b.email = email
	return b
}

func (b *ClientBuilder) WithPassword(rawPassword string) *ClientBuilder {
	b.password = rawPassword
	return b
}

func (b *ClientBuilder) WithName(firstName, lastName string) *ClientBuilder {
	b.firstName = firstName
	b.lastName = lastName
	return b
}

func (b *ClientBuilder) WithCompany(company string) *ClientBuilder {
	b.company = company
	return b
}

func (b *ClientBuilder) WithCountryAndCurrency(country, currency string) *ClientBuilder {
	b.country = country
	b.currency = currency
	return b
}

func (b *ClientBuilder) Build() (*domain.Client, error) {
	if b.email == "" {
		return nil, errors.New("client requires a valid email address")
	}
	if b.password == "" {
		b.password = "DefaultPassword123!"
	}

	passHash, err := auth.HashPassword(b.password)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	return &domain.Client{
		Email:        b.email,
		PasswordHash: passHash,
		FirstName:    b.firstName,
		LastName:     b.lastName,
		Company:      b.company,
		Country:      b.country,
		Currency:     b.currency,
		Status:       b.status,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}
