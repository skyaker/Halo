package handlers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func GetMainKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🍽 Я начал есть"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("✅ Я закончил"),
		),
	)
}
