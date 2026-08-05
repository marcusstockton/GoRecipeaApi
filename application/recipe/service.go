package recipe

import (
	"strings"

	domainrecipe "recipea.com/m/domain/recipe"
)

type CreateRecipeRequest struct {
	Title       string
	Description string
	OwnedBy     uint
	Ingredients []domainrecipe.Ingredient
	Steps       []domainrecipe.Step
}

type UpdateRecipeRequest struct {
	ID          uint
	Title       string
	Description string
	UpdatedBy   uint
	Ingredients []domainrecipe.Ingredient
	Steps       []domainrecipe.Step
}

type RecipeResponse struct {
	ID          uint
	Title       string
	Description string
	OwnedBy     uint
	Ingredients []domainrecipe.Ingredient
	Steps       []domainrecipe.Step
	LikeCount   int
}

type RecipeService struct {
	Repo domainrecipe.Repository
}

func NewRecipeService(repo domainrecipe.Repository) *RecipeService {
	return &RecipeService{Repo: repo}
}

func (s *RecipeService) CreateRecipe(req CreateRecipeRequest) (*RecipeResponse, error) {
	r, err := domainrecipe.NewRecipe(req.Title, req.Description, req.OwnedBy, req.Ingredients, req.Steps)
	if err != nil {
		return nil, err
	}

	existingRecipes, err := s.Repo.List()
	if err != nil {
		return nil, err
	}
	for _, existing := range existingRecipes {
		if existing.OwnedBy == r.OwnedBy && strings.EqualFold(strings.TrimSpace(existing.Title), strings.TrimSpace(r.Title)) {
			return nil, domainrecipe.ErrDuplicateRecipe
		}
	}

	if err := s.Repo.Save(r); err != nil {
		return nil, err
	}
	return &RecipeResponse{ID: r.ID, Title: r.Title, Description: r.Description, OwnedBy: r.OwnedBy, Ingredients: r.Ingredients, Steps: r.Steps, LikeCount: 0}, nil
}

func (s *RecipeService) ListRecipes() ([]RecipeResponse, error) {
	recipes, err := s.Repo.List()
	if err != nil {
		return nil, err
	}
	out := make([]RecipeResponse, 0, len(recipes))
	for _, r := range recipes {
		out = append(out, RecipeResponse{ID: r.ID, Title: r.Title, Description: r.Description, OwnedBy: r.OwnedBy, Ingredients: r.Ingredients, Steps: r.Steps, LikeCount: len(r.Likes)})
	}
	return out, nil
}

func (s *RecipeService) UpdateRecipe(req UpdateRecipeRequest) (*RecipeResponse, error) {
	r, err := s.Repo.FindByID(req.ID)
	if err != nil {
		return nil, err
	}
	if r.OwnedBy != req.UpdatedBy {
		return nil, domainrecipe.ErrNotAuthorized
	}
	if err := r.Update(req.Title, req.Description, req.Ingredients, req.Steps); err != nil {
		return nil, err
	}
	if err := s.Repo.Save(r); err != nil {
		return nil, err
	}
	return &RecipeResponse{ID: r.ID, Title: r.Title, Description: r.Description, OwnedBy: r.OwnedBy, Ingredients: r.Ingredients, Steps: r.Steps, LikeCount: len(r.Likes)}, nil
}

func (s *RecipeService) GetRecipe(id uint) (*RecipeResponse, error) {
	r, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return &RecipeResponse{ID: r.ID, Title: r.Title, Description: r.Description, OwnedBy: r.OwnedBy, Ingredients: r.Ingredients, Steps: r.Steps, LikeCount: len(r.Likes)}, nil
}

func (s *RecipeService) DeleteRecipe(id, requestingChefID uint) error {
	r, err := s.Repo.FindByID(id)
	if err != nil {
		return err
	}
	if r.OwnedBy != requestingChefID {
		return domainrecipe.ErrNotAuthorized
	}
	return s.Repo.Delete(id)
}

func (s *RecipeService) AddLike(recipeID, chefID uint) error {
	if chefID == 0 {
		return domainrecipe.ErrInvalidInput
	}
	recipe, err := s.Repo.FindByID(recipeID)
	if err != nil {
		return err
	}
	if recipe.OwnedBy == chefID {
		return domainrecipe.ErrSelfLikeNotAllowed
	}
	for _, like := range recipe.Likes {
		if like.ChefID == chefID {
			return domainrecipe.ErrDuplicateLike
		}
	}
	return s.Repo.AddLike(recipeID, chefID)
}

func (s *RecipeService) RemoveLike(recipeID, chefID uint) error {
	if chefID == 0 {
		return domainrecipe.ErrInvalidInput
	}
	recipe, err := s.Repo.FindByID(recipeID)
	if err != nil {
		return err
	}
	for _, like := range recipe.Likes {
		if like.ChefID == chefID {
			return s.Repo.RemoveLike(recipeID, chefID)
		}
	}
	if len(recipe.Likes) > 0 {
		return domainrecipe.ErrNotAuthorized
	}
	return domainrecipe.ErrNotFound
}

func (s *RecipeService) AddComment(recipeID, chefID uint, content string, parentCommentID *uint) error {
	if chefID == 0 {
		return domainrecipe.ErrInvalidInput
	}
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return domainrecipe.ErrEmptyComment
	}
	if _, err := s.Repo.FindByID(recipeID); err != nil {
		return err
	}
	if parentCommentID != nil {
		parentComment, err := s.Repo.FindCommentByID(*parentCommentID)
		if err != nil {
			return err
		}
		if parentComment.RecipeID != recipeID {
			return domainrecipe.ErrInvalidInput
		}
	}
	return s.Repo.AddComment(recipeID, chefID, trimmed, parentCommentID)
}

func (s *RecipeService) RemoveComment(commentID, requestingChefID uint) error {
	if requestingChefID == 0 {
		return domainrecipe.ErrInvalidInput
	}
	comment, err := s.Repo.FindCommentByID(commentID)
	if err != nil {
		return err
	}
	if comment.ChefID != requestingChefID {
		return domainrecipe.ErrNotAuthorized
	}
	return s.Repo.RemoveComment(commentID, requestingChefID)
}

func (s *RecipeService) ListComments(recipeID uint) ([]domainrecipe.RecipeComment, error) {
	return s.Repo.ListComments(recipeID)
}
