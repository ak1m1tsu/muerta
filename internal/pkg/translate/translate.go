package translate

import (
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/repositories/models"
)

func RecipeModelToFindDTO(model models.Recipe) dto.FindRecipeDTO {
	steps := make([]dto.FindStepDTO, len(model.Steps))
	for i, step := range model.Steps {
		steps[i] = dto.FindStepDTO{
			ID:    step.ID,
			Name:  step.Name,
			Place: step.Place,
		}
	}
	return dto.FindRecipeDTO{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description,
		Steps:       steps,
	}
}
