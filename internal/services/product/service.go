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
	CreateCategory(ctx context.Context, productId int, categoryId int) (dto.FindProductCategoryDTO, error)
	DeleteCategory(ctx context.Context, productId int, categoryId int) error
	FindProductTips(ctx context.Context, id int) ([]dto.FindTipDTO, error)
	CreateProductTip(ctx context.Context, productID, tipID int) (dto.FindTipDTO, error)
	DeleteProductTip(ctx context.Context, productID, tipID int) error
	services.Counter
}

type productService struct {
	repo repo.ProductRepositorer
}

// CreateProductTip implements ProductServicer
func (s *productService) CreateProductTip(ctx context.Context, productID int, tipID int) (dto.FindTipDTO, error) {
	model, err := s.repo.CreateTip(ctx, productID, tipID)
	if err != nil {
		return dto.FindTipDTO{}, fmt.Errorf("error adding product tip: %w", err)
	}
	return translate.TipModelToFindDTO(&model), nil
}

// DeleteProductTip implements ProductServicer
func (s *productService) DeleteProductTip(ctx context.Context, productID int, tipID int) error {
	err := s.repo.DeleteTip(ctx, productID, tipID)
	if err != nil {
		return fmt.Errorf("error removing product tip: %w", err)
	}
	return nil
}

// FindProductTips implements ProductServicer
func (s *productService) FindProductTips(ctx context.Context, id int) ([]dto.FindTipDTO, error) {
	result, err := s.repo.FindTips(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error finding product tips: %w", err)
	}
	return translate.TipModelsToFindDTOs(result), nil
}

func (s *productService) Count(ctx context.Context) (int, error) {
	count, err := s.repo.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("error counting products: %w", err)
	}
	return count, nil
}

// CreateCategory implements ProductServicer
func (s *productService) CreateCategory(ctx context.Context, productId int, categoryId int) (dto.FindProductCategoryDTO, error) {
	model, err := s.repo.CreateCategory(ctx, productId, categoryId)
	if err != nil {
		return dto.FindProductCategoryDTO{}, fmt.Errorf("error adding product category: %w", err)
	}
	return translate.ProductCategoryModelToFindDTO(&model), nil
}

// DeleteCategory implements ProductServicer
func (s *productService) DeleteCategory(ctx context.Context, productId int, categoryId int) error {
	if err := s.repo.DeleteCategory(ctx, productId, categoryId); err != nil {
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
