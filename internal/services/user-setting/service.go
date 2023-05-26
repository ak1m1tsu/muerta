package usersetting

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/api/router/params"
	"github.com/romankravchuk/muerta/internal/services/utils"
	"github.com/romankravchuk/muerta/internal/storage/postgres/models"
	"github.com/romankravchuk/muerta/internal/storage/postgres/setting"
)

type UserSettingsServicer interface {
	FindSettingByID(ctx context.Context, id int) (params.FindSetting, error)
	FindSettings(ctx context.Context, filter *params.SettingFilter) ([]params.FindSetting, error)
	CreateSetting(ctx context.Context, setting *params.CreateSetting) error
	UpdateSetting(ctx context.Context, id int, setting *params.UpdateSetting) error
	DeleteSetting(ctx context.Context, id int) error
	RestoreSetting(ctx context.Context, id int) (params.FindSetting, error)
	Count(ctx context.Context, filter params.SettingFilter) (int, error)
}

type userSettingsService struct {
	repo setting.SettingsRepositorer
}

// Count implements UserSettingsServicer
func (s *userSettingsService) Count(ctx context.Context, filter params.SettingFilter) (int, error) {
	return s.repo.Count(ctx, models.SettingFilter{Name: filter.Name})
}

func New(repo setting.SettingsRepositorer) UserSettingsServicer {
	return &userSettingsService{repo: repo}
}

func (s *userSettingsService) FindSettingByID(
	ctx context.Context,
	id int,
) (params.FindSetting, error) {
	model, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return params.FindSetting{}, err
	}
	dto := utils.SettingModelToFind(&model)
	return dto, nil
}

func (s *userSettingsService) FindSettings(
	ctx context.Context,
	filter *params.SettingFilter,
) ([]params.FindSetting, error) {
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
	dtos := utils.SettingModelsToFinds(models)
	return dtos, nil
}

func (s *userSettingsService) CreateSetting(
	ctx context.Context,
	payload *params.CreateSetting,
) error {
	model := utils.CreateSettingToModel(payload)
	err := s.repo.Create(ctx, model)
	return err
}

func (s *userSettingsService) UpdateSetting(
	ctx context.Context,
	id int,
	payload *params.UpdateSetting,
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
) (params.FindSetting, error) {
	model, err := s.repo.Restore(ctx, id)
	if err != nil {
		return params.FindSetting{}, fmt.Errorf("restore setting error: %w", err)
	}
	return utils.SettingModelToFind(&model), nil
}
