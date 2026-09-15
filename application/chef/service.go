package chef

import (
	domain "recipea.com/m/domain/chef"
)

type Repository interface {
	Save(chef *domain.Chef) error
	FindByID(id uint) (*domain.Chef, error)
	FindByEmail(email string) (*domain.Chef, error)
	Delete(id uint) error
	List() ([]domain.Chef, error)
}

type TokenProvider interface {
	CreateToken(subject uint) (string, error)
	ParseToken(token string) (uint, error)
}

type ChefService struct {
	Repo  Repository
	Token TokenProvider
}

type CreateChefRequest struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type UpdateChefRequest struct {
	ID        uint
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type ChefResponse struct {
	ID        uint
	FirstName string
	LastName  string
	Email     string
}

func NewChefService(repo Repository, token TokenProvider) *ChefService {
	return &ChefService{Repo: repo, Token: token}
}

func (s *ChefService) CreateChef(req CreateChefRequest) (*ChefResponse, error) {
	chef, err := domain.NewChef(req.FirstName, req.LastName, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	if err := s.Repo.Save(chef); err != nil {
		return nil, err
	}
	return &ChefResponse{ID: chef.ID, FirstName: chef.FirstName, LastName: chef.LastName, Email: chef.Email}, nil
}

func (s *ChefService) ListChefs() ([]ChefResponse, error) {
	chefs, err := s.Repo.List()
	if err != nil {
		return nil, err
	}
	out := make([]ChefResponse, 0, len(chefs))
	for _, c := range chefs {
		out = append(out, ChefResponse{ID: c.ID, FirstName: c.FirstName, LastName: c.LastName, Email: c.Email})
	}
	return out, nil
}

func (s *ChefService) GetChef(id uint) (*ChefResponse, error) {
	chef, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return &ChefResponse{ID: chef.ID, FirstName: chef.FirstName, LastName: chef.LastName, Email: chef.Email}, nil
}

func (s *ChefService) Login(email, password string) (string, *ChefResponse, error) {
	chef, err := s.Repo.FindByEmail(email)
	if err != nil {
		return "", nil, err
	}
	if !chef.VerifyPassword(password) {
		return "", nil, domain.ErrInvalidCredentials
	}
	token, err := s.Token.CreateToken(chef.ID)
	if err != nil {
		return "", nil, err
	}
	return token, &ChefResponse{ID: chef.ID, FirstName: chef.FirstName, LastName: chef.LastName, Email: chef.Email}, nil
}
