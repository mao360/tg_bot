package config

import (
	"os"
	"strconv"
)

type Config struct {
	BotToken             string
	PaymentProviderToken string
	DatabaseURL          string
	LogLevel             string
	Port                 int
	Environment          string
	Debug                bool
}

func Load() *Config {
	return &Config{
		BotToken:             getEnv("TELEGRAM_BOT_TOKEN", ""),
		PaymentProviderToken: getEnv("PAYMENT_PROVIDER_TOKEN", ""),
		DatabaseURL:          getEnv("DATABASE_URL", "postgres://localhost/tg_bot"),
		LogLevel:             getEnv("LOG_LEVEL", "info"),
		Port:                 getEnvAsInt("PORT", 8080),
		Environment:          getEnv("ENVIRONMENT", "development"),
		Debug:                getEnv("DEBUG", "false") == "true",
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
