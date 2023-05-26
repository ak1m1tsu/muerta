package measure

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/api/router/params"
	"github.com/romankravchuk/muerta/internal/services/utils"
	repository "github.com/romankravchuk/muerta/internal/storage/postgres/measure"
	"github.com/romankravchuk/muerta/internal/storage/postgres/models"
)

type MeasureServicer interface {
	FindMeasureByID(ctx context.Context, id int) (params.FindMeasure, error)
	FindMeasures(ctx context.Context, filter *params.MeasureFilter) ([]params.FindMeasure, error)
	CreateMeasure(ctx context.Context, payload *params.CreateMeasure) error
	UpdateMeasure(ctx context.Context, id int, payload *params.UpdateMeasure) error
	DeleteMeasure(ctx context.Context, id int) error
	Count(ctx context.Context, filter params.MeasureFilter) (int, error)
}

type measureService struct {
	repo repository.MeasureRepositorer
}

func (s *measureService) Count(ctx context.Context, filter params.MeasureFilter) (int, error) {
	count, err := s.repo.Count(ctx, models.MeasureFilter{Name: filter.Name})
	if err != nil {
		return 0, fmt.Errorf("error counting measures: %w", err)
	}
	return count, nil
}

// CreateMeasure implements MeasureServicer
func (svc *measureService) CreateMeasure(ctx context.Context, payload *params.CreateMeasure) error {
	model := utils.CreateMeasureToModel(payload)
	if err := svc.repo.Create(ctx, model); err != nil {
		return err
	}
	return nil
}

// DeleteMeasure implements MeasureServicer
func (svc *measureService) DeleteMeasure(ctx context.Context, id int) error {
	if err := svc.repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

// FindMeasureByID implements MeasureServicer
func (svc *measureService) FindMeasureByID(
	ctx context.Context,
	id int,
) (params.FindMeasure, error) {
	model, err := svc.repo.FindByID(ctx, id)
	result := utils.MeasureModelToFind(&model)
	if err != nil {
		return params.FindMeasure{}, err
	}
	return result, nil
}

// FindMeasures implements MeasureServicer
func (svc *measureService) FindMeasures(
	ctx context.Context,
	filter *params.MeasureFilter,
) ([]params.FindMeasure, error) {
	models, err := svc.repo.FindMany(ctx, models.MeasureFilter{
		PageFilter: models.PageFilter{
			Limit:  filter.Limit,
			Offset: filter.Offset,
		},
		Name: filter.Name,
	})
	dtos := utils.MeasureModelsToFinds(models)
	if err != nil {
		return nil, err
	}
	return dtos, nil
}

// UpdateMeasure implements MeasureServicer
func (svc *measureService) UpdateMeasure(
	ctx context.Context,
	id int,
	payload *params.UpdateMeasure,
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

func New(repo repository.MeasureRepositorer) MeasureServicer {
	return &measureService{
		repo: repo,
	}
}
