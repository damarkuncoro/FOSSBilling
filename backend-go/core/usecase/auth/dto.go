package auth

import "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"

type RegisterDTO struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Company      string `json:"company,omitempty"`
	Address1     string `json:"address_1,omitempty"`
	City         string `json:"city,omitempty"`
	Country      string `json:"country,omitempty"`
	Phone        string `json:"phone,omitempty"`
	Currency     string `json:"currency,omitempty"`
	Honeypot     string `json:"website_hp,omitempty"`
	CaptchaToken string `json:"captcha_token,omitempty"`
}

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateProfileDTO struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Gender       string `json:"gender"`
	Birthday     string `json:"birthday"` // YYYY-MM-DD
	Company      string `json:"company"`
	CompanyVat   string `json:"company_vat"`
	CompanyNumber string `json:"company_number"`
	Type         string `json:"type"`
	Address1     string `json:"address_1"`
	Address2     string `json:"address_2"`
	City         string `json:"city"`
	State        string `json:"state"`
	Postcode     string `json:"postcode"`
	Country      string `json:"country"`
	PhoneCC      string `json:"phone_cc"`
	Phone        string `json:"phone"`
	Currency     string `json:"currency"`
	BillingEmail string `json:"billing_email"`
	Custom1      string `json:"custom_1"`
	Custom2      string `json:"custom_2"`
	Custom3      string `json:"custom_3"`
	Custom4      string `json:"custom_4"`
	Custom5      string `json:"custom_5"`
	Custom6      string `json:"custom_6"`
	Custom7      string `json:"custom_7"`
	Custom8      string `json:"custom_8"`
	Custom9      string `json:"custom_9"`
	Custom10     string `json:"custom_10"`
	Custom11     string `json:"custom_11"`
	Custom12     string `json:"custom_12"`
	Custom13     string `json:"custom_13"`
	Custom14     string `json:"custom_14"`
	Custom15     string `json:"custom_15"`
	Custom16     string `json:"custom_16"`
	Custom17     string `json:"custom_17"`
	Custom18     string `json:"custom_18"`
	Custom19     string `json:"custom_19"`
	Custom20     string `json:"custom_20"`
}

type AuthResponse struct {
	Token             string      `json:"token,omitempty"`
	TwoFactorRequired bool        `json:"two_factor_required,omitempty"`
	Client            ClientBrief `json:"client,omitempty"`
}

type TwoFactorSetupResponse struct {
	Secret string `json:"secret"`
	QRURL  string `json:"qr_url"`
}

type ClientBrief struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Currency  string `json:"currency"`
}

type ProfileResponse struct {
	Client  ClientDetail  `json:"client"`
	Balance decimal.Money `json:"balance"`
}

type ClientDetail struct {
	ID            int64   `json:"id"`
	AID           *string `json:"aid,omitempty"`
	Email         string  `json:"email"`
	FirstName     string  `json:"first_name"`
	LastName      string  `json:"last_name"`
	Gender        string  `json:"gender,omitempty"`
	Birthday      string  `json:"birthday,omitempty"` // YYYY-MM-DD
	Company       string  `json:"company"`
	CompanyVat    string  `json:"company_vat,omitempty"`
	CompanyNumber string  `json:"company_number,omitempty"`
	Type          string  `json:"type"`
	Address1      string  `json:"address_1"`
	Address2      string  `json:"address_2"`
	City          string  `json:"city"`
	State         string  `json:"state"`
	Postcode      string  `json:"postcode"`
	Country       string  `json:"country"`
	PhoneCC       string  `json:"phone_cc"`
	Phone         string  `json:"phone"`
	Currency      string  `json:"currency"`
	BillingEmail  string  `json:"billing_email,omitempty"`
	Status        string  `json:"status"`
	Custom1       string  `json:"custom_1,omitempty"`
	Custom2       string  `json:"custom_2,omitempty"`
	Custom3       string  `json:"custom_3,omitempty"`
	Custom4       string  `json:"custom_4,omitempty"`
	Custom5       string  `json:"custom_5,omitempty"`
	Custom6       string  `json:"custom_6,omitempty"`
	Custom7       string  `json:"custom_7,omitempty"`
	Custom8       string  `json:"custom_8,omitempty"`
	Custom9       string  `json:"custom_9,omitempty"`
	Custom10      string  `json:"custom_10,omitempty"`
	Custom11      string  `json:"custom_11,omitempty"`
	Custom12      string  `json:"custom_12,omitempty"`
	Custom13      string  `json:"custom_13,omitempty"`
	Custom14      string  `json:"custom_14,omitempty"`
	Custom15      string  `json:"custom_15,omitempty"`
	Custom16      string  `json:"custom_16,omitempty"`
	Custom17      string  `json:"custom_17,omitempty"`
	Custom18      string  `json:"custom_18,omitempty"`
	Custom19      string  `json:"custom_19,omitempty"`
	Custom20      string  `json:"custom_20,omitempty"`
}
