package product_consumer

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

type ProductConsumerTopics struct {
	ProductTopic       *kafka.Reader
	ProductDeleteTopic *kafka.Reader
}

type ProductConsumerTopicsInterface interface {
	//Product
	ConsumerProductTopic(ctx context.Context) (kafka.Message, error)
	CommitProductTopic(ctx context.Context, message kafka.Message) error

	ConsumerProductDeleteTopic(ctx context.Context) (kafka.Message, error)
	CommitProductDeleteTopic(ctx context.Context, message kafka.Message) error


}

func NewProductConsumerTopics(ProductTopic *kafka.Reader, ProductDeleteTopic *kafka.Reader) ProductConsumerTopicsInterface {
	return &ProductConsumerTopics{
		ProductTopic:       ProductTopic,
		ProductDeleteTopic: ProductDeleteTopic,
	}
}

func (k *ProductConsumerTopics) ConsumerProductTopic(ctx context.Context) (kafka.Message, error) {

	kmessage, err := k.ProductTopic.FetchMessage(ctx)
	if err != nil {
		return kafka.Message{}, err
	}

	return kmessage, nil
}

func (k *ProductConsumerTopics) CommitProductTopic(ctx context.Context, message kafka.Message) error {
	err := k.ProductTopic.CommitMessages(ctx, message)
	if err != nil {
		return err
	}
	return nil
}

func (k *ProductConsumerTopics) ConsumerProductDeleteTopic(ctx context.Context) (kafka.Message, error) {
	kmessage, err := k.ProductDeleteTopic.FetchMessage(ctx)
	if err != nil {
		return kafka.Message{}, err
	}

	return kmessage, nil
}

func (k *ProductConsumerTopics) CommitProductDeleteTopic(ctx context.Context, message kafka.Message) error {
	err := k.ProductDeleteTopic.CommitMessages(ctx, message)
	if err != nil {
		return err
	}
	return nil
}