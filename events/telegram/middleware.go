package telegram

import (
	"database/sql"
	"errors"
	"log"
	"strings"
	"tg-bot/lib/e"
	"tg-bot/repository"
)

func (p *Processor) middleware(metaData Meta, text string) (err error) {
	errMsg := "middleware catch error"
	defer func() { err = e.WrapIfErr(errMsg, err) }()
	log.Println("^^^^^^^^^^^^^^^^^^^^\nmiddleware")
	text = strings.TrimSpace(text)
	chatID := metaData.ChatID
	chat, err := p.storage.GetChat(chatID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if errors.Is(err, sql.ErrNoRows) {
		log.Println("hello errors.Is(err, sql.ErrNoRows)")
		chat = repository.Chat{
			ChatID:   chatID,
			Active:   true,
			Lang:     repository.Ru,
			State:    repository.DefaultState,
			MsgGroup: repository.StartMessageGroup,
			CMD:      false,
		}
		if text == "/start" {
			log.Println("/start")
			err = p.storage.CreateChat(chat)
			if err != nil {
				return err
			}

			p.processing(chat, text, metaData.MessageID)
			return err
		}
		err = p.sendMessage(chat, "/unsubscribed")
		return err
	}

	if !chat.Active {
		if text == "start" {
			err = p.storage.EnableChat(chatID)
			if err != nil {
				return err
			}
			err = p.processing(chat, "/start", metaData.MessageID)
			return err
		}
		err = p.sendMessage(chat, "/unsubscribed")
		return err
	}

	return p.processing(chat, text, metaData.MessageID)
}
