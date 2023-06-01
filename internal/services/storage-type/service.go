package storagetype

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/api/router/params"
	"github.com/romankravchuk/muerta/internal/services/utils"
	"github.com/romankravchuk/muerta/internal/storage/postgres/models"
	repository "github.com/romankravchuk/muerta/internal/storage/postgres/storage-type"
)

type StorageTypeServicer interface {
	FindStorageTypeByID(ctx context.Context, id int) (params.FindStorageType, error)
	FindStorageTypes(
		ctx context.Context,
		filter *params.StorageTypeFilter,
	) ([]params.FindStorageType, error)
	CreateStorageType(ctx context.Context, payload *params.CreateStorageType) error
	UpdateStorageType(ctx context.Context, id int, payload *params.UpdateStorageType) error
	DeleteStorageType(ctx context.Context, id int) error
	FindStorages(ctx context.Context, id int) ([]params.FindStorage, error)
	FindTips(ctx context.Context, id int) ([]params.FindTip, error)
	CreateTip(ctx context.Context, id, tipID int) (params.FindTip, error)
	DeleteTip(ctx context.Context, id, tipID int) error
	Count(ctx context.Context, filter params.StorageTypeFilter) (int, error)
}

type storageTypeService struct {
	repo repository.StorageTypeRepositorer
}

func (s *storageTypeService) Count(
	ctx context.Context,
	filter params.StorageTypeFilter,
) (int, error) {
	count, err := s.repo.Count(ctx, models.StorageTypeFilter{Name: filter.Name})
	if err != nil {
		return 0, fmt.Errorf("error counting storages types: %w", err)
	}
	return count, nil
}

// CreateTip implements StorageTypeServicer
func (s *storageTypeService) CreateTip(
	ctx context.Context,
	id int,
	tipID int,
) (params.FindTip, error) {
	result, err := s.repo.CreateTip(ctx, id, tipID)
	if err != nil {
		return params.FindTip{}, err
	}
	return utils.TipModelToFind(&result), nil
}

// DeleteTip implements StorageTypeServicer
func (s *storageTypeService) DeleteTip(ctx context.Context, id int, tipID int) error {
	if err := s.repo.DeleteTip(ctx, id, tipID); err != nil {
		return err
	}
	return nil
}

// FindStorages implements StorageTypeServicer
func (s *storageTypeService) FindStorages(
	ctx context.Context,
	id int,
) ([]params.FindStorage, error) {
	result, err := s.repo.FindStorages(ctx, id)
	if err != nil {
		return nil, err
	}
	return utils.StorageModelsToFinds(result), nil
}

// FindTips implements StorageTypeServicer
func (s *storageTypeService) FindTips(ctx context.Context, id int) ([]params.FindTip, error) {
	result, err := s.repo.FindTips(ctx, id)
	if err != nil {
		return nil, err
	}
	return utils.TipModelsToFinds(result), nil
}

// CreateStorageType implements StorageTypeServicer
func (svc *storageTypeService) CreateStorageType(
	ctx context.Context,
	payload *params.CreateStorageType,
) error {
	model := utils.CreateStorageTypeToModel(payload)
	if err := svc.repo.Create(ctx, model); err != nil {
		return err
	}
	return nil
}

// DeleteStorageType implements StorageTypeServicer
func (svc *storageTypeService) DeleteStorageType(ctx context.Context, id int) error {
	if err := svc.repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

// FindStorageTypeByID implements StorageTypeServicer
func (svc *storageTypeService) FindStorageTypeByID(
	ctx context.Context,
	id int,
) (params.FindStorageType, error) {
	model, err := svc.repo.FindByID(ctx, id)
	result := utils.StorageTypeModelToFind(&model)
	if err != nil {
		return params.FindStorageType{}, err
	}
	return result, nil
}

// FindStorageTypes implements StorageTypeServicer
func (svc *storageTypeService) FindStorageTypes(
	ctx context.Context,
	filter *params.StorageTypeFilter,
) ([]params.FindStorageType, error) {
	models, err := svc.repo.FindMany(ctx, models.StorageTypeFilter{
		PageFilter: models.PageFilter{Limit: filter.Limit, Offset: filter.Offset},
		Name:       filter.Name,
	})
	dtos := utils.StorageTypeModelsToFinds(models)
	if err != nil {
		return nil, err
	}
	return dtos, nil
}

// UpdateStorageType implements StorageTypeServicer
func (svc *storageTypeService) UpdateStorageType(
	ctx context.Context,
	id int,
	payload *params.UpdateStorageType,
) error {
	model, err := svc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if payload.Name != "" {
		model.Name = payload.Name
	}
	if err := svc.repo.Update(ctx, model); err != nil {
		return err
	}
	return nil
}

func New(repo repository.StorageTypeRepositorer) StorageTypeServicer {
	return &storageTypeService{
		repo: repo,
	}
}
