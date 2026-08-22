package product_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/celio001/product-cqrs/worker/internal/modules/consumer"
	"github.com/celio001/product-cqrs/worker/internal/modules/producer-dlq"
	"github.com/celio001/product-cqrs/worker/internal/modules/product"
	product_respository "github.com/celio001/product-cqrs/worker/internal/modules/product/respository"
	"github.com/celio001/product-cqrs/worker/pkg/logger"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type productService struct {
	productRepo   product_respository.ProductRepositoryInterface
	consumerTopic consumer.ConsumerTopicsInterface
	productDlq    producer.ProducerDlqInterface
}

type ProductServiceInterface interface {
	CreateProductSvc(ctx context.Context)
}

func NewProductService(productRepo product_respository.ProductRepositoryInterface, consumerTopic consumer.ConsumerTopicsInterface, productDlq producer.ProducerDlqInterface) ProductServiceInterface {
	return &productService{
		productRepo:   productRepo,
		consumerTopic: consumerTopic,
		productDlq:    productDlq,
	}
}

func (s *productService) CreateProductSvc(ctx context.Context) {
	
	for {

		p, err := s.consumerTopic.ConsumerProductTopic(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				fmt.Println("Goroutine: Stop signal received. Shutting down worker...")
				return
			}

			logger.Info("error connect topic create product",
				zap.String("error", err.Error()),
				zap.String("event.action", "ERROR_CONNECT_TOPIC_CREATE_PRODUCT"))
			continue
		}

		err = s.createProductWithRetry(ctx, p, 3)
		if err != nil {
			err = s.productDlq.PublishProductDlq(ctx, p, err)
			if err != nil {
				logger.Error("error publish message to DLQ",
					zap.String("error", err.Error()),
					zap.String("event.action", "ERROR_PUBLISH_MESSAGE_DLQ"))
					continue
			}
			err = s.consumerTopic.CommitProductTopic(ctx, p)
			if err != nil {
			logger.Error("error commit message",
				zap.String("error", err.Error()),
				zap.String("event.action", "ERROR_COMMIT_MESSAGE"))
			}
			continue
		}

		err = s.consumerTopic.CommitProductTopic(ctx, p)
		if err != nil {
			logger.Error("error commit message",
				zap.String("error", err.Error()),
				zap.String("event.action", "ERROR_COMMIT_MESSAGE"))
		}
		
	}
}

func (s *productService) createProductWithRetry(ctx context.Context, p kafka.Message, retries int) error {

	var err error
	for i := 1; i <= retries; i++ {

		var prod product.Product

		err = json.Unmarshal(p.Value, &prod)
		if err != nil {
			logger.Info("error connect topic create product",
				zap.String("error", err.Error()),
				zap.String("event.action", "ERROR_CONNECT_TOPIC_CREATE_PRODUCT"))
			return err
		}

		err = s.productRepo.CreateProductRepository(ctx, prod)
		if err == nil {
			logger.Info("product created successfully",
				zap.String("product.name", prod.Name),
				zap.String("event.action", "PRODUCT_CREATED_SUCCESSFULLY"))
			return nil
		}

		logger.Error("error create product",
				zap.String("error", err.Error()),
				zap.String("event.action", "ERROR_CREATE_PRODUCT"))


		time.Sleep(time.Duration(i) * 500 * time.Millisecond)
	}
	return err
}

