package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	product_dto "github.com/celio001/product-command/internal/fiber/v1/product/dto"
	"github.com/celio001/product-command/pkg/logger"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type producerCommand struct {
	w *kafka.Writer
}

type ProducerCommandInterface interface {
	PublishProductCreated(ctx context.Context, p product_dto.CreateProductResponse) error
}

func NewProducerCommand(w *kafka.Writer) ProducerCommandInterface {
	return &producerCommand{w: w}
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

	err = k.w.WriteMessages(ctx, kafka.Message{
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
