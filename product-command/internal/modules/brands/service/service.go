package brands_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/celio001/product-command/internal/modules/brands"
	brandsRepo "github.com/celio001/product-command/internal/modules/brands/repository"
	"github.com/celio001/product-command/internal/modules/producer"
	"github.com/celio001/product-command/pkg/logger"
	"go.uber.org/zap"
)

var (
	ErrCreateBrand = errors.New("error create brand")
)

type brandSvc struct {
	repo      brandsRepo.BrandsRepoInterface
	Kproducer producer.ProducerCommandInterface
}

type BrandSvcInterface interface {
	CreateBrandSvc(ctx context.Context, brand brands.Brand) (brands.Brand, error)
	SoftDeleteBrandSvc(ctx context.Context, id uuid.UUID) error
}

func NewBrandSvc(repo brandsRepo.BrandsRepoInterface, Kproducer producer.ProducerCommandInterface) BrandSvcInterface {
	return &brandSvc{repo: repo, Kproducer: Kproducer}
}

func (b *brandSvc) CreateBrandSvc(ctx context.Context, brand brands.Brand) (brands.Brand, error) {

	tx, err := b.repo.BeginTx(ctx)
	if err != nil {
		logger.Error("failed init transection",
			zap.String("error.message", err.Error()),
			zap.String("error.code", "ERROR_INIT_TRANSECTION"),
		)
		return brands.Brand{}, err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil{
				err = fmt.Errorf("%w: rollback error: %v", err, rbErr)
			}
			return
		}
		err = tx.Commit(ctx)
	}()

	TxRepo := b.repo.WithTx(tx)

	bCreated, err := TxRepo.CreateBrand(ctx, brand)
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
		zap.String("event.action", "create_brand_success"))
	
	if err = b.Kproducer.PublishBrandCreated(ctx, bCreated); err != nil {
		return brands.Brand{}, err
	}

	return bCreated, nil
}

func (b *brandSvc) SoftDeleteBrandSvc(ctx context.Context, id uuid.UUID) error {
	brand, err := b.repo.GetBrandByID(ctx, id)
	if err != nil {
		if errors.Is(err, brandsRepo.ErrBrandNotFound) {
			logger.Error("brand not found to delete",
				zap.String("error.message", err.Error()),
				zap.String("error.code", "ERROR_GET_BRAND"),
			)
			return err
		} else {
			logger.Error("error get brand to delete",
				zap.String("error.message", err.Error()),
				zap.String("error.code", "INTERNAL_ERROR_GET_BRAND"),
				zap.String("brand.id", id.String()),
			)
			return ErrCreateBrand
		}
	}
	err = b.repo.SoftDeleteBrand(ctx, brand.ID)
	if err != nil {
		if errors.Is(err, brandsRepo.ErrBrandNotFound) {
			logger.Error("error delete brand not found",
				zap.String("error.message", err.Error()),
				zap.String("error.code", "ERROR_UUID_BRAND_NOT_FOUND_TO_DELETE"),
				zap.String("brand.id", id.String()),
			)
			return err
		} else {
			logger.Error("internal error delete brand",
				zap.String("error.message", err.Error()),
				zap.String("error.code", "ERROR_INTERNAL_DELETE_BRAND"),
				zap.String("brand.id", id.String()),
			)
			return err
		}
	}
	return nil
}
