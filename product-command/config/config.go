package config

import (
	"os"
	"strings"
)

var config = map[string]string{
	"ENV":       "development",
	"HTTP_PORT": "8081",

	"SERVICE_NAME":    "product-command",
	"SERVICE_VERSION": "1.0.0",

	//Postgres
	"POSTGRES_DB_DSN": "postgres://postgres:postgres@postgres-main:5432/product?sslmode=disable",
	
	//Kafka
	"KAFKA_BROKERS":               "kafka1:9092",
	"KAFKA_PRODUCT_TOPIC":         "product.created",
	"KAFKA_PRODUCT_DELETED_TOPIC": "product.deleted",
	"KAFKA_BRAND_TOPIC":           "brand.created",
	"KAFKA_CATEGORY_TOPIC":        "category.created",

	//Jaeger
	"JAEGER_URL": "jaeger:4317",
}

func GetString(k string) string {
	v := os.Getenv(k)
	if v == "" {
		return config[k]
	}
	return v
}

func GetStrings(k string) []string {
	v := os.Getenv(k)
	if v == "" {
		v := config[k]
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
