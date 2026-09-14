package recipe

import (
	"testing"
)

func validIngredients() []Ingredient {
	return []Ingredient{
		{
			Name:     "Chicken",
			Quantity: "500",
			Unit:     "g",
		},
		{
			Name:     "Rice",
			Quantity: "200",
			Unit:     "g",
		},
	}
}

func validSteps() []Step {
	return []Step{
		{
			Order:  1,
			Action: "Cook the chicken",
		},
		{
			Order:  2,
			Action: "Cook the rice",
		},
	}
}

func TestNewRecipe(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		description string
		ownedBy     uint
		ingredients []Ingredient
		steps       []Step
		wantErr     error
	}{
		{
			name:        "creates valid recipe",
			title:       "Chicken and Rice",
			description: "A simple chicken and rice recipe",
			ownedBy:     1,
			ingredients: validIngredients(),
			steps:       validSteps(),
			wantErr:     nil,
		},
		{
			name:        "rejects empty title",
			title:       "",
			description: "A recipe",
			ownedBy:     1,
			ingredients: validIngredients(),
			steps:       validSteps(),
			wantErr:     ErrInvalidInput,
		},
		{
			name:        "rejects zero owner",
			title:       "Chicken and Rice",
			description: "A recipe",
			ownedBy:     0,
			ingredients: validIngredients(),
			steps:       validSteps(),
			wantErr:     ErrInvalidInput,
		},
		{
			name:        "rejects empty ingredients",
			title:       "Chicken and Rice",
			description: "A recipe",
			ownedBy:     1,
			ingredients: []Ingredient{},
			steps:       validSteps(),
			wantErr:     ErrInvalidInput,
		},
		{
			name:        "rejects nil ingredients",
			title:       "Chicken and Rice",
			description: "A recipe",
			ownedBy:     1,
			ingredients: nil,
			steps:       validSteps(),
			wantErr:     ErrInvalidInput,
		},
		{
			name:        "rejects empty steps",
			title:       "Chicken and Rice",
			description: "A recipe",
			ownedBy:     1,
			ingredients: validIngredients(),
			steps:       []Step{},
			wantErr:     ErrInvalidInput,
		},
		{
			name:        "rejects nil steps",
			title:       "Chicken and Rice",
			description: "A recipe",
			ownedBy:     1,
			ingredients: validIngredients(),
			steps:       nil,
			wantErr:     ErrInvalidInput,
		},
		{
			name:        "rejects ingredient with empty name",
			title:       "Chicken and Rice",
			description: "A recipe",
			ownedBy:     1,
			ingredients: []Ingredient{
				{
					Name:     "",
					Quantity: "500",
				},
			},
			steps:   validSteps(),
			wantErr: ErrInvalidInput,
		},
		{
			name:        "rejects ingredient with empty quantity",
			title:       "Chicken and Rice",
			description: "A recipe",
			ownedBy:     1,
			ingredients: []Ingredient{
				{
					Name:     "Chicken",
					Quantity: "",
				},
			},
			steps:   validSteps(),
			wantErr: ErrInvalidInput,
		},
		{
			name:        "rejects step with empty action",
			title:       "Chicken and Rice",
			description: "A recipe",
			ownedBy:     1,
			ingredients: validIngredients(),
			steps: []Step{
				{
					Order:  1,
					Action: "",
				},
			},
			wantErr: ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recipe, err := NewRecipe(
				tt.title,
				tt.description,
				tt.ownedBy,
				tt.ingredients,
				tt.steps,
			)

			if err != tt.wantErr {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}

			if tt.wantErr != nil {
				if recipe != nil {
					t.Fatal("expected recipe to be nil when an error occurs")
				}
				return
			}

			if recipe == nil {
				t.Fatal("expected recipe to be created")
			}

			if recipe.Title != tt.title {
				t.Errorf("expected title %q, got %q", tt.title, recipe.Title)
			}

			if recipe.Description != tt.description {
				t.Errorf(
					"expected description %q, got %q",
					tt.description,
					recipe.Description,
				)
			}

			if recipe.OwnedBy != tt.ownedBy {
				t.Errorf("expected owner %d, got %d", tt.ownedBy, recipe.OwnedBy)
			}

			if len(recipe.Ingredients) != len(tt.ingredients) {
				t.Errorf(
					"expected %d ingredients, got %d",
					len(tt.ingredients),
					len(recipe.Ingredients),
				)
			}

			if len(recipe.Steps) != len(tt.steps) {
				t.Errorf(
					"expected %d steps, got %d",
					len(tt.steps),
					len(recipe.Steps),
				)
			}
		})
	}
}

