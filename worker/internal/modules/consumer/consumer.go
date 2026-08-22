package consumer

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type consumerTopics struct {
	ProductTopic *kafka.Reader
}

type ConsumerTopicsInterface interface {
	ConsumerProductTopic(ctx context.Context) (kafka.Message, error)
	CommitProductTopic(ctx context.Context, message kafka.Message) error
}

func NewConsumerTopics(ProductTopic *kafka.Reader) ConsumerTopicsInterface {
	return &consumerTopics{
		ProductTopic: ProductTopic,
	}
}

func (k *consumerTopics) ConsumerProductTopic(ctx context.Context) (kafka.Message, error) {

	kmessage, err := k.ProductTopic.FetchMessage(ctx)
	if err != nil {
		return kafka.Message{}, err
	}
	return kmessage, nil
}

func (k *consumerTopics) CommitProductTopic(ctx context.Context, message kafka.Message) error {
	err := k.ProductTopic.CommitMessages(ctx, message)
	if err != nil {
		return err
	}
	return nil
}