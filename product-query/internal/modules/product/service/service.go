package product_service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/celio001/product-cqrs/product-query/internal/modules/product"
	product_repository "github.com/celio001/product-cqrs/product-query/internal/modules/product/repository"
	"github.com/celio001/product-cqrs/product-query/pkg/logger"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type productService struct {
	productRepo product_repository.ProductRepositoryInterface
	rd          *redis.Client
}

type ProductServiceInterface interface {
	GetProductByIDSvc(ctx context.Context, id uuid.UUID) (product.Product, error)
}

func NewProductService(productRepo product_repository.ProductRepositoryInterface, rd *redis.Client) ProductServiceInterface {
	return &productService{productRepo: productRepo, rd: rd}
}

func (s *productService) GetProductByIDSvc(ctx context.Context, id uuid.UUID) (product.Product, error) {

	val, err := s.rd.Get(ctx, id.String()).Result()
	switch {
	case err == redis.Nil:
		logger.Info("product not found in cache",
			zap.String("event.action", "PRODUCT_NOT_FOUND_IN_CACHE"),
			zap.String("product.id", id.String()),
		)
	case err != nil:
		logger.Error("error get product in cache",
			zap.String("error.type", "RedisError"),
			zap.String("error.message", err.Error()),
			zap.String("error.code", "REDIS_ERROR"),
		)
	default:
		var p product.Product
		err := json.Unmarshal([]byte(val), &p)
		if err != nil {
			logger.Error("error unmarshal product in cache",
				zap.String("error.type", "UnmarshalError"),
				zap.String("error.message", err.Error()),
				zap.String("error.code", "UNMARSHAL_ERROR"),
			)
		} else {
			logger.Info("product found in cache",
				zap.String("event.action", "PRODUCT_FOUND_IN_CACHE"),
				zap.String("product.id", id.String()),
			)
			return p, nil
		}
	}

	logger.Info("product not found in cache",
		zap.String("event.action", "PRODUCT_NOT_FOUND_IN_CACHE"),
		zap.String("product.id", id.String()),
	)

	p, err := s.productRepo.GetProductByID(ctx, id)
	if err != nil {
		logger.Error("error get product in repository",
			zap.String("error.type", "RepositoryError"),
			zap.String("error.message", err.Error()),
			zap.String("error.code", "REPOSITORY_ERROR"),
		)
		return product.Product{}, err
	}

	cacheValue, err := json.Marshal(p)
	if err != nil {
		logger.Error("error marshal product for cache",
			zap.String("error.type", "MarshalError"),
			zap.String("error.message", err.Error()),
			zap.String("error.code", "MARSHAL_ERROR"),
		)
		return p, nil
	}

	err = s.rd.Set(ctx, id.String(), cacheValue, 5*time.Minute).Err()
	if err != nil {
		logger.Error("error create cache product",
			zap.String("error.type", "RedisError"),
			zap.String("error.message", err.Error()),
			zap.String("error.code", "REDIS_ERROR"),
		)
	}
	return p, nil
}
