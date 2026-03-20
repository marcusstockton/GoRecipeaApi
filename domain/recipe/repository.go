package recipe

type Repository interface {
	Save(*Recipe) error
	FindByID(uint) (*Recipe, error)
	List() ([]Recipe, error)
	Delete(uint) error
}
