package models

import "errors"

var (
	ErrNotFound       = errors.New("record not found")
	ErrAlreadyExists  = errors.New("record already exists")
	ErrInternal       = errors.New("internal server error")
	ErrInvalidRequest = errors.New("invalid request parameters")

	ErrIncorrectPassword = errors.New("incorrect password")
	ErrInvalidToken      = errors.New("invalid token")
)
