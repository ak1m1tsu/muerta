package user

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/api/router/params"
	"github.com/romankravchuk/muerta/internal/services/utils"
	"github.com/romankravchuk/muerta/internal/storage/postgres/models"
	repo "github.com/romankravchuk/muerta/internal/storage/postgres/user"
)

type UserServicer interface {
	FindUserByID(ctx context.Context, id int) (params.FindUser, error)
	FindUsers(ctx context.Context, filter *params.UserFilter) ([]params.FindUser, error)
	CreateUser(ctx context.Context, payload *params.CreateUser) error
	UpdateUser(ctx context.Context, id int, user *params.UpdateUser) error
	DeleteUser(ctx context.Context, id int) error
	RestoreUser(ctx context.Context, id int) error
	FindSettings(ctx context.Context, id int) ([]params.FindSetting, error)
	UpdateSetting(
		ctx context.Context,
		id int,
		payload *params.UpdateUserSetting,
	) (params.FindSetting, error)
	FindRoles(ctx context.Context, id int) ([]params.FindRole, error)
	AddStorage(
		ctx context.Context,
		id, storageID int,
	) (params.FindStorage, error)
	RemoveStorage(ctx context.Context, id, storageID int) error
	FindStorages(ctx context.Context, id int) ([]params.FindStorage, error)
	FindShelfLives(ctx context.Context, id int) ([]params.FindShelfLife, error)
	CreateShelfLife(
		ctx context.Context,
		id int,
		payload *params.CreateShelfLife,
	) (params.FindShelfLife, error)
	UpdateShelfLife(
		ctx context.Context,
		id int,
		payload *params.UserShelfLife,
	) (params.FindShelfLife, error)
	RestoreShelfLife(
		ctx context.Context,
		id, shelfLifeID int,
	) (params.FindShelfLife, error)
	DeleteShelfLife(ctx context.Context, id, shelfLifeID int) error
	Count(ctx context.Context, filter params.UserFilter) (int, error)
}

type userService struct {
	repo repo.UserStorage
}

// Count implements UserServicer
func (s *userService) Count(ctx context.Context, filter params.UserFilter) (int, error) {
	count, err := s.repo.Count(ctx, models.UserFilter{Name: filter.Name})
	if err != nil {
		return 0, fmt.Errorf("error counting users: %w", err)
	}
	return count, nil
}

func New(repo repo.UserStorage) UserServicer {
	return &userService{
		repo: repo,
	}
}

// CreateShelfLife implements UserServicer
func (svc *userService) CreateShelfLife(
	ctx context.Context,
	id int,
	payload *params.CreateShelfLife,
) (params.FindShelfLife, error) {
	model := utils.CreateShelfLifeToModel(payload)
	createdModel, err := svc.repo.CreateShelfLife(ctx, id, model)
	if err != nil {
		return params.FindShelfLife{}, fmt.Errorf("error creating shelf life: %w", err)
	}
	return utils.ShelfLifeModelToFind(&createdModel), nil
}

// DeleteShelfLife implements UserServicer
func (svc *userService) DeleteShelfLife(
	ctx context.Context,
	id, shelfLifeID int,
) error {
	if err := svc.repo.DeleteShelfLife(ctx, id, shelfLifeID); err != nil {
		return fmt.Errorf("error deleting shelf life: %w", err)
	}
	return nil
}

// FindShelfLives implements UserServicer
func (svc *userService) FindShelfLives(
	ctx context.Context,
	id int,
) ([]params.FindShelfLife, error) {
	models, err := svc.repo.FindShelfLives(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error finding shelf lives: %w", err)
	}
	return utils.ShelfLifeModelsToFinds(models), nil
}

// RestoreShelfLife implements UserServicer
func (svc *userService) RestoreShelfLife(
	ctx context.Context,
	id, shelfLifeID int,
) (params.FindShelfLife, error) {
	model, err := svc.repo.RestoreShelfLife(ctx, id, shelfLifeID)
	if err != nil {
		return params.FindShelfLife{}, fmt.Errorf("error restoring shelf life: %w", err)
	}
	return utils.ShelfLifeModelToFind(&model), nil
}

