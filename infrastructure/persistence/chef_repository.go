package persistence

import (
	"errors"

	"gorm.io/gorm"
	domain "recipea.com/m/domain/chef"
)

type ChefRepository struct {
	DB *gorm.DB
}

func NewChefRepository(db *gorm.DB) *ChefRepository {
	return &ChefRepository{DB: db}
}

func (r *ChefRepository) Save(chef *domain.Chef) error {
	if err := r.DB.Save(chef).Error; err != nil {
		return err
	}
	return nil
}

func (r *ChefRepository) FindByID(id uint) (*domain.Chef, error) {
	var chef domain.Chef
	if err := r.DB.First(&chef, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &chef, nil
}

func (r *ChefRepository) FindByEmail(email string) (*domain.Chef, error) {
	var chef domain.Chef
	if err := r.DB.First(&chef, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &chef, nil
}

func (r *ChefRepository) Delete(id uint) error {
	if err := r.DB.Delete(&domain.Chef{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *ChefRepository) List() ([]domain.Chef, error) {
	var chefs []domain.Chef
	if err := r.DB.Find(&chefs).Error; err != nil {
		return nil, err
	}
	return chefs, nil
}
