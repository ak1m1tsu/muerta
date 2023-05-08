package recipe

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/pkg/translate"
	"github.com/romankravchuk/muerta/internal/repositories/recipes"
)

type RecipeServicer interface {
	CreateRecipe(ctx context.Context, payload *dto.CreateRecipeDTO) error
	FindRecipeByID(ctx context.Context, id int) (dto.FindRecipeDTO, error)
	FindRecipes(ctx context.Context, filter *dto.RecipeFilterDTO) ([]dto.FindRecipeDTO, error)
	UpdateRecipe(ctx context.Context, id int, payload *dto.UpdateRecipeDTO) error
	DeleteRecipe(ctx context.Context, id int) error
	RestoreRecipe(ctx context.Context, id int) error
}

type RecipeService struct {
	repository recipes.RecipesRepositorer
}

func New(repository recipes.RecipesRepositorer) *RecipeService {
	return &RecipeService{repository: repository}
}

func (s *RecipeService) CreateRecipe(ctx context.Context, payload *dto.CreateRecipeDTO) error {
	model := translate.CreateRecipeDTOToModel(payload)
	err := s.repository.Create(ctx, &model)
	return err
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
