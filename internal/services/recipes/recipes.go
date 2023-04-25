package recipes

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/pkg/translate"
	"github.com/romankravchuk/muerta/internal/repositories/models"
	"github.com/romankravchuk/muerta/internal/repositories/recipes"
)

type RecipeServicer interface {
	CreateRecipe(ctx context.Context, payload *dto.CreateRecipeDTO) error
	FindRecipeByID(ctx context.Context, id int) (dto.FindRecipeDTO, error)
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
	steps := make([]models.Step, len(payload.Steps))
	for i, step := range payload.Steps {
		steps[i].ID = step.ID
		steps[i].Place = step.Place
	}
	err := s.repository.Create(ctx, &models.Recipe{
		Name:        payload.Name,
		Description: payload.Description,
		Steps:       steps,
	})
	return err
}

func (s *RecipeService) FindRecipeByID(ctx context.Context, id int) (dto.FindRecipeDTO, error) {
	recipe, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return dto.FindRecipeDTO{}, fmt.Errorf("recipe not found by id: %w", err)
	}
	result := translate.RecipeModelToFindDTO(recipe)
	return result, nil
}
