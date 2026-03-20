package recipe

import domainrecipe "recipea.com/m/domain/recipe"

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
	if err := s.Repo.Save(r); err != nil {
		return nil, err
	}
	return &RecipeResponse{ID: r.ID, Title: r.Title, Description: r.Description, OwnedBy: r.OwnedBy, Ingredients: r.Ingredients, Steps: r.Steps}, nil
}

func (s *RecipeService) ListRecipes() ([]RecipeResponse, error) {
	recipes, err := s.Repo.List()
	if err != nil {
		return nil, err
	}
	out := make([]RecipeResponse, 0, len(recipes))
	for _, r := range recipes {
		out = append(out, RecipeResponse{ID: r.ID, Title: r.Title, Description: r.Description, OwnedBy: r.OwnedBy, Ingredients: r.Ingredients, Steps: r.Steps})
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
	return &RecipeResponse{ID: r.ID, Title: r.Title, Description: r.Description, OwnedBy: r.OwnedBy, Ingredients: r.Ingredients, Steps: r.Steps}, nil
}

func (s *RecipeService) GetRecipe(id uint) (*RecipeResponse, error) {
	r, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return &RecipeResponse{ID: r.ID, Title: r.Title, Description: r.Description, OwnedBy: r.OwnedBy, Ingredients: r.Ingredients, Steps: r.Steps}, nil
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
