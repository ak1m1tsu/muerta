package tip

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/pkg/translate"
	repository "github.com/romankravchuk/muerta/internal/repositories/tip"
	"github.com/romankravchuk/muerta/internal/services"
)

type TipServicer interface {
	FindTipByID(ctx context.Context, id int) (dto.FindTipDTO, error)
	FindTips(ctx context.Context, filter *dto.TipFilterDTO) ([]dto.FindTipDTO, error)
	CreateTip(ctx context.Context, payload *dto.CreateTipDTO) error
	UpdateTip(ctx context.Context, id int, payload *dto.UpdateTipDTO) error
	DeleteTip(ctx context.Context, id int) error
	RestoreTip(ctx context.Context, id int) error
	FindTipStorages(ctx context.Context, id int) ([]dto.FindStorageDTO, error)
	FindTipProducts(ctx context.Context, id int) ([]dto.FindProductDTO, error)
	services.Counter
}

type tipService struct {
	repo repository.TipRepositorer
}

func New(repo repository.TipRepositorer) TipServicer {
	return &tipService{
		repo: repo,
	}
}

func (s *tipService) Count(ctx context.Context) (int, error) {
	count, err := s.repo.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("error counting users: %w", err)
	}
	return count, nil
}

// FindTipProducts implements TipServicer
func (svc *tipService) FindTipProducts(ctx context.Context, id int) ([]dto.FindProductDTO, error) {
	products, err := svc.repo.FindProducts(ctx, id)
	if err != nil {
		return nil, err
	}
	dtos := translate.ProductModelsToFindDTOs(products)
	return dtos, nil
}

// FindTipStorages implements TipServicer
func (sbc *tipService) FindTipStorages(ctx context.Context, id int) ([]dto.FindStorageDTO, error) {
	storages, err := sbc.repo.FindStorages(ctx, id)
	if err != nil {
		return nil, err
	}
	dtos := translate.StorageModelsToFindDTOs(storages)
	return dtos, nil
}

// CreateTip implements TipServicer
func (svc *tipService) CreateTip(ctx context.Context, payload *dto.CreateTipDTO) error {
	model := translate.CreateTipDTOToModel(payload)
	if err := svc.repo.Create(ctx, model); err != nil {
		return err
	}
	return nil
}

// DeleteTip implements TipServicer
func (svc *tipService) DeleteTip(ctx context.Context, id int) error {
	if err := svc.repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

// FindTipByID implements TipServicer
func (svc *tipService) FindTipByID(ctx context.Context, id int) (dto.FindTipDTO, error) {
	model, err := svc.repo.FindByID(ctx, id)
	result := translate.TipModelToFindDTO(&model)
	if err != nil {
		return dto.FindTipDTO{}, err
	}
	return result, nil
}

// FindTips implements TipServicer
func (svc *tipService) FindTips(ctx context.Context, filter *dto.TipFilterDTO) ([]dto.FindTipDTO, error) {
	models, err := svc.repo.FindMany(ctx, filter.Limit, filter.Offset, filter.Description)
	dtos := translate.TipModelsToFindDTOs(models)
	if err != nil {
		return nil, err
	}
	return dtos, nil
}

// RestoreTip implements TipServicer
func (svc *tipService) RestoreTip(ctx context.Context, id int) error {
	if err := svc.repo.Restore(ctx, id); err != nil {
		return err
	}
	return nil
}

// UpdateTip implements TipServicer
func (svc *tipService) UpdateTip(ctx context.Context, id int, payload *dto.UpdateTipDTO) error {
	model, err := svc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if payload.Description != "" {
		model.Description = payload.Description
	}
	if err := svc.repo.Update(ctx, model); err != nil {
		return err
	}
	return nil
}
