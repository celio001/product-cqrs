package brand_publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/celio001/product-command/internal/modules/brands"
	"github.com/celio001/product-command/pkg/logger"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
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

type brandPublisher struct {
	BrandTopic *kafka.Writer
	tracer     trace.Tracer
}

type BrandPublisherInterface interface {
	PublishBrandCreated(ctx context.Context, b brands.Brand) error
}

func NewBrandPublisher(BrandTopic *kafka.Writer, tracer trace.Tracer) BrandPublisherInterface {
	return &brandPublisher{
		BrandTopic: BrandTopic,
		tracer:     tracer,
	}
}

func (k *brandPublisher) PublishBrandCreated(ctx context.Context, b brands.Brand) error {
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
