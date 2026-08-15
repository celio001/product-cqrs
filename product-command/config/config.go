package config

import (
	"os"
	"strings"
)

var config = map[string]string{
	"ENV": "development",

	//Postgres
	//Kafka
	"KAFKA_BROKERS": "localhost:29092",
	"KAFKA_TOPIC":   "product.created",
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
