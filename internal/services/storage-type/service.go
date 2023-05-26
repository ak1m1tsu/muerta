package storagetype

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/pkg/translate"
	"github.com/romankravchuk/muerta/internal/storage/postgres/models"
	repository "github.com/romankravchuk/muerta/internal/storage/postgres/storage-type"
)

type StorageTypeServicer interface {
	FindStorageTypeByID(ctx context.Context, id int) (dto.FindStorageType, error)
	FindStorageTypes(
		ctx context.Context,
		filter *dto.StorageTypeFilter,
	) ([]dto.FindStorageType, error)
	CreateStorageType(ctx context.Context, payload *dto.CreateStorageType) error
	UpdateStorageType(ctx context.Context, id int, payload *dto.UpdateStorageType) error
	DeleteStorageType(ctx context.Context, id int) error
	FindStorages(ctx context.Context, id int) ([]dto.FindStorage, error)
	FindTips(ctx context.Context, id int) ([]dto.FindTip, error)
	CreateTip(ctx context.Context, id, tipID int) (dto.FindTip, error)
	DeleteTip(ctx context.Context, id, tipID int) error
	Count(ctx context.Context, filter dto.StorageTypeFilter) (int, error)
}

type storageTypeService struct {
	repo repository.StorageTypeRepositorer
}

func (s *storageTypeService) Count(
	ctx context.Context,
	filter dto.StorageTypeFilter,
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
) (dto.FindTip, error) {
	result, err := s.repo.CreateTip(ctx, id, tipID)
	if err != nil {
		return dto.FindTip{}, err
	}
	return translate.TipModelToFind(&result), nil
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
) ([]dto.FindStorage, error) {
	result, err := s.repo.FindStorages(ctx, id)
	if err != nil {
		return nil, err
	}
	return translate.StorageModelsToFinds(result), nil
}

// FindTips implements StorageTypeServicer
func (s *storageTypeService) FindTips(ctx context.Context, id int) ([]dto.FindTip, error) {
	result, err := s.repo.FindTips(ctx, id)
	if err != nil {
		return nil, err
	}
	return translate.TipModelsToFinds(result), nil
}

// CreateStorageType implements StorageTypeServicer
func (svc *storageTypeService) CreateStorageType(
	ctx context.Context,
	payload *dto.CreateStorageType,
) error {
	model := translate.CreateStorageTypeToModel(payload)
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
) (dto.FindStorageType, error) {
	model, err := svc.repo.FindByID(ctx, id)
	result := translate.StorageTypeModelToFind(&model)
	if err != nil {
		return dto.FindStorageType{}, err
	}
	return result, nil
}

// FindStorageTypes implements StorageTypeServicer
func (svc *storageTypeService) FindStorageTypes(
	ctx context.Context,
	filter *dto.StorageTypeFilter,
) ([]dto.FindStorageType, error) {
	models, err := svc.repo.FindMany(ctx, models.StorageTypeFilter{
		PageFilter: models.PageFilter{Limit: filter.Limit, Offset: filter.Offset},
		Name:       filter.Name,
	})
	dtos := translate.StorageTypeModelsToFinds(models)
	if err != nil {
		return nil, err
	}
	return dtos, nil
}

// UpdateStorageType implements StorageTypeServicer
func (svc *storageTypeService) UpdateStorageType(
	ctx context.Context,
	id int,
	payload *dto.UpdateStorageType,
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
