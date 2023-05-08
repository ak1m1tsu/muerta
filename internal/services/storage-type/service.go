package storagetype

import (
	"context"

	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/pkg/translate"
	repository "github.com/romankravchuk/muerta/internal/repositories/storage-type"
)

type StorageTypeServicer interface {
	FindStorageTypeByID(ctx context.Context, id int) (dto.FindStorageTypeDTO, error)
	FindStorageTypes(ctx context.Context, filter *dto.StorageTypeFilterDTO) ([]dto.FindStorageTypeDTO, error)
	CreateStorageType(ctx context.Context, payload *dto.CreateStorageTypeDTO) error
	UpdateStorageType(ctx context.Context, id int, payload *dto.UpdateStorageTypeDTO) error
	DeleteStorageType(ctx context.Context, id int) error
}

type storageTypeService struct {
	repo repository.StorageTypeRepositorer
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
