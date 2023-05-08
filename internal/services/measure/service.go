package measure

import (
	"context"

	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/pkg/translate"
	repository "github.com/romankravchuk/muerta/internal/repositories/measure"
)

type MeasureServicer interface {
	FindMeasureByID(ctx context.Context, id int) (dto.FindMeasureDTO, error)
	FindMeasures(ctx context.Context, filter *dto.MeasureFilterDTO) ([]dto.FindMeasureDTO, error)
	CreateMeasure(ctx context.Context, payload *dto.CreateMeasureDTO) error
	UpdateMeasure(ctx context.Context, id int, payload *dto.UpdateMeasureDTO) error
	DeleteMeasure(ctx context.Context, id int) error
	RestoreMeasure(ctx context.Context, id int) error
}

type measureService struct {
	repo repository.MeasureRepositorer
}

// CreateMeasure implements MeasureServicer
func (svc *measureService) CreateMeasure(ctx context.Context, payload *dto.CreateMeasureDTO) error {
	model := translate.CreateMeasureDTOToModel(payload)
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
func (svc *measureService) FindMeasureByID(ctx context.Context, id int) (dto.FindMeasureDTO, error) {
	model, err := svc.repo.FindByID(ctx, id)
	result := translate.MeasureModelToFindDTO(&model)
	if err != nil {
		return dto.FindMeasureDTO{}, err
	}
	return result, nil
}

// FindMeasures implements MeasureServicer
func (svc *measureService) FindMeasures(ctx context.Context, filter *dto.MeasureFilterDTO) ([]dto.FindMeasureDTO, error) {
	models, err := svc.repo.FindMany(ctx, filter.Limit, filter.Offset, filter.Name)
	dtos := translate.MeasureModelsToFindDTOs(models)
	if err != nil {
		return nil, err
	}
	return dtos, nil
}

// RestoreMeasure implements MeasureServicer
func (svc *measureService) RestoreMeasure(ctx context.Context, id int) error {
	if err := svc.repo.Restore(ctx, id); err != nil {
		return err
	}
	return nil
}

// UpdateMeasure implements MeasureServicer
func (svc *measureService) UpdateMeasure(ctx context.Context, id int, payload *dto.UpdateMeasureDTO) error {
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
