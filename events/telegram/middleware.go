package telegram

import (
	"database/sql"
	"errors"
	"strings"
	"tg-bot/lib/e"
	"tg-bot/repository"
)

func (p *Processor) middleware(chatID int64, text, username string) (err error) {
	errMsg := "middleware catch error"
	defer func() { err = e.WrapIfErr(errMsg, err) }()

	text = strings.TrimSpace(text)
	chatInfo, err := p.storage.ChatInfo(chatID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if errors.Is(err, sql.ErrNoRows) {
		if text == "/start" {
			chatInfo = repository.ChatInfo{
				ChatID: chatID,
				Active: true,
				Lang:   repository.Ru,
				State:  repository.DefaultState,
			}
			err = p.storage.SaveChat(chatInfo)
			if err != nil {
				return err
			}

			_, err = p.sendMessage(chatID, "/start", chatInfo.Lang)
			return err
		}
		_, err = p.sendMessage(chatID, "/unsubscribed", repository.Ru)
		return err
	}

	if !chatInfo.Active {
		if text == "start" {
			err = p.storage.EnableChat(chatID)
			if err != nil {
				return err
			}
			_, err = p.sendMessage(chatID, "/start", chatInfo.Lang)
			return err
		}
		_, err = p.sendMessage(chatID, "/unsubscribed", chatInfo.Lang)
		return err
	}

	return p.doCmd(chatInfo, text, username)
}
