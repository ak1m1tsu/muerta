package tip

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/api/router/params"
	"github.com/romankravchuk/muerta/internal/services/utils"
	"github.com/romankravchuk/muerta/internal/storage/postgres/models"
	repository "github.com/romankravchuk/muerta/internal/storage/postgres/tip"
)

type TipServicer interface {
	FindTipByID(ctx context.Context, id int) (params.FindTip, error)
	FindTips(ctx context.Context, filter *params.TipFilter) ([]params.FindTip, error)
	CreateTip(ctx context.Context, payload *params.CreateTip) (params.FindTip, error)
	UpdateTip(ctx context.Context, id int, payload *params.UpdateTip) error
	DeleteTip(ctx context.Context, id int) error
	RestoreTip(ctx context.Context, id int) error
	FindTipStorages(ctx context.Context, id int) ([]params.FindStorage, error)
	FindTipProducts(ctx context.Context, id int) ([]params.FindProduct, error)
	AddProductToTip(ctx context.Context, tipID, productID int) (params.FindProduct, error)
	RemoveProductFromTip(ctx context.Context, tipID, productID int) error
	AddStorageToTip(ctx context.Context, tipID, storageID int) (params.FindStorage, error)
	RemoveStorageFromTip(ctx context.Context, tipID, storageID int) error
	Count(ctx context.Context, filter params.TipFilter) (int, error)
}

type tipService struct {
	repo repository.TipRepositorer
}

// AddProductToTip implements TipServicer
func (s *tipService) AddProductToTip(
	ctx context.Context,
	tipID int,
	productID int,
) (params.FindProduct, error) {
	result, err := s.repo.AddProduct(ctx, tipID, productID)
	if err != nil {
		return params.FindProduct{}, err
	}
	return utils.ProductModelToFind(&result), nil
}

// AddStorageToTip implements TipServicer
func (s *tipService) AddStorageToTip(
	ctx context.Context,
	tipID int,
	storageID int,
) (params.FindStorage, error) {
	result, err := s.repo.AddStorage(ctx, tipID, storageID)
	if err != nil {
		return params.FindStorage{}, err
	}
	return utils.StorageModelToFind(&result), nil
}

// RemoveProductFromTip implements TipServicer
func (s *tipService) RemoveProductFromTip(ctx context.Context, tipID int, productID int) error {
	if err := s.repo.RemoveProduct(ctx, tipID, productID); err != nil {
		return err
	}
	return nil
}

// RemoveStorageFromTip implements TipServicer
func (s *tipService) RemoveStorageFromTip(ctx context.Context, tipID int, storageID int) error {
	if err := s.repo.RemoveStorage(ctx, tipID, storageID); err != nil {
		return err
	}
	return nil
}

func New(repo repository.TipRepositorer) TipServicer {
	return &tipService{
		repo: repo,
	}
}

func (s *tipService) Count(ctx context.Context, filter params.TipFilter) (int, error) {
	count, err := s.repo.Count(ctx, models.TipFilter{Description: filter.Description})
	if err != nil {
		return 0, fmt.Errorf("error counting users: %w", err)
	}
	return count, nil
}

// FindTipProducts implements TipServicer
func (svc *tipService) FindTipProducts(ctx context.Context, id int) ([]params.FindProduct, error) {
	products, err := svc.repo.FindProducts(ctx, id)
	if err != nil {
		return nil, err
	}
	dtos := utils.ProductModelsToFinds(products)
	return dtos, nil
}

// FindTipStorages implements TipServicer
func (sbc *tipService) FindTipStorages(ctx context.Context, id int) ([]params.FindStorage, error) {
	storages, err := sbc.repo.FindStorages(ctx, id)
	if err != nil {
		return nil, err
	}
	dtos := utils.StorageModelsToFinds(storages)
	return dtos, nil
}

// CreateTip implements TipServicer
func (svc *tipService) CreateTip(
	ctx context.Context,
	payload *params.CreateTip,
) (params.FindTip, error) {
	model := utils.CreateTipToModel(payload)
	if err := svc.repo.Create(ctx, &model); err != nil {
		return params.FindTip{}, err
	}
	return utils.TipModelToFind(&model), nil
}

// DeleteTip implements TipServicer
func (svc *tipService) DeleteTip(ctx context.Context, id int) error {
	if err := svc.repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

// FindTipByID implements TipServicer
func (svc *tipService) FindTipByID(ctx context.Context, id int) (params.FindTip, error) {
	model, err := svc.repo.FindByID(ctx, id)
	result := utils.TipModelToFind(&model)
	if err != nil {
		return params.FindTip{}, err
	}
	return result, nil
}

// FindTips implements TipServicer
func (svc *tipService) FindTips(
	ctx context.Context,
	filter *params.TipFilter,
) ([]params.FindTip, error) {
	models, err := svc.repo.FindMany(ctx, models.TipFilter{
		PageFilter: models.PageFilter{
			Limit:  filter.Limit,
			Offset: filter.Offset,
		},
		Description: filter.Description,
	})
	dtos := utils.TipModelsToFinds(models)
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
func (svc *tipService) UpdateTip(ctx context.Context, id int, payload *params.UpdateTip) error {
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
