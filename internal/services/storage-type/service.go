package storagetype

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/pkg/translate"
	repository "github.com/romankravchuk/muerta/internal/repositories/storage-type"
	"github.com/romankravchuk/muerta/internal/services"
)

type StorageTypeServicer interface {
	FindStorageTypeByID(ctx context.Context, id int) (dto.FindStorageTypeDTO, error)
	FindStorageTypes(ctx context.Context, filter *dto.StorageTypeFilterDTO) ([]dto.FindStorageTypeDTO, error)
	CreateStorageType(ctx context.Context, payload *dto.CreateStorageTypeDTO) error
	UpdateStorageType(ctx context.Context, id int, payload *dto.UpdateStorageTypeDTO) error
	DeleteStorageType(ctx context.Context, id int) error
	FindStorages(ctx context.Context, id int) ([]dto.FindStorageDTO, error)
	FindTips(ctx context.Context, id int) ([]dto.FindTipDTO, error)
	CreateTip(ctx context.Context, id, tipID int) (dto.FindTipDTO, error)
	DeleteTip(ctx context.Context, id, tipID int) error
	services.Counter
}

type storageTypeService struct {
	repo repository.StorageTypeRepositorer
}

func (s *storageTypeService) Count(ctx context.Context) (int, error) {
	count, err := s.repo.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("error counting storages types: %w", err)
	}
	return count, nil
}

// CreateTip implements StorageTypeServicer
func (s *storageTypeService) CreateTip(ctx context.Context, id int, tipID int) (dto.FindTipDTO, error) {
	result, err := s.repo.CreateTip(ctx, id, tipID)
	if err != nil {
		return dto.FindTipDTO{}, err
	}
	return translate.TipModelToFindDTO(&result), nil
}

// DeleteTip implements StorageTypeServicer
func (s *storageTypeService) DeleteTip(ctx context.Context, id int, tipID int) error {
	if err := s.repo.DeleteTip(ctx, id, tipID); err != nil {
		return err
	}
	return nil
}

// FindStorages implements StorageTypeServicer
func (s *storageTypeService) FindStorages(ctx context.Context, id int) ([]dto.FindStorageDTO, error) {
	result, err := s.repo.FindStorages(ctx, id)
	if err != nil {
		return nil, err
	}
	return translate.StorageModelsToFindDTOs(result), nil
}

// FindTips implements StorageTypeServicer
func (s *storageTypeService) FindTips(ctx context.Context, id int) ([]dto.FindTipDTO, error) {
	result, err := s.repo.FindTips(ctx, id)
	if err != nil {
		return nil, err
	}
	return translate.TipModelsToFindDTOs(result), nil
}

// CreateStorageType implements StorageTypeServicer
func (svc *storageTypeService) CreateStorageType(ctx context.Context, payload *dto.CreateStorageTypeDTO) error {
	model := translate.CreateStorageTypeDTOToModel(payload)
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
func (svc *storageTypeService) FindStorageTypeByID(ctx context.Context, id int) (dto.FindStorageTypeDTO, error) {
	model, err := svc.repo.FindByID(ctx, id)
	result := translate.StorageTypeModelToFindDTO(&model)
	if err != nil {
		return dto.FindStorageTypeDTO{}, err
	}
	return result, nil
}

// FindStorageTypes implements StorageTypeServicer
func (svc *storageTypeService) FindStorageTypes(ctx context.Context, filter *dto.StorageTypeFilterDTO) ([]dto.FindStorageTypeDTO, error) {
	models, err := svc.repo.FindMany(ctx, filter.Limit, filter.Offset, filter.Name)
	dtos := translate.StorageTypeModelsToFindDTOs(models)
	if err != nil {
		return nil, err
	}
	return dtos, nil
}

// UpdateStorageType implements StorageTypeServicer
func (svc *storageTypeService) UpdateStorageType(ctx context.Context, id int, payload *dto.UpdateStorageTypeDTO) error {
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
