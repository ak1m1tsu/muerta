package storage

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/pkg/translate"
	"github.com/romankravchuk/muerta/internal/repositories/storage"
	"github.com/romankravchuk/muerta/internal/services"
)

type StorageServicer interface {
	FindStorageByID(ctx context.Context, id int) (dto.FindStorageDTO, error)
	FindStorages(ctx context.Context, filter *dto.StorageFilterDTO) ([]dto.FindStorageDTO, error)
	CreateStorage(ctx context.Context, payload *dto.CreateStorageDTO) error
	UpdateStorage(ctx context.Context, id int, payload *dto.UpdateStorageDTO) error
	DeleteStorage(ctx context.Context, id int) error
	RestoreStorage(ctx context.Context, id int) error
	FindTips(ctx context.Context, id int) ([]dto.FindTipDTO, error)
	CreateTip(ctx context.Context, id, tipID int) (dto.FindTipDTO, error)
	DeleteTip(ctx context.Context, id, tipID int) error
	services.Counter
}

type storageService struct {
	repo storage.StorageRepositorer
}

func (s *storageService) Count(ctx context.Context) (int, error) {
	count, err := s.repo.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("error counting storages: %w", err)
	}
	return count, nil
}

// CreateTip implements StorageServicer
func (s *storageService) CreateTip(ctx context.Context, id int, tipID int) (dto.FindTipDTO, error) {
	result, err := s.repo.CreateTip(ctx, id, tipID)
	if err != nil {
		return dto.FindTipDTO{}, fmt.Errorf("failed to create tip: %w", err)
	}
	return translate.TipModelToFindDTO(&result), nil
}

// DeleteTip implements StorageServicer
func (s *storageService) DeleteTip(ctx context.Context, id int, tipID int) error {
	if err := s.repo.DeleteTip(ctx, id, tipID); err != nil {
		return fmt.Errorf("failed to delete tip: %w", err)
	}
	return nil
}

// FindTips implements StorageServicer
func (s *storageService) FindTips(ctx context.Context, id int) ([]dto.FindTipDTO, error) {
	result, err := s.repo.FindTips(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find tips: %w", err)
	}
	return translate.TipModelsToFindDTOs(result), nil
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
