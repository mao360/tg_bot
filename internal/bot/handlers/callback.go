package handlers

import (
	"context"
	"os"
	"strconv"
	"tg_bot/internal/services/payment"
	"tg_bot/internal/services/user"
	"tg_bot/pkg/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackHandler struct {
	userService    *user.Service
	paymentService *payment.Service
	logger         logger.Logger
}

func NewCallbackHandler(userService *user.Service, paymentService *payment.Service, logger logger.Logger) *CallbackHandler {
	return &CallbackHandler{
		userService:    userService,
		paymentService: paymentService,
		logger:         logger,
	}
}

func (h *CallbackHandler) HandleCallback(ctx context.Context, bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) error {
	userID := callbackQuery.From.ID
	data := callbackQuery.Data

	h.logger.Debug("Handling callback", "user_id", userID, "data", data)

	switch data {
	case "mode_reverse":
		h.userService.SetUserState(ctx, userID, "reverse")
		msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "Выбран режим переворота строки. Введите текст для переворота:")
		_, err := bot.Send(msg)
		return err
		
	case "mode_hello":
		h.userService.SetUserState(ctx, userID, "hello")
		msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "Выбран режим 'Hello + строка'. Введите текст:")
		_, err := bot.Send(msg)
		return err
		
	case "buy_premium":
		return h.sendInvoice(bot, callbackQuery.Message.Chat.ID, "premium")
		
	case "buy_pro":
		return h.sendInvoice(bot, callbackQuery.Message.Chat.ID, "pro")
	}

	// Отвечаем на callback query
	callback := tgbotapi.NewCallback(callbackQuery.ID, "")
	_, err := bot.Request(callback)
	return err
}

func (h *CallbackHandler) sendInvoice(bot *tgbotapi.BotAPI, chatID int64, productID string) error {
	product, err := h.paymentService.GetProduct(context.Background(), productID)
	if err != nil {
		h.logger.Error("Failed to get product", "product_id", productID, "error", err)
		msg := tgbotapi.NewMessage(chatID, "❌ Товар не найден.")
		_, err := bot.Send(msg)
		return err
	}

	// Получаем токен провайдера из переменной окружения
	providerToken := os.Getenv("PAYMENT_PROVIDER_TOKEN")
	if providerToken == "" {
		h.logger.Warn("PAYMENT_PROVIDER_TOKEN not set, using test mode")
		providerToken = "TEST_TOKEN"
	}

	// Создаем invoice
	invoice := tgbotapi.NewInvoice(chatID, product.Name, product.Description,
		strconv.FormatInt(product.Price, 10), providerToken, product.Currency,
		"", []tgbotapi.LabeledPrice{{Label: product.Name, Amount: int(product.Price)}})

	// Отправляем invoice
	_, err = bot.Send(invoice)
	if err != nil {
		h.logger.Error("Error sending invoice", "error", err)
		msg := tgbotapi.NewMessage(chatID, "❌ Ошибка создания счета. Попробуйте позже.")
		_, err := bot.Send(msg)
		return err
	}

	return nil
}
