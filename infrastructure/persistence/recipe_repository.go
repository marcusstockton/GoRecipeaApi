package persistence

import (
	"errors"

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
	if err := r.DB.Preload("Ingredients").Preload("Steps").Preload("Likes").Preload("Comments").First(&rec, id).Error; err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *RecipeRepository) List() ([]recipe.Recipe, error) {
	var recs []recipe.Recipe
	if err := r.DB.Preload("Ingredients").Preload("Steps").Preload("Likes").Preload("Comments").Find(&recs).Error; err != nil {
		return nil, err
	}
	return recs, nil
}

func (r *RecipeRepository) Delete(id uint) error {
	return r.DB.Delete(&recipe.Recipe{}, id).Error
}

func (r *RecipeRepository) AddLike(recipeid, chefid uint) error {
	like := &recipe.RecipeLike{
		RecipeID: recipeid,
		ChefID:   chefid,
	}
	return r.DB.Create(like).Error
}

func (r *RecipeRepository) RemoveLike(recipeid, chefid uint) error {
	return r.DB.Where("recipe_id = ? AND chef_id = ?", recipeid, chefid).Delete(&recipe.RecipeLike{}).Error
}

func (r *RecipeRepository) AddComment(recipeid, chefid uint, content string, parentID *uint) error {
	comment := &recipe.RecipeComment{
		RecipeID: recipeid,
		ChefID:   chefid,
		Content:  content,
		ParentID: parentID,
	}
	return r.DB.Create(comment).Error
}

func (r *RecipeRepository) RemoveComment(commentid, chefid uint) error {
	comment := &recipe.RecipeComment{}
	if err := r.DB.First(comment, commentid).Error; err != nil {
		return err
	}
	if comment.ChefID != chefid {
		return recipe.ErrNotAuthorized
	}
	return r.DB.Delete(&recipe.RecipeComment{}, commentid).Error
}

func (r *RecipeRepository) FindCommentByID(commentid uint) (*recipe.RecipeComment, error) {
	var comment recipe.RecipeComment
	if err := r.DB.First(&comment, commentid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, recipe.ErrNotFound
		}
		return nil, err
	}
	return &comment, nil
}

func (r *RecipeRepository) ListComments(recipeid uint) ([]recipe.RecipeComment, error) {
	var comments []recipe.RecipeComment
	if err := r.DB.Where("recipe_id = ?", recipeid).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}
