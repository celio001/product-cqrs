package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	product_dto "github.com/celio001/product-command/internal/fiber/v1/product/dto"
	"github.com/celio001/product-command/internal/modules/brands"
	"github.com/celio001/product-command/internal/modules/categories"
	"github.com/celio001/product-command/pkg/logger"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type producerCommand struct {
	ProductTopic *kafka.Writer
	BrandTopic   *kafka.Writer
	Category     *kafka.Writer
}

type ProducerCommandInterface interface {
	PublishProductCreated(ctx context.Context, p product_dto.CreateProductResponse) error
	PublishBrandCreated(ctx context.Context, b brands.Brand) error
	PublishCategoryCreated(ctx context.Context, c categories.Categories) error
}

func NewProducerCommand(ProductTopic *kafka.Writer, BrandTopic *kafka.Writer, Category *kafka.Writer) ProducerCommandInterface {
	return &producerCommand{
		ProductTopic: ProductTopic,
		BrandTopic:   BrandTopic,
		Category:     Category,
	}
}

func (k *producerCommand) PublishProductCreated(ctx context.Context, p product_dto.CreateProductResponse) error {

	valueBytes, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("falha ao serializar envelope do pedido: %w", err)
	}

	kafkaHeaders := []kafka.Header{
		{Key: "event_type", Value: []byte("product.created")},
		{Key: "trace_id", Value: []byte("1")},
	}

	err = k.ProductTopic.WriteMessages(ctx, kafka.Message{
		Value:   valueBytes,
		Headers: kafkaHeaders,
		Time:    time.Now(),
	})
	if err != nil {
		logger.Error("failed to publish the message product created",
			zap.String("error.message", err.Error()),
			zap.String("error.code", "ERROR_PUBLISH_CREATE_PRODUCT"),
		)
		return err
	}

	return nil
}

func (k *producerCommand) PublishBrandCreated(ctx context.Context, b brands.Brand) error {
	valueBytes, err := json.Marshal(b)
	if err != nil {
		return fmt.Errorf("falha ao serializar envelope do pedido: %w", err)
	}

	kafkaHeaders := []kafka.Header{
		{Key: "event_type", Value: []byte("brand.created")},
		{Key: "trace_id", Value: []byte(uuid.New().String())},
	}

	err = k.BrandTopic.WriteMessages(ctx, kafka.Message{
		Value:   valueBytes,
		Headers: kafkaHeaders,
		Time:    time.Now(),
	})
	if err != nil {
		logger.Error("failed to publish the message brand created",
			zap.String("error.message", err.Error()),
			zap.String("error.code", "ERROR_PUBLISH_CREATE_PRODUCT"),
		)
		return err
	}

	return nil
}

func (k *producerCommand) PublishCategoryCreated(ctx context.Context, c categories.Categories) error {
	valueBytes, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("falha ao serializar envelope do pedido: %w", err)
	}

	kafkaHeaders := []kafka.Header{
		{Key: "event_type", Value: []byte("category.created")},
		{Key: "trace_id", Value: []byte(uuid.New().String())},
	}

	err = k.Category.WriteMessages(ctx, kafka.Message{
		Value:   valueBytes,
		Headers: kafkaHeaders,
		Time:    time.Now(),
	})
	if err != nil {
		logger.Error("failed to publish the message category created",
			zap.String("error.message", err.Error()),
			zap.String("error.code", "ERROR_PUBLISH_CREATE_CATEGORY"),
		)
		return err
	}

	return nil
}
