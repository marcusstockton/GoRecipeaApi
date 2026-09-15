package recipe

import (
	"strings"
	"time"
)

type Ingredient struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	RecipeID    uint   `json:"-"`
	Quantity    string `json:"quantity"`
	Unit        string `json:"unit"`
	Name        string `json:"name"`
	Preparation string `json:"preparation,omitempty"`
}

type RecipeLike struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	RecipeID  uint      `gorm:"not null;uniqueIndex:idx_recipe_like_recipe_chef" json:"recipe_id"`
	ChefID    uint      `gorm:"not null;uniqueIndex:idx_recipe_like_recipe_chef" json:"chef_id"`
	CreatedAt time.Time `json:"created_at"`
}

type RecipeComment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	RecipeID  uint      `json:"recipe_id"`
	ChefID    uint      `json:"chef_id"`
	ParentID  *uint     `json:"parent_id,omitempty"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Step struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	RecipeID uint   `json:"-"`
	Order    int    `json:"order"`
	Action   string `json:"action"`
}

type Recipe struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	Title       string          `gorm:"size:255;not null;uniqueIndex:idx_recipe_owner_title" json:"title"`
	Description string          `json:"description"`
	OwnedBy     uint            `gorm:"not null;uniqueIndex:idx_recipe_owner_title" json:"owned_by"`
	Ingredients []Ingredient    `json:"ingredients" gorm:"foreignKey:RecipeID;constraint:OnDelete:CASCADE"`
	Steps       []Step          `json:"steps" gorm:"foreignKey:RecipeID;constraint:OnDelete:CASCADE"`
	CreatedAt   time.Time       `json:"created_at"`
	Likes       []RecipeLike    `json:"likes" gorm:"foreignKey:RecipeID;constraint:OnDelete:CASCADE"`
	Comments    []RecipeComment `json:"comments" gorm:"foreignKey:RecipeID;constraint:OnDelete:CASCADE"`
}

func NewRecipe(title, description string, ownedBy uint, ingredients []Ingredient, steps []Step) (*Recipe, error) {
	if title == "" || ownedBy == 0 {
		return nil, ErrInvalidInput
	}
	if len(ingredients) == 0 {
		return nil, ErrInvalidInput
	}
	if len(steps) == 0 {
		return nil, ErrInvalidInput
	}
	for _, ing := range ingredients {
		if ing.Name == "" || ing.Quantity == "" {
			return nil, ErrInvalidInput
		}
	}
	for _, st := range steps {
		if st.Action == "" {
			return nil, ErrInvalidInput
		}
	}
	return &Recipe{Title: title, Description: description, OwnedBy: ownedBy, Ingredients: ingredients, Steps: steps}, nil
}

func (r *Recipe) Update(title, description string, ingredients []Ingredient, steps []Step) error {
	if title == "" {
		return ErrInvalidInput
	}
	if len(ingredients) == 0 || len(steps) == 0 {
		return ErrInvalidInput
	}
	if len(ingredients) > 0 {
		for _, ing := range ingredients {
			if ing.Name == "" || ing.Quantity == "" {
				return ErrInvalidInput
			}
		}
	}
	if len(steps) > 0 {
		for _, st := range steps {
			if st.Action == "" {
				return ErrInvalidInput
			}
		}
	}
	r.Title = title
	r.Description = description
	r.Ingredients = ingredients
	r.Steps = steps
	return nil
}

func (r *Recipe) CanUserLike(chefID uint) error {
	if chefID == 0 {
		return ErrInvalidInput
	}
	if r.OwnedBy == chefID {
		return ErrSelfLikeNotAllowed
	}
	for _, like := range r.Likes {
		if like.ChefID == chefID {
			return ErrDuplicateLike
		}
	}
	return nil
}

func (r *Recipe) CanRemoveLike(chefID uint) error {
	if chefID == 0 {
		return ErrInvalidInput
	}
	if len(r.Likes) == 0 {
		return ErrNotFound
	}
	for _, like := range r.Likes {
		if like.ChefID == chefID {
			return nil
		}
	}
	return ErrNotAuthorized
}

func (r *Recipe) CanAddComment(chefID uint, content string, parentCommentID *uint) error {
	if chefID == 0 {
		return ErrInvalidInput
	}
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return ErrEmptyComment
	}
	if parentCommentID != nil {
		for _, comment := range r.Comments {
			if comment.ID == *parentCommentID {
				if comment.RecipeID != r.ID {
					return ErrInvalidInput
				}
				return nil
			}
		}
		return ErrNotFound
	}
	return nil
}
