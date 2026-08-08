package categories_service

import (
	"context"

	"github.com/celio001/product-command/internal/modules/categories"
	categoriesRepo "github.com/celio001/product-command/internal/modules/categories/repository"
	"github.com/celio001/product-command/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type categoriesSvc struct {
	categoriesRepo categoriesRepo.CategoriesInterface
}

type CategoriesSvcInterface interface {
	CreateCategorySvc(ctx context.Context, categorie categories.Categories) (categories.Categories, error)
	SoftDeleteCategory(ctx context.Context, uuid uuid.UUID) error
}

func NewCategoriesSvc(categoriesRepo categoriesRepo.CategoriesInterface) CategoriesSvcInterface {
	return &categoriesSvc{categoriesRepo: categoriesRepo}
}

func (c *categoriesSvc) CreateCategorySvc(ctx context.Context, categorie categories.Categories) (categories.Categories, error) {

	category, err := c.categoriesRepo.CreateCategory(ctx, categorie)
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

	return category, nil
}

func (c *categoriesSvc) SoftDeleteCategory(ctx context.Context, uuid uuid.UUID) error{
	err := c.categoriesRepo.SoftDeleteCategory(ctx, uuid)
	if err != nil{
		return err
	}
	return nil
}