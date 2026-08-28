package consumer

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type KafkaHeaderCarrier []kafka.Header

func (c *KafkaHeaderCarrier) Get(key string) string {
	for _, h := range *c {
		if h.Key == key {
			return string(h.Value)
		}
	}
	return ""
}

func (c *KafkaHeaderCarrier) Set(key string, value string) {
	*c = append(*c, kafka.Header{
		Key:   key,
		Value: []byte(value),
	})
}

func (c *KafkaHeaderCarrier) Keys() []string {
	keys := make([]string, len(*c))
	for i, h := range *c {
		keys[i] = h.Key
	}
	return keys
}

type consumerTopics struct {
	ProductTopic  *kafka.Reader
	BrandTopic    *kafka.Reader
	CategoryTopic *kafka.Reader
}

type ConsumerTopicsInterface interface {
	//Product
	ConsumerProductTopic(ctx context.Context) (kafka.Message, error)
	CommitProductTopic(ctx context.Context, message kafka.Message) error
	//Brand
	ConsumerBrandTopic(ctx context.Context) (kafka.Message, error)
	CommitBrandTopic(ctx context.Context, message kafka.Message) error
}

func NewConsumerTopics(ProductTopic *kafka.Reader, BrandTopic *kafka.Reader) ConsumerTopicsInterface {
	return &consumerTopics{
		ProductTopic: ProductTopic,
		BrandTopic:   BrandTopic,
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

func (k *consumerTopics) ConsumerBrandTopic(ctx context.Context) (kafka.Message, error) {

	kmessage, err := k.BrandTopic.FetchMessage(ctx)
	if err != nil {
		return kafka.Message{}, err
	}

	return kmessage, nil
}

func (k *consumerTopics) CommitBrandTopic(ctx context.Context, message kafka.Message) error {
	err := k.BrandTopic.CommitMessages(ctx, message)
	if err != nil {
		return err
	}
	return nil
}
