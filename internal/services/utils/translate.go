package utils

import (
	"github.com/google/uuid"
	"github.com/romankravchuk/muerta/internal/api/router/params"
	"github.com/romankravchuk/muerta/internal/pkg/auth"
	"github.com/romankravchuk/muerta/internal/storage/postgres/models"
)

func RecipeModelToFind(model *models.Recipe) params.FindRecipe {
	steps := make([]params.FindStep, len(model.Steps))
	for i, step := range model.Steps {
		steps[i] = params.FindStep{
			ID:    step.ID,
			Name:  step.Name,
			Place: step.Place,
		}
	}
	return params.FindRecipe{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description,
		Steps:       steps,
	}
}

func RecipeModelsToFinds(models []models.Recipe) []params.FindRecipe {
	result := make([]params.FindRecipe, len(models))
	for i, recipe := range models {
		result[i] = params.FindRecipe{
			ID:          recipe.ID,
			Name:        recipe.Name,
			Description: recipe.Description,
		}
	}
	return result
}

func UserModelToFind(model *models.User) params.FindUser {
	settings := make([]params.FindSetting, len(model.Settings))
	for i, setting := range model.Settings {
		settings[i] = params.FindSetting{
			ID:       setting.ID,
			Name:     setting.Name,
			Value:    setting.Value,
			Category: setting.Category.Name,
		}
	}
	return params.FindUser{
		ID:        model.ID,
		Name:      model.Name,
		CreatedAt: model.CreatedAt,
		Settings:  settings,
	}
}

func UserModelsToFinds(models []models.User) []params.FindUser {
	result := make([]params.FindUser, len(models))
	for i, user := range models {
		result[i] = params.FindUser{
			ID:        user.ID,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
		}
	}
	return result
}

func CreateUserToModel(dto *params.CreateUser) models.User {
	settings := make([]models.Setting, len(dto.Settings))
	roles := make([]models.Role, len(dto.Roles))
	for i, setting := range dto.Settings {
		settings[i] = models.Setting{
			ID:    setting.ID,
			Value: setting.Value,
		}
	}
	for i, role := range dto.Roles {
		roles[i] = models.Role{ID: role.ID}
	}
	salt := uuid.New().String()
	return models.User{
		Name:     dto.Name,
		Salt:     salt,
		Settings: settings,
		Roles:    roles,
		Password: models.Password{Hash: auth.GenerateHashFromPassword(dto.Password, salt)},
	}
}

func CreateRecipeToModel(dto *params.CreateRecipe) models.Recipe {
	steps := make([]models.Step, len(dto.Steps))
	ingredients := make([]models.RecipeIngredient, len(dto.Ingredients))
	for i, step := range dto.Steps {
		steps[i].ID = step.ID
		steps[i].Place = step.Place
	}
	for i, ingredient := range dto.Ingredients {
		ingredients[i].Product.ID = ingredient.ProductID
		ingredients[i].Measure.ID = ingredient.MeasureID
		ingredients[i].Quantity = ingredient.Quantity
	}
	return models.Recipe{
		User: models.User{
			ID: dto.UserID,
		},
		Name:        dto.Name,
		Description: dto.Description,
		Steps:       steps,
		Ingredients: ingredients,
	}
}

func SettingModelToFind(model *models.Setting) params.FindSetting {
	return params.FindSetting{
		ID:       model.ID,
		Name:     model.Name,
		Category: model.Category.Name,
		Value:    model.Value,
	}
}

func SettingModelsToFinds(models []models.Setting) []params.FindSetting {
	dtos := make([]params.FindSetting, len(models))
	for i, model := range models {
		dtos[i] = SettingModelToFind(&model)
	}
	return dtos
}

func CreateSettingToModel(dto *params.CreateSetting) models.Setting {
	return models.Setting{
		Name: dto.Name,
		Category: models.SettingCategory{
			ID: dto.CategoryID,
		},
	}
}

func StorageModelToFind(model *models.Vault) params.FindStorage {
	return params.FindStorage{
		ID:          model.ID,
		Name:        model.Name,
		Temperature: model.Temperature,
		CreatedAt:   model.CreatedAt,
		Humidity:    model.Humidity,
		Type: params.FindStorageType{
			ID:   model.Type.ID,
			Name: model.Type.Name,
		},
	}
}

