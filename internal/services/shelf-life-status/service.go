package shelflifestatus

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/api/router/params"
	"github.com/romankravchuk/muerta/internal/services/utils"
	"github.com/romankravchuk/muerta/internal/storage/postgres/models"
	repository "github.com/romankravchuk/muerta/internal/storage/postgres/shelf-life-status"
)

type ShelfLifeStatusServicer interface {
	FindShelfLifeStatusByID(ctx context.Context, id int) (params.FindShelfLifeStatus, error)
	FindShelfLifeStatuss(
		ctx context.Context,
		filter *params.ShelfLifeStatusFilter,
	) ([]params.FindShelfLifeStatus, error)
	CreateShelfLifeStatus(ctx context.Context, payload *params.CreateShelfLifeStatus) error
	UpdateShelfLifeStatus(ctx context.Context, id int, payload *params.UpdateShelfLifeStatus) error
	DeleteShelfLifeStatus(ctx context.Context, id int) error
	Count(ctx context.Context, filter params.ShelfLifeStatusFilter) (int, error)
}

type shelfLifeStatusService struct {
	repo repository.ShelfLifeStatusRepositorer
}

func (s *shelfLifeStatusService) Count(
	ctx context.Context,
	filter params.ShelfLifeStatusFilter,
) (int, error) {
	count, err := s.repo.Count(ctx, models.ShelfLifeStatusFilter{Name: filter.Name})
	if err != nil {
		return 0, fmt.Errorf("error counting shelf life statuses: %w", err)
	}
	return count, nil
}

// CreateShelfLifeStatus implements ShelfLifeStatusServicer
func (svc *shelfLifeStatusService) CreateShelfLifeStatus(
	ctx context.Context,
	payload *params.CreateShelfLifeStatus,
) error {
	model := utils.CreateShelfLifeStatusToModel(payload)
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
func (svc *shelfLifeStatusService) FindShelfLifeStatusByID(
	ctx context.Context,
	id int,
) (params.FindShelfLifeStatus, error) {
	model, err := svc.repo.FindByID(ctx, id)
	result := utils.ShelfLifeStatusModelToFind(&model)
	if err != nil {
		return params.FindShelfLifeStatus{}, err
	}
	return result, nil
}

// FindShelfLifeStatuss implements ShelfLifeStatusServicer
func (svc *shelfLifeStatusService) FindShelfLifeStatuss(
	ctx context.Context,
	filter *params.ShelfLifeStatusFilter,
) ([]params.FindShelfLifeStatus, error) {
	models, err := svc.repo.FindMany(ctx, models.ShelfLifeStatusFilter{
		PageFilter: models.PageFilter{
			Limit:  filter.Limit,
			Offset: filter.Offset,
		},
		Name: filter.Name,
	})
	dtos := utils.ShelfLifeStatusModelsToFinds(models)
	if err != nil {
		return nil, err
	}
	return dtos, nil
}

// UpdateShelfLifeStatus implements ShelfLifeStatusServicer
func (svc *shelfLifeStatusService) UpdateShelfLifeStatus(
	ctx context.Context,
	id int,
	payload *params.UpdateShelfLifeStatus,
) error {
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
