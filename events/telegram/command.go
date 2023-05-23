package telegram

import (
	"errors"
	"fmt"
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
	chatMeta = make(map[int64][]record)
	chatCmd  = make(map[int64]string)
)

var ErrUndefinedCommand = errors.New("undefined command")

func (p Processor) processing(chat repository.Chat, text string, messageID int) error {
	if chat.CMD {
		return p.doCmd(chat, text, messageID)
	}

	// if text ==

	return nil
}

func (p Processor) doCmd(chat repository.Chat, text string, messageID int) error {
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
		history := chatMeta[chat.ChatID]
		switch chat.Lang {
		case repository.Kz:
			text = history[len(history)-1].messageGroup.KzMsg.MsgTrigger
		case repository.Ru:
			text = history[len(history)-1].messageGroup.RuMsg.MsgTrigger
		case repository.En:
			text = history[len(history)-1].messageGroup.EnMsg.MsgTrigger
		}
	}

	return p.sendMessage(&chat, text)
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
		err := p.sendMessage(chat, choiceUnavailableLangTrigger)
		if err != nil {
			return fmt.Errorf("language selection error: %w", err)
		}
		return nil
	}

	err := p.storage.ChangeChatLang(chat.ChatID, lang)
	if err != nil {
		return err
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

func (p Processor) oneStepBack(chat repository.Chat) error {
	// var trigger string
	// history := chatMeta[chat.ChatID]
	// if history == nil {
	// 	chatMeta[chat.ChatID] = []record{}
	// 	trigger = "main menu"
	// 	err := p.sendMessage(&chat, trigger)
	// 	if err != nil {
	// 		return fmt.Errorf("events.telegram.oneStepBack failed: %w", err)
	// 	}
	// 	msgGroup, err := p.getMessageGroup(trigger, chat.Lang)
	// 	if err != nil {
	// 		return fmt.Errorf("events.telegram.oneStepBack failed: %w", err)
	// 	}
	// 	// state :=
	// 	r := record{}
	// } else if len(history) == 1 {
	// 	switch chat.Lang {
	// 	case repository.Kz:
	// 		trigger = history[0].messageGroup.KzMsg.MsgTrigger
	// 	case repository.Ru:
	// 		trigger = history[0].messageGroup.RuMsg.MsgTrigger
	// 	case repository.En:
	// 		trigger = history[0].messageGroup.EnMsg.MsgTrigger
	// 	}
	// 	err := p.sendMessage(&chat, trigger)
	// 	if err != nil {
	// 		return fmt.Errorf("events.telegram.oneStepBack failed: %w", err)
	// 	}

	// } else if len(history) < 2 {
	// }
	return nil
}

func (p Processor) sendMessage(chat *repository.Chat, trigger string) error {
	return nil
}

func (p Processor) getMessageAndReply(text, lang string) (repository.Message, tgbotapi.ReplyKeyboardMarkup, error) {
	return repository.Message{}, tgbotapi.ReplyKeyboardMarkup{}, nil
}

func (p *Processor) getFiles(filesInfo []repository.File) ([]repository.File, error) {
	for i, file := range filesInfo {
		content, err := p.fileManager.GetFile(file.FileName)
		if err != nil {
			return nil, fmt.Errorf("can't get files %w", err)
		}
		filesInfo[i].Content = content
	}
	return filesInfo, nil
}
