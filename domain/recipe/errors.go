package recipe

import "errors"

var (
	ErrNotFound           = errors.New("recipe not found")
	ErrInvalidInput       = errors.New("invalid recipe input")
	ErrNotAuthorized      = errors.New("not authorized to modify this recipe")
	ErrDuplicateRecipe    = errors.New("recipe already exists for this chef")
	ErrDuplicateLike      = errors.New("chef already liked this recipe")
	ErrEmptyComment       = errors.New("comment content cannot be empty")
	ErrSelfLikeNotAllowed = errors.New("chefs cannot like their own recipes")
)
