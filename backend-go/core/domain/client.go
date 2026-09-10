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
	ID               int64        `json:"id"`
	AID              *string      `json:"aid,omitempty"`
	GroupID          *int64       `json:"group_id,omitempty"`
	Email            string       `json:"email"`
	PasswordHash     string       `json:"-"`
	Type             string       `json:"type"`
	FirstName        string       `json:"first_name"`
	LastName         string       `json:"last_name"`
	Gender           string       `json:"gender,omitempty"`
	Birthday         *time.Time   `json:"birthday,omitempty"`
	Company          string       `json:"company,omitempty"`
	CompanyVat       string       `json:"company_vat,omitempty"`
	CompanyNumber    string       `json:"company_number,omitempty"`
	Address1         string       `json:"address_1"`
	Address2         string       `json:"address_2,omitempty"`
	City             string       `json:"city"`
	State            string       `json:"state"`
	Postcode         string       `json:"postcode"`
	Country          string       `json:"country"` // ISO 2-letter
	PhoneCC          string       `json:"phone_cc"`
	Phone            string       `json:"phone"`
	Notes            string       `json:"notes,omitempty"`
	Currency         string       `json:"currency"`
	BillingEmail     string       `json:"billing_email,omitempty"`
	ReferredBy       string       `json:"referred_by,omitempty"`
	TaxExempt        bool         `json:"tax_exempt"`
	Status           ClientStatus `json:"status"`
	TwoFactorEnabled bool         `json:"two_factor_enabled"`
	TwoFactorSecret  *string      `json:"-"`
	Custom1          string       `json:"custom_1,omitempty"`
	Custom2          string       `json:"custom_2,omitempty"`
	Custom3          string       `json:"custom_3,omitempty"`
	Custom4          string       `json:"custom_4,omitempty"`
	Custom5          string       `json:"custom_5,omitempty"`
	Custom6          string       `json:"custom_6,omitempty"`
	Custom7          string       `json:"custom_7,omitempty"`
	Custom8          string       `json:"custom_8,omitempty"`
	Custom9          string       `json:"custom_9,omitempty"`
	Custom10         string       `json:"custom_10,omitempty"`
	Custom11         string       `json:"custom_11,omitempty"`
	Custom12         string       `json:"custom_12,omitempty"`
	Custom13         string       `json:"custom_13,omitempty"`
	Custom14         string       `json:"custom_14,omitempty"`
	Custom15         string       `json:"custom_15,omitempty"`
	Custom16         string       `json:"custom_16,omitempty"`
	Custom17         string       `json:"custom_17,omitempty"`
	Custom18         string       `json:"custom_18,omitempty"`
	Custom19         string       `json:"custom_19,omitempty"`
	Custom20         string       `json:"custom_20,omitempty"`
	CreatedAt        time.Time    `json:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at"`
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
