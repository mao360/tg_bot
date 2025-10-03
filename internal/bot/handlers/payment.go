package handlers

import (
	"context"
	"tg_bot/internal/models"
	"tg_bot/internal/services/payment"
	"tg_bot/internal/services/user"
	"tg_bot/pkg/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PaymentHandler struct {
	paymentService *payment.Service
	userService    *user.Service
	logger         logger.Logger
}

func NewPaymentHandler(paymentService *payment.Service, userService *user.Service, logger logger.Logger) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		userService:    userService,
		logger:         logger,
	}
}

func (h *PaymentHandler) HandlePreCheckout(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.PreCheckoutQuery) error {
	// Создаем модель PreCheckoutQuery
	preCheckoutQuery := &models.PreCheckoutQuery{
		ID:             query.ID,
		TotalAmount:    int64(query.TotalAmount),
		Currency:       query.Currency,
		InvoicePayload: query.InvoicePayload,
	}

	// Валидируем платеж
	if err := h.paymentService.ValidatePreCheckout(ctx, preCheckoutQuery); err != nil {
		h.logger.Error("Pre-checkout validation failed", "error", err)
		answer := tgbotapi.NewCallback(query.ID, "Ошибка валидации платежа")
		_, err := bot.Send(answer)
		return err
	}

	h.logger.Debug("Pre-checkout validated successfully", 
		"amount", query.TotalAmount,
		"currency", query.Currency,
		"payload", query.InvoicePayload)

	// Подтверждаем платеж
	answer := tgbotapi.NewCallback(query.ID, "")
	_, err := bot.Send(answer)
	return err
}

func (h *PaymentHandler) HandleSuccessfulPayment(ctx context.Context, bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	payment := message.SuccessfulPayment

	// Логируем детали платежа
	h.logger.Info("=== PAYMENT RECEIVED ===")
	h.logger.Info("User ID: %d", message.From.ID)
	h.logger.Info("Amount: %d %s", payment.TotalAmount, payment.Currency)
	h.logger.Info("Invoice Payload: %s", payment.InvoicePayload)
	h.logger.Info("Provider Payment Charge ID: %s", payment.ProviderPaymentChargeID)
	h.logger.Info("Telegram Payment Charge ID: %s", payment.TelegramPaymentChargeID)
	h.logger.Info("========================")

	// Создаем модель платежа
	paymentModel := &models.Payment{
		UserID:            message.From.ID,
		Amount:            int64(payment.TotalAmount),
		Currency:          payment.Currency,
		ProductID:         payment.InvoicePayload,
		ProviderChargeID:  payment.ProviderPaymentChargeID,
		TelegramChargeID:  payment.TelegramPaymentChargeID,
		Status:            "completed",
	}

	// Обрабатываем платеж
	if err := h.paymentService.ProcessPayment(ctx, paymentModel); err != nil {
		h.logger.Error("Failed to process payment", "error", err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Ошибка обработки платежа. Обратитесь в поддержку.")
		_, err := bot.Send(msg)
		return err
	}

	// Активируем премиум функции
	if err := h.userService.ActivatePremium(ctx, message.From.ID); err != nil {
		h.logger.Error("Failed to activate premium", "error", err)
	}

	// Отправляем подтверждение пользователю
	msg := tgbotapi.NewMessage(message.Chat.ID,
		"✅ Спасибо за покупку! Ваша подписка активирована.")
	_, err := bot.Send(msg)
	return err
}
