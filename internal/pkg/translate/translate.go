package translate

import (
	"github.com/google/uuid"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/pkg/auth"
	"github.com/romankravchuk/muerta/internal/storage/postgres/models"
)

func RecipeModelToFind(model *models.Recipe) dto.FindRecipe {
	steps := make([]dto.FindStep, len(model.Steps))
	for i, step := range model.Steps {
		steps[i] = dto.FindStep{
			ID:    step.ID,
			Name:  step.Name,
			Place: step.Place,
		}
	}
	return dto.FindRecipe{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description,
		Steps:       steps,
	}
}

func RecipeModelsToFinds(models []models.Recipe) []dto.FindRecipe {
	result := make([]dto.FindRecipe, len(models))
	for i, recipe := range models {
		result[i] = dto.FindRecipe{
			ID:          recipe.ID,
			Name:        recipe.Name,
			Description: recipe.Description,
		}
	}
	return result
}

func UserModelToFind(model *models.User) dto.FindUser {
	settings := make([]dto.FindSetting, len(model.Settings))
	for i, setting := range model.Settings {
		settings[i] = dto.FindSetting{
			ID:       setting.ID,
			Name:     setting.Name,
			Value:    setting.Value,
			Category: setting.Category.Name,
		}
	}
	return dto.FindUser{
		ID:        model.ID,
		Name:      model.Name,
		CreatedAt: model.CreatedAt,
		Settings:  settings,
	}
}

func UserModelsToFinds(models []models.User) []dto.FindUser {
	result := make([]dto.FindUser, len(models))
	for i, user := range models {
		result[i] = dto.FindUser{
			ID:        user.ID,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
		}
	}
	return result
}

