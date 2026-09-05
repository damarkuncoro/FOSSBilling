package domain

import (
	"context"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

type ClientStatus string

const (
	ClientStatusActive    ClientStatus = "active"
	ClientStatusSuspended ClientStatus = "suspended"
	ClientStatusCanceled  ClientStatus = "canceled"
)

type Client struct {
	ID           int64        `json:"id"`
	GroupID      *int64       `json:"group_id,omitempty"`
	Email        string       `json:"email"`
	PasswordHash string       `json:"-"`
	FirstName    string       `json:"first_name"`
	LastName     string       `json:"last_name"`
	Company      string       `json:"company,omitempty"`
	Address1     string       `json:"address_1"`
	Address2     string       `json:"address_2,omitempty"`
	City         string       `json:"city"`
	State        string       `json:"state"`
	Postcode     string       `json:"postcode"`
	Country      string       `json:"country"` // ISO 2-letter
	PhoneCC      string       `json:"phone_cc"`
	Phone        string       `json:"phone"`
	Currency     string       `json:"currency"`
	TaxExempt    bool         `json:"tax_exempt"`
	Status       ClientStatus `json:"status"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

type ClientBalanceType string

const (
	BalanceTypeCredit ClientBalanceType = "credit"
	BalanceTypeDebit  ClientBalanceType = "debit"
)

type ClientBalance struct {
	ID          int64             `json:"id"`
	ClientID    int64             `json:"client_id"`
	Type        ClientBalanceType `json:"type"`
	Amount      decimal.Money     `json:"amount"`
	Description string            `json:"description"`
	RelID       *int64            `json:"rel_id,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
}

type ClientRepository interface {
	GetByID(ctx context.Context, id int64) (*Client, error)
	GetByEmail(ctx context.Context, email string) (*Client, error)
	List(ctx context.Context, limit, offset int) ([]*Client, int, error)
	Create(ctx context.Context, client *Client) error
	Update(ctx context.Context, client *Client) error
	Delete(ctx context.Context, id int64) error
	GetBalance(ctx context.Context, clientID int64) (decimal.Money, error)
	AddBalanceTransaction(ctx context.Context, balance *ClientBalance) error
}

