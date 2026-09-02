package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServiceName    string
	ServiceVersion string
	Env            string
	Port           string
	MongoDB        MongoDB
}

type MongoDB struct {
	DSN string
}

func LoadEnvs() *Config {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("error loading .env file: %v", err)
	}

	return &Config{
		ServiceName:    getEnv("SERVICE_NAME", "product-query"),
		ServiceVersion: getEnv("SERVICE_VERSION", "1.0.0"),
		Env:            getEnv("ENV", "development"),
		Port:           getEnv("PORT", "8082"),
		MongoDB: MongoDB{
			DSN: getEnv("MONGO_DB_DSN", "mongodb://localhost:27017"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