// UpdateShelfLife implements UserServicer
func (svc *userService) UpdateShelfLife(
	ctx context.Context,
	id int,
	payload *params.UserShelfLife,
) (params.FindShelfLife, error) {
	model, err := svc.repo.FindShelfLife(ctx, id, payload.ShelfLifeID)
	if err != nil {
		return params.FindShelfLife{}, fmt.Errorf("error finding shelf life: %w", err)
	}
	if payload.MeasureID != 0 {
		model.Measure.ID = payload.MeasureID
	}
	if payload.ProductID != 0 {
		model.Product.ID = payload.ProductID
	}
	if payload.StorageID != 0 {
		model.Storage.ID = payload.StorageID
	}
	if payload.Quantity != 0 {
		model.Quantity = payload.Quantity
	}
	if payload.PurchaseDate != nil {
		model.PurchaseDate = payload.PurchaseDate
	}
	if payload.EndDate != nil {
		model.EndDate = payload.EndDate
	}
	result, err := svc.repo.UpdateShelfLife(ctx, id, model)
	if err != nil {
		return params.FindShelfLife{}, fmt.Errorf("error updating shelf life: %w", err)
	}
	return utils.ShelfLifeModelToFind(&result), nil
}

// RemoveStorage implements UserServicer
func (svc *userService) RemoveStorage(
	ctx context.Context,
	id, storageID int,
) error {
	err := svc.repo.RemoveVault(ctx, id, storageID)
	if err != nil {
		return fmt.Errorf("error creating storage: %w", err)
	}
	return nil
}

// AddStorage implements UserServicer
func (svc *userService) AddStorage(
	ctx context.Context,
	id, storageID int,
) (params.FindStorage, error) {
	model, err := svc.repo.AddVault(ctx, id, storageID)
	if err != nil {
		return params.FindStorage{}, fmt.Errorf("error creating storage: %w", err)
	}
	return utils.StorageModelToFind(&model), nil
}

// FindStorages implements UserServicer
func (svc *userService) FindStorages(ctx context.Context, id int) ([]params.FindStorage, error) {
	result, err := svc.repo.FindVaults(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error finding storages: %w", err)
	}
	return utils.StorageModelsToFinds(result), nil
}

// FindRoles implements UserServicer
func (svc *userService) FindRoles(ctx context.Context, id int) ([]params.FindRole, error) {
	if _, err := svc.repo.FindByID(ctx, id); err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}
	entities, err := svc.repo.FindRoles(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error finding roles: %w", err)
	}
	result := utils.RoleModelsToFindRoles(entities)
	return result, nil
}

// FindSettings implements UserServicer
func (svc *userService) FindSettings(ctx context.Context, id int) ([]params.FindSetting, error) {
	if _, err := svc.repo.FindByID(ctx, id); err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}
	entities, err := svc.repo.FindSettings(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error finding settings: %w", err)
	}
	result := utils.SettingModelsToFinds(entities)
	return result, nil
}

// UpdateSetting implements UserServicer
func (svc *userService) UpdateSetting(
	ctx context.Context,
	id int,
	payload *params.UpdateUserSetting,
) (params.FindSetting, error) {
	if _, err := svc.repo.FindByID(ctx, id); err != nil {
		return params.FindSetting{}, fmt.Errorf("error finding user: %w", err)
	}
	entity := utils.UpdateSettingToModel(payload)
	result, err := svc.repo.UpdateSetting(ctx, id, entity)
	if err != nil {
		return params.FindSetting{}, fmt.Errorf("error updating setting: %w", err)
	}
	dto := utils.SettingModelToFind(&result)
	return dto, nil
}

func (svc *userService) FindUserByID(ctx context.Context, id int) (params.FindUser, error) {
	user, err := svc.repo.FindByID(ctx, id)
	result := utils.UserModelToFind(&user)
	if err != nil {
		return params.FindUser{}, err
	}
	return result, nil
}

func (svc *userService) FindUsers(
	ctx context.Context,
	filter *params.UserFilter,
) ([]params.FindUser, error) {
	users, err := svc.repo.FindMany(ctx, models.UserFilter{
		PageFilter: models.PageFilter{
			Limit:  filter.Limit,
			Offset: filter.Offset,
		},
		Name: filter.Name,
	})
	dtos := utils.UserModelsToFinds(users)
	if err != nil {
		return nil, err
	}
	return dtos, nil
}

func (svc *userService) CreateUser(ctx context.Context, payload *params.CreateUser) error {
	model := utils.CreateUserToModel(payload)
	if err := svc.repo.Create(ctx, model); err != nil {
		return err
	}
	return nil
}

func (svc *userService) UpdateUser(ctx context.Context, id int, user *params.UpdateUser) error {
	oldUser, err := svc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if user.Name != "" {
		oldUser.Name = user.Name
	}
	if err := svc.repo.Update(ctx, oldUser); err != nil {
		return err
	}
	return nil
}

func (svc *userService) DeleteUser(ctx context.Context, id int) error {
	if err := svc.repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

func (svc *userService) RestoreUser(ctx context.Context, id int) error {
	if err := svc.repo.Restore(ctx, id); err != nil {
		return err
	}
	return nil
}
