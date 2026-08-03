package brands_service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/celio001/product-command/internal/brands"
	brandsRepo "github.com/celio001/product-command/internal/brands/repository"
	"github.com/celio001/product-command/pkg/logger"
	"go.uber.org/zap"
)

var (
	ErrCreateBrand = errors.New("error create brand")
)

type brandSvc struct {
	repo brandsRepo.BrandsRepoInterface
}

type BrandSvcInterface interface {
	CreateBrandSvc(ctx context.Context, brand brands.Brand) (brands.Brand, error)
	SoftDeleteBrandSvc(ctx context.Context, id uuid.UUID) error
}

func NewBrandSvc(repo brandsRepo.BrandsRepoInterface) BrandSvcInterface {
	return &brandSvc{repo: repo}
}

func (b *brandSvc) CreateBrandSvc(ctx context.Context, brand brands.Brand) (brands.Brand, error) {

	logger.Info("initiating brand creation",
        zap.String("brand.name", brand.Name),
        zap.String("event.action", "create_brand_start"),
    )

	bCreated, err := b.repo.CreateBrand(ctx, brand)
	if err != nil {
		logger.Error("failed to create brand",
			zap.String("error.message", err.Error()),
			zap.String("error.code", "ERROR_CREATE_BRAND"),
			zap.String("brand.name", brand.Name),
		)
		return brands.Brand{}, err
	}
	logger.Info("brand created successfully",
		zap.String("brand.name", brand.Name),
		zap.String("brand.id", bCreated.ID.String()),
		zap.String("event.action", "create_brand_success"),)

	return bCreated, nil
}

func (b *brandSvc) SoftDeleteBrandSvc(ctx context.Context, id uuid.UUID) error {
	brand, err := b.repo.GetBrandByID(ctx, id)
	if err != nil {
		if errors.Is(err, brandsRepo.ErrBrandNotFound) {
			logger.Error("error get brand not found", zap.String("error", err.Error()))
			return err
		} else {
			logger.Error("error get brand to delete", zap.String("error", err.Error()))
			return ErrCreateBrand
		}
	}
	err = b.repo.SoftDeleteBrand(ctx, brand.ID)
	if err != nil {
		if errors.Is(err, brandsRepo.ErrBrandNotFound) {
			logger.Error("error delete brand not found", zap.String("error", err.Error()))
			return err
		} else {
			logger.Error("error delete brand", zap.String("error", err.Error()))
			return err
		}
	}
	return nil
}
