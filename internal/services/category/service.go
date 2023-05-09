package category

import (
	"context"

	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/pkg/translate"
	repository "github.com/romankravchuk/muerta/internal/repositories/category"
)

type CategoryServicer interface {
	FindCategoryByID(ctx context.Context, id int) (dto.FindCategoryDTO, error)
	FindCategorys(ctx context.Context, filter *dto.CategoryFilterDTO) ([]dto.FindCategoryDTO, error)
	CreateCategory(ctx context.Context, payload *dto.CreateCategoryDTO) error
	UpdateCategory(ctx context.Context, id int, category *dto.UpdateCategoryDTO) error
	DeleteCategory(ctx context.Context, id int) error
	RestoreCategory(ctx context.Context, id int) error
}

type categoryService struct {
	repo repository.CategoryRepositorer
}

// CreateCategory implements CategoryServicer
func (svc *categoryService) CreateCategory(ctx context.Context, payload *dto.CreateCategoryDTO) error {
	model := translate.CreateCategoryDTOToModel(payload)
	if err := svc.repo.Create(ctx, model); err != nil {
		return err
	}
	return nil
}

// DeleteCategory implements CategoryServicer
func (svc *categoryService) DeleteCategory(ctx context.Context, id int) error {
	if err := svc.repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

// FindCategoryByID implements CategoryServicer
func (svc *categoryService) FindCategoryByID(ctx context.Context, id int) (dto.FindCategoryDTO, error) {
	category, err := svc.repo.FindByID(ctx, id)
	if err != nil {
		return dto.FindCategoryDTO{}, err
	}
	dto := translate.CategoryModelToFindDTO(&category)
	return dto, nil
}

// FindCategorys implements CategoryServicer
func (svc *categoryService) FindCategorys(ctx context.Context, filter *dto.CategoryFilterDTO) ([]dto.FindCategoryDTO, error) {
	categories, err := svc.repo.FindMany(ctx, filter.Limit, filter.Offset, filter.Name)
	if err != nil {
		return nil, err
	}
	dtos := translate.CategoryModelsToFindDTOs(categories)
	return dtos, nil
}

// RestoreCategory implements CategoryServicer
func (svc *categoryService) RestoreCategory(ctx context.Context, id int) error {
	if err := svc.repo.Restore(ctx, id); err != nil {
		return err
	}
	return nil
}

// UpdateCategory implements CategoryServicer
func (svc *categoryService) UpdateCategory(ctx context.Context, id int, category *dto.UpdateCategoryDTO) error {
	model, err := svc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if category.Name != "" {
		model.Name = category.Name
	}
	if err := svc.repo.Update(ctx, model); err != nil {
		return err
	}
	return nil
}

func New(repo repository.CategoryRepositorer) CategoryServicer {
	return &categoryService{
		repo: repo,
	}
}