func TestRecipeUpdate(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		description string
		ingredients []Ingredient
		steps       []Step
		wantErr     error
	}{
		{
			name:        "updates valid recipe",
			title:       "Updated Chicken and Rice",
			description: "An updated recipe",
			ingredients: []Ingredient{
				{
					Name:     "Chicken",
					Quantity: "750",
					Unit:     "g",
				},
				{
					Name:     "Rice",
					Quantity: "300",
					Unit:     "g",
				},
			},
			steps: []Step{
				{
					Order:  1,
					Action: "Prepare the chicken",
				},
				{
					Order:  2,
					Action: "Cook everything together",
				},
			},
			wantErr: nil,
		},
		{
			name:        "rejects empty title",
			title:       "",
			description: "Updated description",
			ingredients: validIngredients(),
			steps:       validSteps(),
			wantErr:     ErrInvalidInput,
		},
		{
			name:        "rejects empty ingredients",
			title:       "Updated Recipe",
			description: "Updated description",
			ingredients: []Ingredient{},
			steps:       validSteps(),
			wantErr:     ErrInvalidInput,
		},
		{
			name:        "rejects nil ingredients",
			title:       "Updated Recipe",
			description: "Updated description",
			ingredients: nil,
			steps:       validSteps(),
			wantErr:     ErrInvalidInput,
		},
		{
			name:        "rejects empty steps",
			title:       "Updated Recipe",
			description: "Updated description",
			ingredients: validIngredients(),
			steps:       []Step{},
			wantErr:     ErrInvalidInput,
		},
		{
			name:        "rejects nil steps",
			title:       "Updated Recipe",
			description: "Updated description",
			ingredients: validIngredients(),
			steps:       nil,
			wantErr:     ErrInvalidInput,
		},
		{
			name:        "rejects ingredient with empty name",
			title:       "Updated Recipe",
			description: "Updated description",
			ingredients: []Ingredient{
				{
					Name:     "",
					Quantity: "500",
				},
			},
			steps:   validSteps(),
			wantErr: ErrInvalidInput,
		},
		{
			name:        "rejects ingredient with empty quantity",
			title:       "Updated Recipe",
			description: "Updated description",
			ingredients: []Ingredient{
				{
					Name:     "Chicken",
					Quantity: "",
				},
			},
			steps:   validSteps(),
			wantErr: ErrInvalidInput,
		},
		{
			name:        "rejects step with empty action",
			title:       "Updated Recipe",
			description: "Updated description",
			ingredients: validIngredients(),
			steps: []Step{
				{
					Order:  1,
					Action: "",
				},
			},
			wantErr: ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recipe := &Recipe{
				ID:          1,
				Title:       "Original Recipe",
				Description: "Original description",
				OwnedBy:     1,
				Ingredients: validIngredients(),
				Steps:       validSteps(),
			}

			err := recipe.Update(
				tt.title,
				tt.description,
				tt.ingredients,
				tt.steps,
			)

			if err != tt.wantErr {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}

			if tt.wantErr != nil {
				// The recipe should remain unchanged when validation fails.
				if recipe.Title != "Original Recipe" {
					t.Errorf("recipe title was changed despite validation failure")
				}

				if recipe.Description != "Original description" {
					t.Errorf("recipe description was changed despite validation failure")
				}

				return
			}

			if recipe.Title != tt.title {
				t.Errorf("expected title %q, got %q", tt.title, recipe.Title)
			}

			if recipe.Description != tt.description {
				t.Errorf(
					"expected description %q, got %q",
					tt.description,
					recipe.Description,
				)
			}

			if len(recipe.Ingredients) != len(tt.ingredients) {
				t.Errorf(
					"expected %d ingredients, got %d",
					len(tt.ingredients),
					len(recipe.Ingredients),
				)
			}

			if len(recipe.Steps) != len(tt.steps) {
				t.Errorf(
					"expected %d steps, got %d",
					len(tt.steps),
					len(recipe.Steps),
				)
			}
		})
	}
}

func TestRecipeUpdatePreservesIDAndOwner(t *testing.T) {
	recipe := &Recipe{
		ID:          123,
		Title:       "Original Recipe",
		Description: "Original description",
		OwnedBy:     456,
		Ingredients: validIngredients(),
		Steps:       validSteps(),
	}

	err := recipe.Update(
		"Updated Recipe",
		"Updated description",
		validIngredients(),
		validSteps(),
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if recipe.ID != 123 {
		t.Errorf("expected ID to remain 123, got %d", recipe.ID)
	}

	if recipe.OwnedBy != 456 {
		t.Errorf("expected owner to remain 456, got %d", recipe.OwnedBy)
	}
}

func TestNewRecipeDoesNotSetDatabaseFields(t *testing.T) {
	recipe, err := NewRecipe(
		"Chicken and Rice",
		"A simple recipe",
		1,
		validIngredients(),
		validSteps(),
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if recipe.ID != 0 {
		t.Errorf("expected ID to be zero, got %d", recipe.ID)
	}

	if !recipe.CreatedAt.IsZero() {
		t.Errorf("expected CreatedAt to be zero")
	}
}
