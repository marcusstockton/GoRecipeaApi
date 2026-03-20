package recipe

import "errors"

var (
	ErrNotFound      = errors.New("recipe not found")
	ErrInvalidInput  = errors.New("invalid recipe input")
	ErrNotAuthorized = errors.New("not authorized to modify this recipe")
)
