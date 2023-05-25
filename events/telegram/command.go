package telegram

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"tg-bot/repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	language                     = "/language"
	choiceUnavailableLangTrigger = "choiceUnavailableLang"
	question                     = "/question"
	askQuestionDoneTrigger       = "askQuestionDone"
	undefinedCommandTrigger      = "undefinedCommand"
)

type record struct {
	state        repository.State
	messageGroup repository.MessageGroup
}

var (
	chatCmd             = make(map[int64]string)
	ErrUndefinedCommand = errors.New("undefined command")
)

func (p Processor) processing(chat repository.Chat, text string, messageID int) error {
	if chat.CMD {
		return p.doCmd(chat, text, messageID)
	}
	if text == language || text == question {
		err := p.storage.EnableCmd(chat.ChatID)
		if err != nil {
			return fmt.Errorf("event.telegram.processing failed: %w", err)
		}

		err = p.sendMessage(chat, text)
		if err != nil {
			return fmt.Errorf("event.telegram.processing failed: %w", err)
		}
		chatCmd[chat.ChatID] = text
		return nil
	}

	msgGroup, err := p.getMessageGroup(text, chat.Lang)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("event.telegram.processing failed: %w", err, text)
	} else if !errors.Is(err, sql.ErrNoRows) {
		log.Println(err)
		return p.sendMessage(chat, undefinedCommandTrigger)
	}

	transition, err := p.storage.GetTransition(chat.State.ID, msgGroup.ID)
	if err != nil {
		log.Println(err)
		return p.sendMessage(chat, undefinedCommandTrigger)
	}

	err = p.sendMessageWithData(&chat, msgGroup)
	if err != nil {
		return fmt.Errorf("event.telegram.processing failed: %w", err)
	}

	err = p.storage.ChangeChatState(chat.ChatID, transition.ToState.ID, msgGroup.ID)
	if err != nil {
		return fmt.Errorf("event.telegram.processing failed: %w", err)
	}

	return nil
}

func (p Processor) doCmd(chat repository.Chat, text string, messageID int) error {
	log.Printf("do cmd text: %s", text)
	var err error
	valid := true
	switch chatCmd[chat.ChatID] {
	case language:
		err = p.changeLang(&chat, text)
	case question:
		// TODO something
		fmt.Println(messageID)
	default:
		valid = false
	}

	if err != nil {
		return fmt.Errorf("event.telegram.processing failed: %w", err)
	}

	err = p.storage.DisableCmd(chat.ChatID)
	if err != nil {
		return fmt.Errorf("event.telegram.processing failed: %w", err)
	}

	if valid {
		log.Println("do cmd valid")
		chat.MsgGroup, err = p.getMessageGroupByID(chat.MsgGroup.ID)
		if err != nil {
			return fmt.Errorf("event.telegram.processing failed: %w", err)
		}
		switch chat.Lang {
		case repository.Kz:
			text = chat.MsgGroup.KzMsg.MsgTrigger
		case repository.Ru:
			text = chat.MsgGroup.RuMsg.MsgTrigger
		case repository.En:
			text = chat.MsgGroup.EnMsg.MsgTrigger
		}
	}

	return p.sendMessage(chat, text)
}

func (p Processor) changeLang(chat *repository.Chat, choice string) error {
	var lang string
	switch strings.ToLower(choice) {
	case "қазақ":
		lang = repository.Kz
	case "русский":
		lang = repository.Ru
	case "english":
		lang = repository.En
	default:
		err := p.sendMessage(*chat, choiceUnavailableLangTrigger)
		if err != nil {
			return fmt.Errorf("events.telegram.changeLang: %w", err)
		}
		return nil
	}

	err := p.storage.ChangeChatLang(chat.ChatID, lang)
	if err != nil {
		return fmt.Errorf("events.telegram.changeLang: %w", err)
	}

	chat.Lang = lang
	return nil
}

func (p Processor) getMessageGroup(trigger, lang string) (repository.MessageGroup, error) {
	msg, err := p.storage.GetMessage(trigger, lang)
	if err != nil {
		return repository.MessageGroup{}, fmt.Errorf("events.telegram.getMessageGroup was failed: %w", err)
	}

	msgGroup, err := p.storage.GetMessageGroup(msg.ID, lang)
	if err != nil {
		return repository.MessageGroup{}, fmt.Errorf("events.telegram.getMessageGroup was failed: %w", err)
	}

	msgGroup.KzMsg, err = p.storage.GetMessageByID(msgGroup.KzMsg.ID)
	if err != nil {
		return repository.MessageGroup{}, fmt.Errorf("events.telegram.getMessageGroup was failed: %w", err)
	}
	msgGroup.RuMsg, err = p.storage.GetMessageByID(msgGroup.RuMsg.ID)
	if err != nil {
		return repository.MessageGroup{}, fmt.Errorf("events.telegram.getMessageGroup was failed: %w", err)
	}
	msgGroup.EnMsg, err = p.storage.GetMessageByID(msgGroup.EnMsg.ID)
	if err != nil {
		return repository.MessageGroup{}, fmt.Errorf("events.telegram.getMessageGroup was failed: %w", err)
	}

	return msgGroup, nil
}

