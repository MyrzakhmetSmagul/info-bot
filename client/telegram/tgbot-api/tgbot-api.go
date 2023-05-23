package tgbot_api

import (
	"fmt"
	"log"
	"tg-bot/client/telegram"
	"tg-bot/repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type client struct {
	bot *tgbotapi.BotAPI
}

func New(token string) telegram.Client {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("create new telegram client failed", err)
	}

	return &client{
		bot: bot,
	}
}

func (c client) Updates(offset, limit int) ([]tgbotapi.Update, error) {
	updConf := tgbotapi.UpdateConfig{
		Offset: offset,
		Limit:  limit,
	}
	updates, err := c.bot.GetUpdates(updConf)
	if err != nil {
		return nil, err
	}

	return updates, nil
}

func (c client) SendMessage(chatID int64, message string, replyMarkup tgbotapi.ReplyKeyboardMarkup) error {
	log.Println(chatID)
	msg := tgbotapi.NewMessage(chatID, message)

	var reply interface{}
	if len(replyMarkup.Keyboard) == 0 {
		reply = tgbotapi.ReplyKeyboardRemove{
			RemoveKeyboard: true,
			Selective:      true,
		}
	} else {
		reply = tgbotapi.NewReplyKeyboard(replyMarkup.Keyboard...)
	}

	msg.ReplyMarkup = reply

	_, err := c.bot.Send(msg)

	if err != nil {
		return fmt.Errorf("client.telegram.tgbot-api.SendMessage failed: %w", err)
	}
	return nil
}

func (c client) SendMessageWithFile(chatID int64, fileInfo repository.File, caption string, replyMarkup tgbotapi.ReplyKeyboardMarkup) error {
	msg := make([]tgbotapi.Chattable, 1)
	file := tgbotapi.FileBytes{
		Name:  fileInfo.FileName,
		Bytes: fileInfo.Content,
	}

	var reply interface{}
	if len(replyMarkup.Keyboard) == 0 {
		reply = tgbotapi.ReplyKeyboardRemove{
			RemoveKeyboard: true,
			Selective:      true,
		}
	} else {
		reply = tgbotapi.NewReplyKeyboard(replyMarkup.Keyboard...)
	}

	switch fileInfo.FileType {
	case repository.Photo:
		photo := tgbotapi.NewPhoto(chatID, file)
		photo.Caption = caption
		photo.ReplyMarkup = reply

		msg[0] = photo
	case repository.Video:
		video := tgbotapi.NewVideo(chatID, file)
		video.Caption = caption
		video.ReplyMarkup = reply

		msg[0] = video
	case repository.Document:
		doc := tgbotapi.NewDocument(chatID, file)
		doc.Caption = caption
		doc.ReplyMarkup = reply

		msg[0] = doc
	default:
		return fmt.Errorf("client.telegram.tgbot-api.SendMessageWithFile failed: %w", repository.ErrUndefinedFileType)
	}

	_, err := c.bot.Send(msg[0])
	if err != nil {
		return fmt.Errorf("client.telegram.tgbot-api.SendMessageWithFile failed: %w", err)
	}

	return nil
}

func (c *client) SendMessageWithFiles(chatID int64, filesInfo []repository.File, caption string) error {
	files := make([]interface{}, 0)
	for i, fileInfo := range filesInfo {
		var temp interface{}
		file := tgbotapi.FileBytes{
			Name:  fileInfo.FileName,
			Bytes: fileInfo.Content,
		}
		switch fileInfo.FileType {
		case repository.Photo:
			photo := tgbotapi.NewInputMediaPhoto(file)
			if i == 0 {
				photo.Caption = caption
			}
			temp = photo
		case repository.Video:
			video := tgbotapi.NewInputMediaVideo(file)
			if i == 0 {
				video.Caption = caption
			}
			temp = video
		case repository.Document:
			doc := tgbotapi.NewInputMediaDocument(file)
			if i == 0 {
				doc.Caption = caption
			}
			temp = doc
		default:
			return fmt.Errorf("client.telegram.tgbot-api.SendMessageWithFiles failed: %w", repository.ErrUndefinedFileType)
		}

		files = append(files, temp)
	}

	mediaGroup := tgbotapi.NewMediaGroup(chatID, files)
	if _, err := c.bot.SendMediaGroup(mediaGroup); err != nil {
		return fmt.Errorf("client.telegram.tgbot-api.SendMessageWithFiles failed: %w", err)
	}

	return nil
}
