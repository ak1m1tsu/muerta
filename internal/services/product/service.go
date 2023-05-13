package product

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/pkg/translate"
	repo "github.com/romankravchuk/muerta/internal/repositories/product"
	"github.com/romankravchuk/muerta/internal/services"
)

type ProductServicer interface {
	FindProductByID(ctx context.Context, id int) (dto.FindProductDTO, error)
	FindProducts(ctx context.Context, filter *dto.ProductFilterDTO) ([]dto.FindProductDTO, error)
	CreateProduct(ctx context.Context, payload *dto.CreateProductDTO) error
	UpdateProduct(ctx context.Context, id int, payload *dto.UpdateProductDTO) error
	DeleteProduct(ctx context.Context, id int) error
	RestoreProduct(ctx context.Context, id int) error
	FindProductCategories(ctx context.Context, id int) ([]dto.FindProductCategoryDTO, error)
	FindProductRecipes(ctx context.Context, id int) ([]dto.FindRecipeDTO, error)
	AddProductCategory(ctx context.Context, productId int, categoryId int) (dto.FindProductCategoryDTO, error)
	RemoveProductCategory(ctx context.Context, productId int, categoryId int) error
	services.Counter
}

type productService struct {
	repo repo.ProductRepositorer
}

func (s *productService) Count(ctx context.Context) (int, error) {
	count, err := s.repo.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("error counting products: %w", err)
	}
	return count, nil
}

// AddProductCategory implements ProductServicer
func (s *productService) AddProductCategory(ctx context.Context, productId int, categoryId int) (dto.FindProductCategoryDTO, error) {
	model, err := s.repo.AddProductCategory(ctx, productId, categoryId)
	if err != nil {
		return dto.FindProductCategoryDTO{}, fmt.Errorf("error adding product category: %w", err)
	}
	return translate.ProductCategoryModelToFindDTO(&model), nil
}

// RemoveProductCategory implements ProductServicer
func (s *productService) RemoveProductCategory(ctx context.Context, productId int, categoryId int) error {
	if err := s.repo.RemoveProductCategory(ctx, productId, categoryId); err != nil {
		return fmt.Errorf("error removing product category: %w", err)
	}
	return nil
}

func New(repo repo.ProductRepositorer) ProductServicer {
	return &productService{
		repo: repo,
	}
}

func (svc *productService) FindProductByID(ctx context.Context, id int) (dto.FindProductDTO, error) {
	model, err := svc.repo.FindByID(ctx, id)
	result := translate.ProductModelToFindDTO(&model)
	if err != nil {
		return dto.FindProductDTO{}, err
	}
	return result, nil
}

func (svc *productService) FindProducts(ctx context.Context, filter *dto.ProductFilterDTO) ([]dto.FindProductDTO, error) {
	models, err := svc.repo.FindMany(ctx, filter.Limit, filter.Offset, filter.Name)
	dtos := translate.ProductModelsToFindDTOs(models)
	if err != nil {
		return nil, err
	}
	return dtos, nil
}

func (svc *productService) CreateProduct(ctx context.Context, payload *dto.CreateProductDTO) error {
	model := translate.CreateProductDTOToModel(payload)
	if err := svc.repo.Create(ctx, model); err != nil {
		return err
	}
	return nil
}

func (svc *productService) UpdateProduct(ctx context.Context, id int, payload *dto.UpdateProductDTO) error {
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

func (svc *productService) DeleteProduct(ctx context.Context, id int) error {
	if err := svc.repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

func (svc *productService) RestoreProduct(ctx context.Context, id int) error {
	if err := svc.repo.Restore(ctx, id); err != nil {
		return err
	}
	return nil
}

func (svc *productService) FindProductCategories(ctx context.Context, id int) ([]dto.FindProductCategoryDTO, error) {
	categories, err := svc.repo.FindCategories(ctx, id)
	if err != nil {
		return nil, err
	}
	dtos := translate.CategoryModelsToFindDTOs(categories)
	return dtos, nil
}
func (svc *productService) FindProductRecipes(ctx context.Context, id int) ([]dto.FindRecipeDTO, error) {
	recipes, err := svc.repo.FindRecipes(ctx, id)
	if err != nil {
		return nil, err
	}
	dtos := translate.RecipeModelsToFindDTOs(recipes)
	return dtos, nil
}
