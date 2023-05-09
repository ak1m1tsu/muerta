package product

import (
	"context"

	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/pkg/translate"
	repo "github.com/romankravchuk/muerta/internal/repositories/product"
)

type ProductServicer interface {
	FindProductByID(ctx context.Context, id int) (dto.FindProductDTO, error)
	FindProducts(ctx context.Context, filter *dto.ProductFilterDTO) ([]dto.FindProductDTO, error)
	CreateProduct(ctx context.Context, payload *dto.CreateProductDTO) error
	UpdateProduct(ctx context.Context, id int, payload *dto.UpdateProductDTO) error
	DeleteProduct(ctx context.Context, id int) error
	RestoreProduct(ctx context.Context, id int) error
}

type productService struct {
	repo repo.ProductRepositorer
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
