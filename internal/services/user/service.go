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
}

type userService struct {
	repo repo.UserRepositorer
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

func New(repo repo.UserRepositorer) UserServicer {
	return &userService{
		repo: repo,
	}
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
