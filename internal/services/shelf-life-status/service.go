package shelflifestatus

import (
	"context"

	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/pkg/translate"
	repository "github.com/romankravchuk/muerta/internal/repositories/shelf-life-status"
)

type ShelfLifeStatusServicer interface {
	FindShelfLifeStatusByID(ctx context.Context, id int) (dto.FindShelfLifeStatusDTO, error)
	FindShelfLifeStatuss(ctx context.Context, filter *dto.ShelfLifeStatusFilterDTO) ([]dto.FindShelfLifeStatusDTO, error)
	CreateShelfLifeStatus(ctx context.Context, payload *dto.CreateShelfLifeStatusDTO) error
	UpdateShelfLifeStatus(ctx context.Context, id int, payload *dto.UpdateShelfLifeStatusDTO) error
	DeleteShelfLifeStatus(ctx context.Context, id int) error
}

type shelfLifeStatusService struct {
	repo repository.ShelfLifeStatusRepositorer
}

// CreateShelfLifeStatus implements ShelfLifeStatusServicer
func (svc *shelfLifeStatusService) CreateShelfLifeStatus(ctx context.Context, payload *dto.CreateShelfLifeStatusDTO) error {
	model := translate.CreateShelfLifeStatusDTOToModel(payload)
	if err := svc.repo.Create(ctx, model); err != nil {
		return err
	}
	return nil
}

// DeleteShelfLifeStatus implements ShelfLifeStatusServicer
func (svc *shelfLifeStatusService) DeleteShelfLifeStatus(ctx context.Context, id int) error {
	if err := svc.repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

// FindShelfLifeStatusByID implements ShelfLifeStatusServicer
func (svc *shelfLifeStatusService) FindShelfLifeStatusByID(ctx context.Context, id int) (dto.FindShelfLifeStatusDTO, error) {
	model, err := svc.repo.FindByID(ctx, id)
	result := translate.ShelfLifeStatusModelToFindDTO(&model)
	if err != nil {
		return dto.FindShelfLifeStatusDTO{}, err
	}
	return result, nil
}

// FindShelfLifeStatuss implements ShelfLifeStatusServicer
func (svc *shelfLifeStatusService) FindShelfLifeStatuss(ctx context.Context, filter *dto.ShelfLifeStatusFilterDTO) ([]dto.FindShelfLifeStatusDTO, error) {
	models, err := svc.repo.FindMany(ctx, filter.Limit, filter.Offset, filter.Name)
	dtos := translate.ShelfLifeStatusModelsToFindDTOs(models)
	if err != nil {
		return nil, err
	}
	return dtos, nil
}

// UpdateShelfLifeStatus implements ShelfLifeStatusServicer
func (svc *shelfLifeStatusService) UpdateShelfLifeStatus(ctx context.Context, id int, payload *dto.UpdateShelfLifeStatusDTO) error {
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

func New(repo repository.ShelfLifeStatusRepositorer) ShelfLifeStatusServicer {
	return &shelfLifeStatusService{
		repo: repo,
	}
}