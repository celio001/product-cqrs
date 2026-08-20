package categories_service

import (
	"context"
	"fmt"

	"github.com/celio001/product-command/internal/modules/categories"
	categoriesRepo "github.com/celio001/product-command/internal/modules/categories/repository"
	"github.com/celio001/product-command/internal/modules/producer"
	"github.com/celio001/product-command/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type categoriesSvc struct {
	categoriesRepo categoriesRepo.CategoriesInterface
	Kproducer      producer.ProducerCommandInterface
}

type CategoriesSvcInterface interface {
	CreateCategorySvc(ctx context.Context, categorie categories.Categories) (categories.Categories, error)
	SoftDeleteCategory(ctx context.Context, uuid uuid.UUID) error
}

func NewCategoriesSvc(categoriesRepo categoriesRepo.CategoriesInterface, Kproducer producer.ProducerCommandInterface) CategoriesSvcInterface {
	return &categoriesSvc{
		categoriesRepo: categoriesRepo,
		Kproducer:      Kproducer,
	}
}

func (c *categoriesSvc) CreateCategorySvc(ctx context.Context, categorie categories.Categories) (categories.Categories, error) {

	tx, err := c.categoriesRepo.BeginTx(ctx)
	if err != nil {
		logger.Error("failed init transection",
			zap.String("error.message", err.Error()),
			zap.String("error.code", "ERROR_INIT_TRANSECTION"),
		)
		return categories.Categories{}, err
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

	TxRepo := c.categoriesRepo.WithTx(tx)

	category, err := TxRepo.CreateCategory(ctx, categorie)
	if err != nil {
		logger.Error("failed to create category",
			zap.String("error.message", err.Error()),
			zap.String("error.code", "ERROR_CREATE_CATEGORY"),
			zap.String("category.name", category.Name),
		)
		return categories.Categories{}, err
	}
	logger.Info("category created successfully",
		zap.String("category.name", category.Name),
		zap.String("category.id", category.ID.String()),
		zap.String("event.action", "CATEGORY_CREATED_SUCCESS"))

	err = c.Kproducer.PublishCategoryCreated(ctx, category)
	if err != nil {
		logger.Error("failed to publish the message category created",
			zap.String("error.message", err.Error()),
			zap.String("error.code", "ERROR_PUBLISH_CREATE_CATEGORY"),
		)
		return categories.Categories{}, err
	}

	return category, nil
}

func (c *categoriesSvc) SoftDeleteCategory(ctx context.Context, uuid uuid.UUID) error {
	err := c.categoriesRepo.SoftDeleteCategory(ctx, uuid)
	if err != nil {
		return err
	}
	return nil
}
