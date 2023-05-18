package user

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/pkg/translate"
	"github.com/romankravchuk/muerta/internal/repositories/models"
	repo "github.com/romankravchuk/muerta/internal/repositories/user"
)

type UserServicer interface {
	FindUserByID(ctx context.Context, id int) (dto.FindUser, error)
	FindUsers(ctx context.Context, filter *dto.UserFilter) ([]dto.FindUser, error)
	CreateUser(ctx context.Context, payload *dto.CreateUser) error
	UpdateUser(ctx context.Context, id int, user *dto.UpdateUser) error
	DeleteUser(ctx context.Context, id int) error
	RestoreUser(ctx context.Context, id int) error
	FindSettings(ctx context.Context, id int) ([]dto.FindSetting, error)
	UpdateSetting(
		ctx context.Context,
		id int,
		payload *dto.UpdateUserSetting,
	) (dto.FindSetting, error)
	FindRoles(ctx context.Context, id int) ([]dto.FindRole, error)
	CreateStorage(
		ctx context.Context,
		id int,
		payload *dto.UserStorage,
	) (dto.FindStorage, error)
	DeleteStorage(ctx context.Context, id int, payload *dto.UserStorage) error
	FindStorages(ctx context.Context, id int) ([]dto.FindStorage, error)
	FindShelfLives(ctx context.Context, id int) ([]dto.FindShelfLife, error)
	CreateShelfLife(
		ctx context.Context,
		id int,
		payload *dto.CreateShelfLife,
	) (dto.FindShelfLife, error)
	UpdateShelfLife(
		ctx context.Context,
		id int,
		payload *dto.UserShelfLife,
	) (dto.FindShelfLife, error)
	RestoreShelfLife(
		ctx context.Context,
		id int,
		payload *dto.UserShelfLife,
	) (dto.FindShelfLife, error)
	DeleteShelfLife(ctx context.Context, id int, payload *dto.UserShelfLife) error
	Count(ctx context.Context, filter dto.UserFilter) (int, error)
}

type userService struct {
	repo repo.UserRepositorer
}

// Count implements UserServicer
func (s *userService) Count(ctx context.Context, filter dto.UserFilter) (int, error) {
	count, err := s.repo.Count(ctx, models.UserFilter{Name: filter.Name})
	if err != nil {
		return 0, fmt.Errorf("error counting users: %w", err)
	}
	return count, nil
}

func New(repo repo.UserRepositorer) UserServicer {
	return &userService{
		repo: repo,
	}
}

// CreateShelfLife implements UserServicer
func (svc *userService) CreateShelfLife(
	ctx context.Context,
	id int,
	payload *dto.CreateShelfLife,
) (dto.FindShelfLife, error) {
	model := translate.CreateShelfLifeToModel(payload)
	createdModel, err := svc.repo.CreateShelfLife(ctx, id, model)
	if err != nil {
		return dto.FindShelfLife{}, fmt.Errorf("error creating shelf life: %w", err)
	}
	return translate.ShelfLifeModelToFind(&createdModel), nil
}

// DeleteShelfLife implements UserServicer
func (svc *userService) DeleteShelfLife(
	ctx context.Context,
	id int,
	payload *dto.UserShelfLife,
) error {
	if err := svc.repo.DeleteShelfLife(ctx, id, payload.ShelfLifeID); err != nil {
		return fmt.Errorf("error deleting shelf life: %w", err)
	}
	return nil
}

// FindShelfLives implements UserServicer
func (svc *userService) FindShelfLives(
	ctx context.Context,
	id int,
) ([]dto.FindShelfLife, error) {
	models, err := svc.repo.FindShelfLives(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error finding shelf lives: %w", err)
	}
	return translate.ShelfLifeModelsToFinds(models), nil
}

// RestoreShelfLife implements UserServicer
func (svc *userService) RestoreShelfLife(
	ctx context.Context,
	id int,
	payload *dto.UserShelfLife,
) (dto.FindShelfLife, error) {
	model, err := svc.repo.RestoreShelfLife(ctx, id, payload.ShelfLifeID)
	if err != nil {
		return dto.FindShelfLife{}, fmt.Errorf("error restoring shelf life: %w", err)
	}
	return translate.ShelfLifeModelToFind(&model), nil
}

