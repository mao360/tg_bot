package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Состояния пользователей для отслеживания выбранного режима
var userStates = make(map[int64]string)

func main() {
	// Получаем токен из переменной окружения
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is required")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			handleMessage(bot, update.Message)
		} else if update.CallbackQuery != nil {
			handleCallbackQuery(bot, update.CallbackQuery)
		}
	}
}

func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	userID := message.From.ID
	text := message.Text

	// Обработка команд
	if text == "/start" {
		sendMainMenu(bot, message.Chat.ID)
		return
	}

	// Обработка текста в зависимости от выбранного режима
	if state, exists := userStates[userID]; exists {
		var response string
		switch state {
		case "reverse":
			response = reverseString(text)
		case "hello":
			response = "Hello " + text
		}

		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		bot.Send(msg)

		// Показываем главное меню после обработки
		sendMainMenu(bot, message.Chat.ID)
	} else {
		// Если режим не выбран, показываем главное меню
		sendMainMenu(bot, message.Chat.ID)
	}
}

func handleCallbackQuery(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	userID := callbackQuery.From.ID
	data := callbackQuery.Data

	switch data {
	case "mode_reverse":
		userStates[userID] = "reverse"
		msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "Выбран режим переворота строки. Введите текст для переворота:")
		bot.Send(msg)
	case "mode_hello":
		userStates[userID] = "hello"
		msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "Выбран режим 'Hello + строка'. Введите текст:")
		bot.Send(msg)
	}

	// Отвечаем на callback query
	callback := tgbotapi.NewCallback(callbackQuery.ID, "")
	bot.Request(callback)
}

func sendMainMenu(bot *tgbotapi.BotAPI, chatID int64) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Режим 1: Перевернуть строку", "mode_reverse"),
			tgbotapi.NewInlineKeyboardButtonData("Режим 2: Hello + строка", "mode_hello"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, "Выберите режим работы:")
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)

}
