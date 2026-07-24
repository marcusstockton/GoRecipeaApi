package chef

import "errors"

var (
	ErrNotFound           = errors.New("chef not found")
	ErrInvalidInput       = errors.New("invalid chef input")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
