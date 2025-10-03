package bot

import (
	"context"
	"tg_bot/internal/bot/handlers"
	"tg_bot/internal/config"
	"tg_bot/internal/services/payment"
	"tg_bot/internal/services/user"
	"tg_bot/pkg/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	config          *config.Config
	logger          logger.Logger
	botAPI          *tgbotapi.BotAPI
	messageHandler  *handlers.MessageHandler
	callbackHandler *handlers.CallbackHandler
	paymentHandler  *handlers.PaymentHandler
}

func New(cfg *config.Config, logger logger.Logger) *Bot {
	// Инициализируем сервисы
	userService := user.NewService(logger)
	paymentService := payment.NewService(logger)

	// Инициализируем обработчики
	messageHandler := handlers.NewMessageHandler(userService, paymentService, logger)
	callbackHandler := handlers.NewCallbackHandler(userService, paymentService, logger)
	paymentHandler := handlers.NewPaymentHandler(paymentService, userService, logger)

	return &Bot{
		config:          cfg,
		logger:          logger,
		messageHandler:  messageHandler,
		callbackHandler: callbackHandler,
		paymentHandler:  paymentHandler,
	}
}

func (b *Bot) Run() error {
	// Создаем бота
	bot, err := tgbotapi.NewBotAPI(b.config.BotToken)
	if err != nil {
		return err
	}

	b.botAPI = bot
	b.botAPI.Debug = b.config.Debug

	b.logger.Info("Bot authorized", "username", bot.Self.UserName)

	// Настраиваем обновления
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	b.logger.Info("Bot started, listening for updates...")

	// Обрабатываем обновления
	for update := range updates {
		ctx := context.Background()

		if update.Message != nil {
			// Обрабатываем успешный платеж
			if update.Message.SuccessfulPayment != nil {
				if err := b.paymentHandler.HandleSuccessfulPayment(ctx, bot, update.Message); err != nil {
					b.logger.Error("Failed to handle successful payment", "error", err)
				}
				continue
			}

			// Обрабатываем обычные сообщения
			if err := b.messageHandler.HandleMessage(ctx, bot, update.Message); err != nil {
				b.logger.Error("Failed to handle message", "error", err)
			}
		} else if update.CallbackQuery != nil {
			if err := b.callbackHandler.HandleCallback(ctx, bot, update.CallbackQuery); err != nil {
				b.logger.Error("Failed to handle callback", "error", err)
			}
		} else if update.PreCheckoutQuery != nil {
			if err := b.paymentHandler.HandlePreCheckout(ctx, bot, update.PreCheckoutQuery); err != nil {
				b.logger.Error("Failed to handle pre-checkout", "error", err)
			}
		}
	}

	return nil
}
