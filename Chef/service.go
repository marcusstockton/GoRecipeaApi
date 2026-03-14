package Chef

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	domain "recipea.com/m/shared"
)

type Service struct {
	DB *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

func (s *Service) GetAllChefs() ([]domain.Chef, error) {
	var chefs []domain.Chef
	if err := s.DB.Find(&chefs).Error; err != nil {
		return nil, err
	}
	return chefs, nil
}

func (s *Service) CreateChef(input domain.Chef) (domain.Chef, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.Chef{}, err
	}
	chef := domain.Chef{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  string(hashedPassword),
	}
	if err := s.DB.Create(&chef).Error; err != nil {
		return domain.Chef{}, err
	}
	return chef, nil
}

func (s *Service) GetChefByID(id string) (domain.Chef, error) {
	var chef domain.Chef
	if err := s.DB.First(&chef, id).Error; err != nil {
		return domain.Chef{}, err
	}
	return chef, nil
}

func (s *Service) FindChefByEmail(email string) (domain.Chef, error) {
	var chef domain.Chef
	if err := s.DB.First(&chef, "email = ?", email).Error; err != nil {
		return domain.Chef{}, err
	}
	return chef, nil
}

func (s *Service) UpdateChef(id string, updates map[string]any) (domain.Chef, error) {
	var chef domain.Chef
	if err := s.DB.First(&chef, id).Error; err != nil {
		return domain.Chef{}, err
	}
	if err := s.DB.Model(&chef).Updates(updates).Error; err != nil {
		return domain.Chef{}, err
	}
	if err := s.DB.First(&chef, id).Error; err != nil {
		return domain.Chef{}, err
	}
	return chef, nil
}

func (s *Service) DeleteChef(id string) error {
	if err := s.DB.Delete(&domain.Chef{}, id).Error; err != nil {
		return err
	}
	return nil
}