func CreateUserToModel(dto *dto.CreateUser) models.User {
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

func CreateRecipeToModel(dto *dto.CreateRecipe) models.Recipe {
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

func SettingModelToFind(model *models.Setting) dto.FindSetting {
	return dto.FindSetting{
		ID:       model.ID,
		Name:     model.Name,
		Category: model.Category.Name,
		Value:    model.Value,
	}
}

func SettingModelsToFinds(models []models.Setting) []dto.FindSetting {
	dtos := make([]dto.FindSetting, len(models))
	for i, model := range models {
		dtos[i] = SettingModelToFind(&model)
	}
	return dtos
}

func CreateSettingToModel(dto *dto.CreateSetting) models.Setting {
	return models.Setting{
		Name: dto.Name,
		Category: models.SettingCategory{
			ID: dto.CategoryID,
		},
	}
}

func StorageModelToFind(model *models.Vault) dto.FindStorage {
	return dto.FindStorage{
		ID:          model.ID,
		Name:        model.Name,
		Temperature: model.Temperature,
		CreatedAt:   model.CreatedAt,
		Humidity:    model.Humidity,
		Type: dto.FindStorageType{
			ID:   model.Type.ID,
			Name: model.Type.Name,
		},
	}
}

func StorageModelsToFinds(models []models.Vault) []dto.FindStorage {
	dtos := make([]dto.FindStorage, len(models))
	for i, model := range models {
		dtos[i] = StorageModelToFind(&model)
	}
	return dtos
}

func CreateStorageToModel(dto *dto.CreateStorage) models.Vault {
	return models.Vault{
		Name:        dto.Name,
		Temperature: dto.Temperature,
		Humidity:    dto.Humidity,
		Type: models.StorageType{
			ID: dto.TypeID,
		},
	}
}

func ProductModelToFind(model *models.Product) dto.FindProduct {
	return dto.FindProduct{
		ID:   model.ID,
		Name: model.Name,
	}
}

func ProductModelsToFinds(models []models.Product) []dto.FindProduct {
	dtos := make([]dto.FindProduct, len(models))
	for i, model := range models {
		dtos[i] = ProductModelToFind(&model)
	}
	return dtos
}

func CreateProductToModel(dto *dto.CreateProduct) models.Product {
	return models.Product{
		Name: dto.Name,
	}
}

func RoleModelToFindRole(model *models.Role) dto.FindRole {
	return dto.FindRole{
		ID:   model.ID,
		Name: model.Name,
	}
}

func RoleModelsToFindRoles(models []models.Role) []dto.FindRole {
	roles := make([]dto.FindRole, len(models))
	for i, model := range models {
		roles[i] = RoleModelToFindRole(&model)
	}
	return roles
}

func CreateRoleToModel(dto *dto.CreateRole) models.Role {
	return models.Role{
		Name: dto.Name,
	}
}

func CategoryModelsToFinds(model []models.ProductCategory) []dto.FindProductCategory {
	categories := make([]dto.FindProductCategory, len(model))
	for i, category := range model {
		categories[i] = ProductCategoryModelToFind(&category)
	}
	return categories
}

func ProductCategoryModelToFind(model *models.ProductCategory) dto.FindProductCategory {
	return dto.FindProductCategory{
		ID:        model.ID,
		Name:      model.Name,
		CreatedAt: model.CreatedAt,
	}
}

func CreateCategoryToModel(dto *dto.CreateProductCategory) models.ProductCategory {
	return models.ProductCategory{
		Name: dto.Name,
	}
}

func CreateTipToModel(dto *dto.CreateTip) models.Tip {
	return models.Tip{
		Description: dto.Description,
	}
}

func TipModelToFind(model *models.Tip) dto.FindTip {
	return dto.FindTip{
		ID:          model.ID,
		Description: model.Description,
	}
}

func TipModelsToFinds(models []models.Tip) []dto.FindTip {
	dtos := make([]dto.FindTip, len(models))
	for i, model := range models {
		dtos[i] = TipModelToFind(&model)
	}
	return dtos
}

func CreateMeasureToModel(dto *dto.CreateMeasure) models.Measure {
	return models.Measure{
		Name: dto.Name,
	}
}

func MeasureModelToFind(model *models.Measure) dto.FindMeasure {
	return dto.FindMeasure{
		ID:   model.ID,
		Name: model.Name,
	}
}

func MeasureModelsToFinds(models []models.Measure) []dto.FindMeasure {
	dtos := make([]dto.FindMeasure, len(models))
	for i, model := range models {
		dtos[i] = MeasureModelToFind(&model)
	}
	return dtos
}

func CreateStorageTypeToModel(dto *dto.CreateStorageType) models.StorageType {
	return models.StorageType{
		Name: dto.Name,
	}
}

func StorageTypeModelToFind(model *models.StorageType) dto.FindStorageType {
	return dto.FindStorageType{
		ID:   model.ID,
		Name: model.Name,
	}
}

func StorageTypeModelsToFinds(models []models.StorageType) []dto.FindStorageType {
	dtos := make([]dto.FindStorageType, len(models))
	for i, model := range models {
		dtos[i] = StorageTypeModelToFind(&model)
	}
	return dtos
}

func CreateShelfLifeToModel(dto *dto.CreateShelfLife) models.ShelfLife {
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

func ShelfLifeModelToFind(model *models.ShelfLife) dto.FindShelfLife {
	return dto.FindShelfLife{
		ID: model.ID,
		Product: dto.FindProduct{
			ID:   model.Product.ID,
			Name: model.Product.Name,
		},
		Storage: dto.FindStorage{
			ID:          model.Storage.ID,
			Name:        model.Storage.Name,
			Temperature: model.Storage.Temperature,
			Humidity:    model.Storage.Humidity,
			CreatedAt:   model.Storage.CreatedAt,
		},
		Measure: dto.FindMeasure{
			ID:   model.Measure.ID,
			Name: model.Measure.Name,
		},
		Quantity:     model.Quantity,
		PurchaseDate: model.PurchaseDate,
		EndDate:      model.EndDate,
	}
}

func ShelfLifeModelsToFinds(models []models.ShelfLife) []dto.FindShelfLife {
	dtos := make([]dto.FindShelfLife, len(models))
	for i, model := range models {
		dtos[i] = ShelfLifeModelToFind(&model)
	}
	return dtos
}

func CreateShelfLifeStatusToModel(dto *dto.CreateShelfLifeStatus) models.ShelfLifeStatus {
	return models.ShelfLifeStatus{
		Name: dto.Name,
	}
}

func ShelfLifeStatusModelToFind(model *models.ShelfLifeStatus) dto.FindShelfLifeStatus {
	return dto.FindShelfLifeStatus{
		ID:   model.ID,
		Name: model.Name,
	}
}

func ShelfLifeStatusModelsToFinds(models []models.ShelfLifeStatus) []dto.FindShelfLifeStatus {
	dtos := make([]dto.FindShelfLifeStatus, len(models))
	for i, model := range models {
		dtos[i] = ShelfLifeStatusModelToFind(&model)
	}
	return dtos
}

func SignUpToModel(payload *dto.SignUp) models.User {
	return models.User{
		Name: payload.Name,
	}
}

func CreateIngredientToModel(payload *dto.CreateIngredient) models.RecipeIngredient {
	return models.RecipeIngredient{
		Product:  models.Product{ID: payload.ProductID},
		Measure:  models.Measure{ID: payload.MeasureID},
		Quantity: payload.Quantity,
	}
}

func UpdateIngredientToModel(payload *dto.UpdateIngredient) models.RecipeIngredient {
	return models.RecipeIngredient{
		Product:  models.Product{ID: payload.ProductID},
		Measure:  models.Measure{ID: payload.MeasureID},
		Quantity: payload.Quantity,
	}
}

func RecipeIngredientModelToFind(entity *models.RecipeIngredient) dto.FindRecipeIngredient {
	return dto.FindRecipeIngredient{
		Product: dto.FindProduct{
			ID:   entity.Product.ID,
			Name: entity.Product.Name,
		},
		Measure: dto.FindMeasure{
			ID:   entity.Measure.ID,
			Name: entity.Measure.Name,
		},
		Quantity: entity.Quantity,
	}
}

func RecipeIngredientModelsToFinds(
	entities []models.RecipeIngredient,
) []dto.FindRecipeIngredient {
	payloads := make([]dto.FindRecipeIngredient, len(entities))
	for i, entity := range entities {
		payloads[i] = RecipeIngredientModelToFind(&entity)
	}
	return payloads
}

func UpdateSettingToModel(payload *dto.UpdateUserSetting) models.Setting {
	return models.Setting{
		Value: payload.Value,
	}
}

func UserStorageToModel(payload *dto.UserStorage) models.Vault {
	return models.Vault{
		ID: payload.StorageID,
	}
}

func StepModelToFind(model models.Step) dto.FindStep {
	return dto.FindStep{
		ID:    model.ID,
		Name:  model.Name,
		Place: model.Place,
	}
}

func StepModelsToFinds(models []models.Step) []dto.FindStep {
	payloads := make([]dto.FindStep, len(models))
	for i, model := range models {
		payloads[i] = StepModelToFind(model)
	}
	return payloads
}

func CreateToStepModel(dto *dto.CreateStep) models.Step {
	return models.Step{
		Name: dto.Name,
	}
}
