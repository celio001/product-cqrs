package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Configs struct {
	ServiceName    string
	ServiceVersion string
	Env            string
	KafkaCfg       KafkaConfig
	Redis          Redis
	MongoDB        MongoDB
	JaegerConfig   JaegerConfig
}

type KafkaConfig struct {
	KafkaBrokers       []string
	ProductTopic       string
	ProductDeleteTopic string
	ProdctDlqTopic     string
	BrandTopic         string
	BrandDqlTopic      string
	CategoryTopic      string
	CategoryDlqTopic   string
}

type Redis struct {
	ADDR     string
}
type MongoDB struct {
	DSN string
}

type JaegerConfig struct {
	URL string
}

func LoadEnvs() *Configs {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("error loading .env file: %v", err)
	}

	return &Configs{
		ServiceName:    getEnv("SERVICE_NAME", "worker-sync"),
		ServiceVersion: getEnv("SERVICE_VERSION", "1.0.0"),
		Env:            getEnv("ENV", "development"),
		KafkaCfg: KafkaConfig{
			KafkaBrokers:       kafkaGetStrings("KAFKA_BROKERS"),
			ProductTopic:       getEnv("KAFKA_PRODUCT_TOPIC", "product.created"),
			ProductDeleteTopic: getEnv("KAFKA_PRODUCT_DELETE_TOPIC", "product.deleted"),
			ProdctDlqTopic:     getEnv("KAFKA_PRODUCT_DLQ_TOPIC", "product.dlq"),
			BrandTopic:         getEnv("KAFKA_BRAND_TOPIC", "brand.created"),
			BrandDqlTopic:      getEnv("KAFKA_BRAND_DLQ_TOPIC", "brand.dlq"),
			CategoryTopic:      getEnv("KAFKA_CATEGORY_TOPIC", "category.created"),
			CategoryDlqTopic:   getEnv("KAFKA_CATEGORY__DLQ_TOPIC", "category.dlq"),
		},
		Redis: Redis{
			ADDR:     getEnv("REDIS_HOST", "localhost:6379"),
		},
		MongoDB: MongoDB{
			DSN: getEnv("MONGO_DB_DSN", "mongodb://mongo:27017"),
		},
		JaegerConfig: JaegerConfig{
			URL: getEnv("JAEGER_URL", "jaeger:4317"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func kafkaGetStrings(k string) []string {
	v := os.Getenv(k)
	if v == "" {
		return nil
	}
	parts := strings.Split(v, ",")
	var out []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
