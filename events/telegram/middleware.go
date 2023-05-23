package telegram

import (
	"database/sql"
	"errors"
	"strings"
	"tg-bot/lib/e"
	"tg-bot/repository"
)

func (p *Processor) middleware(metaData Meta, text string) (err error) {
	errMsg := "middleware catch error"
	defer func() { err = e.WrapIfErr(errMsg, err) }()
	text = strings.TrimSpace(text)
	chatID := metaData.ChatID
	chat, err := p.storage.GetChat(chatID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if errors.Is(err, sql.ErrNoRows) {
		if text == "/start" {
			chat = repository.Chat{
				ChatID: chatID,
				Active: true,
				Lang:   repository.Ru,
				State:  repository.DefaultState,
				CMD:    false,
			}
			chatMeta[chatID] = chatInfo{
				cmd:    "",
				record: []record{},
			}
			err = p.storage.CreateChat(chat)
			if err != nil {
				delete(chatMeta, chatID)
				return err
			}

			err = p.sendMessage(&chat, "/start")
			return err
		}
		err = p.sendMessage(&chat, "/unsubscribed")
		return err
	}

	if !chat.Active {
		if text == "start" {
			err = p.storage.EnableChat(chatID)
			if err != nil {
				return err
			}
			err = p.sendMessage(&chat, "/start")
			return err
		}
		err = p.sendMessage(&chat, "/unsubscribed")
		return err
	}

	return p.processing(chat, text, metaData.MessageID)
}