// UpdateShelfLife implements UserServicer
func (svc *userService) UpdateShelfLife(
	ctx context.Context,
	id int,
	payload *dto.UserShelfLife,
) (dto.FindShelfLife, error) {
	model, err := svc.repo.FindShelfLife(ctx, id, payload.ShelfLifeID)
	if err != nil {
		return dto.FindShelfLife{}, fmt.Errorf("error finding shelf life: %w", err)
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
		return dto.FindShelfLife{}, fmt.Errorf("error updating shelf life: %w", err)
	}
	return translate.ShelfLifeModelToFind(&result), nil
}

// DeleteStorage implements UserServicer
func (svc *userService) DeleteStorage(
	ctx context.Context,
	id int,
	payload *dto.UserStorage,
) error {
	model := translate.UserStorageToModel(payload)
	err := svc.repo.DeleteStorage(ctx, id, model)
	if err != nil {
		return fmt.Errorf("error creating storage: %w", err)
	}
	return nil
}

// CreateStorage implements UserServicer
func (svc *userService) CreateStorage(
	ctx context.Context,
	id int,
	payload *dto.UserStorage,
) (dto.FindStorage, error) {
	model := translate.UserStorageToModel(payload)
	model, err := svc.repo.CreateStorage(ctx, id, model)
	if err != nil {
		return dto.FindStorage{}, fmt.Errorf("error creating storage: %w", err)
	}
	return translate.StorageModelToFind(&model), nil
}

// FindStorages implements UserServicer
func (svc *userService) FindStorages(ctx context.Context, id int) ([]dto.FindStorage, error) {
	result, err := svc.repo.FindStorages(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error finding storages: %w", err)
	}
	return translate.StorageModelsToFinds(result), nil
}

// FindRoles implements UserServicer
func (svc *userService) FindRoles(ctx context.Context, id int) ([]dto.FindRole, error) {
	if _, err := svc.repo.FindByID(ctx, id); err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}
	entities, err := svc.repo.FindRoles(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error finding roles: %w", err)
	}
	result := translate.RoleModelsToFindRoles(entities)
	return result, nil
}

// FindSettings implements UserServicer
func (svc *userService) FindSettings(ctx context.Context, id int) ([]dto.FindSetting, error) {
	if _, err := svc.repo.FindByID(ctx, id); err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}
	entities, err := svc.repo.FindSettings(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error finding settings: %w", err)
	}
	result := translate.SettingModelsToFinds(entities)
	return result, nil
}

// UpdateSetting implements UserServicer
func (svc *userService) UpdateSetting(
	ctx context.Context,
	id int,
	payload *dto.UpdateUserSetting,
) (dto.FindSetting, error) {
	if _, err := svc.repo.FindByID(ctx, id); err != nil {
		return dto.FindSetting{}, fmt.Errorf("error finding user: %w", err)
	}
	entity := translate.UpdateSettingToModel(payload)
	result, err := svc.repo.UpdateSetting(ctx, id, entity)
	if err != nil {
		return dto.FindSetting{}, fmt.Errorf("error updating setting: %w", err)
	}
	dto := translate.SettingModelToFind(&result)
	return dto, nil
}

func (svc *userService) FindUserByID(ctx context.Context, id int) (dto.FindUser, error) {
	user, err := svc.repo.FindByID(ctx, id)
	result := translate.UserModelToFind(&user)
	if err != nil {
		return dto.FindUser{}, err
	}
	return result, nil
}

func (svc *userService) FindUsers(
	ctx context.Context,
	filter *dto.UserFilter,
) ([]dto.FindUser, error) {
	users, err := svc.repo.FindMany(ctx, models.UserFilter{
		PageFilter: models.PageFilter{
			Limit:  filter.Limit,
			Offset: filter.Offset,
		},
		Name: filter.Name,
	})
	dtos := translate.UserModelsToFinds(users)
	if err != nil {
		return nil, err
	}
	return dtos, nil
}

func (svc *userService) CreateUser(ctx context.Context, payload *dto.CreateUser) error {
	model := translate.CreateUserToModel(payload)
	if err := svc.repo.Create(ctx, model); err != nil {
		return err
	}
	return nil
}

func (svc *userService) UpdateUser(ctx context.Context, id int, user *dto.UpdateUser) error {
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
