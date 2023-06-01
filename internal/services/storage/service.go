package storage

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/api/router/params"
	"github.com/romankravchuk/muerta/internal/services/utils"
	"github.com/romankravchuk/muerta/internal/storage/postgres/models"
	"github.com/romankravchuk/muerta/internal/storage/postgres/storage"
)

type StorageServicer interface {
	FindStorageByID(ctx context.Context, id int) (params.FindStorage, error)
	FindStorages(ctx context.Context, filter *params.StorageFilter) ([]params.FindStorage, error)
	CreateStorage(ctx context.Context, payload *params.CreateStorage) error
	UpdateStorage(ctx context.Context, id int, payload *params.UpdateStorage) error
	DeleteStorage(ctx context.Context, id int) error
	RestoreStorage(ctx context.Context, id int) error
	FindTips(ctx context.Context, id int) ([]params.FindTip, error)
	CreateTip(ctx context.Context, id, tipID int) (params.FindTip, error)
	DeleteTip(ctx context.Context, id, tipID int) error
	FindShelfLives(ctx context.Context, id int) ([]params.FindShelfLife, error)
	Count(ctx context.Context, filter params.StorageFilter) (int, error)
}

type storageService struct {
	repo storage.StorageRepositorer
}

func (s *storageService) FindShelfLives(
	ctx context.Context,
	id int,
) ([]params.FindShelfLife, error) {
	result, err := s.repo.FindShelfLives(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find shelf lives: %w", err)
	}
	return utils.ShelfLifeModelsToFinds(result), nil
}

func (s *storageService) Count(ctx context.Context, filter params.StorageFilter) (int, error) {
	count, err := s.repo.Count(ctx, models.StorageFilter{Name: filter.Name})
	if err != nil {
		return 0, fmt.Errorf("error counting storages: %w", err)
	}
	return count, nil
}

// CreateTip implements StorageServicer
func (s *storageService) CreateTip(ctx context.Context, id int, tipID int) (params.FindTip, error) {
	result, err := s.repo.CreateTip(ctx, id, tipID)
	if err != nil {
		return params.FindTip{}, fmt.Errorf("failed to create tip: %w", err)
	}
	return utils.TipModelToFind(&result), nil
}

// DeleteTip implements StorageServicer
func (s *storageService) DeleteTip(ctx context.Context, id int, tipID int) error {
	if err := s.repo.DeleteTip(ctx, id, tipID); err != nil {
		return fmt.Errorf("failed to delete tip: %w", err)
	}
	return nil
}

// FindTips implements StorageServicer
func (s *storageService) FindTips(ctx context.Context, id int) ([]params.FindTip, error) {
	result, err := s.repo.FindTips(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find tips: %w", err)
	}
	return utils.TipModelsToFinds(result), nil
}

func New(repo storage.StorageRepositorer) StorageServicer {
	return &storageService{
		repo: repo,
	}
}

func (s *storageService) FindStorageByID(ctx context.Context, id int) (params.FindStorage, error) {
	model, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return params.FindStorage{}, fmt.Errorf("failed to find storage: %w", err)
	}
	dto := utils.StorageModelToFind(&model)
	return dto, nil
}

func (s *storageService) FindStorages(
	ctx context.Context,
	filter *params.StorageFilter,
) ([]params.FindStorage, error) {
	models, err := s.repo.FindMany(ctx, models.StorageFilter{
		PageFilter: models.PageFilter{Limit: filter.Limit, Offset: filter.Offset},
		Name:       filter.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find storages: %w", err)
	}
	dtos := utils.StorageModelsToFinds(models)
	return dtos, nil
}

func (s *storageService) CreateStorage(ctx context.Context, payload *params.CreateStorage) error {
	model := utils.CreateStorageToModel(payload)
	if err := s.repo.Create(ctx, &model); err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}
	return nil
}

func (s *storageService) UpdateStorage(
	ctx context.Context,
	id int,
	payload *params.UpdateStorage,
) error {
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
