package recipe

import "time"

type Ingredient struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	RecipeID    uint   `json:"-"`
	Quantity    string `json:"quantity"`
	Unit        string `json:"unit"`
	Name        string `json:"name"`
	Preparation string `json:"preparation,omitempty"`
}

type Step struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	RecipeID uint   `json:"-"`
	Order    int    `json:"order"`
	Action   string `json:"action"`
}

type Recipe struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	OwnedBy     uint         `json:"owned_by"`
	Ingredients []Ingredient `json:"ingredients" gorm:"foreignKey:RecipeID;constraint:OnDelete:CASCADE"`
	Steps       []Step       `json:"steps" gorm:"foreignKey:RecipeID;constraint:OnDelete:CASCADE"`
	CreatedAt   time.Time    `json:"created_at"`
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
