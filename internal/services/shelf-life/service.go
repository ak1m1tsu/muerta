package shelflife

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/api/router/params"
	"github.com/romankravchuk/muerta/internal/services/utils"
	"github.com/romankravchuk/muerta/internal/storage/postgres/models"
	repository "github.com/romankravchuk/muerta/internal/storage/postgres/shelf-life"
)

type ShelfLifeServicer interface {
	FindShelfLifeByID(ctx context.Context, id int) (params.FindShelfLife, error)
	FindShelfLifes(ctx context.Context, filter *params.ShelfLifeFilter) ([]params.FindShelfLife, error)
	CreateShelfLife(ctx context.Context, payload *params.CreateShelfLife) error
	UpdateShelfLife(ctx context.Context, id int, payload *params.UpdateShelfLife) error
	DeleteShelfLife(ctx context.Context, id int) error
	RestoreShelfLife(ctx context.Context, id int) error
	FindShelfLifeStatuses(ctx context.Context, id int) ([]params.FindShelfLifeStatus, error)
	CreateShelfLifeStatus(ctx context.Context, id, status int) (params.FindShelfLifeStatus, error)
	DeleteShelfLifeStatus(ctx context.Context, id, status int) error
	Count(ctx context.Context, filter params.ShelfLifeFilter) (int, error)
}

type shelfLifeSerivce struct {
	repo repository.ShelfLifeRepositorer
}

func (s *shelfLifeSerivce) Count(ctx context.Context, filter params.ShelfLifeFilter) (int, error) {
	count, err := s.repo.Count(ctx, models.ShelfLifeFilter{})
	if err != nil {
		return 0, fmt.Errorf("error counting shelf lives: %w", err)
	}
	return count, nil
}

// CreateShelfLifeStatus implements ShelfLifeServicer
func (s *shelfLifeSerivce) CreateShelfLifeStatus(ctx context.Context, id, status int) (params.FindShelfLifeStatus, error) {
	model, err := s.repo.CreateStatus(ctx, id, status)
	if err != nil {
		return params.FindShelfLifeStatus{}, err
	}
	return utils.ShelfLifeStatusModelToFind(&model), nil
}

// DeleteShelfLifeStatus implements ShelfLifeServicer
func (s *shelfLifeSerivce) DeleteShelfLifeStatus(ctx context.Context, id, status int) error {
	if err := s.repo.DeleteStatus(ctx, id, status); err != nil {
		return err
	}
	return nil
}

// FindShelfLifeStatuses implements ShelfLifeServicer
func (s *shelfLifeSerivce) FindShelfLifeStatuses(ctx context.Context, id int) ([]params.FindShelfLifeStatus, error) {
	models, err := s.repo.FindStatuses(ctx, id)
	if err != nil {
		return nil, err
	}
	return utils.ShelfLifeStatusModelsToFinds(models), nil
}

// CreateShelfLife implements ShelfLifeServicer
func (svc *shelfLifeSerivce) CreateShelfLife(ctx context.Context, payload *params.CreateShelfLife) error {
	model := utils.CreateShelfLifeToModel(payload)
	if err := svc.repo.Create(ctx, model); err != nil {
		return err
	}
	return nil
}

// DeleteShelfLife implements ShelfLifeServicer
func (svc *shelfLifeSerivce) DeleteShelfLife(ctx context.Context, id int) error {
	if err := svc.repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

// FindShelfLifeByID implements ShelfLifeServicer
func (svc *shelfLifeSerivce) FindShelfLifeByID(ctx context.Context, id int) (params.FindShelfLife, error) {
	model, err := svc.repo.FindByID(ctx, id)
	result := utils.ShelfLifeModelToFind(&model)
	if err != nil {
		return params.FindShelfLife{}, err
	}
	return result, nil
}

// FindShelfLifes implements ShelfLifeServicer
func (svc *shelfLifeSerivce) FindShelfLifes(ctx context.Context, filter *params.ShelfLifeFilter) ([]params.FindShelfLife, error) {
	models, err := svc.repo.FindMany(ctx, models.ShelfLifeFilter{
		PageFilter: models.PageFilter{
			Limit:  filter.Limit,
			Offset: filter.Offset,
		},
	})
	fmt.Printf("%+v\n", models)
	dtos := utils.ShelfLifeModelsToFinds(models)
	if err != nil {
		return nil, err
	}
	return dtos, nil
}

// RestoreShelfLife implements ShelfLifeServicer
func (svc *shelfLifeSerivce) RestoreShelfLife(ctx context.Context, id int) error {
	if err := svc.repo.Restore(ctx, id); err != nil {
		return err
	}
	return nil
}

// UpdateShelfLife implements ShelfLifeServicer
func (svc *shelfLifeSerivce) UpdateShelfLife(
	ctx context.Context,
	id int,
	payload *params.UpdateShelfLife,
) error {
	model, err := svc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if payload.PurchaseDate != nil {
		model.PurchaseDate = payload.PurchaseDate
	}
	if payload.EndDate != nil {
		model.EndDate = payload.EndDate
	}
	if payload.Quantity > 0 {
		model.Quantity = payload.Quantity
	}
	if payload.MeasureID != 0 {
		model.Measure.ID = payload.MeasureID
	}
	if payload.StorageID != 0 {
		model.Storage.ID = payload.StorageID
	}
	if payload.ProductID != 0 {
		model.Product.ID = payload.ProductID
	}
	if err := svc.repo.Update(ctx, model); err != nil {
		return err
	}
	return nil
}

func New(repo repository.ShelfLifeRepositorer) ShelfLifeServicer {
	return &shelfLifeSerivce{
		repo: repo,
	}
}
