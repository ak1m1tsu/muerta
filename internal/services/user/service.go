package user

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/pkg/translate"
	repo "github.com/romankravchuk/muerta/internal/repositories/user"
)

type UserServicer interface {
	FindUserByID(ctx context.Context, id int) (dto.FindUserDTO, error)
	FindUsers(ctx context.Context, filter *dto.UserFilterDTO) ([]dto.FindUserDTO, error)
	CreateUser(ctx context.Context, payload *dto.CreateUserDTO) error
	UpdateUser(ctx context.Context, id int, user *dto.UpdateUserDTO) error
	DeleteUser(ctx context.Context, id int) error
	RestoreUser(ctx context.Context, id int) error
	FindSettings(ctx context.Context, id int) ([]dto.FindSettingDTO, error)
	UpdateSetting(ctx context.Context, id int, payload *dto.UpdateUserSettingDTO) (dto.FindSettingDTO, error)
	FindRoles(ctx context.Context, id int) ([]dto.FindRoleDTO, error)
	CreateStorage(ctx context.Context, id int, payload *dto.UserStorageDTO) (dto.FindStorageDTO, error)
	DeleteStorage(ctx context.Context, id int, payload *dto.UserStorageDTO) error
	FindStorages(ctx context.Context, id int) ([]dto.FindStorageDTO, error)
	FindShelfLives(ctx context.Context, id int) ([]dto.FindShelfLifeDTO, error)
	CreateShelfLife(ctx context.Context, id int, payload *dto.CreateShelfLifeDTO) (dto.FindShelfLifeDTO, error)
	UpdateShelfLife(ctx context.Context, id int, payload *dto.UserShelfLifeDTO) (dto.FindShelfLifeDTO, error)
	RestoreShelfLife(ctx context.Context, id int, payload *dto.UserShelfLifeDTO) (dto.FindShelfLifeDTO, error)
	DeleteShelfLife(ctx context.Context, id int, payload *dto.UserShelfLifeDTO) error
}

type userService struct {
	repo repo.UserRepositorer
}

func New(repo repo.UserRepositorer) UserServicer {
	return &userService{
		repo: repo,
	}
}

// CreateShelfLife implements UserServicer
func (svc *userService) CreateShelfLife(ctx context.Context, id int, payload *dto.CreateShelfLifeDTO) (dto.FindShelfLifeDTO, error) {
	model := translate.CreateShelfLifeDTOToModel(payload)
	createdModel, err := svc.repo.CreateShelfLife(ctx, id, model)
	if err != nil {
		return dto.FindShelfLifeDTO{}, fmt.Errorf("error creating shelf life: %w", err)
	}
	return translate.ShelfLifeModelToFindDTO(&createdModel), nil
}

// DeleteShelfLife implements UserServicer
func (svc *userService) DeleteShelfLife(ctx context.Context, id int, payload *dto.UserShelfLifeDTO) error {
	if err := svc.repo.DeleteShelfLife(ctx, id, payload.ShelfLifeID); err != nil {
		return fmt.Errorf("error deleting shelf life: %w", err)
	}
	return nil
}

// FindShelfLives implements UserServicer
func (svc *userService) FindShelfLives(ctx context.Context, id int) ([]dto.FindShelfLifeDTO, error) {
	models, err := svc.repo.FindShelfLives(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error finding shelf lives: %w", err)
	}
	return translate.ShelfLifeModelsToFindDTOs(models), nil
}

// RestoreShelfLife implements UserServicer
func (svc *userService) RestoreShelfLife(ctx context.Context, id int, payload *dto.UserShelfLifeDTO) (dto.FindShelfLifeDTO, error) {
	model, err := svc.repo.RestoreShelfLife(ctx, id, payload.ShelfLifeID)
	if err != nil {
		return dto.FindShelfLifeDTO{}, fmt.Errorf("error restoring shelf life: %w", err)
	}
	return translate.ShelfLifeModelToFindDTO(&model), nil
}

