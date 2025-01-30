package constants

import "errors"

var (
	ErrInternal           = errors.New("internal Error")
	ErrNotFound           = errors.New("resource not found")
	ErrNil                = errors.New("nil value")
	ErrDuplicate          = errors.New("duplicate resource in use")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
