package errors

import "errors"

var (
	ErrNotFound          = errors.New("resource not found")
	ErrUnauthorized      = errors.New("unauthorized access")
	ErrForbidden         = errors.New("access forbidden")
	ErrInvalidInput      = errors.New("invalid input payload")
	ErrDuplicateEntry    = errors.New("record already exists")
	ErrInsufficientFunds = errors.New("insufficient balance")
	ErrOrderExpired      = errors.New("order has expired or cannot be renewed")
	ErrProvisioning      = errors.New("service provisioning failed")
	ErrInternal          = errors.New("internal server error")
)
