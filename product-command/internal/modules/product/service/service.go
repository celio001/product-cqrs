package product_service

import (
	"context"
	"errors"
	"fmt"

	product_dto "github.com/celio001/product-command/internal/fiber/v1/product/dto"
	brands_repository "github.com/celio001/product-command/internal/modules/brands/repository"
	categories_repository "github.com/celio001/product-command/internal/modules/categories/repository"
	"github.com/celio001/product-command/internal/modules/fiscal"
	fiscal_repository "github.com/celio001/product-command/internal/modules/fiscal/repository"
	"github.com/celio001/product-command/internal/modules/inventory"
	inventory_repository "github.com/celio001/product-command/internal/modules/inventory/repository"
	"github.com/celio001/product-command/internal/modules/product"
	product_repository "github.com/celio001/product-command/internal/modules/product/repository"
	"github.com/celio001/product-command/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type productSvc struct {
	fiscalRepo     fiscal_repository.FiscalRepositoryInterface
	productRepo    product_repository.ProductRepoInterface
	InventoryRep   inventory_repository.InventoryRepoInterface
	CategoriesRepo categories_repository.CategoriesInterface
	BrandRepo      brands_repository.BrandsRepoInterface
}

type ProductSvcInterface interface {
	CreateProductSvc(ctx context.Context, p product.Product, i inventory.Inventory, f fiscal.FiscalData) (resp product_dto.CreateProductResponse, err error)
	SoftDeleteProductSvc(ctx context.Context, id uuid.UUID) error
}

func NewProductSvc(productRepo product_repository.ProductRepoInterface, fiscalRepo fiscal_repository.FiscalRepositoryInterface, InventoryRep inventory_repository.InventoryRepoInterface, CategoriesRepo categories_repository.CategoriesInterface, BrandRepo brands_repository.BrandsRepoInterface) ProductSvcInterface {
	return &productSvc{
		fiscalRepo:     fiscalRepo,
		productRepo:    productRepo,
		InventoryRep:   InventoryRep,
		CategoriesRepo: CategoriesRepo,
		BrandRepo:      BrandRepo,
	}
}

func (s *productSvc) CreateProductSvc(ctx context.Context, p product.Product, i inventory.Inventory, f fiscal.FiscalData) (resp product_dto.CreateProductResponse, err error) {

	c, err := s.CategoriesRepo.GetCategoryByID(ctx, p.CategoryID)
	if err != nil {
		return product_dto.CreateProductResponse{}, err
	}
	p.CategoryID = c.ID

	b, err := s.BrandRepo.GetBrandByID(ctx, p.BrandID)
	if err != nil {
		return product_dto.CreateProductResponse{}, err
	}
	p.BrandID = b.ID

	tx, err := s.productRepo.BeginTx(ctx)
	if err != nil {
		return product_dto.CreateProductResponse{}, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				err = fmt.Errorf("%w: rollback error: %v", err, rbErr)
			}
			return
		}
		err = tx.Commit(ctx)
	}()

	productRepoTx := s.productRepo.WithTx(tx)
	inventoryRepoTx := s.InventoryRep.WithTx(tx)
	fiscaRepoTx := s.fiscalRepo.WithTx(tx)

	product, err := productRepoTx.CreateProductRepo(ctx, p)
	if err != nil {
		return product_dto.CreateProductResponse{}, err
	}

	i.ProductID = product.ID
	inventory, err := inventoryRepoTx.CreateInventoryRepo(ctx, i)
	if err != nil {
		return product_dto.CreateProductResponse{}, err
	}

	f.ProductId = product.ID
	fiscal, err := fiscaRepoTx.CreateFiscalData(ctx, f)
	if err != nil {
		return product_dto.CreateProductResponse{}, err
	}

	resp = s.productResponse(product, inventory, fiscal)

	return resp, nil
}

func (s *productSvc) SoftDeleteProductSvc(ctx context.Context, id uuid.UUID) error {
	err := s.productRepo.SoftDeleteProduct(ctx, id)
	if err != nil {
		if errors.Is(err, product_repository.ErrProductNotFound) {
			logger.Error("error delete product not found",
				zap.String("error.message", err.Error()),
				zap.String("error.code", "ERROR_UUID_PRODUCT_NOT_FOUND_TO_DELETE"),
				zap.String("product.id", id.String()),
			)
		}
		return err
	}

	return nil
}

func (s *productSvc) productResponse(p product.Product, i inventory.Inventory, f fiscal.FiscalData) product_dto.CreateProductResponse {
	return product_dto.CreateProductResponse{
		ID:                  p.ID,
		BrandID:             p.BrandID,
		CategoryID:          p.CategoryID,
		Name:                p.Name,
		Sku:                 p.Sku,
		BarCodeEan:          p.BarCodeEan,
		ShortDescription:    p.ShortDescription,
		DetailedDescription: p.DetailedDescription,
		UnitOfMeasure:       p.UnitOfMeasure,
		CostPrice:           p.CostPrice,
		SalePrice:           p.SalePrice,
		PromotionalPrice:    p.PromotionalPrice,
		GrossWeight:         p.GrossWeight,
		NetWeight:           p.NetWeight,
		Height:              p.Height,
		Width:               p.Width,
		Length:              p.Length,
		Status:              p.Status,
		CreatedAt:           p.CreatedAt,
		UpdatedAt:           p.UpdatedAt,
		Stock: product_dto.InventoryProductResponse{
			ID:                i.ID,
			ProductID:         p.ID,
			LocationAisle:     i.LocationAisle,
			QuantityAvailable: i.QuantityAvailable,
			MinimumStock:      i.MinimumStock,
			MaximumStock:      i.MaximumStock,
			UpdatedAt:         i.UpdatedAt,
		},
		Fiscal: product_dto.FiscalProductResponse{
			ID:         f.ID,
			ProductID:  p.ID,
			NcmCode:    f.NcmCode,
			CestCode:   f.CestCode,
			OriginCode: f.OriginCode,
			IcmsRate:   f.IcmsRate,
			PisRate:    f.PisRate,
			CofinsRate: f.CofinsRate,
			IpiRate:    f.IpiRate,
			UpdatedAt:  f.UpdatedAt,
		},
	}
}
