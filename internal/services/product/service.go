package product

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/api/router/params"
	"github.com/romankravchuk/muerta/internal/services/utils"
	"github.com/romankravchuk/muerta/internal/storage/postgres/models"
	repo "github.com/romankravchuk/muerta/internal/storage/postgres/product"
)

type ProductServicer interface {
	FindProductByID(ctx context.Context, id int) (params.FindProduct, error)
	FindProducts(ctx context.Context, filter *params.ProductFilter) ([]params.FindProduct, error)
	CreateProduct(ctx context.Context, payload *params.CreateProduct) error
	UpdateProduct(ctx context.Context, id int, payload *params.UpdateProduct) error
	DeleteProduct(ctx context.Context, id int) error
	RestoreProduct(ctx context.Context, id int) error
	FindProductCategories(ctx context.Context, id int) ([]params.FindProductCategory, error)
	FindProductRecipes(ctx context.Context, id int) ([]params.FindRecipe, error)
	CreateCategory(
		ctx context.Context,
		productId int,
		categoryId int,
	) (params.FindProductCategory, error)
	DeleteCategory(ctx context.Context, productId int, categoryId int) error
	FindProductTips(ctx context.Context, id int) ([]params.FindTip, error)
	CreateProductTip(ctx context.Context, productID, tipID int) (params.FindTip, error)
	DeleteProductTip(ctx context.Context, productID, tipID int) error
	Count(ctx context.Context, filter params.ProductFilter) (int, error)
}

type productService struct {
	repo repo.ProductRepositorer
}

// CreateProductTip implements ProductServicer
func (s *productService) CreateProductTip(
	ctx context.Context,
	productID int,
	tipID int,
) (params.FindTip, error) {
	model, err := s.repo.CreateTip(ctx, productID, tipID)
	if err != nil {
		return params.FindTip{}, fmt.Errorf("error adding product tip: %w", err)
	}
	return utils.TipModelToFind(&model), nil
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
func (s *productService) FindProductTips(ctx context.Context, id int) ([]params.FindTip, error) {
	result, err := s.repo.FindTips(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error finding product tips: %w", err)
	}
	return utils.TipModelsToFinds(result), nil
}

func (s *productService) Count(ctx context.Context, filter params.ProductFilter) (int, error) {
	count, err := s.repo.Count(ctx, models.ProductFilter{Name: filter.Name})
	if err != nil {
		return 0, fmt.Errorf("error counting products: %w", err)
	}
	return count, nil
}

// CreateCategory implements ProductServicer
func (s *productService) CreateCategory(
	ctx context.Context,
	productId int,
	categoryId int,
) (params.FindProductCategory, error) {
	model, err := s.repo.CreateCategory(ctx, productId, categoryId)
	if err != nil {
		return params.FindProductCategory{}, fmt.Errorf("error adding product category: %w", err)
	}
	return utils.ProductCategoryModelToFind(&model), nil
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

func (svc *productService) FindProductByID(
	ctx context.Context,
	id int,
) (params.FindProduct, error) {
	model, err := svc.repo.FindByID(ctx, id)
	result := utils.ProductModelToFind(&model)
	if err != nil {
		return params.FindProduct{}, err
	}
	return result, nil
}

func (svc *productService) FindProducts(
	ctx context.Context,
	filter *params.ProductFilter,
) ([]params.FindProduct, error) {
	models, err := svc.repo.FindMany(ctx, models.ProductFilter{
		PageFilter: models.PageFilter{
			Limit:  filter.Limit,
			Offset: filter.Offset,
		},
		Name: filter.Name,
	})
	dtos := utils.ProductModelsToFinds(models)
	if err != nil {
		return nil, err
	}
	return dtos, nil
}

func (svc *productService) CreateProduct(ctx context.Context, payload *params.CreateProduct) error {
	model := utils.CreateProductToModel(payload)
	if err := svc.repo.Create(ctx, model); err != nil {
		return err
	}
	return nil
}

func (svc *productService) UpdateProduct(
	ctx context.Context,
	id int,
	payload *params.UpdateProduct,
) error {
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

func (svc *productService) FindProductCategories(
	ctx context.Context,
	id int,
) ([]params.FindProductCategory, error) {
	categories, err := svc.repo.FindCategories(ctx, id)
	if err != nil {
		return nil, err
	}
	dtos := utils.CategoryModelsToFinds(categories)
	return dtos, nil
}

func (svc *productService) FindProductRecipes(
	ctx context.Context,
	id int,
) ([]params.FindRecipe, error) {
	recipes, err := svc.repo.FindRecipes(ctx, id)
	if err != nil {
		return nil, err
	}
	dtos := utils.RecipeModelsToFinds(recipes)
	return dtos, nil
}
