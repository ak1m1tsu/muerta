package step

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/pkg/translate"
	repository "github.com/romankravchuk/muerta/internal/repositories/step"
	"github.com/romankravchuk/muerta/internal/services"
)

type StepServicer interface {
	FindSteps(ctx context.Context, filter *dto.StepFilterDTO) ([]dto.FindStepDTO, error)
	CreateStep(ctx context.Context, payload *dto.CreateStepDTO) (dto.FindStepDTO, error)
	FindStep(ctx context.Context, id int) (dto.FindStepDTO, error)
	UpdateStep(ctx context.Context, id int, payload *dto.UpdateStepDTO) (dto.FindStepDTO, error)
	DeleteStep(ctx context.Context, id int) error
	RestoreStep(ctx context.Context, id int) (dto.FindStepDTO, error)
	services.Counter
}

type stepService struct {
	repo repository.StepRepositorer
}

// CreateStep implements StepServicer
func (s *stepService) CreateStep(ctx context.Context, payload *dto.CreateStepDTO) (dto.FindStepDTO, error) {
	model := translate.CreateDTOToStepModel(payload)
	if err := s.repo.Create(ctx, &model); err != nil {
		return dto.FindStepDTO{}, fmt.Errorf("error creating step: %w", err)
	}
	return translate.StepModelToFindDTO(model), nil
}

// DeleteStep implements StepServicer
func (s *stepService) DeleteStep(ctx context.Context, id int) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("error deleting step: %w", err)
	}
	return nil
}

// FindStep implements StepServicer
func (s *stepService) FindStep(ctx context.Context, id int) (dto.FindStepDTO, error) {
	model, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return dto.FindStepDTO{}, fmt.Errorf("error finding step: %w", err)
	}
	return translate.StepModelToFindDTO(model), nil
}

// FindSteps implements StepServicer
func (s *stepService) FindSteps(ctx context.Context, filter *dto.StepFilterDTO) ([]dto.FindStepDTO, error) {
	model, err := s.repo.FindMany(ctx, filter.Limit, filter.Offset, filter.Name)
	if err != nil {
		return nil, fmt.Errorf("error finding steps: %w", err)
	}
	return translate.StepModelsToFindDTOs(model), nil
}

// RestoreStep implements StepServicer
func (s *stepService) RestoreStep(ctx context.Context, id int) (dto.FindStepDTO, error) {
	model, err := s.repo.Restore(ctx, id)
	if err != nil {
		return dto.FindStepDTO{}, fmt.Errorf("error restoring step: %w", err)
	}
	return translate.StepModelToFindDTO(model), nil
}

func (s *stepService) Count(ctx context.Context) (int, error) {
	count, err := s.repo.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("error counting steps: %w", err)
	}
	return count, nil
}

// UpdateStep implements StepServicer
func (s *stepService) UpdateStep(ctx context.Context, id int, payload *dto.UpdateStepDTO) (dto.FindStepDTO, error) {
	model, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return dto.FindStepDTO{}, fmt.Errorf("error finding step: %w", err)
	}
	if payload.Name != "" {
		model.Name = payload.Name
	}
	if err := s.repo.Update(ctx, id, model); err != nil {
		return dto.FindStepDTO{}, fmt.Errorf("error updating step: %w", err)
	}
	return translate.StepModelToFindDTO(model), nil
}

func New(repo repository.StepRepositorer) StepServicer {
	return &stepService{
		repo: repo,
	}
}
