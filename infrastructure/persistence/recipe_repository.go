package persistence

import (
	"gorm.io/gorm"
	"recipea.com/m/domain/recipe"
)

type RecipeRepository struct {
	DB *gorm.DB
}

func NewRecipeRepository(db *gorm.DB) *RecipeRepository {
	return &RecipeRepository{DB: db}
}

func (r *RecipeRepository) Save(recipe *recipe.Recipe) error {
	if err := r.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(recipe).Error; err != nil {
		return err
	}
	return nil
}

func (r *RecipeRepository) FindByID(id uint) (*recipe.Recipe, error) {
	var rec recipe.Recipe
	if err := r.DB.Preload("Ingredients").Preload("Steps").First(&rec, id).Error; err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *RecipeRepository) List() ([]recipe.Recipe, error) {
	var recs []recipe.Recipe
	if err := r.DB.Preload("Ingredients").Preload("Steps").Find(&recs).Error; err != nil {
		return nil, err
	}
	return recs, nil
}

func (r *RecipeRepository) Delete(id uint) error {
	return r.DB.Delete(&recipe.Recipe{}, id).Error
}
