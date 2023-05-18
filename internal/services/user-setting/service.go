package usersetting

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/pkg/translate"
	"github.com/romankravchuk/muerta/internal/repositories/models"
	"github.com/romankravchuk/muerta/internal/repositories/setting"
)

type UserSettingsServicer interface {
	FindSettingByID(ctx context.Context, id int) (dto.FindSetting, error)
	FindSettings(ctx context.Context, filter *dto.SettingFilter) ([]dto.FindSetting, error)
	CreateSetting(ctx context.Context, setting *dto.CreateSetting) error
	UpdateSetting(ctx context.Context, id int, setting *dto.UpdateSetting) error
	DeleteSetting(ctx context.Context, id int) error
	RestoreSetting(ctx context.Context, id int) (dto.FindSetting, error)
	Count(ctx context.Context, filter dto.SettingFilter) (int, error)
}

type userSettingsService struct {
	repo setting.SettingsRepositorer
}

// Count implements UserSettingsServicer
func (s *userSettingsService) Count(ctx context.Context, filter dto.SettingFilter) (int, error) {
	return s.repo.Count(ctx, models.SettingFilter{Name: filter.Name})
}

func New(repo setting.SettingsRepositorer) UserSettingsServicer {
	return &userSettingsService{repo: repo}
}

func (s *userSettingsService) FindSettingByID(
	ctx context.Context,
	id int,
) (dto.FindSetting, error) {
	model, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return dto.FindSetting{}, err
	}
	dto := translate.SettingModelToFind(&model)
	return dto, nil
}

func (s *userSettingsService) FindSettings(
	ctx context.Context,
	filter *dto.SettingFilter,
) ([]dto.FindSetting, error) {
	models, err := s.repo.FindMany(ctx, models.SettingFilter{
		PageFilter: models.PageFilter{
			Limit:  filter.Limit,
			Offset: filter.Offset,
		},
		Name: filter.Name,
	})
	if err != nil {
		return nil, err
	}
	dtos := translate.SettingModelsToFinds(models)
	return dtos, nil
}

func (s *userSettingsService) CreateSetting(
	ctx context.Context,
	payload *dto.CreateSetting,
) error {
	model := translate.CreateSettingToModel(payload)
	err := s.repo.Create(ctx, model)
	return err
}

func (s *userSettingsService) UpdateSetting(
	ctx context.Context,
	id int,
	payload *dto.UpdateSetting,
) error {
	model, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("setting not found: %w", err)
	}
	if payload.Name != "" {
		model.Name = payload.Name
	}
	if payload.CategoryID != 0 {
		model.Category.ID = payload.CategoryID
	}
	if err := s.repo.Update(ctx, model); err != nil {
		return fmt.Errorf("update setting error: %w", err)
	}
	return nil
}

func (s *userSettingsService) DeleteSetting(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("delete setting error: %w", err)
	}
	return nil
}

func (s *userSettingsService) RestoreSetting(
	ctx context.Context,
	id int,
) (dto.FindSetting, error) {
	model, err := s.repo.Restore(ctx, id)
	if err != nil {
		return dto.FindSetting{}, fmt.Errorf("restore setting error: %w", err)
	}
	return translate.SettingModelToFind(&model), nil
}
