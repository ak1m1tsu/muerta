package recipe

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/pkg/translate"
	recipes "github.com/romankravchuk/muerta/internal/repositories/recipe"
)

type RecipeServicer interface {
	CreateRecipe(ctx context.Context, payload *dto.CreateRecipeDTO) error
	FindRecipeByID(ctx context.Context, id int) (dto.FindRecipeDTO, error)
	FindRecipes(ctx context.Context, filter *dto.RecipeFilterDTO) ([]dto.FindRecipeDTO, error)
	UpdateRecipe(ctx context.Context, id int, payload *dto.UpdateRecipeDTO) error
	DeleteRecipe(ctx context.Context, id int) error
	RestoreRecipe(ctx context.Context, id int) error
	FindRecipeIngredients(ctx context.Context, id int) ([]dto.FindRecipeIngredientDTO, error)
	CreateRecipeIngredient(ctx context.Context, id int, payload *dto.CreateRecipeIngredientDTO) (dto.FindRecipeIngredientDTO, error)
	UpdateRecipeIngredient(ctx context.Context, id int, payload *dto.UpdateRecipeIngredientDTO) (dto.FindRecipeIngredientDTO, error)
	DeleteRecipeIngredient(ctx context.Context, id int, payload *dto.DeleteRecipeIngredientDTO) error
}

type RecipeService struct {
	repository recipes.RecipesRepositorer
}

// CreateRecipeIngredient implements RecipeServicer
func (s *RecipeService) CreateRecipeIngredient(ctx context.Context, id int, payload *dto.CreateRecipeIngredientDTO) (dto.FindRecipeIngredientDTO, error) {
	model := translate.CreateRecipeIngredientDTOToModel(payload)
	ingredient, err := s.repository.CreateRecipeIngredient(ctx, id, &model)
	if err != nil {
		return dto.FindRecipeIngredientDTO{}, fmt.Errorf("create recipe ingredient: %w", err)
	}
	dto := translate.RecipeIngredientModelToFindDTO(&ingredient)
	return dto, nil
}

// DeleteRecipeIngredient implements RecipeServicer
func (s *RecipeService) DeleteRecipeIngredient(ctx context.Context, id int, payload *dto.DeleteRecipeIngredientDTO) error {
	if err := s.repository.DeleteRecipeIngredient(ctx, id, payload.ProductID); err != nil {
		return fmt.Errorf("delete recipe ingredient: %w", err)
	}
	return nil
}

// FindRecipeIngredients implements RecipeServicer
func (s *RecipeService) FindRecipeIngredients(ctx context.Context, id int) ([]dto.FindRecipeIngredientDTO, error) {
	ingredients, err := s.repository.FindIngredients(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ingredients not found: %w", err)
	}
	result := translate.RecipeIngredientModelsToFindDTOs(ingredients)
	return result, nil
}

// UpdateRecipeIngredient implements RecipeServicer
func (s *RecipeService) UpdateRecipeIngredient(ctx context.Context, id int, payload *dto.UpdateRecipeIngredientDTO) (dto.FindRecipeIngredientDTO, error) {
	model := translate.UpdateRecipeIngredientDTOToModel(payload)
	ingredient, err := s.repository.UpdateRecipeIngredient(ctx, id, &model)
	if err != nil {
		return dto.FindRecipeIngredientDTO{}, fmt.Errorf("update recipe ingredient: %w", err)
	}
	dto := translate.RecipeIngredientModelToFindDTO(&ingredient)
	return dto, nil
}

func New(repository recipes.RecipesRepositorer) RecipeServicer {
	return &RecipeService{repository: repository}
}

func (s *RecipeService) CreateRecipe(ctx context.Context, payload *dto.CreateRecipeDTO) error {
	model := translate.CreateRecipeDTOToModel(payload)
	if err := s.repository.Create(ctx, &model); err != nil {
		return fmt.Errorf("create recipe: %w", err)
	}
	return nil
}

func (s *RecipeService) FindRecipeByID(ctx context.Context, id int) (dto.FindRecipeDTO, error) {
	recipe, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return dto.FindRecipeDTO{}, fmt.Errorf("recipe not found by id: %w", err)
	}
	result := translate.RecipeModelToFindDTO(&recipe)
	return result, nil
}

func (s *RecipeService) FindRecipes(ctx context.Context, filter *dto.RecipeFilterDTO) ([]dto.FindRecipeDTO, error) {
	recipes, err := s.repository.FindMany(ctx, filter.Limit, filter.Offset, filter.Name)
	if err != nil {
		return nil, fmt.Errorf("recipes not found: %w", err)
	}
	result := translate.RecipeModelsToFindDTOs(recipes)
	return result, nil
}

func (s *RecipeService) UpdateRecipe(ctx context.Context, id int, payload *dto.UpdateRecipeDTO) error {
	recipe, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("recipe not found by id: %w", err)
	}
	if payload.Name != "" {
		recipe.Name = payload.Name
	}
	if payload.Description != "" {
		recipe.Description = payload.Description
	}
	if err := s.repository.Update(ctx, &recipe); err != nil {
		return fmt.Errorf("update recipe: %w", err)
	}
	return nil
}

func (s *RecipeService) DeleteRecipe(ctx context.Context, id int) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete recipe: %w", err)
	}
	return nil
}

func (s *RecipeService) RestoreRecipe(ctx context.Context, id int) error {
	if err := s.repository.Restore(ctx, id); err != nil {
		return fmt.Errorf("restore recipe: %w", err)
	}
	return nil
}
