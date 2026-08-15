package kafka

import (
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
)

func NewKafkaProducer(brokers []string, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:  kafka.TCP(brokers...),
		Topic: topic,
		Balancer: &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		BatchSize:    100,                   
		BatchTimeout: 50 * time.Millisecond, 
		BatchBytes:   1048576,              
		MaxAttempts:  5,                
		WriteTimeout: 10 * time.Second, 
		ReadTimeout:  10 * time.Second,
		Compression: compress.Snappy,
		Async: false,
	}
}
