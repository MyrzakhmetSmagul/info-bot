package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tg-bot/repository"
)

type Client interface {
	Updates(offset, limit int) ([]tgbotapi.Update, error)

	SendMessage(chatID int64,
		message string,
		replyMarkup tgbotapi.ReplyKeyboardMarkup) error

	SendMessageWithFile(chatID int64,
		fileInfo repository.File,
		caption string,
		replyMarkup tgbotapi.ReplyKeyboardMarkup) error

	SendMessageWithFiles(chatID int64,
		filesInfo []repository.File, caption string) error
}
