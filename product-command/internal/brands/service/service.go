package service

import (
	"context"

	"github.com/celio001/product-command/internal/brands"
	"github.com/celio001/product-command/internal/brands/repository"
	"github.com/celio001/product-command/pkg/logger"
	"go.uber.org/zap"
)

type brandSvc struct {
	repo repository.BrandsRepoInterface
}

type BrandSvcInterface interface {
	CreateBrandSvc(ctx context.Context, brand brands.Brand) (brands.Brand, error)
}

func NewBrandSvc(repo repository.BrandsRepoInterface) BrandSvcInterface {
	return &brandSvc{repo: repo}
}

func (b *brandSvc) CreateBrandSvc(ctx context.Context, brand brands.Brand) (brands.Brand, error) {
	bCreated, err := b.repo.CreateBrand(ctx, brand)
	if err != nil {
		logger.Error("error create brand", zap.String("error", err.Error()))
		return brands.Brand{}, err
	}

	return bCreated, nil
}
