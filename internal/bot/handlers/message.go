package handlers

import (
	"context"
	"tg_bot/internal/models"
	"tg_bot/internal/services/payment"
	"tg_bot/internal/services/user"
	"tg_bot/pkg/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageHandler struct {
	userService    *user.Service
	paymentService *payment.Service
	logger         logger.Logger
}

func NewMessageHandler(userService *user.Service, paymentService *payment.Service, logger logger.Logger) *MessageHandler {
	return &MessageHandler{
		userService:    userService,
		paymentService: paymentService,
		logger:         logger,
	}
}

func (h *MessageHandler) HandleMessage(ctx context.Context, bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	// Получаем или создаем пользователя
	user, err := h.userService.GetOrCreateUser(
		ctx,
		message.From.ID,
		message.From.UserName,
		message.From.FirstName,
		message.From.LastName,
	)
	if err != nil {
		h.logger.Error("Failed to get user", "error", err)
		return err
	}

	// Обрабатываем команды
	switch message.Text {
	case "/start":
		return h.handleStart(bot, message.Chat.ID, user)
	case "/buy":
		return h.handleBuy(bot, message.Chat.ID, user)
	default:
		return h.handleText(bot, message, user)
	}
}

func (h *MessageHandler) handleStart(bot *tgbotapi.BotAPI, chatID int64, user *models.User) error {
	// Отправляем главное меню
	return h.sendMainMenu(bot, chatID, user.IsPremium)
}

func (h *MessageHandler) handleBuy(bot *tgbotapi.BotAPI, chatID int64, user *models.User) error {
	// Отправляем меню покупок
	return h.sendPaymentMenu(bot, chatID)
}

func (h *MessageHandler) handleText(bot *tgbotapi.BotAPI, message *tgbotapi.Message, user *models.User) error {
	// Проверяем состояние пользователя
	state, exists := h.userService.GetUserState(context.Background(), user.TelegramID)
	if !exists {
		// Если режим не выбран, показываем главное меню
		return h.sendMainMenu(bot, message.Chat.ID, user.IsPremium)
	}

	var response string
	switch state {
	case "reverse":
		response = h.reverseString(message.Text)
	case "hello":
		response = "Hello " + message.Text
	default:
		response = "Неизвестный режим"
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	_, err := bot.Send(msg)
	if err != nil {
		h.logger.Error("Failed to send message", "error", err)
		return err
	}

	// Показываем главное меню после обработки
	return h.sendMainMenu(bot, message.Chat.ID, user.IsPremium)
}

func (h *MessageHandler) sendMainMenu(bot *tgbotapi.BotAPI, chatID int64, isPremium bool) error {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Режим 1: Перевернуть строку", "mode_reverse"),
			tgbotapi.NewInlineKeyboardButtonData("Режим 2: Hello + строка", "mode_hello"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💳 Купить подписку", "buy_premium"),
		),
	)

	text := "Выберите режим работы или купите подписку:"
	if isPremium {
		text += "\n\n✨ У вас активна премиум подписка!"
	}

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	_, err := bot.Send(msg)
	return err
}

func (h *MessageHandler) sendPaymentMenu(bot *tgbotapi.BotAPI, chatID int64) error {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💎 Премиум (299₽)", "buy_premium"),
			tgbotapi.NewInlineKeyboardButtonData("🚀 Pro (799₽)", "buy_pro"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, "Выберите подписку:")
	msg.ReplyMarkup = keyboard
	_, err := bot.Send(msg)
	return err
}

func (h *MessageHandler) reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
