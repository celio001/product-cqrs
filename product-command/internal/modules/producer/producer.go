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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type kafkaHeaderCarrier struct {
	headers *[]kafka.Header
}

func (c kafkaHeaderCarrier) Get(key string) string {
	for _, header := range *c.headers {
		if header.Key == key {
			return string(header.Value)
		}
	}
	return ""
}

func (c kafkaHeaderCarrier) Set(key, value string) {
	for i := range *c.headers {
		if (*c.headers)[i].Key == key {
			(*c.headers)[i].Value = []byte(value)
			return
		}
	}
	*c.headers = append(*c.headers, kafka.Header{Key: key, Value: []byte(value)})
}

func (c kafkaHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(*c.headers))
	for _, header := range *c.headers {
		keys = append(keys, header.Key)
	}
	return keys
}

type producerCommand struct {
	ProductTopic *kafka.Writer
	BrandTopic   *kafka.Writer
	Category     *kafka.Writer
	tracer       trace.Tracer
}

type ProducerCommandInterface interface {
	PublishProductCreated(ctx context.Context, p product_dto.CreateProductResponse) error
	PublishProductDeleted(ctx context.Context, id uuid.UUID) error

	PublishBrandCreated(ctx context.Context, b brands.Brand) error
	PublishCategoryCreated(ctx context.Context, c categories.Categories) error
}

func NewProducerCommand(ProductTopic *kafka.Writer, BrandTopic *kafka.Writer, Category *kafka.Writer, tracer trace.Tracer) ProducerCommandInterface {
	return &producerCommand{
		ProductTopic: ProductTopic,
		BrandTopic:   BrandTopic,
		Category:     Category,
		tracer:       tracer,
	}
}

func (k *producerCommand) PublishProductCreated(ctx context.Context, p product_dto.CreateProductResponse) error {

	ctx, span := k.tracer.Start(ctx, "kafka.produce.event-submitted")

	defer span.End()

	span.SetAttributes(
		attribute.String("messaging.system", "kafka"),
		attribute.String("messaging.destination.name", k.ProductTopic.Topic),
		attribute.String("product.id", p.ID.String()),
	)

	valueBytes, err := json.Marshal(p)
	if err != nil {
		span.SetStatus(codes.Error, "failed to serialize product message")
		span.RecordError(err)
		return fmt.Errorf("falha ao serializar envelope do pedido: %w", err)
	}

	kafkaHeaders := []kafka.Header{
		{Key: "event_type", Value: []byte("product.created")},
	}
	// Inject the trace context into the Kafka headers
	otel.GetTextMapPropagator().Inject(ctx, kafkaHeaderCarrier{headers: &kafkaHeaders})

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
		span.SetStatus(codes.Error, "failed to publish product message")
		span.RecordError(err)
		return err
	}

	return nil
}

func (k *producerCommand) PublishProductDeleted(ctx context.Context, id uuid.UUID) error {
	ctx, span := k.tracer.Start(ctx, "kafka.produce.event-submitted")
	
	defer span.End()
	span.SetAttributes(
		attribute.String("messaging.system", "kafka"),
		attribute.String("messaging.destination.name", k.ProductTopic.Topic),
		attribute.String("product.id", id.String()),
	)
	
	kafkaHeaders := []kafka.Header{
		{Key: "event_type", Value: []byte("product.deleted")},
	}
	otel.GetTextMapPropagator().Inject(ctx, kafkaHeaderCarrier{headers: &kafkaHeaders})

	err := k.ProductTopic.WriteMessages(ctx, kafka.Message{
		Headers: kafkaHeaders,
		Time:    time.Now(),
	})
	if err != nil {
		logger.Error("failed to publish the message product deleted",
			zap.String("error.message", err.Error()),
			zap.String("error.code", "ERROR_PUBLISH_DELETE_PRODUCT"),
		)
		span.SetStatus(codes.Error, "failed to publish product message")
		span.RecordError(err)
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