func StorageModelsToFinds(models []models.Vault) []params.FindStorage {
	dtos := make([]params.FindStorage, len(models))
	for i, model := range models {
		dtos[i] = StorageModelToFind(&model)
	}
	return dtos
}

func CreateStorageToModel(dto *params.CreateStorage) models.Vault {
	return models.Vault{
		Name:        dto.Name,
		Temperature: dto.Temperature,
		Humidity:    dto.Humidity,
		Type: models.StorageType{
			ID: dto.TypeID,
		},
	}
}

func ProductModelToFind(model *models.Product) params.FindProduct {
	return params.FindProduct{
		ID:   model.ID,
		Name: model.Name,
	}
}

func ProductModelsToFinds(models []models.Product) []params.FindProduct {
	dtos := make([]params.FindProduct, len(models))
	for i, model := range models {
		dtos[i] = ProductModelToFind(&model)
	}
	return dtos
}

func CreateProductToModel(dto *params.CreateProduct) models.Product {
	return models.Product{
		Name: dto.Name,
	}
}

func RoleModelToFindRole(model *models.Role) params.FindRole {
	return params.FindRole{
		ID:   model.ID,
		Name: model.Name,
	}
}

func RoleModelsToFindRoles(models []models.Role) []params.FindRole {
	roles := make([]params.FindRole, len(models))
	for i, model := range models {
		roles[i] = RoleModelToFindRole(&model)
	}
	return roles
}

func CreateRoleToModel(dto *params.CreateRole) models.Role {
	return models.Role{
		Name: dto.Name,
	}
}

func CategoryModelsToFinds(model []models.ProductCategory) []params.FindProductCategory {
	categories := make([]params.FindProductCategory, len(model))
	for i, category := range model {
		categories[i] = ProductCategoryModelToFind(&category)
	}
	return categories
}

func ProductCategoryModelToFind(model *models.ProductCategory) params.FindProductCategory {
	return params.FindProductCategory{
		ID:        model.ID,
		Name:      model.Name,
		CreatedAt: model.CreatedAt,
	}
}

func CreateCategoryToModel(dto *params.CreateProductCategory) models.ProductCategory {
	return models.ProductCategory{
		Name: dto.Name,
	}
}

func CreateTipToModel(dto *params.CreateTip) models.Tip {
	return models.Tip{
		Description: dto.Description,
	}
}

func TipModelToFind(model *models.Tip) params.FindTip {
	return params.FindTip{
		ID:          model.ID,
		Description: model.Description,
	}
}

func TipModelsToFinds(models []models.Tip) []params.FindTip {
	dtos := make([]params.FindTip, len(models))
	for i, model := range models {
		dtos[i] = TipModelToFind(&model)
	}
	return dtos
}

func CreateMeasureToModel(dto *params.CreateMeasure) models.Measure {
	return models.Measure{
		Name: dto.Name,
	}
}

func MeasureModelToFind(model *models.Measure) params.FindMeasure {
	return params.FindMeasure{
		ID:   model.ID,
		Name: model.Name,
	}
}

func MeasureModelsToFinds(models []models.Measure) []params.FindMeasure {
	dtos := make([]params.FindMeasure, len(models))
	for i, model := range models {
		dtos[i] = MeasureModelToFind(&model)
	}
	return dtos
}

func CreateStorageTypeToModel(dto *params.CreateStorageType) models.StorageType {
	return models.StorageType{
		Name: dto.Name,
	}
}

func StorageTypeModelToFind(model *models.StorageType) params.FindStorageType {
	return params.FindStorageType{
		ID:   model.ID,
		Name: model.Name,
	}
}

func StorageTypeModelsToFinds(models []models.StorageType) []params.FindStorageType {
	dtos := make([]params.FindStorageType, len(models))
	for i, model := range models {
		dtos[i] = StorageTypeModelToFind(&model)
	}
	return dtos
}

func CreateShelfLifeToModel(dto *params.CreateShelfLife) models.ShelfLife {
	return models.ShelfLife{
		Product: models.Product{
			ID: dto.ProductID,
		},
		Storage: models.Vault{
			ID: dto.StorageID,
		},
		Measure: models.Measure{
			ID: dto.MeasureID,
		},
		User: models.User{
			ID: dto.UserID,
		},
		Quantity:     dto.Quantity,
		PurchaseDate: dto.PurchaseDate,
		EndDate:      dto.EndDate,
	}
}

