package storage

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/pkg/translate"
	"github.com/romankravchuk/muerta/internal/repositories/storage"
)

type StorageServicer interface {
	FindStorageByID(ctx context.Context, id int) (dto.FindStorageDTO, error)
	FindStorages(ctx context.Context, filter *dto.StorageFilterDTO) ([]dto.FindStorageDTO, error)
	CreateStorage(ctx context.Context, payload *dto.CreateStorageDTO) error
	UpdateStorage(ctx context.Context, id int, payload *dto.UpdateStorageDTO) error
	DeleteStorage(ctx context.Context, id int) error
	RestoreStorage(ctx context.Context, id int) error
}

type storageService struct {
	repo storage.StorageRepositorer
}

func New(repo storage.StorageRepositorer) StorageServicer {
	return &storageService{
		repo: repo,
	}
}

func (s *storageService) FindStorageByID(ctx context.Context, id int) (dto.FindStorageDTO, error) {
	model, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return dto.FindStorageDTO{}, fmt.Errorf("failed to find storage: %w", err)
	}
	dto := translate.StorageModelToFindDTO(&model)
	return dto, nil
}

func (s *storageService) FindStorages(ctx context.Context, filter *dto.StorageFilterDTO) ([]dto.FindStorageDTO, error) {
	models, err := s.repo.FindMany(ctx, filter.Limit, filter.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find storages: %w", err)
	}
	dtos := translate.StorageModelsToFindDTOs(models)
	return dtos, nil
}

func (s *storageService) CreateStorage(ctx context.Context, payload *dto.CreateStorageDTO) error {
	model := translate.CreateStorageDTOToModel(payload)
	if err := s.repo.Create(ctx, &model); err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}
	return nil
}

func (s *storageService) UpdateStorage(ctx context.Context, id int, payload *dto.UpdateStorageDTO) error {
	model, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find storage: %w", err)
	}
	fmt.Printf("%+v\n", payload)
	if payload.Name != "" {
		model.Name = payload.Name
	}
	if payload.Temperature != 0.0 {
		model.Temperature = payload.Temperature
	}
	if payload.Humidity != 0.0 {
		model.Humidity = payload.Humidity
	}
	if payload.TypeID != 0 {
		model.Type.ID = payload.TypeID
	}
	fmt.Printf("%+v\n", model)
	if err := s.repo.Update(ctx, &model); err != nil {
		return fmt.Errorf("failed to update storage: %w", err)
	}
	return nil
}

func (s *storageService) DeleteStorage(ctx context.Context, id int) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete storage: %w", err)
	}
	return nil
}

func (s *storageService) RestoreStorage(ctx context.Context, id int) error {
	if err := s.repo.Restore(ctx, id); err != nil {
		return fmt.Errorf("failed to restore storage: %w", err)
	}
	return nil
}
