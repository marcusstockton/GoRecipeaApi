package recipe

import (
	"errors"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	domainrecipe "recipea.com/m/domain/recipe"
	"recipea.com/m/infrastructure/persistence"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	if err := db.AutoMigrate(&domainrecipe.Recipe{}, &domainrecipe.Ingredient{}, &domainrecipe.Step{}); err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}

	return db
}

func TestCreateRecipeRejectsDuplicateTitlePerChef(t *testing.T) {
	db := setupTestDB(t)
	service := NewRecipeService(persistence.NewRecipeRepository(db))

	req := CreateRecipeRequest{
		Title:       "Spaghetti Bolognese",
		Description: "A classic pasta",
		OwnedBy:     10,
		Ingredients: []domainrecipe.Ingredient{{Quantity: "2", Unit: "cups", Name: "pasta"}},
		Steps:       []domainrecipe.Step{{Order: 1, Action: "Boil pasta"}},
	}

	if _, err := service.CreateRecipe(req); err != nil {
		t.Fatalf("first create failed: %v", err)
	}

	_, err := service.CreateRecipe(req)
	if !errors.Is(err, domainrecipe.ErrDuplicateRecipe) {
		t.Fatalf("expected duplicate recipe error, got %v", err)
	}
}

func TestCreateRecipeAllowsSameTitleForDifferentChefs(t *testing.T) {
	db := setupTestDB(t)
	service := NewRecipeService(persistence.NewRecipeRepository(db))

	firstReq := CreateRecipeRequest{
		Title:       "Spaghetti Bolognese",
		Description: "A classic pasta",
		OwnedBy:     10,
		Ingredients: []domainrecipe.Ingredient{{Quantity: "2", Unit: "cups", Name: "pasta"}},
		Steps:       []domainrecipe.Step{{Order: 1, Action: "Boil pasta"}},
	}
	secondReq := firstReq
	secondReq.OwnedBy = 20

	if _, err := service.CreateRecipe(firstReq); err != nil {
		t.Fatalf("first create failed: %v", err)
	}

	if _, err := service.CreateRecipe(secondReq); err != nil {
		t.Fatalf("second create should be allowed for different chef: %v", err)
	}
}
