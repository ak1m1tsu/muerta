package shelflife

import (
	"context"

	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/pkg/translate"
	repository "github.com/romankravchuk/muerta/internal/repositories/shelf-life"
)

type ShelfLifeServicer interface {
	FindShelfLifeByID(ctx context.Context, id int) (dto.FindShelfLifeDTO, error)
	FindShelfLifes(ctx context.Context, filter *dto.ShelfLifeFilterDTO) ([]dto.FindShelfLifeDTO, error)
	CreateShelfLife(ctx context.Context, payload *dto.CreateShelfLifeDTO) error
	UpdateShelfLife(ctx context.Context, id int, payload *dto.UpdateShelfLifeDTO) error
	DeleteShelfLife(ctx context.Context, id int) error
	RestoreShelfLife(ctx context.Context, id int) error
}

type shelfLifeSerivce struct {
	repo repository.ShelfLifeRepositorer
}

// CreateShelfLife implements ShelfLifeServicer
func (svc *shelfLifeSerivce) CreateShelfLife(ctx context.Context, payload *dto.CreateShelfLifeDTO) error {
	model := translate.CreateShelfLifeDTOToModel(payload)
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
func (svc *shelfLifeSerivce) FindShelfLifeByID(ctx context.Context, id int) (dto.FindShelfLifeDTO, error) {
	model, err := svc.repo.FindByID(ctx, id)
	result := translate.ShelfLifeModelToFindDTO(&model)
	if err != nil {
		return dto.FindShelfLifeDTO{}, err
	}
	return result, nil
}

// FindShelfLifes implements ShelfLifeServicer
func (svc *shelfLifeSerivce) FindShelfLifes(ctx context.Context, filter *dto.ShelfLifeFilterDTO) ([]dto.FindShelfLifeDTO, error) {
	models, err := svc.repo.FindMany(ctx, filter.Limit, filter.Offset)
	dtos := translate.ShelfLifeModelsToFindDTOs(models)
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
func (svc *shelfLifeSerivce) UpdateShelfLife(ctx context.Context, id int, payload *dto.UpdateShelfLifeDTO) error {
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