// UpdateShelfLife implements UserServicer
func (svc *userService) UpdateShelfLife(ctx context.Context, id int, payload *dto.UserShelfLifeDTO) (dto.FindShelfLifeDTO, error) {
	model, err := svc.repo.FindShelfLife(ctx, id, payload.ShelfLifeID)
	if err != nil {
		return dto.FindShelfLifeDTO{}, fmt.Errorf("error finding shelf life: %w", err)
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
		return dto.FindShelfLifeDTO{}, fmt.Errorf("error updating shelf life: %w", err)
	}
	return translate.ShelfLifeModelToFindDTO(&result), nil
}

// DeleteStorage implements UserServicer
func (svc *userService) DeleteStorage(ctx context.Context, id int, payload *dto.UserStorageDTO) error {
	model := translate.UserStorageDTOToModel(payload)
	err := svc.repo.DeleteStorage(ctx, id, model)
	if err != nil {
		return fmt.Errorf("error creating storage: %w", err)
	}
	return nil
}

// CreateStorage implements UserServicer
func (svc *userService) CreateStorage(ctx context.Context, id int, payload *dto.UserStorageDTO) (dto.FindStorageDTO, error) {
	model := translate.UserStorageDTOToModel(payload)
	model, err := svc.repo.CreateStorage(ctx, id, model)
	if err != nil {
		return dto.FindStorageDTO{}, fmt.Errorf("error creating storage: %w", err)
	}
	return translate.StorageModelToFindDTO(&model), nil
}

// FindStorages implements UserServicer
func (svc *userService) FindStorages(ctx context.Context, id int) ([]dto.FindStorageDTO, error) {
	result, err := svc.repo.FindStorages(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error finding storages: %w", err)
	}
	return translate.StorageModelsToFindDTOs(result), nil
}

// FindRoles implements UserServicer
func (svc *userService) FindRoles(ctx context.Context, id int) ([]dto.FindRoleDTO, error) {
	if _, err := svc.repo.FindByID(ctx, id); err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}
	entities, err := svc.repo.FindRoles(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error finding roles: %w", err)
	}
	result := translate.RoleModelsToFindRoleDTOs(entities)
	return result, nil
}

// FindSettings implements UserServicer
func (svc *userService) FindSettings(ctx context.Context, id int) ([]dto.FindSettingDTO, error) {
	if _, err := svc.repo.FindByID(ctx, id); err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}
	entities, err := svc.repo.FindSettings(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error finding settings: %w", err)
	}
	result := translate.SettingModelsToFindDTOs(entities)
	return result, nil
}

// UpdateSetting implements UserServicer
func (svc *userService) UpdateSetting(ctx context.Context, id int, payload *dto.UpdateUserSettingDTO) (dto.FindSettingDTO, error) {
	if _, err := svc.repo.FindByID(ctx, id); err != nil {
		return dto.FindSettingDTO{}, fmt.Errorf("error finding user: %w", err)
	}
	entity := translate.UpdateSettingDTOToModel(payload)
	result, err := svc.repo.UpdateSetting(ctx, id, entity)
	if err != nil {
		return dto.FindSettingDTO{}, fmt.Errorf("error updating setting: %w", err)
	}
	dto := translate.SettingModelToFindDTO(&result)
	return dto, nil
}

func (svc *userService) FindUserByID(ctx context.Context, id int) (dto.FindUserDTO, error) {
	user, err := svc.repo.FindByID(ctx, id)
	result := translate.UserModelToFindDTO(&user)
	if err != nil {
		return dto.FindUserDTO{}, err
	}
	return result, nil
}

func (svc *userService) FindUsers(ctx context.Context, filter *dto.UserFilterDTO) ([]dto.FindUserDTO, error) {
	users, err := svc.repo.FindMany(ctx, filter.Limit, filter.Offset, filter.Name)
	dtos := translate.UserModelsToFindDTOs(users)
	if err != nil {
		return nil, err
	}
	return dtos, nil
}

func (svc *userService) CreateUser(ctx context.Context, payload *dto.CreateUserDTO) error {
	model := translate.CreateUserDTOToModel(payload)
	if err := svc.repo.Create(ctx, model); err != nil {
		return err
	}
	return nil
}

func (svc *userService) UpdateUser(ctx context.Context, id int, user *dto.UpdateUserDTO) error {
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
