package recipe

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/pkg/translate"
	"github.com/romankravchuk/muerta/internal/repositories/models"
	recipes "github.com/romankravchuk/muerta/internal/repositories/recipe"
)

type RecipeServicer interface {
	CreateRecipe(ctx context.Context, payload *dto.CreateRecipe) error
	FindRecipeByID(ctx context.Context, id int) (dto.FindRecipe, error)
	FindRecipes(ctx context.Context, filter *dto.RecipeFilter) ([]dto.FindRecipe, error)
	UpdateRecipe(ctx context.Context, id int, payload *dto.UpdateRecipe) error
	DeleteRecipe(ctx context.Context, id int) error
	RestoreRecipe(ctx context.Context, id int) error
	FindRecipeIngredients(ctx context.Context, id int) ([]dto.FindRecipeIngredient, error)
	CreateIngredient(
		ctx context.Context,
		id int,
		payload *dto.CreateIngredient,
	) (dto.FindRecipeIngredient, error)
	UpdateIngredient(
		ctx context.Context,
		id int,
		payload *dto.UpdateIngredient,
	) (dto.FindRecipeIngredient, error)
	DeleteIngredient(ctx context.Context, id int, payload *dto.DeleteIngredient) error
	FindRecipeSteps(ctx context.Context, recipeID int) ([]dto.FindStep, error)
	CreateRecipeStep(ctx context.Context, recipeID, stepID, place int) (dto.FindStep, error)
	DeleteRecipeStep(ctx context.Context, recipeID, stepID, place int) error
	Count(ctx context.Context, filter dto.RecipeFilter) (int, error)
}

type recipeService struct {
	repo recipes.RecipesRepositorer
}

func (s *recipeService) Count(ctx context.Context, filter dto.RecipeFilter) (int, error) {
	count, err := s.repo.Count(ctx, models.RecipeFilter{Name: filter.Name})
	if err != nil {
		return 0, fmt.Errorf("error counting recipes: %w", err)
	}
	return count, nil
}

// CreateRecipeStep implements RecipeServicer
func (s *recipeService) CreateRecipeStep(
	ctx context.Context,
	recipeID int,
	stepID int,
	place int,
) (dto.FindStep, error) {
	model, err := s.repo.CreateStep(ctx, recipeID, stepID, place)
	if err != nil {
		return dto.FindStep{}, fmt.Errorf("create step: %w", err)
	}
	return translate.StepModelToFind(model), nil
}

// DeleteRecipeStep implements RecipeServicer
func (s *recipeService) DeleteRecipeStep(
	ctx context.Context,
	recipeID int,
	stepID int,
	place int,
) error {
	if err := s.repo.DeleteStep(ctx, recipeID, stepID, place); err != nil {
		return fmt.Errorf("delete step: %w", err)
	}
	return nil
}

// FindRecipeSteps implements RecipeServicer
func (s *recipeService) FindRecipeSteps(
	ctx context.Context,
	recipeID int,
) ([]dto.FindStep, error) {
	entities, err := s.repo.FindSteps(ctx, recipeID)
	if err != nil {
		return nil, fmt.Errorf("steps not found: %w", err)
	}
	return translate.StepModelsToFinds(entities), nil
}

// CreateIngredient implements RecipeServicer
func (s *recipeService) CreateIngredient(
	ctx context.Context,
	id int,
	payload *dto.CreateIngredient,
) (dto.FindRecipeIngredient, error) {
	model := translate.CreateIngredientToModel(payload)
	ingredient, err := s.repo.CreateIngredient(ctx, id, &model)
	if err != nil {
		return dto.FindRecipeIngredient{}, fmt.Errorf("create recipe ingredient: %w", err)
	}
	dto := translate.RecipeIngredientModelToFind(&ingredient)
	return dto, nil
}

// DeleteIngredient implements RecipeServicer
func (s *recipeService) DeleteIngredient(
	ctx context.Context,
	id int,
	payload *dto.DeleteIngredient,
) error {
	if err := s.repo.DeleteIngredient(ctx, id, payload.ProductID); err != nil {
		return fmt.Errorf("delete recipe ingredient: %w", err)
	}
	return nil
}

// FindRecipeIngredients implements RecipeServicer
func (s *recipeService) FindRecipeIngredients(
	ctx context.Context,
	id int,
) ([]dto.FindRecipeIngredient, error) {
	ingredients, err := s.repo.FindIngredients(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ingredients not found: %w", err)
	}
	result := translate.RecipeIngredientModelsToFinds(ingredients)
	return result, nil
}

// UpdateIngredient implements RecipeServicer
func (s *recipeService) UpdateIngredient(
	ctx context.Context,
	id int,
	payload *dto.UpdateIngredient,
) (dto.FindRecipeIngredient, error) {
	model := translate.UpdateIngredientToModel(payload)
	ingredient, err := s.repo.UpdateIngredient(ctx, id, &model)
	if err != nil {
		return dto.FindRecipeIngredient{}, fmt.Errorf("update recipe ingredient: %w", err)
	}
	dto := translate.RecipeIngredientModelToFind(&ingredient)
	return dto, nil
}

func New(repository recipes.RecipesRepositorer) RecipeServicer {
	return &recipeService{repo: repository}
}

func (s *recipeService) CreateRecipe(ctx context.Context, payload *dto.CreateRecipe) error {
	model := translate.CreateRecipeToModel(payload)
	if err := s.repo.Create(ctx, &model); err != nil {
		return fmt.Errorf("create recipe: %w", err)
	}
	return nil
}

func (s *recipeService) FindRecipeByID(ctx context.Context, id int) (dto.FindRecipe, error) {
	recipe, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return dto.FindRecipe{}, fmt.Errorf("recipe not found by id: %w", err)
	}
	result := translate.RecipeModelToFind(&recipe)
	return result, nil
}

func (s *recipeService) FindRecipes(
	ctx context.Context,
	filter *dto.RecipeFilter,
) ([]dto.FindRecipe, error) {
	recipes, err := s.repo.FindMany(ctx, models.RecipeFilter{
		PageFilter: models.PageFilter{
			Limit:  filter.Limit,
			Offset: filter.Offset,
		},
		Name: filter.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("recipes not found: %w", err)
	}
	result := translate.RecipeModelsToFinds(recipes)
	return result, nil
}

func (s *recipeService) UpdateRecipe(
	ctx context.Context,
	id int,
	payload *dto.UpdateRecipe,
) error {
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
