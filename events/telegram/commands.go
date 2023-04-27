package telegram

import (
	"database/sql"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
	"tg-bot/lib/e"
	"tg-bot/repository"
)

const (
	choiceUnavailableLangTrigger = "choiceUnavailableLang"
	askQuestionDoneTrigger       = "askQuestionDone"
	undefinedCommandTrigger      = "undefinedCommand"
)

var errUndefinedCommand = errors.New("undefined command")

func (p *Processor) doCmd(chatInfo repository.ChatInfo, text, username string) (err error) {
	errMsg := fmt.Sprintf("can't do cmd %s from %s", text, username)
	defer func() { err = e.WrapIfErr(errMsg, err) }()

	log.Printf("got new command '%s' from '%s'\n", text, username)
	var state repository.State
	switch chatInfo.State {
	case repository.ChangeLangState:
		state, err = p.changeLang(chatInfo, text)
	case repository.AskQuestionState:
		state, err = p.sendMessage(chatInfo.ChatID, askQuestionDoneTrigger, chatInfo.Lang)
	default:
		state, err = p.sendMessage(chatInfo.ChatID, text, chatInfo.Lang)
	}

	if err != nil {
		if errors.Is(err, errUndefinedCommand) {
			state = chatInfo.State
		} else {
			return err
		}
	}

	if !state.Equals(chatInfo.State) {
		prev := chatInfo.State
		chatInfo.State = state
		err = p.storage.ChangeState(chatInfo)
		if err != nil {
			chatInfo.State = prev
			return err
		}
	}
	return nil
}

func (p *Processor) changeLang(info repository.ChatInfo, choice string) (repository.State, error) {
	var newLang repository.Lang
	switch strings.ToLower(choice) {
	case "қазақ":
		newLang = repository.Kz
	case "русский":
		newLang = repository.Ru
	case "english":
		newLang = repository.En
	default:
		_, err := p.sendMessage(info.ChatID, choiceUnavailableLangTrigger, info.Lang)
		if err != nil {
			return info.State, fmt.Errorf("language selection error: %w", err)
		}
		return info.State, nil
	}

	prev := info.Lang
	info.Lang = newLang
	err := p.storage.ChangeLang(info)
	if err != nil {
		info.Lang = prev
		return info.State, err
	}

	state, err := p.sendMessage(info.ChatID, "/start", newLang)
	if err != nil {
		info.Lang = prev
		return info.State, fmt.Errorf("language selection error: %w", err)
	}

	return state, nil
}

func (p *Processor) sendMessage(chatID int64, trigger string, lang repository.Lang) (state repository.State, err error) {
	errMsg := fmt.Sprintf("can't send message \nINFO: trigger'%s' error", trigger)
	defer func() { err = e.WrapIfErr(errMsg, err) }()

	msg, replyMarkup, err := p.getMessageAndReply(trigger, lang)
	log.Println(replyMarkup)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) && trigger != undefinedCommandTrigger {
			_, err = p.sendMessage(chatID, undefinedCommandTrigger, lang)
			if err != nil {
				return repository.DefaultState, err
			}
			return repository.DefaultState, errUndefinedCommand
		}
		return repository.DefaultState, fmt.Errorf("can't send message: %w", err)
	}

	filesInfo, err := p.storage.GetFilesInfoOfMessage(msg.ID)
	if err != nil {
		return repository.DefaultState, err
	}

	if filesInfo == nil || len(filesInfo) == 0 {
		err = p.tg.SendMessage(chatID, msg.Text, replyMarkup)
	} else {
		filesInfo, err = p.getFiles(filesInfo)
		if err != nil {
			return repository.DefaultState, err
		}
		if len(filesInfo) == 1 {
			err = p.tg.SendMessageWithFile(chatID, filesInfo[0], msg.Text, replyMarkup)
		} else {
			err = p.tg.SendMessageWithFiles(chatID, filesInfo, msg.Text)
		}
	}

	if err != nil {
		return repository.DefaultState, fmt.Errorf("can't send message: %w", err)
	}

	return msg.State, nil
}

func (p *Processor) getMessageAndReply(text string, lang repository.Lang) (*repository.MessageWithLang, *tgbotapi.ReplyKeyboardMarkup, error) {
	msg, err := p.storage.GetLangMessage(text, lang)
	if err != nil {
		return nil, nil, fmt.Errorf("can't get message: %w", err)
	}

	replyMarkup, err := p.storage.GetReplyMarkup(msg.ID, lang)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return msg, nil, nil
		}
		return nil, nil, fmt.Errorf("can't get keyboard buttons: %w", err)
	}

	return msg, replyMarkup, nil
}

func (p *Processor) getFiles(filesInfo []repository.FileInfo) ([]repository.FileInfo, error) {
	for i, fileInfo := range filesInfo {
		content, err := p.fileManager.GetFile(fileInfo.Name)
		if err != nil {
			return nil, fmt.Errorf("can't get files %w", err)
		}
		filesInfo[i].Content = content
	}
	return filesInfo, nil
}