func ShelfLifeModelToFind(model *models.ShelfLife) params.FindShelfLife {
	return params.FindShelfLife{
		ID: model.ID,
		Product: params.FindProduct{
			ID:   model.Product.ID,
			Name: model.Product.Name,
		},
		Storage: params.FindStorage{
			ID:          model.Storage.ID,
			Name:        model.Storage.Name,
			Temperature: model.Storage.Temperature,
			Humidity:    model.Storage.Humidity,
			CreatedAt:   model.Storage.CreatedAt,
		},
		Measure: params.FindMeasure{
			ID:   model.Measure.ID,
			Name: model.Measure.Name,
		},
		Quantity:     model.Quantity,
		PurchaseDate: model.PurchaseDate,
		EndDate:      model.EndDate,
	}
}

func ShelfLifeModelsToFinds(models []models.ShelfLife) []params.FindShelfLife {
	dtos := make([]params.FindShelfLife, len(models))
	for i, model := range models {
		dtos[i] = ShelfLifeModelToFind(&model)
	}
	return dtos
}

func CreateShelfLifeStatusToModel(dto *params.CreateShelfLifeStatus) models.ShelfLifeStatus {
	return models.ShelfLifeStatus{
		Name: dto.Name,
	}
}

func ShelfLifeStatusModelToFind(model *models.ShelfLifeStatus) params.FindShelfLifeStatus {
	return params.FindShelfLifeStatus{
		ID:   model.ID,
		Name: model.Name,
	}
}

func ShelfLifeStatusModelsToFinds(models []models.ShelfLifeStatus) []params.FindShelfLifeStatus {
	dtos := make([]params.FindShelfLifeStatus, len(models))
	for i, model := range models {
		dtos[i] = ShelfLifeStatusModelToFind(&model)
	}
	return dtos
}

func SignUpToModel(payload *params.SignUp) models.User {
	return models.User{
		Name: payload.Name,
	}
}

func CreateIngredientToModel(payload *params.CreateIngredient) models.RecipeIngredient {
	return models.RecipeIngredient{
		Product:  models.Product{ID: payload.ProductID},
		Measure:  models.Measure{ID: payload.MeasureID},
		Quantity: payload.Quantity,
	}
}

func UpdateIngredientToModel(payload *params.UpdateIngredient) models.RecipeIngredient {
	return models.RecipeIngredient{
		Product:  models.Product{ID: payload.ProductID},
		Measure:  models.Measure{ID: payload.MeasureID},
		Quantity: payload.Quantity,
	}
}

func RecipeIngredientModelToFind(entity *models.RecipeIngredient) params.FindRecipeIngredient {
	return params.FindRecipeIngredient{
		Product: params.FindProduct{
			ID:   entity.Product.ID,
			Name: entity.Product.Name,
		},
		Measure: params.FindMeasure{
			ID:   entity.Measure.ID,
			Name: entity.Measure.Name,
		},
		Quantity: entity.Quantity,
	}
}

func RecipeIngredientModelsToFinds(
	entities []models.RecipeIngredient,
) []params.FindRecipeIngredient {
	payloads := make([]params.FindRecipeIngredient, len(entities))
	for i, entity := range entities {
		payloads[i] = RecipeIngredientModelToFind(&entity)
	}
	return payloads
}

func UpdateSettingToModel(payload *params.UpdateUserSetting) models.Setting {
	return models.Setting{
		Value: payload.Value,
	}
}

func UserStorageToModel(payload *params.UserStorage) models.Vault {
	return models.Vault{
		ID: payload.StorageID,
	}
}

func StepModelToFind(model models.Step) params.FindStep {
	return params.FindStep{
		ID:    model.ID,
		Name:  model.Name,
		Place: model.Place,
	}
}

func StepModelsToFinds(models []models.Step) []params.FindStep {
	payloads := make([]params.FindStep, len(models))
	for i, model := range models {
		payloads[i] = StepModelToFind(model)
	}
	return payloads
}

func CreateToStepModel(dto *params.CreateStep) models.Step {
	return models.Step{
		Name: dto.Name,
	}
}
