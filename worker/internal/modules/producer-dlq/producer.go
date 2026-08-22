package producer

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type producerCommand struct {
	ProductDlq *kafka.Writer
}

type ProducerDlqInterface interface {
	PublishProductDlq(ctx context.Context, msg kafka.Message, cause error) error
}

func NewProducerDlq(ProductTopic *kafka.Writer) ProducerDlqInterface {
	return &producerCommand{
		ProductDlq: ProductTopic,
	}
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

	return k.ProductDlq.WriteMessages(ctx, dlqMsg)
}
