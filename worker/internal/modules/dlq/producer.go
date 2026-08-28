package dlq

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type producerCommand struct {
	productDlq *kafka.Writer
	brandDlq   *kafka.Writer
}

type ProducerDlqInterface interface {
	PublishBrandDlq(ctx context.Context, msg kafka.Message, cause error) error
	PublishProductDlq(ctx context.Context, msg kafka.Message, cause error) error
}

func NewProducerDlq(productDlq *kafka.Writer, brandTopic *kafka.Writer) ProducerDlqInterface {
	return &producerCommand{
		productDlq: productDlq,
		brandDlq:   brandTopic,
	}
}

func (k *producerCommand) PublishBrandDlq(ctx context.Context, msg kafka.Message, cause error) error {
	headers := append(msg.Headers,
		kafka.Header{Key: "x-original-topic", Value: []byte(msg.Topic)},
		kafka.Header{Key: "x-original-partition", Value: fmt.Appendf(nil, "%d", msg.Partition)},
		kafka.Header{Key: "x-original-offset", Value: fmt.Appendf(nil, "%d", msg.Offset)},
		kafka.Header{Key: "x-exception-message", Value: []byte(cause.Error())},
		kafka.Header{Key: "x-failed-at", Value: []byte(time.Now().Format(time.RFC3339))},
	)

	dlqMsg := kafka.Message{
		Key:     msg.Key,
		Value:   msg.Value,
		Headers: headers,
	}

	return k.brandDlq.WriteMessages(ctx, dlqMsg)
}

func (k *producerCommand) PublishProductDlq(ctx context.Context, msg kafka.Message, cause error) error {
	headers := append(msg.Headers,
		kafka.Header{Key: "x-original-topic", Value: []byte(msg.Topic)},
		kafka.Header{Key: "x-original-partition", Value: fmt.Appendf(nil, "%d", msg.Partition)},
		kafka.Header{Key: "x-original-offset", Value: fmt.Appendf(nil, "%d", msg.Offset)},
		kafka.Header{Key: "x-exception-message", Value: []byte(cause.Error())},
		kafka.Header{Key: "x-failed-at", Value: []byte(time.Now().Format(time.RFC3339))},
	)

	dlqMsg := kafka.Message{
		Key:     msg.Key,
		Value:   msg.Value,
		Headers: headers,
	}

	return k.productDlq.WriteMessages(ctx, dlqMsg)
}
