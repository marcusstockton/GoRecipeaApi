package recipe

type Repository interface {
	Save(*Recipe) error
	FindByID(uint) (*Recipe, error)
	List() ([]Recipe, error)
	Delete(uint) error
	AddLike(uint, uint) error
	RemoveLike(uint, uint) error
	AddComment(uint, uint, string, *uint) error
	RemoveComment(uint, uint) error
	FindCommentByID(uint) (*RecipeComment, error)
	ListComments(uint) ([]RecipeComment, error)
}
