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
	languageTrigger         = "/language"
	unavailableLangTrigger  = "choiceUnavailableLang"
	questionTrigger         = "/question"
	questionDoneTrigger     = "askQuestionDone"
	undefinedCommandTrigger = "undefinedCommand"
)

var errUndefinedCommand = errors.New("undefined command")

func (p *Processor) processing(chatInfo repository.ChatInfo, text, username string) (err error) {
	errMsg := fmt.Sprintf("processing %s from %s was failed", text, username)
	defer func() { err = e.WrapIfErr(errMsg, err) }()

	log.Printf("got new command '%s' from '%s'\n", text, username)

	if text == repository.LanguageState.Name || text == repository.QuestionState.Name {
		if err = p.storage.EnableCMD(chatInfo.ChatID); err != nil {
			return err
		}
	} else {
		if err = p.storage.DisableCMD(chatInfo.ChatID); err != nil {
			return err
		}
	}

	state, err := p.doCmd(chatInfo, text)

	if err != nil {
		if errors.Is(err, errUndefinedCommand) {
			state = chatInfo.State
		} else {
			return err
		}
	}

	if !(state == chatInfo.State) {
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

func (p *Processor) doCmd(info repository.ChatInfo, text string) (state repository.State, err error) {
	switch info.State.ID {
	case repository.LanguageState.ID:
		state, err = p.changeLang(info, text)
		if err != nil {
			return
		}
	case repository.QuestionState.ID:
		state, err = p.askQuestion(info, text)
		if err != nil {
			return
		}
	default:
		state, err = p.sendMessage(&info, text)
		if err != nil {
			return
		}
	}

	return info.PrevState, nil
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
		_, err := p.sendMessage(&info, unavailableLangTrigger)
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

	state, err := p.sendMessage(&info, info.PrevState.Name)
	if err != nil {
		info.Lang = prev
		return info.State, fmt.Errorf("language selection error: %w", err)
	}

	return state, nil
}

func (p *Processor) askQuestion(info repository.ChatInfo, question string) (repository.State, error) {
	questionInfo := repository.QuestionInfo{
		ChatID:   info.ChatID,
		Question: question,
	}
	err := p.storage.CreateQuestion(&questionInfo)
	if err != nil {
		return info.PrevState, err
	}

	_, err = p.sendMessage(&info, questionDoneTrigger)
	return info.PrevState, nil
}

func (p *Processor) sendMessage(chatInfo *repository.ChatInfo, trigger string) (state repository.State, err error) {
	errMsg := fmt.Sprintf("can't send message \nINFO: trigger'%s' error", trigger)
	defer func() { err = e.WrapIfErr(errMsg, err) }()

	msg, replyMarkup, err := p.getMessageAndReply(chatInfo.PrevState, trigger, chatInfo.Lang)
	log.Println(replyMarkup)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) && trigger != undefinedCommandTrigger {
			_, err = p.sendMessage(chatInfo, undefinedCommandTrigger)
			if err != nil {
				return repository.State{}, err
			}
			return repository.State{}, errUndefinedCommand
		}
		return repository.State{}, fmt.Errorf("can't send message: %w", err)
	}

	filesInfo, err := p.storage.GetFilesInfoOfMessage(msg.ID)
	if err != nil {
		return repository.State{}, err
	}

	if filesInfo == nil || len(filesInfo) == 0 {
		err = p.tg.SendMessage(chatInfo.ChatID, msg.Text, replyMarkup)
	} else {
		if len(filesInfo) == 1 {
			err = p.tg.SendMessageWithFile(chatInfo.ChatID, filesInfo[0], msg.Text, replyMarkup)
		} else {
			err = p.tg.SendMessageWithFiles(chatInfo.ChatID, filesInfo, msg.Text)
		}
	}

	if err != nil {
		return repository.State{}, fmt.Errorf("can't send message: %w", err)
	}

	return msg.State, nil
}

func (p *Processor) getMessageAndReply(prevState repository.State, text string, lang repository.Lang) (*repository.Message, *tgbotapi.ReplyKeyboardMarkup, error) {
	msg, err := p.storage.GetMessage(text, lang, prevState)
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

func (p *Processor) getFiles(filesInfo []*repository.FileInfo) error {
	for i, fileInfo := range filesInfo {
		content, err := p.fileManager.GetFile(fileInfo.Name)
		if err != nil {
			return fmt.Errorf("can't get files %w", err)
		}
		filesInfo[i].Content = content
	}
	return nil
}
