package step

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/pkg/translate"
	"github.com/romankravchuk/muerta/internal/storage/postgres/models"
	repository "github.com/romankravchuk/muerta/internal/storage/postgres/step"
)

type StepServicer interface {
	FindSteps(ctx context.Context, filter *dto.StepFilter) ([]dto.FindStep, error)
	CreateStep(ctx context.Context, payload *dto.CreateStep) (dto.FindStep, error)
	FindStep(ctx context.Context, id int) (dto.FindStep, error)
	UpdateStep(ctx context.Context, id int, payload *dto.UpdateStep) (dto.FindStep, error)
	DeleteStep(ctx context.Context, id int) error
	RestoreStep(ctx context.Context, id int) (dto.FindStep, error)
	Count(ctx context.Context, filter dto.StepFilter) (int, error)
}

type stepService struct {
	repo repository.StepRepositorer
}

// CreateStep implements StepServicer
func (s *stepService) CreateStep(
	ctx context.Context,
	payload *dto.CreateStep,
) (dto.FindStep, error) {
	model := translate.CreateToStepModel(payload)
	if err := s.repo.Create(ctx, &model); err != nil {
		return dto.FindStep{}, fmt.Errorf("error creating step: %w", err)
	}
	return translate.StepModelToFind(model), nil
}

// DeleteStep implements StepServicer
func (s *stepService) DeleteStep(ctx context.Context, id int) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("error deleting step: %w", err)
	}
	return nil
}

// FindStep implements StepServicer
func (s *stepService) FindStep(ctx context.Context, id int) (dto.FindStep, error) {
	model, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return dto.FindStep{}, fmt.Errorf("error finding step: %w", err)
	}
	return translate.StepModelToFind(model), nil
}

// FindSteps implements StepServicer
func (s *stepService) FindSteps(
	ctx context.Context,
	filter *dto.StepFilter,
) ([]dto.FindStep, error) {
	model, err := s.repo.FindMany(ctx, models.StepFilter{
		PageFilter: models.PageFilter{
			Limit:  filter.Limit,
			Offset: filter.Offset,
		},
		Name: filter.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("error finding steps: %w", err)
	}
	return translate.StepModelsToFinds(model), nil
}

// RestoreStep implements StepServicer
func (s *stepService) RestoreStep(ctx context.Context, id int) (dto.FindStep, error) {
	model, err := s.repo.Restore(ctx, id)
	if err != nil {
		return dto.FindStep{}, fmt.Errorf("error restoring step: %w", err)
	}
	return translate.StepModelToFind(model), nil
}

func (s *stepService) Count(ctx context.Context, filter dto.StepFilter) (int, error) {
	count, err := s.repo.Count(ctx, models.StepFilter{Name: filter.Name})
	if err != nil {
		return 0, fmt.Errorf("error counting steps: %w", err)
	}
	return count, nil
}

// UpdateStep implements StepServicer
func (s *stepService) UpdateStep(
	ctx context.Context,
	id int,
	payload *dto.UpdateStep,
) (dto.FindStep, error) {
	model, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return dto.FindStep{}, fmt.Errorf("error finding step: %w", err)
	}
	if payload.Name != "" {
		model.Name = payload.Name
	}
	if err := s.repo.Update(ctx, id, model); err != nil {
		return dto.FindStep{}, fmt.Errorf("error updating step: %w", err)
	}
	return translate.StepModelToFind(model), nil
}

func New(repo repository.StepRepositorer) StepServicer {
	return &stepService{
		repo: repo,
	}
}
