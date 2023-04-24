package recipes

import (
	"context"

	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/repositories/models"
	"github.com/romankravchuk/muerta/internal/repositories/recipes"
)

type RecipeServicer interface {
	CreateRecipe(ctx context.Context, payload *dto.CreateRecipeDTO) error
	// UpdateRecipe(ctx context.Context)
	// DeleteRecipe(ctx context.Context)
	// FindRecipeByID(ctx context.Context)
	// FindRecipeByName(ctx context.Context)
	// FindManyRecipes(ctx context.Context)
}

type RecipeService struct {
	repository recipes.RecipesRepositorer
}

func New(repository recipes.RecipesRepositorer) *RecipeService {
	return &RecipeService{repository: repository}
}

func (s *RecipeService) CreateRecipe(ctx context.Context, payload *dto.CreateRecipeDTO) error {
	steps := make([]models.RecipeStep, len(payload.Steps))
	for i, step := range payload.Steps {
		steps[i].StepID = step.ID
		steps[i].Place = step.Place
	}
	err := s.repository.Create(ctx, &models.Recipe{
		Name:        payload.Name,
		Description: payload.Description,
		Steps:       steps,
	})
	return err
}
