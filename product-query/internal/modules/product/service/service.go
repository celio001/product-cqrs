package product_service

import (
	"context"

	"github.com/celio001/product-cqrs/product-query/internal/modules/product"
	product_repository "github.com/celio001/product-cqrs/product-query/internal/modules/product/repository"
	"github.com/google/uuid"
)

type productService struct {
	productRepo product_repository.ProductRepositoryInterface
}

type ProductServiceInterface interface {
	GetProductByIDSvc(ctx context.Context, id uuid.UUID) (product.Product, error)
}

func NewProductService(productRepo product_repository.ProductRepositoryInterface) ProductServiceInterface {
	return &productService{productRepo: productRepo}
}

func (s *productService) GetProductByIDSvc(ctx context.Context, id uuid.UUID) (product.Product, error) {
	return s.productRepo.GetProductByID(ctx, id)
}
