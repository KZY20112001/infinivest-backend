package global

import "errors"

var (
	ErrInternal  = errors.New("internal Error")
	ErrNotFound  = errors.New("resource not found")
	ErrNil       = errors.New("nil value")
	ErrDuplicate = errors.New("duplicate resource in use")
)
