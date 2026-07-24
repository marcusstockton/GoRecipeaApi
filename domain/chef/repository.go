package chef

type Repository interface {
	Save(chef *Chef) error
	FindByID(id uint) (*Chef, error)
	FindByEmail(email string) (*Chef, error)
	Delete(id uint) error
	List() ([]Chef, error)
}
