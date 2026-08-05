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

	if err := db.AutoMigrate(&domainrecipe.Recipe{}, &domainrecipe.Ingredient{}, &domainrecipe.Step{}, &domainrecipe.RecipeLike{}, &domainrecipe.RecipeComment{}); err != nil {
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

func TestRecipeLikeHasUniqueConstraintPerRecipeAndChef(t *testing.T) {
	db := setupTestDB(t)
	repo := persistence.NewRecipeRepository(db)

	recipeModel := &domainrecipe.Recipe{
		Title:       "Mango Smoothie",
		Description: "Refreshing drink",
		OwnedBy:     1,
		Ingredients: []domainrecipe.Ingredient{{Quantity: "1", Unit: "cup", Name: "mango"}},
		Steps:       []domainrecipe.Step{{Order: 1, Action: "Blend"}},
	}
	if err := repo.Save(recipeModel); err != nil {
		t.Fatalf("save recipe failed: %v", err)
	}

	if err := repo.AddLike(recipeModel.ID, 7); err != nil {
		t.Fatalf("first like should succeed: %v", err)
	}
	if err := repo.AddLike(recipeModel.ID, 7); err == nil {
		t.Fatalf("expected duplicate like to fail because of the unique constraint")
	}
}

func TestRecipeResponsesIncludeLikesButNotComments(t *testing.T) {
	db := setupTestDB(t)
	service := NewRecipeService(persistence.NewRecipeRepository(db))

	created, err := service.CreateRecipe(CreateRecipeRequest{
		Title:       "Berry Tart",
		Description: "Sweet pastry",
		OwnedBy:     5,
		Ingredients: []domainrecipe.Ingredient{{Quantity: "1", Unit: "box", Name: "berries"}},
		Steps:       []domainrecipe.Step{{Order: 1, Action: "Bake"}},
	})
	if err != nil {
		t.Fatalf("create recipe failed: %v", err)
	}

	if err := service.AddLike(created.ID, 8); err != nil {
		t.Fatalf("add like failed: %v", err)
	}

	got, err := service.GetRecipe(created.ID)
	if err != nil {
		t.Fatalf("get recipe failed: %v", err)
	}
	if got.LikeCount != 1 {
		t.Fatalf("expected recipe response to include one like, got %d", got.LikeCount)
	}
}

func TestChefCannotLikeTheirOwnRecipeButCanCommentAndReply(t *testing.T) {
	db := setupTestDB(t)
	service := NewRecipeService(persistence.NewRecipeRepository(db))

	created, err := service.CreateRecipe(CreateRecipeRequest{
		Title:       "Honey Cake",
		Description: "Sweet",
		OwnedBy:     20,
		Ingredients: []domainrecipe.Ingredient{{Quantity: "1", Unit: "cup", Name: "honey"}},
		Steps:       []domainrecipe.Step{{Order: 1, Action: "Bake"}},
	})
	if err != nil {
		t.Fatalf("create recipe failed: %v", err)
	}

	if err := service.AddLike(created.ID, 20); !errors.Is(err, domainrecipe.ErrSelfLikeNotAllowed) {
		t.Fatalf("expected self-like to be forbidden, got %v", err)
	}
	if err := service.AddComment(created.ID, 20, "Made by me", nil); err != nil {
		t.Fatalf("owner should be able to comment on own recipe: %v", err)
	}

	comments, err := service.ListComments(created.ID)
	if err != nil {
		t.Fatalf("list comments failed: %v", err)
	}
	if len(comments) != 1 {
		t.Fatalf("expected one comment, got %d", len(comments))
	}
	if err := service.AddComment(created.ID, 21, "Replying to the owner", &comments[0].ID); err != nil {
		t.Fatalf("another chef should be able to reply to the comment: %v", err)
	}
	comments, err = service.ListComments(created.ID)
	if err != nil {
		t.Fatalf("list comments after reply failed: %v", err)
	}
	if len(comments) != 2 {
		t.Fatalf("expected two comments after reply, got %d", len(comments))
	}
}

func TestLikeAndCommentLifecycleRespectsChefOwnership(t *testing.T) {
	db := setupTestDB(t)
	service := NewRecipeService(persistence.NewRecipeRepository(db))

	created, err := service.CreateRecipe(CreateRecipeRequest{
		Title:       "Lemon Tart",
		Description: "Bright citrus dessert",
		OwnedBy:     10,
		Ingredients: []domainrecipe.Ingredient{{Quantity: "2", Unit: "cups", Name: "flour"}},
		Steps:       []domainrecipe.Step{{Order: 1, Action: "Bake it"}},
	})
	if err != nil {
		t.Fatalf("create recipe failed: %v", err)
	}

	if err := service.AddLike(created.ID, 11); err != nil {
		t.Fatalf("first like failed: %v", err)
	}
	if err := service.AddLike(created.ID, 11); !errors.Is(err, domainrecipe.ErrDuplicateLike) {
		t.Fatalf("expected duplicate like error, got %v", err)
	}
	if err := service.RemoveLike(created.ID, 12); !errors.Is(err, domainrecipe.ErrNotAuthorized) {
		t.Fatalf("expected authorization error when removing another chef's like, got %v", err)
	}
	if err := service.RemoveLike(created.ID, 11); err != nil {
		t.Fatalf("removing own like failed: %v", err)
	}

	if err := service.AddComment(created.ID, 11, "Very nice recipe", nil); err != nil {
		t.Fatalf("add comment failed: %v", err)
	}
	comments, err := service.ListComments(created.ID)
	if err != nil {
		t.Fatalf("list comments failed: %v", err)
	}
	if len(comments) != 1 {
		t.Fatalf("expected one comment, got %d", len(comments))
	}
	if err := service.RemoveComment(comments[0].ID, 12); !errors.Is(err, domainrecipe.ErrNotAuthorized) {
		t.Fatalf("expected authorization error for another chef removing comment, got %v", err)
	}
	if err := service.RemoveComment(comments[0].ID, 11); err != nil {
		t.Fatalf("removing own comment failed: %v", err)
	}
}
