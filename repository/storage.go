package repository

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Repository interface {
	Chat
	MessageKeyboardStorage
}

type Chat interface {
	ChatInfo(chatID int64) (ChatInfo, error)
	SaveChat(info ChatInfo) error
	EnableChat(chatID int64) error
	DisableChat(chatID int64) error
	ChangeLang(info ChatInfo) error
	ChangeState(info ChatInfo) error
	EnableCMD(chatID int64) error
	DisableCMD(chatID int64) error
}

type MessageKeyboardStorage interface {
	GetMessage(trigger string, lang Lang, prevState State) (*Message, error)
	CreateMessage(msg Message) error
	UpdateMessage(msg *Message, text string) (err error)
	GetReplyMarkup(msgID int64, lang Lang) (*tgbotapi.ReplyKeyboardMarkup, error)
	GetKeyboardButtons(messageID int64, lang Lang) ([]tgbotapi.KeyboardButton, error)
	CreateKeyboard(keyboard *Keyboard) error
	CreateKeyboards(keyboards []*Keyboard) error
	CreateReplyMarkup(messageID, keyboardID int64) error
	AddFileForMessage(fileInfo FileInfo) error
	GetFilesInfoOfMessage(messageID int64) ([]FileInfo, error)
}
