package translate

import (
	"github.com/google/uuid"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/pkg/auth"
	"github.com/romankravchuk/muerta/internal/repositories/models"
)

func RecipeModelToFindDTO(model *models.Recipe) dto.FindRecipeDTO {
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

func RecipeModelsToFindDTOs(models []models.Recipe) []dto.FindRecipeDTO {
	result := make([]dto.FindRecipeDTO, len(models))
	for i, recipe := range models {
		result[i] = dto.FindRecipeDTO{
			ID:          recipe.ID,
			Name:        recipe.Name,
			Description: recipe.Description,
		}
	}
	return result
}

func UserModelToFindDTO(model *models.User) dto.FindUserDTO {
	settings := make([]dto.FindSettingDTO, len(model.Settings))
	for i, setting := range model.Settings {
		settings[i] = dto.FindSettingDTO{
			ID:       setting.ID,
			Name:     setting.Name,
			Value:    setting.Value,
			Category: setting.Category.Name,
		}
	}
	return dto.FindUserDTO{
		ID:        model.ID,
		Name:      model.Name,
		CreatedAt: model.CreatedAt,
		Settings:  settings,
	}
}

func UserModelsToFindDTOs(models []models.User) []dto.FindUserDTO {
	result := make([]dto.FindUserDTO, len(models))
	for i, user := range models {
		result[i] = dto.FindUserDTO{
			ID:        user.ID,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
		}
	}
	return result
}

func CreateUserDTOToModel(dto *dto.CreateUserDTO) models.User {
	settings := make([]models.Setting, len(dto.Settings))
	for i, setting := range dto.Settings {
		settings[i] = models.Setting{
			ID:    setting.ID,
			Value: setting.Value,
		}
	}
	salt := uuid.New().String()
	return models.User{
		ID:       dto.ID,
		Name:     dto.Name,
		Salt:     salt,
		Settings: settings,
		Password: models.Password{
			Hash: auth.GenerateHashFromPassword(dto.Password, salt),
		},
	}
}

func CreateRecipeDTOToModel(dto *dto.CreateRecipeDTO) models.Recipe {
	steps := make([]models.Step, len(dto.Steps))
	for i, step := range dto.Steps {
		steps[i].ID = step.ID
		steps[i].Place = step.Place
	}
	return models.Recipe{
		Name:        dto.Name,
		Description: dto.Description,
		Steps:       steps,
	}
}

func SettingModelToFindDTO(model *models.Setting) dto.FindSettingDTO {
	return dto.FindSettingDTO{
		ID:       model.ID,
		Name:     model.Name,
		Category: model.Category.Name,
	}
}

func SettingModelsToFindDTOs(models []models.Setting) []dto.FindSettingDTO {
	dtos := make([]dto.FindSettingDTO, len(models))
	for i, model := range models {
		dtos[i] = SettingModelToFindDTO(&model)
	}
	return dtos
}

func CreateSettingDTOToModel(dto *dto.CreateSettingDTO) models.Setting {
	return models.Setting{
		Name: dto.Name,
		Category: models.Category{
			ID: dto.CategoryID,
		},
	}
}

func StorageModelToFindDTO(model *models.Storage) dto.FindStorageDTO {
	return dto.FindStorageDTO{
		ID:          model.ID,
		Name:        model.Name,
		Temperature: model.Temperature,
		CreatedAt:   model.CreatedAt,
		Humidity:    model.Humidity,
		TypeName:    model.Type.Name,
	}
}

func StorageModelsToFindDTOs(models []models.Storage) []dto.FindStorageDTO {
	dtos := make([]dto.FindStorageDTO, len(models))
	for i, model := range models {
		dtos[i] = StorageModelToFindDTO(&model)
	}
	return dtos
}

func CreateStorageDTOToModel(dto *dto.CreateStorageDTO) models.Storage {
	return models.Storage{
		Name:        dto.Name,
		Temperature: dto.Temperature,
		Humidity:    dto.Humidity,
		Type: models.Type{
			ID: dto.TypeID,
		},
	}
}
