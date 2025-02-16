package constants

import "errors"

var (
	ErrInternal           = errors.New("internal Error")
	ErrNil                = errors.New("nil value")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
