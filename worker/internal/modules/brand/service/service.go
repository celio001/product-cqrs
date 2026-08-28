package brand_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/celio001/product-cqrs/worker/internal/modules/brand"
	brand_repository "github.com/celio001/product-cqrs/worker/internal/modules/brand/repository"
	"github.com/celio001/product-cqrs/worker/internal/modules/consumer"
	"github.com/celio001/product-cqrs/worker/internal/modules/dlq"
	"github.com/celio001/product-cqrs/worker/pkg/logger"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type brandService struct {
	brandRepo     brand_repository.BrandRepositoryInterface
	consumerTopic consumer.ConsumerTopicsInterface
	brandDlq      dlq.ProducerDlqInterface
	tracer        trace.Tracer
}

type BrandServiceInterface interface {
	CreateBrandSvc(ctx context.Context)
}

func NewBrandService(brandRepo brand_repository.BrandRepositoryInterface, consumerTopic consumer.ConsumerTopicsInterface, brandDql dlq.ProducerDlqInterface, tracer trace.Tracer) BrandServiceInterface {
	return &brandService{
		brandRepo:     brandRepo,
		consumerTopic: consumerTopic,
		brandDlq:      brandDql,
		tracer:        tracer,
	}
}

func (s *brandService) CreateBrandSvc(ctx context.Context) {

	for {

		p, err := s.consumerTopic.ConsumerBrandTopic(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				fmt.Println("Goroutine: Stop signal received. Shutting down worker...")
				return
			}

			logger.Info("error connect topic create brand",
				zap.String("error", err.Error()),
				zap.String("event.action", "ERROR_CONNECT_TOPIC_CREATE_BRAND"))
			continue
		}

		s.processBrandMessage(ctx, p)
	}
}

func (s *brandService) processBrandMessage(ctx context.Context, b kafka.Message) {
	carrier := (*consumer.KafkaHeaderCarrier)(&b.Headers)
	parentCtx := otel.GetTextMapPropagator().Extract(ctx, carrier)

	messageCtx, span := s.tracer.Start(parentCtx, "kafka.consume.brand-registered", trace.WithSpanKind(trace.SpanKindConsumer))
	defer span.End()

	err := s.createBrandWithRetry(messageCtx, b, 3)
	if err != nil {
		span.SetStatus(codes.Error, "brand creation failed")
		span.RecordError(err)

		if dlqErr := s.brandDlq.PublishBrandDlq(messageCtx, b, err); dlqErr != nil {
			logger.Error("error publish message to DLQ",
				zap.String("error", dlqErr.Error()),
				zap.String("event.action", "ERROR_PUBLISH_MESSAGE_DLQ"))
			return
		}

		err = s.consumerTopic.CommitBrandTopic(messageCtx, b)
		if err != nil {
			logger.Error("error commit message",
				zap.String("error", err.Error()),
				zap.String("event.action", "ERROR_COMMIT_MESSAGE"))
		}
		return
	}

	err = s.consumerTopic.CommitBrandTopic(messageCtx, b)
	if err != nil {
		logger.Error("error commit message",
			zap.String("error", err.Error()),
			zap.String("event.action", "ERROR_COMMIT_MESSAGE"))
	}

}

func (s *brandService) createBrandWithRetry(ctx context.Context, p kafka.Message, retries int) error {

	var err error
	for i := 1; i <= retries; i++ {

		var b brand.Brand

		err = json.Unmarshal(p.Value, &b)
		if err != nil {
			logger.Info("error connect topic create brand",
				zap.String("error", err.Error()),
				zap.String("event.action", "ERROR_CONNECT_TOPIC_CREATE_BRAND"))
			return err
		}

		err = s.brandRepo.CreateBrandRepository(ctx, b)
		if err == nil {
			logger.Info("brand created successfully",
				zap.String("brand.name", b.Name),
				zap.String("event.action", "BRAND_CREATED_SUCCESSFULLY"))
			return nil
		}

		logger.Error("error create brand",
			zap.String("error", err.Error()),
			zap.String("event.action", "ERROR_CREATE_BRAND"))

		time.Sleep(time.Duration(i) * 500 * time.Millisecond)
	}
	return err
}
