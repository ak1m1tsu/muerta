package usersetting

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/pkg/translate"
	"github.com/romankravchuk/muerta/internal/repositories/setting"
)

type UserSettingsServicer interface {
	FindSettingByID(ctx context.Context, id int) (dto.FindSettingDTO, error)
	FindSettings(ctx context.Context, filter *dto.SettingFilterDTO) ([]dto.FindSettingDTO, error)
	CreateSetting(ctx context.Context, setting *dto.CreateSettingDTO) error
	UpdateSetting(ctx context.Context, id int, setting *dto.UpdateSettingDTO) error
	DeleteSetting(ctx context.Context, id int) error
	RestoreSetting(ctx context.Context, id int) error
}

type userSettingsService struct {
	repo setting.SettingsRepositorer
}

func New(repo setting.SettingsRepositorer) UserSettingsServicer {
	return &userSettingsService{repo: repo}
}

func (s *userSettingsService) FindSettingByID(ctx context.Context, id int) (dto.FindSettingDTO, error) {
	model, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return dto.FindSettingDTO{}, err
	}
	dto := translate.SettingModelToFindDTO(&model)
	return dto, nil
}

func (s *userSettingsService) FindSettings(ctx context.Context, filter *dto.SettingFilterDTO) ([]dto.FindSettingDTO, error) {
	models, err := s.repo.FindMany(ctx, filter.Limit, filter.Offset, filter.Name)
	if err != nil {
		return nil, err
	}
	dtos := translate.SettingModelsToFindDTOs(models)
	return dtos, nil
}

func (s *userSettingsService) CreateSetting(ctx context.Context, payload *dto.CreateSettingDTO) error {
	model := translate.CreateSettingDTOToModel(payload)
	err := s.repo.Create(ctx, model)
	return err
}

func (s *userSettingsService) UpdateSetting(ctx context.Context, id int, payload *dto.UpdateSettingDTO) error {
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
	return nil
}

func (s *userSettingsService) RestoreSetting(ctx context.Context, id int) error {
	return nil
}
