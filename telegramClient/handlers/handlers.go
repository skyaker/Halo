package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func HandleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	switch update.Message.Text {
	case "/start":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Я помогу отслеживать твои действия.")
		msg.ReplyMarkup = GetMainKeyboard()
		bot.Send(msg)

	case "🍽 Я начал есть":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Окей, записал, что ты начал есть.")
		bot.Send(msg)
		// Тут запись в БД

	case "✅ Я закончил":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Хорошо, записал, что ты закончил.")
		bot.Send(msg)
		// Тут обновление записи в БД
	}
}
