package categories_publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/celio001/product-command/internal/modules/categories"
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

type producerCommand struct {
	Category *kafka.Writer
	tracer   trace.Tracer
}

type CategoriesPublisherInterface interface {
	PublishCategoryCreated(ctx context.Context, c categories.Categories) error
}

func NewCategoryPublisher(Category *kafka.Writer, tracer trace.Tracer) CategoriesPublisherInterface {
	return &producerCommand{
		Category: Category,
		tracer:   tracer,
	}
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
