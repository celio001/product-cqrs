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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type productService struct {
	productRepo   product_respository.ProductRepositoryInterface
	consumerTopic consumer.ConsumerTopicsInterface
	productDlq    producer.ProducerDlqInterface
	tracer        trace.Tracer
}

type ProductServiceInterface interface {
	CreateProductSvc(ctx context.Context)
}

func NewProductService(productRepo product_respository.ProductRepositoryInterface, consumerTopic consumer.ConsumerTopicsInterface, productDlq producer.ProducerDlqInterface, tracer trace.Tracer) ProductServiceInterface {
	return &productService{
		productRepo:   productRepo,
		consumerTopic: consumerTopic,
		productDlq:    productDlq,
		tracer:        tracer,
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

		s.processProductMessage(ctx, p)
	}
}

func (s *productService) processProductMessage(ctx context.Context, p kafka.Message) {
	carrier := (*consumer.KafkaHeaderCarrier)(&p.Headers)
	parentCtx := otel.GetTextMapPropagator().Extract(ctx, carrier)

	messageCtx, span := s.tracer.Start(parentCtx, "kafka.consume.product-registered", trace.WithSpanKind(trace.SpanKindConsumer))
	defer span.End()

	err := s.createProductWithRetry(messageCtx, p, 3)
	if err != nil {
		span.SetStatus(codes.Error, "product creation failed")
		span.RecordError(err)

		if dlqErr := s.productDlq.PublishProductDlq(messageCtx, p, err); dlqErr != nil {
			logger.Error("error publish message to DLQ",
				zap.String("error", dlqErr.Error()),
				zap.String("event.action", "ERROR_PUBLISH_MESSAGE_DLQ"))
			return
		}

		err = s.consumerTopic.CommitProductTopic(messageCtx, p)
		if err != nil {
			logger.Error("error commit message",
				zap.String("error", err.Error()),
				zap.String("event.action", "ERROR_COMMIT_MESSAGE"))
		}
		return
	}

	err = s.consumerTopic.CommitProductTopic(messageCtx, p)
	if err != nil {
		logger.Error("error commit message",
			zap.String("error", err.Error()),
			zap.String("event.action", "ERROR_COMMIT_MESSAGE"))
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
