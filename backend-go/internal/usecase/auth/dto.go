package auth

import "github.com/fossbilling/backend-go/pkg/decimal"

type RegisterDTO struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Company   string `json:"company,omitempty"`
	Address1  string `json:"address_1,omitempty"`
	City      string `json:"city,omitempty"`
	Country   string `json:"country,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Currency  string `json:"currency,omitempty"`
}

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateProfileDTO struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Company   string `json:"company"`
	Address1  string `json:"address_1"`
	Address2  string `json:"address_2"`
	City      string `json:"city"`
	State     string `json:"state"`
	Postcode  string `json:"postcode"`
	Country   string `json:"country"`
	PhoneCC   string `json:"phone_cc"`
	Phone     string `json:"phone"`
	Currency  string `json:"currency"`
}

type AuthResponse struct {
	Token  string      `json:"token"`
	Client ClientBrief `json:"client"`
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
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Company   string `json:"company"`
	Address1  string `json:"address_1"`
	Address2  string `json:"address_2"`
	City      string `json:"city"`
	State     string `json:"state"`
	Postcode  string `json:"postcode"`
	Country   string `json:"country"`
	PhoneCC   string `json:"phone_cc"`
	Phone     string `json:"phone"`
	Currency  string `json:"currency"`
	Status    string `json:"status"`
}