func (p Processor) getMessageGroupByID(msgGroupID int) (repository.MessageGroup, error) {
	msgGroup, err := p.storage.GetMessageGroupByID(msgGroupID)
	if err != nil {
		return repository.MessageGroup{}, fmt.Errorf("events.telegram.getMessageGroupByID was failed: %w", err)
	}

	msgGroup.KzMsg, err = p.storage.GetMessageByID(msgGroup.KzMsg.ID)
	if err != nil {
		return repository.MessageGroup{}, fmt.Errorf("events.telegram.getMessageGroupByID was failed: %w", err)
	}
	msgGroup.RuMsg, err = p.storage.GetMessageByID(msgGroup.RuMsg.ID)
	if err != nil {
		return repository.MessageGroup{}, fmt.Errorf("events.telegram.getMessageGroupByID was failed: %w", err)
	}
	msgGroup.EnMsg, err = p.storage.GetMessageByID(msgGroup.EnMsg.ID)
	if err != nil {
		return repository.MessageGroup{}, fmt.Errorf("events.telegram.getMessageGroupByID was failed: %w", err)
	}

	return msgGroup, nil
}

func (p Processor) sendMessageWithData(chat *repository.Chat, msgGroup repository.MessageGroup) error {
	var text string
	switch chat.Lang {
	case repository.Kz:
		text = msgGroup.KzMsg.Text
	case repository.Ru:
		text = msgGroup.RuMsg.Text
	case repository.En:
		text = msgGroup.EnMsg.Text
	}
	rkb, err := p.getReply(msgGroup.ID, chat.Lang)
	if err != nil {
		return fmt.Errorf("events.telegram.sendMessage failed: %w", err)
	}
	log.Println("rkb:", rkb)
	files, err := p.getFilesOfMsgGroup(msgGroup.ID)
	if err != nil {
		return fmt.Errorf("events.telegram.sendMessage failed: %w", err)
	}

	if files == nil || len(files) == 0 {
		err = p.tg.SendMessage(chat.ChatID, text, rkb)
	} else if len(files) == 1 {
		log.Println(len(files))
		err = p.tg.SendMessageWithFile(chat.ChatID, files[0], text, rkb)
	} else {
		log.Println(len(files))

		err = p.tg.SendMessageWithFiles(chat.ChatID, files, text)
	}
	if err != nil {
		return fmt.Errorf("events.telegram.sendMessage was failed: %w", err)
	}
	return nil
}

func (p Processor) sendMessage(chat repository.Chat, trigger string) error {
	msgGroup, err := p.getMessageGroup(trigger, chat.Lang)
	if err != nil {
		return fmt.Errorf("events.telegram.sendMessage failed: %w,\ntrigger: %s", err, trigger)
	}

	var text string
	switch chat.Lang {
	case repository.Kz:
		text = msgGroup.KzMsg.Text
	case repository.Ru:
		text = msgGroup.RuMsg.Text
	case repository.En:
		text = msgGroup.EnMsg.Text
	}
	rkb, err := p.getReply(msgGroup.ID, chat.Lang)
	if err != nil {
		return fmt.Errorf("events.telegram.sendMessage failed: %w", err)
	}
	log.Println("rkb:", rkb)
	files, err := p.getFilesOfMsgGroup(msgGroup.ID)
	if err != nil {
		return fmt.Errorf("events.telegram.sendMessage failed: %w", err)
	}

	if files == nil || len(files) == 0 {
		err = p.tg.SendMessage(chat.ChatID, text, rkb)
	} else if len(files) == 1 {
		log.Println(len(files))
		err = p.tg.SendMessageWithFile(chat.ChatID, files[0], text, rkb)
	} else {
		log.Println(len(files))

		err = p.tg.SendMessageWithFiles(chat.ChatID, files, text)
	}
	if err != nil {
		return fmt.Errorf("events.telegram.sendMessage was failed: %w", err)
	}
	return nil
}

func (p Processor) getReply(msgGroupID int, lang string) (tgbotapi.ReplyKeyboardMarkup, error) {
	rms, err := p.storage.GetReplyMarkupsOfState(msgGroupID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return tgbotapi.ReplyKeyboardMarkup{}, fmt.Errorf("events.telegram.getReply failed: %w", err)
	}
	log.Println("getReply:", rms)
	if rms == nil {
		return tgbotapi.ReplyKeyboardMarkup{}, nil
	}
	rkb := tgbotapi.ReplyKeyboardMarkup{}
	keyboard := []tgbotapi.KeyboardButton{}
	for i := 0; i < len(rms); i++ {
		msgGroup, err := p.storage.GetMessageGroupByID(rms[i].ReplyMessageGroupID)
		if err != nil {
			return tgbotapi.ReplyKeyboardMarkup{}, fmt.Errorf("events.telegram.getReply failed: %w", err)
		}
		var id int
		switch lang {
		case repository.Kz:
			id = msgGroup.KzMsg.ID
		case repository.Ru:
			id = msgGroup.RuMsg.ID
		case repository.En:
			id = msgGroup.EnMsg.ID
		}
		msg, err := p.storage.GetMessageByID(id)
		if err != nil {
			return tgbotapi.ReplyKeyboardMarkup{}, fmt.Errorf("events.telegram.getReply failed: %w", err)
		}
		keyboard = append(keyboard, tgbotapi.KeyboardButton{Text: msg.Text})
		if len(keyboard) == 3 || i+1 == len(rms) {
			rkb.Keyboard = append(rkb.Keyboard, keyboard)
			keyboard = []tgbotapi.KeyboardButton{}
		}
	}

	return rkb, nil
}

func (p *Processor) getFilesOfMsgGroup(msgGroupID int) ([]repository.File, error) {
	files, err := p.storage.GetFilesOfMsgGroup(msgGroupID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("events.getFilesOfMsgGroup failed: %w", err)
	}
	if files == nil {
		return nil, nil
	}
	for i := 0; i < len(files); i++ {
		content, err := p.fileManager.GetFile(files[i].FileName)
		if err != nil {
			return nil, fmt.Errorf("can't get files %w", err)
		}
		files[i].Content = content
	}
	return files, nil
}
