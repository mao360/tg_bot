package main

import (
	"log"

	"tg_bot/internal/bot"
	"tg_bot/internal/config"
	"tg_bot/pkg/logger"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	// Инициализируем логгер
	logger := logger.New(cfg.LogLevel)

	// Создаем и запускаем бота
	botApp := bot.New(cfg, *logger)

	if err := botApp.Run(); err != nil {
		log.Fatal("Failed to start bot:", err)
	}
}
