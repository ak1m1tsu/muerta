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
		Category: models.SettingCategory{
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
		Type: models.StorageType{
			ID: dto.TypeID,
		},
	}
}

func ProductModelToFindDTO(model *models.Product) dto.FindProductDTO {
	return dto.FindProductDTO{
		ID:   model.ID,
		Name: model.Name,
	}
}

func ProductModelsToFindDTOs(models []models.Product) []dto.FindProductDTO {
	dtos := make([]dto.FindProductDTO, len(models))
	for i, model := range models {
		dtos[i] = ProductModelToFindDTO(&model)
	}
	return dtos
}

func CreateProductDTOToModel(dto *dto.CreateProductDTO) models.Product {
	return models.Product{
		Name: dto.Name,
	}
}

func RoleModelToFindRoleDTO(model *models.Role) dto.FindRoleDTO {
	return dto.FindRoleDTO{
		ID:   model.ID,
		Name: model.Name,
	}
}

func RoleModelsToFindRoleDTOs(models []models.Role) []dto.FindRoleDTO {
	roles := make([]dto.FindRoleDTO, len(models))
	for i, model := range models {
		roles[i] = RoleModelToFindRoleDTO(&model)
	}
	return roles
}

func CreateRoleDTOToModel(dto *dto.CreateRoleDTO) models.Role {
	return models.Role{
		Name: dto.Name,
	}
}

func CategoryModelsToFindDTOs(model []models.ProductCategory) []dto.FindProductCategoryDTO {
	categories := make([]dto.FindProductCategoryDTO, len(model))
	for i, category := range model {
		categories[i] = CategoryModelToFindDTO(&category)
	}
	return categories
}

func CategoryModelToFindDTO(model *models.ProductCategory) dto.FindProductCategoryDTO {
	return dto.FindProductCategoryDTO{
		ID:        model.ID,
		Name:      model.Name,
		CreatedAt: model.CreatedAt,
	}
}

func CreateCategoryDTOToModel(dto *dto.CreateProductCategoryDTO) models.ProductCategory {
	return models.ProductCategory{
		Name: dto.Name,
	}
}

func CreateTipDTOToModel(dto *dto.CreateTipDTO) models.Tip {
	return models.Tip{
		Description: dto.Description,
	}
}

func TipModelToFindTipDTO(model *models.Tip) dto.FindTipDTO {
	return dto.FindTipDTO{
		ID:          model.ID,
		Description: model.Description,
	}
}

func TipModelsToFindTipDTOs(models []models.Tip) []dto.FindTipDTO {
	dtos := make([]dto.FindTipDTO, len(models))
	for i, model := range models {
		dtos[i] = TipModelToFindTipDTO(&model)
	}
	return dtos
}

func CreateMeasureDTOToModel(dto *dto.CreateMeasureDTO) models.Measure {
	return models.Measure{
		Name: dto.Name,
	}
}

func MeasureModelToFindDTO(model *models.Measure) dto.FindMeasureDTO {
	return dto.FindMeasureDTO{
		ID:   model.ID,
		Name: model.Name,
	}
}

func MeasureModelsToFindDTOs(models []models.Measure) []dto.FindMeasureDTO {
	dtos := make([]dto.FindMeasureDTO, len(models))
	for i, model := range models {
		dtos[i] = MeasureModelToFindDTO(&model)
	}
	return dtos
}

func CreateStorageTypeDTOToModel(dto *dto.CreateStorageTypeDTO) models.StorageType {
	return models.StorageType{
		Name: dto.Name,
	}
}
func StorageTypeModelToFindDTO(model *models.StorageType) dto.FindStorageTypeDTO {
	return dto.FindStorageTypeDTO{
		ID:   model.ID,
		Name: model.Name,
	}
}
func StorageTypeModelsToFindDTOs(models []models.StorageType) []dto.FindStorageTypeDTO {
	dtos := make([]dto.FindStorageTypeDTO, len(models))
	for i, model := range models {
		dtos[i] = StorageTypeModelToFindDTO(&model)
	}
	return dtos
}

func CreateShelfLifeDTOToModel(dto *dto.CreateShelfLifeDTO) models.ShelfLife {
	return models.ShelfLife{
		Product: models.Product{
			ID: dto.ProductID,
		},
		Storage: models.Storage{
			ID: dto.StorageID,
		},
		Measure: models.Measure{
			ID: dto.MeasureID,
		},
		Quantity:     dto.Quantity,
		PurchaseDate: dto.PurchaseDate,
		EndDate:      dto.EndDate,
	}
}
func ShelfLifeModelToFindDTO(model *models.ShelfLife) dto.FindShelfLifeDTO {
	return dto.FindShelfLifeDTO{
		ID: model.ID,
		Product: dto.FindProductDTO{
			ID:   model.Product.ID,
			Name: model.Product.Name,
		},
		Storage: dto.FindStorageDTO{
			ID:          model.Storage.ID,
			Name:        model.Storage.Name,
			Temperature: model.Storage.Temperature,
			Humidity:    model.Storage.Humidity,
			TypeName:    model.Storage.Type.Name,
			CreatedAt:   model.Storage.CreatedAt,
		},
		Measure: dto.FindMeasureDTO{
			ID:   model.Measure.ID,
			Name: model.Measure.Name,
		},
		Quantity:     model.Quantity,
		PurchaseDate: model.PurchaseDate,
		EndDate:      model.EndDate,
	}
}
func ShelfLifeModelsToFindDTOs(models []models.ShelfLife) []dto.FindShelfLifeDTO {
	dtos := make([]dto.FindShelfLifeDTO, len(models))
	for i, model := range models {
		dtos[i] = ShelfLifeModelToFindDTO(&model)
	}
	return dtos
}
func CreateShelfLifeStatusDTOToModel(dto *dto.CreateShelfLifeStatusDTO) models.ShelfLifeStatus {
	return models.ShelfLifeStatus{
		Name: dto.Name,
	}
}
func ShelfLifeStatusModelToFindDTO(model *models.ShelfLifeStatus) dto.FindShelfLifeStatusDTO {
	return dto.FindShelfLifeStatusDTO{
		ID:   model.ID,
		Name: model.Name,
	}
}
func ShelfLifeStatusModelsToFindDTOs(models []models.ShelfLifeStatus) []dto.FindShelfLifeStatusDTO {
	dtos := make([]dto.FindShelfLifeStatusDTO, len(models))
	for i, model := range models {
		dtos[i] = ShelfLifeStatusModelToFindDTO(&model)
	}
	return dtos
}
func SignUpDTOToModel(payload *dto.SignUpDTO) models.User {
	return models.User{
		Name: payload.Name,
	}
}
