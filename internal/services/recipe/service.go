package recipe

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/pkg/translate"
	recipes "github.com/romankravchuk/muerta/internal/repositories/recipe"
	"github.com/romankravchuk/muerta/internal/services"
)

type RecipeServicer interface {
	CreateRecipe(ctx context.Context, payload *dto.CreateRecipeDTO) error
	FindRecipeByID(ctx context.Context, id int) (dto.FindRecipeDTO, error)
	FindRecipes(ctx context.Context, filter *dto.RecipeFilterDTO) ([]dto.FindRecipeDTO, error)
	UpdateRecipe(ctx context.Context, id int, payload *dto.UpdateRecipeDTO) error
	DeleteRecipe(ctx context.Context, id int) error
	RestoreRecipe(ctx context.Context, id int) error
	FindRecipeIngredients(ctx context.Context, id int) ([]dto.FindRecipeIngredientDTO, error)
	CreateIngredient(ctx context.Context, id int, payload *dto.CreateIngredientDTO) (dto.FindRecipeIngredientDTO, error)
	UpdateIngredient(ctx context.Context, id int, payload *dto.UpdateIngredientDTO) (dto.FindRecipeIngredientDTO, error)
	DeleteIngredient(ctx context.Context, id int, payload *dto.DeleteIngredientDTO) error
	FindRecipeSteps(ctx context.Context, recipeID int) ([]dto.FindStepDTO, error)
	CreateRecipeStep(ctx context.Context, recipeID, stepID, place int) (dto.FindStepDTO, error)
	DeleteRecipeStep(ctx context.Context, recipeID, stepID, place int) error
	services.Counter
}

type recipeService struct {
	repo recipes.RecipesRepositorer
}

func (s *recipeService) Count(ctx context.Context) (int, error) {
	count, err := s.repo.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("error counting recipes: %w", err)
	}
	return count, nil
}

// CreateRecipeStep implements RecipeServicer
func (s *recipeService) CreateRecipeStep(ctx context.Context, recipeID int, stepID int, place int) (dto.FindStepDTO, error) {
	model, err := s.repo.CreateStep(ctx, recipeID, stepID, place)
	if err != nil {
		return dto.FindStepDTO{}, fmt.Errorf("create step: %w", err)
	}
	return translate.StepModelToFindDTO(model), nil
}

// DeleteRecipeStep implements RecipeServicer
func (s *recipeService) DeleteRecipeStep(ctx context.Context, recipeID int, stepID int, place int) error {
	if err := s.repo.DeleteStep(ctx, recipeID, stepID, place); err != nil {
		return fmt.Errorf("delete step: %w", err)
	}
	return nil
}

// FindRecipeSteps implements RecipeServicer
func (s *recipeService) FindRecipeSteps(ctx context.Context, recipeID int) ([]dto.FindStepDTO, error) {
	entities, err := s.repo.FindSteps(ctx, recipeID)
	if err != nil {
		return nil, fmt.Errorf("steps not found: %w", err)
	}
	return translate.StepModelsToFindDTOs(entities), nil
}

// CreateIngredient implements RecipeServicer
func (s *recipeService) CreateIngredient(ctx context.Context, id int, payload *dto.CreateIngredientDTO) (dto.FindRecipeIngredientDTO, error) {
	model := translate.CreateIngredientDTOToModel(payload)
	ingredient, err := s.repo.CreateIngredient(ctx, id, &model)
	if err != nil {
		return dto.FindRecipeIngredientDTO{}, fmt.Errorf("create recipe ingredient: %w", err)
	}
	dto := translate.RecipeIngredientModelToFindDTO(&ingredient)
	return dto, nil
}

// DeleteIngredient implements RecipeServicer
func (s *recipeService) DeleteIngredient(ctx context.Context, id int, payload *dto.DeleteIngredientDTO) error {
	if err := s.repo.DeleteIngredient(ctx, id, payload.ProductID); err != nil {
		return fmt.Errorf("delete recipe ingredient: %w", err)
	}
	return nil
}

// FindRecipeIngredients implements RecipeServicer
func (s *recipeService) FindRecipeIngredients(ctx context.Context, id int) ([]dto.FindRecipeIngredientDTO, error) {
	ingredients, err := s.repo.FindIngredients(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ingredients not found: %w", err)
	}
	result := translate.RecipeIngredientModelsToFindDTOs(ingredients)
	return result, nil
}

// UpdateIngredient implements RecipeServicer
func (s *recipeService) UpdateIngredient(ctx context.Context, id int, payload *dto.UpdateIngredientDTO) (dto.FindRecipeIngredientDTO, error) {
	model := translate.UpdateIngredientDTOToModel(payload)
	ingredient, err := s.repo.UpdateIngredient(ctx, id, &model)
	if err != nil {
		return dto.FindRecipeIngredientDTO{}, fmt.Errorf("update recipe ingredient: %w", err)
	}
	dto := translate.RecipeIngredientModelToFindDTO(&ingredient)
	return dto, nil
}

func New(repository recipes.RecipesRepositorer) RecipeServicer {
	return &recipeService{repo: repository}
}

func (s *recipeService) CreateRecipe(ctx context.Context, payload *dto.CreateRecipeDTO) error {
	fmt.Printf("%#v\n", payload)
	model := translate.CreateRecipeDTOToModel(payload)
	fmt.Printf("%#v\n", model)
	if err := s.repo.Create(ctx, &model); err != nil {
		return fmt.Errorf("create recipe: %w", err)
	}
	return nil
}

func (s *recipeService) FindRecipeByID(ctx context.Context, id int) (dto.FindRecipeDTO, error) {
	recipe, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return dto.FindRecipeDTO{}, fmt.Errorf("recipe not found by id: %w", err)
	}
	result := translate.RecipeModelToFindDTO(&recipe)
	return result, nil
}

func (s *recipeService) FindRecipes(ctx context.Context, filter *dto.RecipeFilterDTO) ([]dto.FindRecipeDTO, error) {
	recipes, err := s.repo.FindMany(ctx, filter.Limit, filter.Offset, filter.Name)
	if err != nil {
		return nil, fmt.Errorf("recipes not found: %w", err)
	}
	result := translate.RecipeModelsToFindDTOs(recipes)
	return result, nil
}

func (s *recipeService) UpdateRecipe(ctx context.Context, id int, payload *dto.UpdateRecipeDTO) error {
	recipe, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("recipe not found by id: %w", err)
	}
	if payload.Name != "" {
		recipe.Name = payload.Name
	}
	if payload.Description != "" {
		recipe.Description = payload.Description
	}
	if err := s.repo.Update(ctx, &recipe); err != nil {
		return fmt.Errorf("update recipe: %w", err)
	}
	return nil
}

func (s *recipeService) DeleteRecipe(ctx context.Context, id int) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete recipe: %w", err)
	}
	return nil
}

func (s *recipeService) RestoreRecipe(ctx context.Context, id int) error {
	if err := s.repo.Restore(ctx, id); err != nil {
		return fmt.Errorf("restore recipe: %w", err)
	}
	return nil
}
