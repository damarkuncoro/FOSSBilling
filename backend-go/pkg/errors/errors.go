package errors

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("forbidden")
	ErrInvalidInput  = errors.New("invalid input")
	ErrDuplicate     = errors.New("exists")
	ErrNoFunds       = errors.New("insufficient funds")
	ErrExpired       = errors.New("expired")
	ErrInternal      = errors.New("internal error")
)
