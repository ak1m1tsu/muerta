package category

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/pkg/translate"
	repository "github.com/romankravchuk/muerta/internal/storage/postgres/category"
	"github.com/romankravchuk/muerta/internal/storage/postgres/models"
)

type CategoryServicer interface {
	FindCategoryByID(ctx context.Context, id int) (dto.FindProductCategory, error)
	FindCategorys(
		ctx context.Context,
		filter *dto.ProductCategoryFilter,
	) ([]dto.FindProductCategory, error)
	CreateCategory(ctx context.Context, payload *dto.CreateProductCategory) error
	UpdateCategory(ctx context.Context, id int, category *dto.UpdateProductCategory) error
	DeleteCategory(ctx context.Context, id int) error
	RestoreCategory(ctx context.Context, id int) error
	Count(ctx context.Context, filter dto.ProductCategoryFilter) (int, error)
}

type categoryService struct {
	repo repository.CategoryRepositorer
}

func (s *categoryService) Count(
	ctx context.Context,
	filter dto.ProductCategoryFilter,
) (int, error) {
	count, err := s.repo.Count(ctx, models.ProductCategoryFilter{Name: filter.Name})
	if err != nil {
		return 0, fmt.Errorf("error counting users: %w", err)
	}
	return count, nil
}

// CreateCategory implements CategoryServicer
func (svc *categoryService) CreateCategory(
	ctx context.Context,
	payload *dto.CreateProductCategory,
) error {
	model := translate.CreateCategoryToModel(payload)
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
func (svc *categoryService) FindCategoryByID(
	ctx context.Context,
	id int,
) (dto.FindProductCategory, error) {
	category, err := svc.repo.FindByID(ctx, id)
	if err != nil {
		return dto.FindProductCategory{}, err
	}
	dto := translate.ProductCategoryModelToFind(&category)
	return dto, nil
}

// FindCategorys implements CategoryServicer
func (svc *categoryService) FindCategorys(
	ctx context.Context,
	filter *dto.ProductCategoryFilter,
) ([]dto.FindProductCategory, error) {
	categories, err := svc.repo.FindMany(ctx, models.ProductCategoryFilter{
		PageFilter: models.PageFilter{
			Limit:  filter.Limit,
			Offset: filter.Offset,
		},
		Name: filter.Name,
	})
	if err != nil {
		return nil, err
	}
	dtos := translate.CategoryModelsToFinds(categories)
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
func (svc *categoryService) UpdateCategory(
	ctx context.Context,
	id int,
	category *dto.UpdateProductCategory,
) error {
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
