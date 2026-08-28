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
	MongoDB        MongoDB
	JaegerConfig   JaegerConfig
}

type KafkaConfig struct {
	KafkaBrokers     []string
	ProductTopic     string
	ProdctDlqTopic   string
	BrandTopic       string
	BrandDqlTopic    string
	CategoryTopic    string
	CategoryDlqTopic string
}

type MongoDB struct {
	DSN string
}

type JaegerConfig struct {
	URL string
}

func LoadEnvs() *Configs {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return &Configs{
		ServiceName:    os.Getenv("SERVICE_NAME"),
		ServiceVersion: os.Getenv("SERVICE_VERSION"),
		Env:            os.Getenv("ENV"),
		KafkaCfg: KafkaConfig{
			KafkaBrokers:     kafkaGetStrings("KAFKA_BROKERS"),
			ProductTopic:     os.Getenv("KAFKA_PRODUCT_TOPIC"),
			ProdctDlqTopic:   os.Getenv("KAFKA_PRODUCT_DLQ_TOPIC"),
			BrandTopic:       os.Getenv("KAFKA_BRAND_TOPIC"),
			BrandDqlTopic:    os.Getenv("KAFKA_BRAND_DLQ_TOPIC"),
			CategoryTopic:    os.Getenv("KAFKA_CATEGORY_TOPIC"),
			CategoryDlqTopic: os.Getenv("KAFKA_CATEGORY__DLQ_TOPIC"),
		},
		MongoDB: MongoDB{
			DSN: os.Getenv("MONGO_DB_DSN"),
		},
		JaegerConfig: JaegerConfig{
			URL: os.Getenv("JAEGER_URL"),
		},
	}
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
